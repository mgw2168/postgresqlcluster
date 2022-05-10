package cluster

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
	"time"
)

func CreatePgCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	switch pg.Spec.PgVersion {
	case "12":
		pg.Spec.CCPImageTag = "centos8-12.7-3.0-4.7.1"
	case "13":
		pg.Spec.CCPImageTag = "centos8-13.3-3.0-4.7.1"
	case "14":
		pg.Spec.CCPImageTag = "debian-14.2-3.1-2.1.1"
	}

	if pg.Spec.BackrestStorageType == "s3" && pg.HasValidS3Conf() {
		if err = CreateManagedResource(pg); err != nil {
			klog.Errorf("create managed resource error: %s", err)
			return err
		}
	}

	var resp pkg.CreateClusterResponse

	clusterReq := &pkg.CreatePgCluster{
		ClientVersion:   pkg.ClientVersion,
		Name:            pg.Spec.Name,
		Namespace:       pg.Spec.Namespace,
		SyncReplication: pg.Spec.SyncReplication,
		CCPImage:        "radondb-postgres-gis-ha",
		// configured at operator crds, so not need
		//CCPImagePrefix:  "docker.io/radondb",
		CCPImageTag:           pg.Spec.CCPImageTag,
		ReplicaCount:          pg.Spec.ReplicaCount,
		CPULimit:              pg.Spec.CPULimit,
		CPURequest:            pg.Spec.CPURequest,
		MemoryLimit:           pg.Spec.MemoryLimit,
		MemoryRequest:         pg.Spec.MemoryRequest,
		Database:              pg.Spec.Database,
		Username:              pg.Spec.Username,
		Password:              pg.Spec.Password,
		StorageConfig:         pg.Spec.StorageConfig,
		ReplicaStorageConfig:  pg.Spec.StorageConfig,
		BackrestStorageConfig: pg.Spec.StorageConfig,
		BackrestPVCSize:       pg.Spec.PVCSize,
		PVCSize:               pg.Spec.PVCSize,
		AutofailFlag:          true,
	}

	if pg.Spec.BackrestStorageType == "s3" && pg.HasValidS3Conf() {
		repoPath := fmt.Sprintf("%s-%s", pg.Name, time.Now().Format("20060102-150405"))
		pg.Status.BackrestRepoPath = repoPath

		clusterReq.BackrestStorageType = pg.Spec.BackrestStorageType

		plainKey, _ := base64.StdEncoding.DecodeString(pg.Spec.BackrestS3Key)
		plainKeySecret, _ := base64.StdEncoding.DecodeString(pg.Spec.BackrestS3KeySecret)

		clusterReq.BackrestS3Key = string(plainKey)
		clusterReq.BackrestS3KeySecret = string(plainKeySecret)

		clusterReq.BackrestS3Bucket = pg.Spec.BackrestS3Bucket
		clusterReq.BackrestS3Region = pg.Spec.BackrestS3Region
		clusterReq.BackrestS3Endpoint = pg.Spec.BackrestS3Endpoint
		clusterReq.BackrestS3URIStyle = pg.Spec.BackrestS3URIStyle
		clusterReq.BackrestS3VerifyTLS = pkg.UpdateBackrestS3VerifyTLSDisable
		clusterReq.BackrestRepoPath = fmt.Sprintf("/%s", repoPath)

		if pg.Spec.RestoreFrom != "" {
			clusterReq.PGDataSource.RestoreFrom = GetRestoreFromName(pg)
			clusterReq.PGDataSource.Namespace = pg.Namespace
			if pg.Spec.RestoreTarget != "" {
				clusterReq.PGDataSource.RestoreOpts = fmt.Sprintf(`--repo-type=%s --type=time --target='%s'`, pg.Spec.BackrestStorageType, pg.Spec.RestoreTarget)
			} else {
				clusterReq.PGDataSource.RestoreOpts = fmt.Sprintf(`--repo-type=%s`, pg.Spec.BackrestStorageType)
			}
		}
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

	models.MergeCondition(pg, pkg.CreateCluster, resp.Status)

	return
}
