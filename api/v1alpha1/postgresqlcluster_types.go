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
	// Deprecated
	Action string `json:"action,omitempty"`

	// Deprecated: move to ClusterConfig
	SyncReplication *bool `json:"syncReplication,omitempty"`

	// Deprecated
	Database string `json:"database,omitempty"`
	// Deprecated
	Username string `json:"username,omitempty"`
	// Deprecated
	Password string `json:"password,omitempty"`

	// Deprecated
	ClusterName []string `json:"clusterName,omitempty"`
	// Deprecated
	AutoFail int `json:"autofail,omitempty"`
	// Deprecated
	Startup bool `json:"startup,omitempty"`
	// Deprecated
	Shutdown bool `json:"shutdown,omitempty"`

	// ******** delete cluster
	// Deprecated
	Selector string `json:"selector,omitempty"`
	// Deprecated
	AllFlag bool `json:"allFlag,omitempty"`
	// Deprecated
	DeleteBackups bool `json:"deleteBackups,omitempty"`
	// Deprecated
	DeleteData bool `json:"deleteData,omitempty"`

	// ******** scale cluster
	// Deprecated
	NodeLabel string `json:"nodeLabel,omitempty"`
	// Deprecated
	ServiceType string `json:"serviceType,omitempty"`
	// Deprecated
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// Deprecated
	ScaleDownDeleteData string `json:"delete-data,omitempty"`

	// ******** restart cluster
	// Deprecated
	Restart bool `json:"restart,omitempty"`
	// Deprecated
	Targets []string `json:"targets,omitempty"`

	// Deprecated
	RollingUpdate bool `json:"rollingUpdate,omitempty"`
	// Deprecated
	PodAnnotation []map[string]string `json:"podAnnotation,omitempty"`

	// Deprecated
	ManagedUser bool `json:"managedUser,omitempty"`
	// Deprecated
	PasswordAgeDays int `json:"passwordAgeDays,omitempty"`
	// Deprecated
	PasswordLength int `json:"passwordLength,omitempty"`
	// Deprecated
	PasswordType string `json:"passwordType,omitempty"`

	// ******** show user
	// Deprecated
	ShowSystemAccounts bool `json:"showSystemAccounts,omitempty"`

	// ******** update user
	// Deprecated
	SetSystemAccountPassword bool `json:"setSystemAccountPassword,omitempty"`

	// ******** The above fields are deprecated or unused ********

	// ******** create cluster params ********
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	CCPImage      string `json:"ccpImage,omitempty"`
	CCPImageTag   string `json:"ccpImageTag,omitempty"`
	PgVersion     string `json:"pgVersion,omitempty"`
	ReplicaCount  int    `json:"replicaCount,omitempty"`
	CPULimit      string `json:"cpuLimit,omitempty"`
	CPURequest    string `json:"cpuRequest,omitempty"`
	MemoryLimit   string `json:"memoryLimit,omitempty"`
	MemoryRequest string `json:"memoryRequest,omitempty"`

	StorageConfig string `json:"storageConfig,omitempty"`

	// ******** update cluster
	PVCSize string `json:"pvcSize,omitempty"`

	// ******** scale down
	ReplicaName string `json:"replicaName,omitempty"`

	// ******** create user
	Users []User `json:"users,omitempty"`

	// pgconf
	ClusterConfig string `json:"postgresqlParams,omitempty"`

	// like "1653299584|full","1653299584|incr", use timestamp to mark change
	PerformBackup string `json:"performBackup,omitempty"`
	// like "1653299584|backup1|backup2|backup3", use timestamp to mark change
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
