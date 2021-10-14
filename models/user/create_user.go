package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func CreatePgUser(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.CreateUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	createUserReq := &pkg.CreateUserRequest{
		Clusters:        clusterName,
		ClientVersion:   pg.Spec.ClientVersion,
		ManagedUser:     pg.Spec.ManagedUser,
		Namespace:       pg.Spec.Namespace,
		Username:        pg.Spec.Username,
		Password:        pg.Spec.Password,
		PasswordAgeDays: 86400,
		PasswordLength:  8,
		PasswordType:    "md5",
	}
	respByte, err := pkg.Call("POST", pkg.CreateUserPath, createUserReq)
	if err != nil {
		klog.Errorf("call create user error: %s", err.Error())
		return
	}
	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return
	}

	if resp.Code == pkg.Ok {
		pg.Status.State = v1alpha1.Success
	} else {
		pg.Status.State = v1alpha1.Failed
	}

	res, ok := pg.Status.Condition[v1alpha1.CreateUser]
	if ok {
		res.Code = resp.Code
		res.Msg = resp.Msg
	} else {
		pg.Status.Condition = map[string]v1alpha1.ApiResult{
			v1alpha1.CreateUser: {
				Code: resp.Code,
				Msg:  resp.Msg,
			}}
	}
	return
}
