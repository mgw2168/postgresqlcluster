package backup

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func ShowBackup(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var backups []v1alpha1.PgBackup

	url := pkg.BackrestPath + "/" + pg.Name + "?version=" + pkg.ClientVersion + "&selector=" + "" + "&namespace=" + pg.Namespace
	var resp pkg.ShowBackrestResponse

	respByte, err := pkg.Call("GET", url, nil)
	if err != nil {
		klog.Errorf("call backrest backup error: %s", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	for i, _ := range resp.Items {
		item := resp.Items[i]
		backupInfos := item.Info

		for j, _ := range backupInfos {
			info := backupInfos[j]
			for k, _ := range info.Backups {
				bi := info.Backups[k]
				backup := v1alpha1.PgBackup{
					Type:           bi.Type,
					Name:           bi.Label,
					StorageType:    item.StorageType,
					StartTime:      bi.Timestamp.Start,
					EndTime:        bi.Timestamp.Stop,
					StartArchive:   bi.Archive.Start,
					StopArchive:    bi.Archive.Stop,
					DatabaseSize:   bi.Info.Size,
					RepositorySize: bi.Info.Repository.Size,
				}
				backups = append(backups, backup)
			}
		}
	}

	models.MergeCondition(pg, pkg.ShowBackup, resp.Status)
	pg.Status.Backups = backups

	return nil
}
