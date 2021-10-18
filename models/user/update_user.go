package user

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func UpdatePgUser(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.UpdateUserResponse
	var clusterName []string
	clusterName = append(clusterName, pg.Spec.Name)
	updateUserReq := &pkg.UpdateUserRequest{
		Clusters:                 clusterName,
		ClientVersion:            "4.7.1",
		Namespace:                pg.Spec.Namespace,
		Username:                 pg.Spec.Username,
		Password:                 pg.Spec.Password,
		PasswordAgeDays:          86400,
		PasswordLength:           8,
		PasswordType:             "md5",
		SetSystemAccountPassword: pg.Spec.SetSystemAccountPassword,
	}
	respByte, err := pkg.Call("POST", pkg.UpdateUserPath, updateUserReq)
	if err != nil {
		klog.Errorf("call update user error: %s", err.Error())
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

	//res, ok := pg.Status.Condition[v1alpha1.UpdateUser]
	//if ok {
	//	res.Code = resp.Code
	//	res.Msg = resp.Msg
	//} else {
	//	pg.Status.Condition = map[string]v1alpha1.ApiResult{
	//		v1alpha1.UpdateUser: {
	//			Code: resp.Code,
	//			Msg:  resp.Msg,
	//		}}
	//}
	flag := true
	for _, res := range pg.Status.Condition {
		if res.Api == v1alpha1.UpdateUser {
			flag = false
			res.Code = resp.Code
			res.Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  v1alpha1.UpdateUser,
			Code: resp.Code,
			Msg:  resp.Msg,
			Data: "",
		})
	}
	return
}
