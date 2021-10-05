package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) listPgUser(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.ShowUserResponse
	listUserReq := &pgcluster.ShowUserRequest{
		AllFlag:            pg.Spec.AllFlag,
		Clusters:           pg.Spec.ClusterName,
		ClientVersion:      pg.Spec.ClientVersion,
		Namespace:          pg.Spec.Namespace,
		Selector:           pg.Spec.Selector,
		ShowSystemAccounts: pg.Spec.ShowSystemAccounts,
	}
	respByte, err := request.Call("POST", request.ShowUserPath, listUserReq)
	if err != nil {
		klog.Errorf("call create user error: %s", err.Error())
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
