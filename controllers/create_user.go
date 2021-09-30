package controllers

import (
	"context"
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) createPgUser(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.CreateUserResponse
	createUserReq := &pgcluster.CreateUserRequest{
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
	respByte, err := request.Call("POST", request.CreateUserPath, createUserReq)
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
