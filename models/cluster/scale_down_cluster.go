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
			"?version=4.7.1"+
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

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.ScaleDownCluster {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.ScaleDownCluster,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}
	return
}
