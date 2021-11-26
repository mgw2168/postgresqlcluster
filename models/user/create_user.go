package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func CreatePgUser(pg *v1alpha1.PostgreSQLCluster, username, passwd string) (err error) {
	var resp pkg.CreateUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	createUserReq := &pkg.CreateUserRequest{
		Clusters:        clusterName,
		ClientVersion:   "4.7.1",
		ManagedUser:     pg.Spec.ManagedUser,
		Namespace:       pg.Spec.Namespace,
		Username:        username,
		Password:        passwd,
		PasswordAgeDays: 86400,
		PasswordLength:  8,
		PasswordType:    "md5",
	}
	klog.Infof("params: %+v", createUserReq)
	respByte, err := pkg.Call("POST", pkg.CreateUserPath, createUserReq)
	if err != nil {
		klog.Errorf("call create user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.CreateUser {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.CreateUser,
			Code: resp.Code,
			Msg:  resp.Msg,
		})
	}
	return
}
