package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func DeletePgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.DeleteClusterResponse
	clusterReq := &pkg.DeleteClusterRequest{
		Clustername:   pg.Spec.Name,
		ClientVersion: "4.7.1",
		Namespace:     pg.Spec.Namespace,
		DeleteBackups: false,
		DeleteData:    false,
	}
	klog.Infof("params: %+v", clusterReq)
	respByte, err := pkg.Call("POST", pkg.DeleteClusterPath, clusterReq)
	if err != nil {
		klog.Errorf("call delete cluster error: %s", err.Error())
		return
	}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		klog.Errorf("delete cluster json unmarshal error: %s", err.Error())
		return
	}
	if resp.Code == pkg.Ok {
		pg.Status.State = pkg.Success
	} else {
		pg.Status.State = pkg.Failed
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.DeleteCluster {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.DeleteCluster,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}
	return
}
