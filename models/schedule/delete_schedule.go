package schedule

import (
	"encoding/json"
	"fmt"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

const nameFormat = "%s-pgbackrest-%s"

func DeleteSchedule(pg *v1alpha1.PostgreSQLCluster, backupType string) (err error) {
	scheduleName := fmt.Sprintf(nameFormat, pg.Name, backupType)

	var resp pkg.DeleteScheduleResponse
	// request create schedule
	deleteScheduleReq := &pkg.DeleteScheduleRequest{
		Namespace:    pg.Namespace,
		ScheduleName: scheduleName,
		ClusterName:  pg.Name,
	}

	klog.Infof("params: %+v", deleteScheduleReq)
	respByte, err := pkg.Call("POST", pkg.ScheduleDeletePath, deleteScheduleReq)
	if err != nil {
		klog.Errorf("call delete schedule error: %s", err.Error())
		return
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	models.MergeCondition(pg, pkg.DeleteSchedule, resp.Status)

	return nil
}
