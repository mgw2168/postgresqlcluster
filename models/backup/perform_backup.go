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

func PerformBackup(pg *v1alpha1.PostgreSQLCluster) (err error) {
	parts := strings.Split(pg.Spec.PerformBackup, "|")
	if len(parts) < 2 {
		klog.Errorf("invalid backup arguments")
		return fmt.Errorf("invalid backup arguments: %s", pg.Spec.PerformBackup)
	}

	if parts[1] != "full" && parts[1] != "incr" && parts[1] != "diff" {
		klog.Errorf("%s of backup is not supported", parts[1])
		return nil
	}

	var resp pkg.CreateClusterResponse

	// take a backup
	backupReq := &pkg.CreateBackrestBackupRequest{
		Namespace:           pg.Namespace,
		Args:                []string{pg.Name},
		Selector:            "",
		BackupOpts:          fmt.Sprintf(`--type=%s`, parts[1]),
		BackrestStorageType: pg.Spec.BackrestStorageType,
	}

	klog.Infof("params: %+v", backupReq)
	respByte, err := pkg.Call("POST", pkg.BackrestBackupPath, backupReq)
	if err != nil {
		klog.Errorf("call backrest backup error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	models.MergeCondition(pg, pkg.PerformBackup, resp.Status)

	return
}
