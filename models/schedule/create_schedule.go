package schedule

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func CreateSchedule(pg *v1alpha1.PostgreSQLCluster, backupType string) (err error) {
	cronExpression := ""

	switch backupType {
	case "full":
		cronExpression = pg.Spec.FullBackupSchedule
	case "incr":
		cronExpression = pg.Spec.IncrBackupSchedule
	case "diff":
		cronExpression = pg.Spec.DiffBackupSchedule
	}

	var resp pkg.CreateScheduleResponse
	// request create schedule
	scheduleReq := &pkg.CreateScheduleRequest{
		ClusterName:    pg.Name,
		Namespace:      pg.Namespace,
		Schedule:       cronExpression,
		ScheduleType:   "pgbackrest",
		PGBackRestType: backupType,
	}

	klog.Infof("params: %+v", scheduleReq)
	respByte, err := pkg.Call("POST", pkg.SchedulePath, scheduleReq)
	if err != nil {
		klog.Errorf("call create schedule error: %s", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	models.MergeCondition(pg, pkg.CreateSchedule, resp.Status)

	return nil
}
