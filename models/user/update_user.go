package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func UpdatePgUser(pg *v1alpha1.PostgreSQLCluster, username, passwd string, isSuperUser bool) (err error) {
	var resp pkg.UpdateUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	updateUserReq := &pkg.UpdateUserRequest{
		Clusters:                 clusterName,
		ClientVersion:            pkg.ClientVersion,
		Namespace:                pg.Spec.Namespace,
		Username:                 username,
		Password:                 passwd,
		PasswordAgeDays:          86400,
		PasswordLength:           8,
		Superuser:                isSuperUser,
		PasswordType:             "md5",
		SetSystemAccountPassword: false,
	}
	klog.Infof("params: %+v", updateUserReq)
	respByte, err := pkg.Call("POST", pkg.UpdateUserPath, updateUserReq)
	if err != nil {
		klog.Errorf("call update user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	models.MergeCondition(pg, pkg.UpdateUser, resp.Status)

	return
}
