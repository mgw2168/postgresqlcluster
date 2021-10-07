package controllers

import (
	"context"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) restartCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.RestartResponse
	restartReq := &pkg.RestartRequest{
		Namespace:     pg.Spec.Namespace,
		ClusterName:   pg.Spec.Name,
		RollingUpdate: pg.Spec.RollingUpdate,
		Targets:       pg.Spec.Targets,
		ClientVersion: pg.Spec.ClientVersion,
	}
	respByte, err := pkg.Call("POST", pkg.RestartClusterPath, restartReq)
	if err != nil {
		klog.Errorf("call restart cluster error: %s", err.Error())
		return
	}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		klog.Errorf("restart cluster json unmarshal error: %s", err.Error())
		return
	}
	if resp.Code == pkg.Ok {
		// update cluster status
		pg.Status.PostgreSQLClusterState = v1alpha1.Created
		pg.Status.Condition = append(pg.Status.Condition, string(respByte))
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.PostgreSQLClusterState = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
