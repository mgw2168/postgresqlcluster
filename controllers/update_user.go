package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) updatePgUser(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.UpdateUserResponse
	updateUserReq := &pkg.UpdateUserRequest{
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
	respByte, err := pkg.Call("POST", pkg.UpdateUserPath, updateUserReq)
	if err != nil {
		klog.Errorf("call update user error: %s", err.Error())
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
