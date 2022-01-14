package eventhandler

import (
	"github.com/kubesphere/k8sclient"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func WhenStorageClassCreated(createEvent event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
	freshStorageClass := createEvent.Object.(*storagev1.StorageClass)

	var Pgo models.PgoConfig
	_, err := Pgo.GetConfig(k8sclient.GetKubernetesClient(), pkg.PgoNamespace)
	if err != nil {
		klog.Error(err)
	}

	_, err = Pgo.UpdateCm(k8sclient.GetKubernetesClient(), pkg.PgoNamespace, freshStorageClass)
	if err != nil {
		klog.Errorf("update configmap error: %s", err)
	}

	klog.Infof("New StorageClass Added to Store: %s", freshStorageClass.Name)
}
