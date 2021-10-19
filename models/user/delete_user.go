package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func DeletePgUser(pg *v1alpha1.PostgreSQLCluster, username string) (err error) {
	var resp pkg.DeleteUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	deleteUserReq := &pkg.DeleteUserRequest{
		ClientVersion: "4.7.1",
		Clusters:      clusterName,
		Namespace:     pg.Spec.Namespace,
		Username:      username,
	}
	respByte, err := pkg.Call("POST", pkg.DeleteUserPath, deleteUserReq)
	if err != nil {
		klog.Errorf("call delete user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}
	if resp.Code == pkg.Ok {
		pg.Status.State = v1alpha1.Success
	} else {
		pg.Status.State = v1alpha1.Failed
	}

	flag := true
	for _, res := range pg.Status.Condition {
		if res.Api == v1alpha1.DeleteUser {
			flag = false
			res.Code = resp.Code
			res.Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  v1alpha1.DeleteUser,
			Code: resp.Code,
			Msg:  resp.Msg,
			Data: "",
		})
	}
	return
}
