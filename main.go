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

package main

import (
	"context"
	"flag"
	"os"
	"time"

	sv1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	pgclusterv1alpha1 "github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers"
	"github.com/kubesphere/models"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(pgclusterv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8088", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8089", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "d61cefd2.kubesphere.io",
	})

	err = controllers.Add(mgr)
	if err != nil {
		setupLog.Error(err, "add controller error")
		os.Exit(1)
	}

	if err = (&controllers.PostgreSQLClusterReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "PostgreSQLCluster")
		os.Exit(1)
	}
	//storageclass逻辑
	// create a client for kube resources
	clintset, err := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		os.Exit(1)
	}
	sharedInformers := informers.NewSharedInformerFactory(clintset, time.Minute)
	class := sharedInformers.Storage().V1().StorageClasses()
	informerSc := class.Informer()
	// informerScLister := class.Lister()
	var Pgo models.PgoConfig
	if _, err := Pgo.GetConfig(clintset, "pgo"); err != nil {
		klog.Error(err)
	}
	informerSc.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(*sv1.StorageClass)
			_, err := Pgo.UpdateCm(clintset, "pgo", mObj)
			if err != nil {
				klog.Errorf("update configmap error: %s", err)
			}
			klog.Infof("New StorageClass Added to Store: %s", mObj.Name)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
		},
		DeleteFunc: func(obj interface{}) {
		},
	})
	stopCh := make(chan struct{})
	if err := mgr.Add(manager.RunnableFunc(func(context.Context) error {
		sharedInformers.Start(stopCh)
		sharedInformers.WaitForCacheSync(stopCh)
		return nil
	})); err != nil {
		setupLog.Error(err, "unable to set up sc informer")
		os.Exit(1)
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
