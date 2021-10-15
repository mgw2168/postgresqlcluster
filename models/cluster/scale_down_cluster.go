package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func ScaleDownPgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.ScaleDownResponse
	respByte, err := pkg.Call("GET",
		pkg.ScaleDownClusterPath+
			pg.Spec.Name+
			"?version="+pg.Spec.ClientVersion+
			"&namespace="+pg.Spec.Namespace+
			"&replica-name="+pg.Spec.ReplicaName+
			"&delete-data=false",
		nil)
	if err != nil {
		klog.Errorf("call scale down cluster error: ", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("scale down cluster json unmarshal error: ", err.Error())
		return
	}

	if resp.Code == pkg.Ok {
		pg.Status.State = v1alpha1.Success
	} else {
		pg.Status.State = v1alpha1.Failed
	}
	res, ok := pg.Status.Condition[v1alpha1.ScaleDownCluster]
	if ok {
		res.Code = resp.Code
		res.Msg = resp.Msg
	} else {
		pg.Status.Condition = map[string]v1alpha1.ApiResult{
			v1alpha1.ScaleDownCluster: {
				Code: resp.Code,
				Msg:  resp.Msg,
			}}
	}
	return
}
