package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func DeletePgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	// delete managed resource if spec->restoreFrom is managed by dmp
	DeleteManagedResource(pg)

	var resp pkg.DeleteClusterResponse
	clusterReq := &pkg.DeleteClusterRequest{
		Clustername:   pg.Spec.Name,
		ClientVersion: pkg.ClientVersion,
		Namespace:     pg.Spec.Namespace,
		DeleteBackups: false,
		DeleteData:    false,
	}
	klog.Infof("params: %+v", clusterReq)
	respByte, err := pkg.Call("POST", pkg.DeleteClusterPath, clusterReq)
	if err != nil {
		klog.Errorf("call delete cluster error: %s", err.Error())
		return
	}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		klog.Errorf("delete cluster json unmarshal error: %s", err.Error())
		return
	}

	models.MergeCondition(pg, pkg.DeleteCluster, resp.Status)

	return
}
