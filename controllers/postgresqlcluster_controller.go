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
	"github.com/kubesphere/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PostgreSQLClusterResource struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func AppendChannel() {
	for {
		res := <-PgClusterResourceChan
		fmt.Println("======= todo delete resource in channel.....")
		fmt.Println(res.Name, res.Namespace)
	}
}

var (
	PgClusterResourceChan chan PostgreSQLClusterResource
)

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

	if pgCluster.Status.PostgreSQLClusterState == "" {
		pgCluster.Status.PostgreSQLClusterState = v1alpha1.Creating
		err := r.Status().Update(ctx, pgCluster)
		return reconcile.Result{}, err
	}

	if pgCluster.Status.PostgreSQLClusterState == v1alpha1.Creating {
		if pgCluster.Spec.Action == v1alpha1.CreateCluster {
			if err := r.createPgCluster(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.UpdateCluster {
			if err := r.updatePgCluster(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.ScaleCluster {
			if err := r.scalePgCluster(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.ScaleDownCluster {
			if err := r.scaleDownPgCluster(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.DeleteCluster {
			if err := r.deletePgCluster(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.RestartCluster {
			if err := r.restartCluster(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.CreateUser {
			if err := r.createPgUser(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.UpdateUser {
			if err := r.updatePgUser(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.DeleteUser {
			if err := r.DeletePgUser(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		} else if pgCluster.Spec.Action == v1alpha1.ShowUser {
			if err := r.listPgUser(ctx, pgCluster); err != nil {
				klog.Error(err.Error())
				return ctrl.Result{}, err
			}
		}
		PgClusterResourceChan <- PostgreSQLClusterResource{
			Namespace: pgCluster.Spec.Namespace,
			Name:      pgCluster.Spec.Name,
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgreSQLClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// 定义个channel用来存储
	PgClusterResourceChan = make(chan PostgreSQLClusterResource, 10)
	go AppendChannel()
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PostgreSQLCluster{}).
		Complete(r)
}
