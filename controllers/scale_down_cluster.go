package controllers

import (
	"context"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func (r *PostgreSQLClusterReconciler) scaleDownPgCluster(ctx context.Context, pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.ScaleDownResponse
	//scaleReq := &pgcluster.ScaleDownRequest{
	//	Name:          pg.Spec.Name,
	//	ClientVersion: pg.Spec.ClientVersion,
	//	Namespace:     pg.Spec.Namespace,
	//	DeleteData:    pg.Spec.DeleteData,
	//}
	var delete_data string
	if pg.Spec.DeleteData {
		delete_data = "true"
	} else {
		delete_data = "false"
	}
	respByte, err := pkg.Call("GET",
		pkg.ScaleDownClusterPath+
			pg.Spec.Name+
			"?version="+pg.Spec.ClientVersion+
			"&namespace="+pg.Spec.Namespace+
			"&replica-name="+pg.Spec.ReplicaName+
			"&delete-data="+delete_data,
		nil)
	if err != nil {
		klog.Errorf("call scale down cluster error: ", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("scale down cluster json unmarshal error: ", err.Error())
	}

	if resp.Code == pkg.Ok {
		// update cluster status
		pg.Status.PostgreSQLClusterState = v1alpha1.Scaled
		pg.Status.Condition = append(pg.Status.Condition, string(respByte))
		err = r.Status().Update(ctx, pg)
	} else {
		pg.Status.PostgreSQLClusterState = v1alpha1.Failed
		err = r.Status().Update(ctx, pg)
	}
	return
}
