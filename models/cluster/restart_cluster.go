package cluster

import (
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models"
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
		ClientVersion: pkg.ClientVersion,
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

	models.MergeCondition(pg, pkg.RestartCluster, resp.Status)

	pg.Spec.Restart = false
	return
}
