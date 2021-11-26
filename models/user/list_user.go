package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func ListPgUser(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.ShowUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	listUserReq := &pkg.ShowUserRequest{
		Clusters:           clusterName,
		ClientVersion:      "4.7.1",
		Namespace:          pg.Spec.Namespace,
		ShowSystemAccounts: pg.Spec.ShowSystemAccounts,
	}
	klog.Infof("params: %+v", listUserReq)
	respByte, err := pkg.Call("POST", pkg.ShowUserPath, listUserReq)
	if err != nil {
		klog.Errorf("call list user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.ShowUser {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.ShowUser,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}
	return
}
