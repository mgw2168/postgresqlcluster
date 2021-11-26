package cluster

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func CreatePgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.CreateClusterResponse
	if pg.Spec.PgVersion == "12" {
		pg.Spec.CCPImageTag = "centos8-12.7-3.0-4.7.1"
	} else if pg.Spec.PgVersion == "13" {
		pg.Spec.CCPImageTag = "centos8-13.3-3.0-4.7.1"
	}
	clusterReq := &pkg.CreatePgCluster{
		ClientVersion:   "4.7.1",
		Name:            pg.Spec.Name,
		Namespace:       pg.Spec.Namespace,
		SyncReplication: pg.Spec.SyncReplication,
		CCPImage:        "radondb-postgres-gis-ha",
		CCPImagePrefix:  "docker.io/radondb",
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
		AutofailFlag:    true,
	}
	klog.Infof("params: %+v", clusterReq)
	respByte, err := pkg.Call("POST", pkg.CreateClusterPath, clusterReq)
	if err != nil {
		klog.Errorf("call create cluster error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.CreateCluster {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.CreateCluster,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}
	return
}
