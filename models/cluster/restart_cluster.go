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
	if resp.Code == pkg.Ok {
		// update cluster status
		pg.Status.State = v1alpha1.Success
	} else {
		pg.Status.State = v1alpha1.Failed
	}

	flag := true
	for _, res := range pg.Status.Condition {
		if res.Api == v1alpha1.RestartCluster {
			flag = false
			res.Code = resp.Code
			res.Msg = resp.Msg
			break
		}
	}
	if flag {
		pg.Status.Condition = append(pg.Status.Condition, v1alpha1.ApiResult{
			Api:  v1alpha1.RestartCluster,
			Code: resp.Code,
			Msg:  resp.Msg,
			Data: "",
		})
	}
	pg.Spec.Restart = false
	return
}
