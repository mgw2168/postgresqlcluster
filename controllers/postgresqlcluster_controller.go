/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	v1 "github.com/kubesphere/api/v1"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/eventhandler"
	"github.com/kubesphere/k8sclient"
	"github.com/kubesphere/models/backup"
	"github.com/kubesphere/models/cluster"
	"github.com/kubesphere/pkg"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

type PostgreSQLClusterResource struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// PostgreSQLClusterReconciler reconciles a PostgreSQLCluster object
type PostgreSQLClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	checkTime = 20 * time.Second
)

//+kubebuilder:rbac:groups=pgcluster.radondb.com,resources=postgresqlclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pgcluster.radondb.com,resources=postgresqlclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pgcluster.radondb.com,resources=postgresqlclusters/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=configmaps;secrets;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list

func (r *PostgreSQLClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	pgCluster := &v1alpha1.PostgreSQLCluster{}
	if err := r.Get(ctx, req.NamespacedName, pgCluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if pkg.IsFree(pgCluster.Name, pgCluster.Namespace) {
		klog.Infof("cluster:%s in namespace:%s is updated by reconciler", pgCluster.Name, pgCluster.Namespace)
		return ctrl.Result{RequeueAfter: checkTime}, r.updateState(pgCluster)
	}

	klog.Infof("cluster:%s in namespace:%s is updated by event handler", pgCluster.Name, pgCluster.Namespace)
	return ctrl.Result{RequeueAfter: checkTime}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgreSQLClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	if r.Scheme == nil {
		r.Scheme = mgr.GetScheme()
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(&predicate.Funcs{
			// we only handler event we care about, filter by name and namespace of resource
			UpdateFunc: eventhandler.UpdateEventFilter,
		}).
		For(&v1alpha1.PostgreSQLCluster{}).
		Watches(&source.Kind{Type: &storagev1.StorageClass{}}, handler.Funcs{
			// when a new storage class added, we also need to update it to pgcluster operator
			CreateFunc: eventhandler.WhenStorageClassCreated,
		}).
		Watches(&source.Kind{Type: &corev1.ConfigMap{}}, handler.Funcs{
			UpdateFunc: eventhandler.WhenConfigMapUpdated,
		}).
		Watches(&source.Kind{Type: &v1alpha1.PostgreSQLCluster{}}, handler.Funcs{
			CreateFunc: func(createEvent event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
				pg := createEvent.Object.(*v1alpha1.PostgreSQLCluster)

				pkg.Lock(pg.Name, pg.Namespace)
				defer pkg.UnLock(pg.Name, pg.Namespace)

				if pg.Status.State == "" {
					if err := cluster.CreatePgCluster(pg); err != nil {
						klog.Errorf("create Pgcluster resource error: %s", err)
					}
					if err := r.updateState(pg); err != nil {
						klog.Errorf("update PostgreSQLCluster state error: %s", err)
					}
				}
			},

			UpdateFunc: func(updateEvent event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
				oldCluster := updateEvent.ObjectOld.(*v1alpha1.PostgreSQLCluster)
				newCluster := updateEvent.ObjectNew.(*v1alpha1.PostgreSQLCluster)

				pkg.Lock(newCluster.Name, newCluster.Namespace)
				defer pkg.UnLock(newCluster.Name, newCluster.Namespace)

				err := doUpdateCluster(oldCluster, newCluster)
				if err != nil {
					klog.Errorf("update cluster error: %s", err)
				}
				if err = r.updateState(newCluster); err != nil {
					klog.Errorf("update PostgreSQLCluster state error: %s", err)
				}
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
				pg := deleteEvent.Object.(*v1alpha1.PostgreSQLCluster)
				err := cluster.DeletePgCluster(pg)
				if err != nil {
					klog.Errorf("delete cluster error: %s", err)
				}
			},
		}).
		Complete(r)
}

func (r *PostgreSQLClusterReconciler) isInBackup(pg *v1alpha1.PostgreSQLCluster) bool {
	k8s := k8sclient.GetKubernetesClient()
	backupJobName := "backrest-backup-%s"
	scheduleBackupJobName := "%s-full-sch-backup"

	job, err := k8s.BatchV1().Jobs(pg.Namespace).Get(context.TODO(), fmt.Sprintf(backupJobName, pg.Name), metav1.GetOptions{})
	if err != nil {
		return false
	}

	if job.Status.Active > 0 {
		return true
	}

	job, err = k8s.BatchV1().Jobs(pg.Namespace).Get(context.TODO(), fmt.Sprintf(scheduleBackupJobName, pg.Name), metav1.GetOptions{})
	if err != nil {
		return false
	}
	if job.Status.Active > 0 {
		return true
	}

	return false
}

func (r *PostgreSQLClusterReconciler) updateState(pg *v1alpha1.PostgreSQLCluster) error {
	pgc := &v1.Pgcluster{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: pg.Namespace,
		Name:      pg.Name,
	}, pgc)
	if err != nil {
		klog.Errorf("get pgcluster resource error: %s", err)
	}

	if pgc.Status.State != "" {
		pg.Status.State = string(pgc.Status.State)

		if pgc.Status.State != StatusBootstrapped && pgc.Status.State != StatusProcessed {
			if r.isInBackup(pg) {
				pg.Status.State = StatusInBackup
			}
		}
	}

	backup.ShowBackup(pg)

	return r.Status().Update(context.TODO(), pg)
}
