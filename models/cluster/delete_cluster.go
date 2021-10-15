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
		ClientVersion: pg.Spec.ClientVersion,
		Namespace:     pg.Spec.Namespace,
		DeleteBackups: false,
		DeleteData:    false,
	}
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
		pg.Status.State = v1alpha1.Success
	} else {
		pg.Status.State = v1alpha1.Failed
	}

	res, ok := pg.Status.Condition[v1alpha1.DeleteCluster]
	if ok {
		res.Code = resp.Code
		res.Msg = resp.Msg
	} else {
		pg.Status.Condition = map[string]v1alpha1.ApiResult{
			v1alpha1.DeleteCluster: {
				Code: resp.Code,
				Msg:  resp.Msg,
			}}
	}
	return
}
