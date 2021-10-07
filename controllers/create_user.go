package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) createPgUser(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.CreateUserResponse
	createUserReq := &pkg.CreateUserRequest{
		AllFlag:         pg.Spec.AllFlag,
		Clusters:        pg.Spec.ClusterName,
		ClientVersion:   pg.Spec.ClientVersion,
		ManagedUser:     pg.Spec.ManagedUser,
		Namespace:       pg.Spec.Namespace,
		Password:        pg.Spec.Password,
		PasswordAgeDays: pg.Spec.PasswordAgeDays,
		PasswordLength:  pg.Spec.PasswordLength,
		PasswordType:    pg.Spec.PasswordType,
		Selector:        pg.Spec.Selector,
		Username:        pg.Spec.Username,
	}
	respByte, err := pkg.Call("POST", pkg.CreateUserPath, createUserReq)
	if err != nil {
		klog.Errorf("call create user error: %s", err.Error())
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
