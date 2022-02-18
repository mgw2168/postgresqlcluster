package controllers

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models/cluster"
	"k8s.io/klog/v2"
)

func doUpdateCluster(oldObj, newObj *v1alpha1.PostgreSQLCluster) (err error) {
	// update pvc
	if oldObj.Spec.PVCSize != newObj.Spec.PVCSize && oldObj.Spec.PVCSize != "" {
		err = cluster.UpdatePgCluster(newObj, true)
		if err != nil {
			klog.Errorf("update pvc error: %s", err)
		}
	}

	// update cpu and memory
	if oldObj.Spec.CPURequest != newObj.Spec.CPURequest || oldObj.Spec.CPULimit != newObj.Spec.CPULimit ||
		oldObj.Spec.MemoryLimit != newObj.Spec.MemoryLimit ||
		oldObj.Spec.MemoryRequest != newObj.Spec.MemoryRequest {
		err = cluster.UpdatePgCluster(newObj, false)
		if err != nil {
			klog.Errorf("update cpu and memory error: %s", err.Error())
		}
	}

	// scale up
	if oldObj.Spec.ReplicaCount != newObj.Spec.ReplicaCount && newObj.Spec.ReplicaCount > oldObj.Spec.ReplicaCount {
		replicaCount := newObj.Spec.ReplicaCount - oldObj.Spec.ReplicaCount
		err = cluster.ScaleUpPgCluster(newObj, replicaCount)
		if err != nil {
			klog.Errorf("scale up error: %s", err.Error())
		}
	}

	// scale down
	if oldObj.Spec.ReplicaName != newObj.Spec.ReplicaName && newObj.Spec.ReplicaName != "" {
		err = cluster.ScaleDownPgCluster(newObj)
		if err != nil {
			klog.Errorf("scale down error: %s", err.Error())
		}
	}

	// restart
	if oldObj.Spec.Restart {
		err = cluster.RestartCluster(newObj)
		if err != nil {
			klog.Errorf("restart cluster error: %s", err.Error())
		}
	}

	// update users
	return doUpdateUsers(oldObj, newObj)
}
