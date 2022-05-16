package schedule

import "github.com/kubesphere/api/v1alpha1"

func UpdateSchedule(pg *v1alpha1.PostgreSQLCluster, backupType string) (err error) {
	err = DeleteSchedule(pg, backupType)
	if err != nil {
		return
	}

	err = CreateSchedule(pg, backupType)
	if err != nil {
		return
	}

	return nil
}
