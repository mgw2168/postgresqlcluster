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
	sliceutil "github.com/kubesphere/controllers/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kubesphere/api/v1alpha1"
)

var (
	pgClusterFinalizer = "finalizers.radondb.com/pgcluster"
	decUnstructured    = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
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

	if pgCluster.Status.State == "" {
		pgCluster.Status.State = v1alpha1.Creating
		err := r.Status().Update(ctx, pgCluster)
		return ctrl.Result{}, err
	}

	if pgCluster.DeletionTimestamp.IsZero() {
		// The object is not being deleted
		if !sliceutil.HasString(pgCluster.Finalizers, pgClusterFinalizer) {
			pgCluster.Finalizers = append(pgCluster.Finalizers, pgClusterFinalizer)
			err := r.Update(ctx, pgCluster)
			return ctrl.Result{}, err
		}
	} else {
		// the object is not being deleted
		if sliceutil.HasString(pgCluster.Finalizers, pgClusterFinalizer) {
			// todo do delete
			err := r.Update(ctx, pgCluster)
			return reconcile.Result{}, err
		}
	}

	// install postgresql cluster
	if pgCluster.Status.State == v1alpha1.Creating {
		if err := r.installPostgreSQLCluster(ctx, pgCluster); err != nil {
			klog.Error(err.Error())
		}
	} else if pgCluster.Status.Version != pgCluster.Spec.ClientVersion {
		if err := r.deleteCluster(ctx, pgCluster); err != nil {

		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgreSQLClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Client == nil {
		r.Client = mgr.GetClient()
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PostgreSQLCluster{}).
		Complete(r)
}

func (r *PostgreSQLClusterReconciler) installPostgreSQLCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	err = r.createPgCluster(ctx, pg)
	if err != nil {
		klog.Errorf("install pg cluster error: %s", err.Error())
		return err
	}
	return err
}

func getUnstructuredObjStatus(obj *unstructured.Unstructured) string {
	var clusterStatus string
	statusMap, ok := obj.Object["status"].(map[string]interface{})
	if ok {
		clusterStatus, ok = statusMap["state"].(string)
		if ok {
			return clusterStatus
		} else {
			clusterStatus = v1alpha1.ClusterStatusUnknown
		}
	} else {
		clusterStatus = v1alpha1.ClusterStatusUnknown
	}
	return clusterStatus
}
