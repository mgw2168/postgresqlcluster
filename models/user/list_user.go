package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func ListPgUser(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.ShowUserResponse
	listUserReq := &pkg.ShowUserRequest{
		Clusters:           pg.Spec.ClusterName,
		ClientVersion:      pg.Spec.ClientVersion,
		Namespace:          pg.Spec.Namespace,
		ShowSystemAccounts: pg.Spec.ShowSystemAccounts,
	}
	respByte, err := pkg.Call("POST", pkg.ShowUserPath, listUserReq)
	if err != nil {
		klog.Errorf("call list user error: %s", err.Error())
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

	res, ok := pg.Status.Condition[v1alpha1.ShowUser]
	if ok {
		res.Code = resp.Code
		res.Msg = resp.Msg
	} else {
		pg.Status.Condition = map[string]v1alpha1.ApiResult{
			v1alpha1.ShowUser: {
				Code: resp.Code,
				Msg:  resp.Msg,
			}}
	}

	return
}
