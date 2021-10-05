package controllers

import (
	"context"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) deletePgCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.DeleteClusterResponse
	clusterReq := &pgcluster.DeleteClusterRequest{
		Clustername:   pg.Spec.Name,
		Selector:      pg.Spec.Selector,
		ClientVersion: pg.Spec.ClientVersion,
		Namespace:     pg.Spec.Namespace,
		AllFlag:       pg.Spec.AllFlag,
		DeleteBackups: pg.Spec.DeleteBackups,
		DeleteData:    pg.Spec.DeleteData,
	}
	respByte, err := request.Call("POST", request.DeleteClusterPath, clusterReq)
	if err != nil {
		klog.Errorf("call delete cluster error: %s", err.Error())
		return
	}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		klog.Errorf("delete cluster json unmarshal error: %s", err.Error())
		return
	}
	if resp.Code == request.Ok {
		// update cluster status
		pg.Status.PostgreSQLClusterState = v1alpha1.Deleted
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.PostgreSQLClusterState = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
