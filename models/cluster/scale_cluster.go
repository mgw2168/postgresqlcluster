package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func ScaleUpPgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.ClusterScaleResponse
	scaleReq := &pkg.ClusterScaleRequest{
		Name:          pg.Spec.Name,
		ClientVersion: pg.Spec.ClientVersion,
		Namespace:     pg.Spec.Namespace,
		CCPImageTag:   pg.Spec.CCPImageTag,
		ReplicaCount:  pg.Spec.ReplicaCount,
		StorageConfig: pg.Spec.StorageConfig,
	}
	respByte, err := pkg.Call("POST", pkg.ScaleClusterPath+pg.Spec.Name, scaleReq)
	if err != nil {
		klog.Errorf("call scale cluster error: ", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("scale cluster json unmarshal error: ", err.Error())
		return
	}

	if resp.Code == pkg.Ok {
		pg.Status.State = v1alpha1.Success
	} else {
		pg.Status.State = v1alpha1.Failed
	}

	res, ok := pg.Status.Condition[v1alpha1.ScaleCluster]
	if ok {
		res.Code = resp.Code
		res.Msg = resp.Msg
	} else {
		pg.Status.Condition = map[string]v1alpha1.ApiResult{
			v1alpha1.ScaleCluster: {
				Code: resp.Code,
				Msg:  resp.Msg,
			}}
	}
	return
}