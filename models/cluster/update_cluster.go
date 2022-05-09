package cluster

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func UpdatePgCluster(pg *v1alpha1.PostgreSQLCluster, pvc bool) (err error) {
	var pvcSize string
	if pvc {
		pvcSize = pg.Spec.PVCSize
	}
	var resp pkg.UpdateClusterResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	updateReq := &pkg.UpdateClusterRequest{
		Clustername:   clusterName,
		ClientVersion: pkg.ClientVersion,
		Namespace:     pg.Spec.Namespace,
		Autofail:      1,
		CPULimit:      pg.Spec.CPULimit,
		CPURequest:    pg.Spec.CPURequest,
		MemoryLimit:   pg.Spec.MemoryLimit,
		MemoryRequest: pg.Spec.MemoryRequest,
		PVCSize:       pvcSize,
		Startup:       true,
	}
	klog.Infof("params: %+v", updateReq)
	respByte, err := pkg.Call("POST", pkg.UpdateClusterPath, updateReq)
	if err != nil {
		klog.Errorf("call update cluster error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	models.MergeCondition(pg, pkg.UpdateCluster, resp.Status)

	return
}
