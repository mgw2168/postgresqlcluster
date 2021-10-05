package controllers

import (
	"context"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) scalePgCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.ClusterScaleResponse
	scaleReq := &pgcluster.ClusterScaleRequest{
		Name:          pg.Spec.Name,
		ClientVersion: pg.Spec.ClientVersion,
		Namespace:     pg.Spec.Namespace,
		CCPImageTag:   pg.Spec.CCPImageTag,
		NodeLabel:     pg.Spec.NodeLabel,
		ReplicaCount:  pg.Spec.ReplicaCount,
		ServiceType:   pg.Spec.ServiceType,
		StorageConfig: pg.Spec.StorageConfig,
		Tolerations:   pg.Spec.Tolerations,
	}
	respByte, err := request.Call("POST", request.ScaleClusterPath+pg.Spec.Name, scaleReq)
	if err != nil {
		klog.Errorf("call scale cluster error: ", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("scale cluster json unmarshal error: ", err.Error())
	}

	if resp.Code == request.Ok {
		// update cluster status
		pg.Status.PostgreSQLClusterState = v1alpha1.Scaled
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.PostgreSQLClusterState = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
