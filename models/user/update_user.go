package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func UpdatePgUser(pg *v1alpha1.PostgreSQLCluster, username, passwd string) (err error) {
	var resp pkg.UpdateUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	updateUserReq := &pkg.UpdateUserRequest{
		Clusters:                 clusterName,
		ClientVersion:            "4.7.1",
		Namespace:                pg.Spec.Namespace,
		Username:                 username,
		Password:                 passwd,
		PasswordAgeDays:          86400,
		PasswordLength:           8,
		PasswordType:             "md5",
		SetSystemAccountPassword: pg.Spec.SetSystemAccountPassword,
	}
	klog.Infof("params: %+v", updateUserReq)
	respByte, err := pkg.Call("POST", pkg.UpdateUserPath, updateUserReq)
	if err != nil {
		klog.Errorf("call update user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.UpdateUser {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.UpdateUser,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}
	return
}
