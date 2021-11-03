package pkg

import v1 "k8s.io/api/core/v1"

type Status struct {
	// status code
	Code string
	// status message
	Msg string
}

// ******** create cluster
type CreatePgCluster struct {
	ClientVersion   string
	Name            string
	Namespace       string
	SyncReplication *bool
	CCPImage        string
	CCPImageTag     string
	CCPImagePrefix  string
	ReplicaCount    int
	CPULimit        string
	CPURequest      string
	MemoryLimit     string
	MemoryRequest   string
	Database        string
	Username        string
	Password        string
	StorageConfig   string
	PVCSize         string
	AutofailFlag    bool
}

type CreateClusterResponse struct {
	Result CreateClusterDetail `json:"result"`
	Status `json:"status"`
}

type CreateClusterDetail struct {
	Database   string
	Name       string
	Users      []CreateClusterDetailUser
	WorkflowID string
}

type CreateClusterDetailUser struct {
	Password string
	Username string
}

// ********** update cluster
type UpdateClusterAutofailStatus int

type UpdateClusterRequest struct {
	Clustername   []string
	ClientVersion string
	Namespace     string
	AllFlag       bool
	Autofail      int
	CPULimit      string
	CPURequest    string
	MemoryLimit   string
	MemoryRequest string
	PVCSize       string
	Startup       bool
	Shutdown      bool
	Tolerations   []v1.Toleration `json:"tolerations"`
}

type UpdateClusterResponse struct {
	Results []string
	Status
}

// ********** delete cluster
type DeleteClusterRequest struct {
	Clustername   string
	Selector      string
	ClientVersion string
	Namespace     string
	AllFlag       bool
	DeleteBackups bool
	DeleteData    bool
}

type DeleteClusterResponse struct {
	Results []string
	Status
}

// ********* show cluster
type ShowClusterRequest struct {
	Clustername   string `json:"clustername"`
	Selector      string `json:"selector"`
	Ccpimagetag   string `json:"ccpimagetag"`
	ClientVersion string `json:"clientversion"`
	Namespace     string `json:"namespace"`
	AllFlag       bool   `json:"allflag"`
}

// ******** scale cluster
type ClusterScaleRequest struct {
	Name          string
	ClientVersion string
	Namespace     string
	CCPImageTag   string
	NodeLabel     string
	ReplicaCount  int
	ServiceType   string
	StorageConfig string
	Tolerations   []v1.Toleration `json:"tolerations"`
}
type ClusterScaleResponse struct {
	Results []string
	Status
}

// ******** scale down
type ScaleDownRequest struct {
	Name          string
	ClientVersion string
	Namespace     string
	ReplicaName   string `json:"replica-name"`
	DeleteData    bool   `json:"delete-data"`
}

type ScaleDownResponse struct {
	Results []string
	Status
}

// ******** restart cluster
type RestartRequest struct {
	Namespace     string
	ClusterName   string
	RollingUpdate bool
	Targets       []string
	ClientVersion string
}

type RestartResponse struct {
	Result RestartDetail
	Status
}

type RestartDetail struct {
	ClusterName  string
	Instances    []InstanceDetail
	Error        bool
	ErrorMessage string
}

type InstanceDetail struct {
	InstanceName string
	Error        bool
	ErrorMessage string
}

// ******** create user
type CreateUserRequest struct {
	AllFlag         bool
	Clusters        []string
	ClientVersion   string
	ManagedUser     bool
	Namespace       string
	Password        string
	PasswordAgeDays int
	PasswordLength  int
	// PasswordType is one of "md5" or "scram-sha-256", defaults to "md5"
	PasswordType string
	Selector     string
	Username     string
}

type CreateUserResponse struct {
	Results []UserResponseDetail
	Status
}

type UserResponseDetail struct {
	ClusterName  string
	Error        bool
	ErrorMessage string
	Password     string
	Username     string
	ValidUntil   string
}

// ******** delete user
type DeleteUserRequest struct {
	AllFlag       bool
	ClientVersion string
	Clusters      []string
	Namespace     string
	Selector      string
	Username      string
}

type DeleteUserResponse struct {
	Results []UserResponseDetail
	Status
}

// ******** update user
type UpdateUserRequest struct {
	ClientVersion            string
	Namespace                string
	AllFlag                  bool
	Selector                 string
	Clusters                 []string
	Username                 string
	Password                 string
	PasswordAgeDays          int
	PasswordLength           int
	PasswordType             string
	SetSystemAccountPassword bool
}

type UpdateUserResponse struct {
	Results []UserResponseDetail
	Status
}

// ******** show user
type ShowUserRequest struct {
	AllFlag            bool
	Clusters           []string
	ClientVersion      string
	Expired            int
	Namespace          string
	Selector           string
	ShowSystemAccounts bool
}

type ShowUserResponse struct {
	Results []UserResponseDetail
	Status
}
