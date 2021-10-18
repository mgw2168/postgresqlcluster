package controllers

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models/cluster"
	"github.com/kubesphere/models/user"
	"k8s.io/klog/v2"
)

func doUpdateCluster(oldObj, newObj *v1alpha1.PostgreSQLCluster) (err error) {
	// update pvc
	if oldObj.Spec.PVCSize != newObj.Spec.PVCSize {
		err = cluster.UpdatePgCluster(newObj)
		if err != nil {
			klog.Errorf("update pvc error: %s", err)
		}
	}
	// update cpu and memory
	if oldObj.Spec.CPURequest != newObj.Spec.CPURequest || oldObj.Spec.CPULimit != newObj.Spec.CPULimit ||
		oldObj.Spec.MemoryLimit != newObj.Spec.MemoryLimit ||
		oldObj.Spec.MemoryRequest != newObj.Spec.MemoryRequest {
		err = cluster.UpdatePgCluster(newObj)
		if err != nil {
			klog.Errorf("update cpu and memory error: %s", err.Error())
		}
	}
	// scale up
	if oldObj.Spec.ReplicaCount != newObj.Spec.ReplicaCount && newObj.Spec.ReplicaCount > oldObj.Spec.ReplicaCount {
		err = cluster.ScaleUpPgCluster(newObj)
		if err != nil {
			klog.Errorf("scale up error: %s", err.Error())
		}
	}

	// scale down
	if oldObj.Spec.ReplicaName != newObj.Spec.ReplicaName {
		err = cluster.ScaleDownPgCluster(newObj)
		if err != nil {
			klog.Errorf("scale down error: %s", err.Error())
		}
	}

	// restart todo TODO
	if oldObj.Spec.Restart {
		err = cluster.RestartCluster(newObj)
		if err != nil {
			klog.Errorf("restart cluster error: %s", err.Error())
		}
	}

	// create user
	//if oldObj.Spec.Username != newObj.Spec.Username || oldObj.Spec.Password != newObj.Spec.Password && newObj.Spec.Password != "" {
	if len(newObj.Spec.Users) > 0 {
		err = user.CreatePgUser(newObj)
		if err != nil {
			klog.Errorf("create pg user error: %s", err.Error())
		}
	}

	// delete user
	if oldObj.Spec.Username != newObj.Spec.Username && newObj.Spec.Password == "" {
		err = user.DeletePgUser(newObj)
		if err != nil {
			klog.Errorf("update password error: %s", err.Error())
		}
	}

	// update user password
	if oldObj.Spec.Username == newObj.Spec.Username && oldObj.Spec.Password != newObj.Spec.Password {
		err = user.UpdatePgUser(newObj)
		if err != nil {
			klog.Errorf("update password error: %s", err.Error())
		}
	}
	// list user
	if oldObj.Spec.Username != newObj.Spec.Username && newObj.Spec.Password == "" {
		err = user.ListPgUser(newObj)
		if err != nil {
			klog.Errorf("update password error: %s", err.Error())
		}
	}
	return nil
}
