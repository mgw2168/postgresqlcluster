package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) updatePgCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.UpdateClusterResponse
	updateReq := &pgcluster.UpdateClusterRequest{
		Clustername:   pg.Spec.ClusterName,
		ClientVersion: pg.Spec.ClientVersion,
		Namespace:     pg.Spec.Namespace,
		AllFlag:       pg.Spec.AllFlag,
		Autofail:      pg.Spec.AutoFail,
		CPULimit:      pg.Spec.CPULimit,
		CPURequest:    pg.Spec.CPURequest,
		MemoryLimit:   pg.Spec.MemoryLimit,
		MemoryRequest: pg.Spec.MemoryRequest,
		PVCSize:       pg.Spec.PVCSize,
		Startup:       pg.Spec.Startup,
		Shutdown:      pg.Spec.Shutdown,
		Tolerations:   pg.Spec.Tolerations,
	}
	respByte, err := request.Call("POST", request.UpdateClusterPath, updateReq)
	if err != nil {
		klog.Errorf("call update cluster error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
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
