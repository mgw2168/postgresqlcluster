package controllers

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models/user"
	"k8s.io/klog/v2"
)

type clusterUsers struct {
	Users []v1alpha1.User
}

func GetUserList(obj *v1alpha1.PostgreSQLCluster) clusterUsers {
	users := obj.Spec.Users
	if obj.Spec.Username != "" {
		users = append(users, v1alpha1.User{UserName: obj.Spec.Username, Password: obj.Spec.Password})
	}
	return clusterUsers{Users: users}
}

func (in clusterUsers) IsExisted(user v1alpha1.User) (int, bool) {
	for i, u := range in.Users {
		if u.UserName == user.UserName {
			return i, true
		}
	}
	return -1, false
}

func (in clusterUsers) NeedUpdate(index int, user v1alpha1.User) bool {
	if in.Users[index].UserName == user.UserName && in.Users[index].Password != user.Password {
		return true
	}
	return false
}

func doUpdateUsers(oldObj, newObj *v1alpha1.PostgreSQLCluster) error {
	oldUserList := GetUserList(oldObj)
	newUserList := GetUserList(newObj)

	for _, u := range newUserList.Users {
		if u.UserName == "postgres" || u.UserName == "ccp_monitoring" || u.UserName == "primaryuser" {
			continue
		}
		// delete users which are not in Spec.Users/Spec.Username
		if _, existed := oldUserList.IsExisted(u); !existed {
			klog.Infof("user: %s will be created", u.UserName)
			err := user.CreatePgUser(newObj, u.UserName, u.Password)
			if err != nil {
				klog.Errorf("create user error: %s", err.Error())
			}
		}
	}

	// update users in Spec.Users/Spec.Username into PGC
	for _, u := range oldUserList.Users {
		if u.UserName == "postgres" || u.UserName == "ccp_monitoring" || u.UserName == "primaryuser" {
			continue
		}
		index, existed := newUserList.IsExisted(u)
		if !existed {
			klog.Infof("user: %s will be deleted", u.UserName)
			err := user.DeletePgUser(newObj, u.UserName)
			if err != nil {
				klog.Errorf("delete user error: %s", err.Error())
			}
			continue
		}
		if newUserList.NeedUpdate(index, u) {
			klog.Infof("user: %s will be updated", u.UserName)
			err := user.UpdatePgUser(newObj, u.UserName, u.Password)
			if err != nil {
				klog.Errorf("update password error: %s", err.Error())
			}
		}
	}

	return nil
}
