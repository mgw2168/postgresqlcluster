package controllers

import (
	"context"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) restartCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.RestartResponse
	restartReq := &pgcluster.RestartRequest{
		Namespace:     pg.Spec.Namespace,
		ClusterName:   pg.Spec.Name,
		RollingUpdate: pg.Spec.RollingUpdate,
		Targets:       pg.Spec.Targets,
		ClientVersion: pg.Spec.ClientVersion,
	}
	respByte, err := request.Call("POST", request.RestartClusterPath, restartReq)
	if err != nil {
		klog.Errorf("call restart cluster error: %s", err.Error())
		return
	}
	err = json.Unmarshal(respByte, resp)
	if err != nil {
		klog.Errorf("restart cluster json unmarshal error: %s", err.Error())
		return
	}
	if resp.Code == request.Ok {
		// update cluster status
		pg.Status.PostgreSQLClusterState = v1alpha1.Created
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.PostgreSQLClusterState = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
