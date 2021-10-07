package v1alpha1

const (
	Creating = "creating"
	Deleted  = "deleted"
	Failed   = "failed"
	Created  = "created"
	Scaled   = "scaled"

	// action of the cluster or user
	CreateCluster    = "create_cluster"
	DeleteCluster    = "delete_cluster"
	UpdateCluster    = "update_cluster"
	ScaleCluster     = "scale_cluster"
	ScaleDownCluster = "scaledown_cluster"
	RestartCluster   = "restart_cluster"
	CreateUser       = "create_user"
	DeleteUser       = "delete_user"
	UpdateUser       = "update_user"
	ShowUser         = "show_user"

	// cluster status
	ClusterStatusUnknown = "unknown"
	// 状态更新异常
)
