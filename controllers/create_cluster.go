package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) createPgCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.CreateClusterResponse
	clusterReq := &pkg.CreatePgCluster{
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
	respByte, err := pkg.Call("POST", pkg.CreateClusterPath, clusterReq)
	if err != nil {
		klog.Errorf("call create cluster error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	if resp.Code == pkg.Ok {
		// update cluster status
		pg.Status.PostgreSQLClusterState = v1alpha1.Created
		// append result to status.condition
		pg.Status.Condition = append(pg.Status.Condition, string(respByte))
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.PostgreSQLClusterState = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
