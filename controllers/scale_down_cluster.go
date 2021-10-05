package controllers

import (
	"context"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/controllers/pgcluster"
	"github.com/kubesphere/controllers/request"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) scaleDownPgCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pgcluster.ScaleDownResponse
	//scaleReq := &pgcluster.ScaleDownRequest{
	//	Name:          pg.Spec.Name,
	//	ClientVersion: pg.Spec.ClientVersion,
	//	Namespace:     pg.Spec.Namespace,
	//	DeleteData:    pg.Spec.DeleteData,
	//}
	respByte, err := request.Call("GET",
		request.ScaleDownClusterPath+
			pg.Spec.Name+
			"?version="+pg.Spec.ClientVersion+
			"&namespace="+pg.Spec.Namespace+
			"&replica-name="+pg.Spec.ReplicaName+
			"&delete-data="+pg.Spec.ScaleDownDeleteData,
		nil)
	if err != nil {
		klog.Errorf("call scale down cluster error: ", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, resp); err != nil {
		klog.Errorf("scale down cluster json unmarshal error: ", err.Error())
	}

	if resp.Code == request.Ok {
		// update cluster status
		pg.Status.PostgreSQLClusterState = v1alpha1.Scaled
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.PostgreSQLClusterState = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
