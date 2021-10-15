package cluster

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func CreatePgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.CreateClusterResponse
	clusterReq := &pkg.CreatePgCluster{
		ClientVersion:   pg.Spec.ClientVersion,
		Name:            pg.Spec.Name,
		Namespace:       pg.Spec.Namespace,
		SyncReplication: pg.Spec.SyncReplication,
		CCPImage:        pg.Spec.CCPImage,
		CCPImageTag:     pg.Spec.CCPImageTag,
		ReplicaCount:    pg.Spec.ReplicaCount,
		CPULimit:        pg.Spec.CPULimit,
		CPURequest:      pg.Spec.CPURequest,
		MemoryLimit:     pg.Spec.MemoryLimit,
		MemoryRequest:   pg.Spec.MemoryRequest,
		Database:        pg.Spec.Database,
		Username:        pg.Spec.Username,
		Password:        pg.Spec.Password,
		StorageConfig:   pg.Spec.StorageConfig,
		PVCSize:         pg.Spec.PVCSize,
	}
	respByte, err := pkg.Call("POST", pkg.CreateClusterPath, clusterReq)
	if err != nil {
		klog.Errorf("call create cluster error: %s", err.Error())
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

	res, ok := pg.Status.Condition[v1alpha1.CreateCluster]
	if ok {
		res.Code = resp.Code
		res.Msg = resp.Msg
	} else {
		pg.Status.Condition = map[string]v1alpha1.ApiResult{
			v1alpha1.CreateCluster: {
				Code: resp.Code,
				Msg:  resp.Msg,
			}}
	}
	return
}