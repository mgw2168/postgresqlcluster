package pkg

const (
	// api authorize
	UserName = "admin"
	PassWord = "examplepassword"
	// namespace of pg operator
	PgoNamespace = "dmp-system"
	// port of the postgresql operator api server
	srvPort = "8443"
	// IP svc name.namespace
	//IP = "http://127.0.0.1"
	IP = "http://postgres-operator." + PgoNamespace
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

	// backrestbackup
	BackrestBackupPath = HostPath + "/backrestbackup"

	// backrest
	BackrestPath = HostPath + "/backrest"

	// schedule
	SchedulePath       = HostPath + "/schedule"
	ScheduleDeletePath = HostPath + "/scheduledelete"

	// status
	Ok      = "ok"
	Error   = "error"
	Failed  = "failed"
	Success = "success"

	// client version
	ClientVersion = "2.1.1"

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
	PerformBackup    = "perform_backup"
	ShowBackup       = "show_backup"
	DeleteBackup     = "delete_backup"
	CreateSchedule   = "create_schedule"
	DeleteSchedule   = "delete_schedule"
)
