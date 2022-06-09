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

	ReplicaStorageConfig  string
	BackrestStorageConfig string
	BackrestPVCSize       string
	AutofailFlag          bool

	BackrestStorageType string
	BackrestS3Key       string
	BackrestS3KeySecret string
	BackrestS3Bucket    string
	BackrestS3Region    string
	BackrestS3Endpoint  string
	BackrestS3URIStyle  string
	BackrestS3VerifyTLS UpdateBackrestS3VerifyTLS
	BackrestRepoPath    string

	PGDataSource PGDataSourceSpec
}

type PGDataSourceSpec struct {
	Namespace   string `json:"namespace"`
	RestoreFrom string `json:"restoreFrom"`
	RestoreOpts string `json:"restoreOpts"`
}

// UpdateBackrestS3VerifyTLS defines the types for updating the S3 TLS verification configuration
type UpdateBackrestS3VerifyTLS int

// set the different values around updating the S3 TLS verification configuration
const (
	UpdateBackrestS3VerifyTLSDoNothing UpdateBackrestS3VerifyTLS = iota
	UpdateBackrestS3VerifyTLSEnable
	UpdateBackrestS3VerifyTLSDisable
)

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
	Superuser    bool
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
	Superuser                bool
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

// backrestbackup
type CreateBackrestBackupRequest struct {
	Namespace           string
	Args                []string
	Selector            string
	BackupOpts          string
	BackrestStorageType string
}

type CreateBackrestBackupResponse struct {
	Results []string
	Status
}

type ShowBackrestDetail struct {
	Name        string
	Info        []PgBackRestInfo
	StorageType string
}

type PgBackRestInfo struct {
	Archives []PgBackRestInfoArchive `json:"archive"`
	Backups  []PgBackRestInfoBackup  `json:"backup"`
	Cipher   string                  `json:"cipher"`
	DBs      []PgBackRestInfoDB      `json:"db"`
	Name     string                  `json:"name"`
	Status   PgBackRestInfoStatus    `json:"status"`
}

type PgBackRestInfoDB struct {
	ID       int    `json:"id"`
	SystemID int64  `json:"system-id,omitempty"`
	Version  string `json:"version,omitempty"`
}

type PgBackRestInfoArchive struct {
	DB  PgBackRestInfoDB `json:"db"`
	ID  string           `json:"id"`
	Max string           `json:"max"`
	Min string           `json:"min"`
}

type PgBackRestInfoStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ShowBackrestResponse struct {
	Items []ShowBackrestDetail
	Status
}

type PgBackRestInfoBackup struct {
	Archive   PgBackRestInfoBackupArchive   `json:"archive"`
	Backrest  PgBackRestInfoBackupBackrest  `json:"backrest"`
	Database  PgBackRestInfoDB              `json:"database"`
	Info      PgBackRestInfoBackupInfo      `json:"info"`
	Label     string                        `json:"label"`
	Prior     string                        `json:"prior"`
	Reference []string                      `json:"reference"`
	Timestamp PgBackRestInfoBackupTimestamp `json:"timestamp"`
	Type      string                        `json:"type"`
}

type PgBackRestInfoBackupTimestamp struct {
	Start int64 `json:"start"`
	Stop  int64 `json:"stop"`
}

type PgBackRestInfoBackupArchive struct {
	Start string `json:"start"`
	Stop  string `json:"stop"`
}

type PgBackRestInfoBackupBackrest struct {
	Format  int    `json:"format"`
	Version string `json:"version"`
}

type PgBackRestInfoBackupInfo struct {
	Delta      int64                              `json:"delta"`
	Repository PgBackRestInfoBackupInfoRepository `json:"repository"`
	Size       int64                              `json:"size"`
}

type PgBackRestInfoBackupInfoRepository struct {
	Delta int64 `json:"delta"`
	Size  int64 `json:"size"`
}

type DeleteBackrestBackupRequest struct {
	// ClientVersion represents the version of the client that is making the API
	// request
	ClientVersion string
	// ClusterName is the name of the pgcluster of which we want to delete the
	// backup from
	ClusterName string
	// Namespace isthe namespace that the cluster is in
	Namespace string
	// Target is the nane of the backup to be deleted
	Target string
}

type DeleteBackrestBackupResponse struct {
	Status
}

// schedule
type CreateScheduleRequest struct {
	ClusterName         string
	Name                string
	Namespace           string
	Schedule            string
	ScheduleType        string
	Selector            string
	PGBackRestType      string
	BackrestStorageType string
	PVCName             string
	ScheduleOptions     string
	StorageConfig       string
	PolicyName          string
	Database            string
	Secret              string
}

type CreateScheduleResponse struct {
	Results []string
	Status
}

type DeleteScheduleRequest struct {
	Namespace    string
	ScheduleName string
	ClusterName  string
	Selector     string
}

type DeleteScheduleResponse struct {
	Results []string
	Status
}
