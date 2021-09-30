package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) updatePgUser(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.UpdateUserResponse
	updateUserReq := &pgcluster.UpdateUserRequest{
		AllFlag:                  pg.Spec.AllFlag,
		ClientVersion:            pg.Spec.ClientVersion,
		Clusters:                 pg.Spec.ClusterName,
		Namespace:                pg.Spec.Namespace,
		Password:                 pg.Spec.Password,
		PasswordAgeDays:          pg.Spec.PasswordAgeDays,
		PasswordLength:           pg.Spec.PasswordLength,
		PasswordType:             pg.Spec.PasswordType,
		Selector:                 pg.Spec.Selector,
		SetSystemAccountPassword: pg.Spec.SetSystemAccountPassword,
		Username:                 pg.Spec.Username,
	}
	respByte, err := request.Call("POST", request.UpdateUserPath, updateUserReq)
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
