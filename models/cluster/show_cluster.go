package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/pkg"
	"k8s.io/klog/v2"
)

func ShowClusterConfig(pg *v1alpha1.PostgreSQLCluster) (pkg.PGClusterInfo, error) {
	var resp pkg.ShowClusterResponse

	showClusterReq := &pkg.ShowClusterRequest{
		Clustername:   pg.Spec.Name,
		ClientVersion: "4.7.1",
		Namespace:     pg.Spec.Namespace,
		AllFlag:       true,
	}

	respByte, err := pkg.Call("POST", pkg.ShowClusterPath, showClusterReq)
	if err != nil {
		klog.Errorf("call show cluster error: ", err.Error())
		return pkg.PGClusterInfo{}, err
	}

	if err = json.Unmarshal(respByte, &resp); err != nil {
		klog.Errorf("json unmarshal error: %s; data: %s", err, respByte)
		return pkg.PGClusterInfo{}, err
	}

	if len(resp.Results) == 0 {
		return pkg.PGClusterInfo{}, fmt.Errorf("no pgcluster:%s found", pg.Spec.Name)
	}

	return resp.Results[0], nil
}
