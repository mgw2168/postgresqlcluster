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
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models/cluster"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
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

//+kubebuilder:rbac:groups=pgcluster.kubesphere.io,resources=postgresqlclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pgcluster.kubesphere.io,resources=postgresqlclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pgcluster.kubesphere.io,resources=postgresqlclusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PostgreSQLCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *PostgreSQLClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	pgCluster := &v1alpha1.PostgreSQLCluster{}
	if err := r.Get(ctx, req.NamespacedName, pgCluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return ctrl.Result{}, nil
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
		For(&v1alpha1.PostgreSQLCluster{}).
		Complete(r)
}

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("postgresqlCluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	reconcileObj := r.(*PostgreSQLClusterReconciler)
	// Watch for changes to PostgreSQLCluster
	err = c.Watch(&source.Kind{Type: &v1alpha1.PostgreSQLCluster{}}, &handler.Funcs{
		CreateFunc: func(event event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
			pg := event.Object.(*v1alpha1.PostgreSQLCluster)
			if pg.Status.State == "" {
				if err := cluster.CreatePgCluster(pg); err != nil {
					klog.Errorf("create Pgcluster resource error: %s", err)
				}
				err = reconcileObj.Status().Update(context.TODO(), pg)
				if err != nil {
					if errors.IsConflict(err) {
						return
					}
					klog.Errorf("update Pgcluster status error: %s", err)
				}
			}
		},

		UpdateFunc: func(updateEvent event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
			oldCluster := updateEvent.ObjectOld.(*v1alpha1.PostgreSQLCluster)
			newCluster := updateEvent.ObjectNew.(*v1alpha1.PostgreSQLCluster)

			err := doUpdateCluster(oldCluster, newCluster)
			if err != nil {
				klog.Errorf("update cluster error: %s", err)
			}

			err = reconcileObj.Status().Update(context.TODO(), newCluster)
			if err != nil {

				klog.Errorf("update Pgcluster status error: %s", err)
			}
		},
		DeleteFunc: func(deleteEvent event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
			pg := deleteEvent.Object.(*v1alpha1.PostgreSQLCluster)
			err := cluster.DeletePgCluster(pg)
			if err != nil {
				klog.Errorf("delete cluster error: %s", err)
			}
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &PostgreSQLClusterReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
}
