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

	if !reflect.DeepEqual(newObj.Spec.Users, oldObj.Spec.Users) {
		// create user
		for _, newUser := range newObj.Spec.Users {
			if newUser.UserName == "postgres" {
				continue
			}
			if !models.InSlice(oldObj, newUser.UserName) {
				err = user.CreatePgUser(newObj, newUser.UserName, newUser.Password)
				if err != nil {
					klog.Errorf("create user error: %s", err.Error())
				}
			}
		}

		// delete user
		for _, oldUser := range oldObj.Spec.Users {
			if oldUser.UserName == "postgres" {
				continue
			}
			if !models.InSlice(newObj, oldUser.UserName) {
				err = user.DeletePgUser(newObj, oldUser.UserName)
				if err != nil {
					klog.Errorf("delete user error: %s", err.Error())
				}
			}
			// update user
			for _, newUser := range newObj.Spec.Users {
				if newUser.UserName == "postgres" {
					continue
				}
				if oldUser.UserName == newUser.UserName && oldUser.Password != newUser.Password {
					err = user.UpdatePgUser(newObj, newUser.UserName, newUser.Password)
					if err != nil {
						klog.Errorf("update password error: %s", err.Error())
					}
				}
			}
		}
	}
	// delete user
	if oldObj.Spec.Username != "" && newObj.Spec.Username == "" {
		// delete user
		err = user.DeletePgUser(newObj, oldObj.Spec.Username)
		if err != nil {
			klog.Errorf("delete user error: %s", err.Error())
		}
	}

	// update user passwd
	if oldObj.Spec.Username != "" && oldObj.Spec.Username == newObj.Spec.Username && oldObj.Spec.Password != newObj.Spec.Password {
		// update password
		err = user.UpdatePgUser(newObj, newObj.Spec.Username, newObj.Spec.Password)
		if err != nil {
			klog.Errorf("update password error: %s", err.Error())
		}
	}

	// delete default user and create new user
	if oldObj.Spec.Username != "" && oldObj.Spec.Username != newObj.Spec.Username {
		err = user.DeletePgUser(newObj, oldObj.Spec.Username)
		if err != nil {
			klog.Errorf("delete user error: %s", err.Error())
		}
		// update password
		err = user.CreatePgUser(newObj, newObj.Spec.Username, newObj.Spec.Password)
		if err != nil {
			klog.Errorf("update password error: %s", err.Error())
		}
	}
	return nil
}
