package backup

import (
	"encoding/json"
	"fmt"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
	"strings"
)

func DeleteBackup(pg *v1alpha1.PostgreSQLCluster) (err error) {
	parts := strings.Split(pg.Spec.BackupToDelete, "|")
	if len(parts) < 2 {
		klog.Errorf("invalid delete backup arguments")
		return fmt.Errorf("invalid delete backup arguments: %s", pg.Spec.PerformBackup)
	}

	var resp pkg.DeleteBackrestBackupResponse

	for i := 1; i < len(parts); i++ {
		// take a backup
		deleteBackupReq := &pkg.DeleteBackrestBackupRequest{
			ClientVersion: pkg.ClientVersion,
			ClusterName:   pg.Name,
			Namespace:     pg.Namespace,
			Target:        parts[i],
		}

		klog.Infof("params: %+v", deleteBackupReq)
		respByte, err := pkg.Call("DELETE", pkg.BackrestPath, deleteBackupReq)
		if err != nil {
			klog.Errorf("call delete backrest backup error: %s", err.Error())
			continue
		}
		if err = json.Unmarshal(respByte, &resp); err != nil {
			klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
			continue
		}

		if resp.Status.Code == "ok" {
			pg.Status.BackupDeletingQueue = append(pg.Status.BackupDeletingQueue, parts[i])
		}

		models.MergeCondition(pg, pkg.DeleteBackup, resp.Status)
	}

	return nil
}
