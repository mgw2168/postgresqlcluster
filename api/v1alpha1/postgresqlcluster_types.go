/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PostgreSQLClusterSpec defines the desired state of PostgreSQLCluster
type PostgreSQLClusterSpec struct {
	Action string `json:"action,omitempty"`
	// ******** create cluster params ********
	//ClientVersion string `json:"clientVersion,omitempty"`
	// cluster name
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	SyncReplication *bool  `json:"syncReplication,omitempty"`
	CCPImage        string `json:"ccpImage,omitempty"`
	CCPImageTag     string `json:"ccpImageTag,omitempty"`
	PgVersion       string `json:"pgVersion,omitempty"`
	ReplicaCount    int    `json:"replicaCount,omitempty"`
	CPULimit        string `json:"cpuLimit,omitempty"`
	CPURequest      string `json:"cpuRequest,omitempty"`
	MemoryLimit     string `json:"memoryLimit,omitempty"`
	MemoryRequest   string `json:"memoryRequest,omitempty"`
	Database        string `json:"database,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	StorageConfig   string `json:"storageConfig,omitempty"`

	// ******** update cluster
	ClusterName []string `json:"clusterName,omitempty"`
	AutoFail    int      `json:"autofail,omitempty"`
	PVCSize     string   `json:"pvcSize,omitempty"`
	Startup     bool     `json:"startup,omitempty"`
	Shutdown    bool     `json:"shutdown,omitempty"`

	// ******** delete cluster
	//ClusterName   string `json:"clusterName"`
	Selector      string `json:"selector,omitempty"`
	AllFlag       bool   `json:"allFlag,omitempty"`
	DeleteBackups bool   `json:"deleteBackups,omitempty"`
	DeleteData    bool   `json:"deleteData,omitempty"`

	// ******** scale cluster
	NodeLabel   string          `json:"nodeLabel,omitempty"`
	ServiceType string          `json:"serviceType,omitempty"`
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// ******** scale down
	ReplicaName         string `json:"replicaName,omitempty"`
	ScaleDownDeleteData string `json:"delete-data,omitempty"`
	// ******** restart cluster
	Restart       bool                `json:"restart,omitempty"`
	RollingUpdate bool                `json:"rollingUpdate,omitempty"`
	Targets       []string            `json:"targets,omitempty"`
	PodAnnotation []map[string]string `json:"podAnnotation,omitempty"`

	// ******** create user
	Users           []User `json:"users,omitempty"`
	ManagedUser     bool   `json:"managedUser,omitempty"`
	PasswordAgeDays int    `json:"passwordAgeDays,omitempty"`
	PasswordLength  int    `json:"passwordLength,omitempty"`
	PasswordType    string `json:"passwordType,omitempty"`

	// ******** update user
	SetSystemAccountPassword bool `json:"setSystemAccountPassword,omitempty"`

	// ******** show user
	ShowSystemAccounts bool `json:"showSystemAccounts,omitempty"`

	ClusterConfig string `json:"postgresqlParams,omitempty"`

	// TODO add comment
	PerformBackup string `json:"performBackup,omitempty"`
	// TODO add comment
	BackupToDelete string `json:"backupToDelete,omitempty"`

	FullBackupSchedule string `json:"fullBackupSchedule,omitempty"`
	DiffBackupSchedule string `json:"diffBackupSchedule,omitempty"`
	IncrBackupSchedule string `json:"incrBackupSchedule,omitempty"`

	BackrestStorageType string `json:"backrestStorageType,omitempty"`

	BackrestS3Key       string `json:"backrestS3Key,omitempty"`
	BackrestS3KeySecret string `json:"backrestS3KeySecret,omitempty"`
	BackrestS3Bucket    string `json:"backrestS3Bucket,omitempty"`
	BackrestS3Endpoint  string `json:"backrestS3Endpoint,omitempty"`
	BackrestS3Region    string `json:"backrestS3Region,omitempty"`
	BackrestS3URIStyle  string `json:"backrestS3URIStyle,omitempty"`

	RestoreFrom   string `json:"restoreFrom,omitempty"`
	RestoreTarget string `json:"restoreTarget,omitempty"`
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// PostgreSQLClusterStatus defines the observed state of PostgreSQLCluster
type PostgreSQLClusterStatus struct {
	Backups             []PgBackup  `json:"backups,omitempty"`
	Condition           []ApiResult `json:"condition,omitempty"`
	BackupDeletingQueue []string    `json:"backupDeletingQueue,omitempty"`
	BackrestRepoPath    string      `json:"backrestRepoPath,omitempty"`
	State               string      `json:"state,omitempty"`
}

type PgBackup struct {
	Type            string   `json:"type"`
	Name            string   `json:"name"`
	RepoPath        string   `json:"repoPath"`
	StorageType     string   `json:"storageType"`
	StartTime       int64    `json:"startTime"`
	EndTime         int64    `json:"endTime"`
	StartArchive    string   `json:"startArchive"`
	StopArchive     string   `json:"stopArchive"`
	DatabaseSize    int64    `json:"databaseSize"`
	RepositorySize  int64    `json:"repositorySize"`
	BackupReference []string `json:"backupReference,omitempty"`
}

// ApiResult defines the result of pg operator ApiServer
type ApiResult struct {
	Api        string `json:"api,omitempty"`
	Code       string `json:"code,omitempty"`
	Msg        string `json:"msg,omitempty"`
	Data       string `json:"data,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,namespace=pgo,resources=deployments,verbs=get;list;watch;patch
//+kubebuilder:rbac:groups=core,namespace=pgo,resources=configmaps,verbs=get;list;watch;patch

// PostgreSQLCluster is the Schema for the postgresqlclusters API
type PostgreSQLCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgreSQLClusterSpec   `json:"spec,omitempty"`
	Status PostgreSQLClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PostgreSQLClusterList contains a list of PostgreSQLClusterSpec
type PostgreSQLClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgreSQLCluster `json:"items"`
}

type PgStorageSpec struct {
	Name               string `json:"name"`
	StorageClass       string `json:"storageclass"`
	AccessMode         string `json:"accessmode"`
	Size               string `json:"size"`
	StorageType        string `json:"storagetype"`
	SupplementalGroups string `json:"supplementalgroups"`
	MatchLabels        string `json:"matchLabels"`
}

func init() {
	SchemeBuilder.Register(&PostgreSQLCluster{}, &PostgreSQLClusterList{})
}
