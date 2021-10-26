package pkg

const (
	// api authorize
	UserName = "admin"
	PassWord = "examplepassword"
	// port of the postgresql operator api server
	srvPort = "8443"
	// IP svc name.namespace
	//IP = "http://139.198.21.143"
	IP = "http://postgres-operator.pgo"
	// Cluster path
	HostPath           = IP + ":" + srvPort
	CreateClusterPath  = HostPath + "/clusters"
	DeleteClusterPath  = HostPath + "/clustersdelete"
	ShowClusterPath    = HostPath + "/showclusters"
	UpdateClusterPath  = HostPath + "/clustersupdate"
	RestartClusterPath = HostPath + "/restart"

	// scale cluster
	ScaleClusterPath     = HostPath + "/clusters/scale/"
	ScaleDownClusterPath = HostPath + "/scaledown/"

	// User path
	CreateUserPath = HostPath + "/usercreate"
	DeleteUserPath = HostPath + "/userdelete"
	UpdateUserPath = HostPath + "/userupdate"
	ShowUserPath   = HostPath + "/usershow"

	// status
	Ok      = "ok"
	Error   = "error"
	Failed  = "failed"
	Success = "success"

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
)
