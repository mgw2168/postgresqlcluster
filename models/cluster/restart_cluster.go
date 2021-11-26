package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

// todo annotation
func RestartCluster(pg *v1alpha1.PostgreSQLCluster) (err error) {
	var resp pkg.RestartResponse
	restartReq := &pkg.RestartRequest{
		Namespace:     pg.Spec.Namespace,
		ClusterName:   pg.Spec.Name,
		RollingUpdate: true,
		Targets:       pg.Spec.Targets,
		ClientVersion: "4.7.1",
	}
	klog.Infof("params: %+v", restartReq)
	respByte, err := pkg.Call("POST", pkg.RestartClusterPath, restartReq)
	if err != nil {
		klog.Errorf("call restart cluster error: %s", err.Error())
		return
	}
	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		klog.Errorf("restart cluster json unmarshal error: %s", err.Error())
		return
	}

	flag := true
	for i, _ := range pg.Status.Condition {
		if pg.Status.Condition[i].Api == pkg.RestartCluster {
			flag = false
			pg.Status.Condition[i].Code = resp.Code
			pg.Status.Condition[i].Msg = resp.Msg
			pg.Status.Condition[i].Data = resp.Result.ErrorMessage
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  pkg.RestartCluster,
			Code: resp.Code,
			Msg:  resp.Msg,
			Data: resp.Result.ErrorMessage,
		})
	}
	pg.Spec.Restart = false
	return
}
