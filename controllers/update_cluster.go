package controllers

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/models/cluster"
	"github.com/kubesphere/models/user"
	"k8s.io/klog/v2"
	"reflect"
)

func doUpdateCluster(oldObj, newObj *v1alpha1.PostgreSQLCluster) (err error) {
	// update pvc
	if oldObj.Spec.PVCSize != newObj.Spec.PVCSize && oldObj.Spec.PVCSize != "" {
		err = cluster.UpdatePgCluster(newObj, true)
		if err != nil {
			klog.Errorf("update pvc error: %s", err)
		}
		klog.Error("update pvc error===")
	}
	// update cpu and memory
	if oldObj.Spec.CPURequest != newObj.Spec.CPURequest || oldObj.Spec.CPULimit != newObj.Spec.CPULimit ||
		oldObj.Spec.MemoryLimit != newObj.Spec.MemoryLimit ||
		oldObj.Spec.MemoryRequest != newObj.Spec.MemoryRequest {
		err = cluster.UpdatePgCluster(newObj, false)
		if err != nil {
			klog.Errorf("update cpu and memory error: %s", err.Error())
		}
		klog.Error("update cpu and memory error")
	}
	// scale up
	if oldObj.Spec.ReplicaCount != newObj.Spec.ReplicaCount && newObj.Spec.ReplicaCount > oldObj.Spec.ReplicaCount {
		err = cluster.ScaleUpPgCluster(newObj)
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

	if len(newObj.Spec.Users) > 0 && !reflect.DeepEqual(newObj.Spec.Users, oldObj.Spec.Users) {
		// create user
		for _, newUser := range newObj.Spec.Users {
			if !models.InSlice(oldObj, newUser.UserName) {
				err = user.CreatePgUser(newObj, newUser.UserName, newUser.Password)
				if err != nil {
					klog.Errorf("create user error: %s", err.Error())
				}
			}
		}

		// delete user
		for _, oldUser := range oldObj.Spec.Users {
			if !models.InSlice(newObj, oldUser.UserName) {
				err = user.DeletePgUser(newObj, oldUser.UserName)
				if err != nil {
					klog.Errorf("delete user error: %s", err.Error())
				}
			}
			// update user
			for _, newUser := range newObj.Spec.Users {
				if oldUser.UserName == newUser.UserName && oldUser.Password != newUser.Password {
					err = user.UpdatePgUser(newObj, newUser.UserName, newUser.Password)
					if err != nil {
						klog.Errorf("update password error: %s", err.Error())
					}
				}
			}
		}
	}

	// list user
	if oldObj.Spec.Username != newObj.Spec.Username && newObj.Spec.Password == "" {
		err = user.ListPgUser(newObj)
		if err != nil {
			klog.Errorf("list user error: %s", err.Error())
		}
	}
	return nil
}
