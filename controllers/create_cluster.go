package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) createPgCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.CreateClusterResponse
	clusterReq := &pgcluster.CreatePgCluster{
		ClientVersion:   pg.Spec.ClientVersion,
		Name:            pg.Spec.Name,
		Namespace:       pg.Spec.Namespace,
		SyncReplication: pg.Spec.SyncReplication,
		CCPImage:        pg.Spec.CCPImage,
		CCPImageTag:     pg.Spec.CCPImageTag,
		ReplicaCount:    pg.Spec.ReplicaCount,
		CPULimit:        pg.Spec.CPULimit,
		CPURequest:      pg.Spec.CPURequest,
		MemoryLimit:     pg.Spec.MemoryLimit,
		MemoryRequest:   pg.Spec.MemoryRequest,
		Database:        pg.Spec.Database,
		Username:        pg.Spec.Username,
		Password:        pg.Spec.Password,
	}
	respByte, err := request.Call("POST", request.CreateClusterPath, clusterReq)
	if err != nil {
		klog.Errorf("call create cluster error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}
	if resp.Code == request.Ok {
		// update cluster status
		pg.Status.State = v1alpha1.Created
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.State = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
