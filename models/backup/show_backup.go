package backup

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
	"sort"
)

func ShowBackup(pg *v1alpha1.PostgreSQLCluster) (err error) {
	backups := make(map[string]v1alpha1.PgBackup)
	deletingBackups := make(map[string]string)

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
					Type:            bi.Type,
					Name:            bi.Label,
					StorageType:     item.StorageType,
					StartTime:       bi.Timestamp.Start,
					EndTime:         bi.Timestamp.Stop,
					StartArchive:    bi.Archive.Start,
					StopArchive:     bi.Archive.Stop,
					DatabaseSize:    bi.Info.Size,
					RepositorySize:  bi.Info.Repository.Size,
					RepoPath:        pg.Status.BackrestRepoPath,
					BackupReference: bi.Reference,
				}
				backups[backup.Name] = backup
			}
		}
	}

	for _, b := range pg.Status.BackupDeletingQueue {
		if _, ok := backups[b]; ok {
			deletingBackups[b] = b

			// do not show backups which is being deleting
			delete(backups, b)
		}
	}

	// update BackupDeletingQueue to remove backups which already deleting success.
	pg.Status.BackupDeletingQueue = nil
	for k, _ := range deletingBackups {
		pg.Status.BackupDeletingQueue = append(pg.Status.BackupDeletingQueue, k)
	}

	if len(backups) > 1 {
		pg.Status.Backups = nil
		for _, v := range backups {
			pg.Status.Backups = append(pg.Status.Backups, v)
		}
		sort.Slice(pg.Status.Backups, func(i, j int) bool {
			return pg.Status.Backups[i].StartTime < pg.Status.Backups[j].StartTime
		})

		// hidde the first backup which can not delete
		pg.Status.Backups = pg.Status.Backups[1:]
	} else {
		pg.Status.Backups = nil
	}

	models.MergeCondition(pg, pkg.ShowBackup, resp.Status)

	return nil
}
