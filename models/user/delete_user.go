package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func DeletePgUser(pg *v1alpha1.PostgreSQLCluster, username string) (err error) {
	var resp pkg.DeleteUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	deleteUserReq := &pkg.DeleteUserRequest{
		ClientVersion: pkg.ClientVersion,
		Clusters:      clusterName,
		Namespace:     pg.Spec.Namespace,
		Username:      username,
	}
	klog.Infof("params: %+v", deleteUserReq)
	respByte, err := pkg.Call("POST", pkg.DeleteUserPath, deleteUserReq)
	if err != nil {
		klog.Errorf("call delete user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	models.MergeCondition(pg, pkg.DeleteUser, resp.Status)

	return
}
