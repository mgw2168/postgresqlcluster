package eventhandler

import (
	"github.com/kubesphere/k8sclient"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func WhenConfigMapUpdated(updateEvent event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
	cm := updateEvent.ObjectNew.(*corev1.ConfigMap)

	var Pgo models.PgoConfig
	_, err := Pgo.GetConfig(k8sclient.GetKubernetesClient(), pkg.PgoNamespace)
	if err != nil {
		klog.Error(err)
	}

	_, err = Pgo.UpdateCmInformer(k8sclient.GetKubernetesClient(), pkg.PgoNamespace, cm)
	if err != nil {
		klog.Errorf("update configmap error: ", err)
	}
	klog.Infof("update cm %s success!", models.CustomConfigMapName)
}
