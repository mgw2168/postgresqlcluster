package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func GetExistedUserList(pg *v1alpha1.PostgreSQLCluster) (pkg.ShowUserResponse, error) {
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
		return resp, err
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return resp, err
	}

	return resp, nil
}
