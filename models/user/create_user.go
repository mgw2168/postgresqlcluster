package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func CreatePgUser(pg *v1alpha1.PostgreSQLCluster, username, passwd string, isSuperUser bool) (err error) {
	var resp pkg.CreateUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	createUserReq := &pkg.CreateUserRequest{
		Clusters:        clusterName,
		ClientVersion:   pkg.ClientVersion,
		ManagedUser:     true,
		Namespace:       pg.Spec.Namespace,
		Username:        username,
		Password:        passwd,
		Superuser:       isSuperUser,
		PasswordAgeDays: 86400,
		PasswordLength:  8,
		PasswordType:    "md5",
	}
	klog.Infof("params: %+v", createUserReq)
	respByte, err := pkg.Call("POST", pkg.CreateUserPath, createUserReq)
	if err != nil {
		klog.Errorf("call create user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	models.MergeCondition(pg, pkg.CreateUser, resp.Status)

	return
}
