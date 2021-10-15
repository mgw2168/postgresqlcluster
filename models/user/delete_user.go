package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func DeletePgUser(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.DeleteUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	deleteUserReq := &pkg.DeleteUserRequest{
		ClientVersion: pg.Spec.ClientVersion,
		Clusters:      clusterName,
		Namespace:     pg.Spec.Namespace,
		Username:      pg.Spec.Username,
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

	res, ok := pg.Status.Condition[v1alpha1.DeleteUser]
	if ok {
		res.Code = resp.Code
		res.Msg = resp.Msg
	} else {
		pg.Status.Condition = map[string]v1alpha1.ApiResult{
			v1alpha1.DeleteUser: {
				Code: resp.Code,
				Msg:  resp.Msg,
			}}
	}
	return
}
