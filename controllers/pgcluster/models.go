package pgcluster

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
	ReplicaCount    int
	CPULimit        string
	CPURequest      string
	MemoryLimit     string
	MemoryRequest   string
	Database        string
	Username        string
	Password        string
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
	Autofail      UpdateClusterAutofailStatus
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
	// Name of the cluster to show
	// required: true
	Clustername string `json:"clustername"`
	// Selector of the cluster to show
	Selector string `json:"selector"`
	// Image tag of the cluster
	Ccpimagetag string `json:"ccpimagetag"`
	// Version of API client
	// required: true
	ClientVersion string `json:"clientversion"`
	// Namespace to search
	// required: true
	Namespace string `json:"namespace"`
	// Shows all clusters
	AllFlag bool `json:"allflag"`
}
