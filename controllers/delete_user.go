package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) DeletePgUser(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.DeleteUserResponse
	deleteUserReq := &pkg.DeleteUserRequest{
		AllFlag:       pg.Spec.AllFlag,
		ClientVersion: pg.Spec.ClientVersion,
		Clusters:      pg.Spec.ClusterName,
		Namespace:     pg.Spec.Namespace,
		Selector:      pg.Spec.Selector,
		Username:      pg.Spec.Username,
	}
	respByte, err := pkg.Call("POST", pkg.DeleteUserPath, deleteUserReq)
	if err != nil {
		klog.Errorf("call delete user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
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
