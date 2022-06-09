package models

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"time"
)

func MergeCondition(pg *v1alpha1.PostgreSQLCluster, action string, status pkg.Status) {
	// update condition to the latest status
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == action {
			pg.Status.Condition[i].UpdateTime = time.Now().Format(time.RFC3339)
			pg.Status.Condition[i].Code = status.Code
			pg.Status.Condition[i].Msg = status.Msg
			return
		}
	}

	// append new condition
	pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
		UpdateTime: time.Now().Format(time.RFC3339),
		Api:        action,
		Code:       status.Code,
		Msg:        status.Msg,
	})
}

func InSlice(pgCluster *v1alpha1.PostgreSQLCluster, username string) bool {
	for _, u := range pgCluster.Spec.Users {
		if u.UserName == username {
			return true
		}
	}
	return false
}
