package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

func ScaleDownPgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.ScaleDownResponse
	respByte, err := pkg.Call("GET",
		pkg.ScaleDownClusterPath+
			pg.Spec.Name+
			"?version="+pkg.ClientVersion+
			"&namespace="+pg.Spec.Namespace+
			"&replica-name="+pg.Spec.ReplicaName+
			"&delete-data=false",
		nil)
	if err != nil {
		klog.Errorf("call scale down cluster error: ", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("scale down cluster json unmarshal error: ", err.Error())
		return
	}

	models.MergeCondition(pg, pkg.ScaleDownCluster, resp.Status)

	return
}
