package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func ScaleUpPgCluster(pg *v1alpha1.PostgreSQLCluster, replicaCount int) (err error) {
	var resp pkg.ClusterScaleResponse
	scaleReq := &pkg.ClusterScaleRequest{
		Name:          pg.Spec.Name,
		ClientVersion: pkg.ClientVersion,
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

	models.MergeCondition(pg, pkg.ScaleCluster, resp.Status)

	return
}
