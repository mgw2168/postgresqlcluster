package k8sclient

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

var k8sClient kubernetes.Interface

func init() {
	var err error
	k8sClient, err = kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		klog.Error(err, "setup k8s client failed")
		os.Exit(-1)
	}
}

func GetKubernetesClient() kubernetes.Interface {
	return k8sClient
}
