package eventhandler

import (
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func UpdateEventFilter(updateEvent event.UpdateEvent) bool {
	cm, ok := updateEvent.ObjectNew.(*corev1.ConfigMap)
	if ok {
		if cm.GetNamespace() != pkg.PgoNamespace || cm.GetName() != models.CustomConfigMapName {
			return false
		}
		klog.Infof("cm:%s change event we care about", cm.GetName())
	}
	return true
}
