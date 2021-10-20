package cluster

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func UpdatePgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.UpdateClusterResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	updateReq := &pkg.UpdateClusterRequest{
		Clustername:   clusterName,
		ClientVersion: "4.7.1",
		Namespace:     pg.Spec.Namespace,
		Autofail:      0,
		CPULimit:      pg.Spec.CPULimit,
		CPURequest:    pg.Spec.CPURequest,
		MemoryLimit:   pg.Spec.MemoryLimit,
		MemoryRequest: pg.Spec.MemoryRequest,
		PVCSize:       pg.Spec.PVCSize,
		Startup:       true,
	}
	respByte, err := pkg.Call("POST", pkg.UpdateClusterPath, updateReq)
	if err != nil {
		klog.Errorf("call update cluster error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}
	if resp.Code == pkg.Ok {
		pg.Status.State = v1alpha1.Success
	} else {
		pg.Status.State = v1alpha1.Failed
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == v1alpha1.UpdateCluster {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  v1alpha1.UpdateCluster,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}

	return
}
