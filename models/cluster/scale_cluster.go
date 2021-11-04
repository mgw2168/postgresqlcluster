package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func ScaleUpPgCluster(pg *v1alpha1.PostgreSQLCluster, replicaCount int) (err error) {
	var resp pkg.ClusterScaleResponse
	scaleReq := &pkg.ClusterScaleRequest{
		Name:          pg.Spec.Name,
		ClientVersion: "4.7.1",
		Namespace:     pg.Spec.Namespace,
		CCPImageTag:   pg.Spec.CCPImageTag,
		ReplicaCount:  replicaCount,
		StorageConfig: pg.Spec.StorageConfig,
	}
	klog.Infof("params: %+v", scaleReq)
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
		pg.Status.State = pkg.Success
	} else {
		pg.Status.State = pkg.Failed
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.ScaleCluster {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.ScaleCluster,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}
	return
}
