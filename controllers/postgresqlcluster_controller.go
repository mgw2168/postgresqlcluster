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
	"github.com/kubesphere/eventhandler"
	"github.com/kubesphere/models/cluster"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (r *PostgreSQLClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.Infof("start reconcile postgresqlcluster.")

	_ = log.FromContext(ctx)
	fakePGC := v1alpha1.PostgreSQLCluster{}
	if err := r.Get(ctx, req.NamespacedName, &fakePGC); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// when postgresqlcluster is deleted
	if !fakePGC.ObjectMeta.DeletionTimestamp.IsZero() {
		//delete cluster
		if err := cluster.DeletePgCluster(&fakePGC); err != nil {
			return ctrl.Result{}, err
		}
		fakePGC.ObjectMeta.Finalizers = []string{}
		return ctrl.Result{}, r.Client.Update(ctx, &fakePGC)
	}

	// create cluster
	if fakePGC.ObjectMeta.DeletionTimestamp.IsZero() && len(fakePGC.ObjectMeta.Finalizers) == 0 {
		if err := cluster.CreatePgCluster(&fakePGC); err != nil {
			klog.Errorf("create Pgcluster resource error: %s", err)
			return ctrl.Result{}, err
		}
		fakePGC.ObjectMeta.Finalizers = append(fakePGC.ObjectMeta.Finalizers, "finalizers.postgresqlcluster.radondb.com/created")
		return r.updateSpecThenReturn(ctx, &fakePGC)
	}

	// waiting cluster initialized done
	if fakePGC.ObjectMeta.DeletionTimestamp.IsZero() && len(fakePGC.ObjectMeta.Finalizers) == 1 && fakePGC.Status.State != "pgcluster Initialized" {
		return r.fetchStateThenReturn(ctx, &fakePGC)
	}

	// when postgresqlcluster is ready
	if fakePGC.ObjectMeta.DeletionTimestamp.IsZero() && len(fakePGC.ObjectMeta.Finalizers) == 1 && fakePGC.Status.State == "pgcluster Initialized" {
		fakePGC.ObjectMeta.Finalizers = append(fakePGC.ObjectMeta.Finalizers, "finalizers.postgresqlcluster.radondb.com/initialized")
		return r.updateSpecThenReturn(ctx, &fakePGC)
	}

	// handler cpu/mem/replica change
	shouldBreakThisTerm, err := r.compareWithPGCluster(ctx, &fakePGC)
	if err != nil {
		return ctrl.Result{}, err
	}
	if shouldBreakThisTerm {
		return r.updateSpecThenReturn(ctx, &fakePGC)
	}

	// handle user info change
	if err := r.updateUserList(&fakePGC); err != nil {
		return ctrl.Result{}, err
	}

	return r.fetchStateThenReturn(ctx, &fakePGC)
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
		Complete(r)
}
