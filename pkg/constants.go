package pkg

const (
	// api authorize
	UserName = "admin"
	PassWord = "examplepassword"
	// port of the postgresql operator api server
	srvPort = "8443"
	// IP svcname.namespace
	IP = "http://139.198.21.143"
	//IP = "http://postgres-operator.pgo"
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

	// Ok status
	Ok = "ok"

	// Error code string
	Error = "error"
)
