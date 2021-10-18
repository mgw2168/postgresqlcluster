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
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// PostgreSQLClusterStatus defines the observed state of PostgreSQLCluster
type PostgreSQLClusterStatus struct {
	Condition []ApiResult `json:"condition,omitempty"`
	State     string      `json:"state,omitempty"`
}

// ApiResult defines the result of pg operator ApiServer
type ApiResult struct {
	Api  string `json:"api,omitempty"`
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data string `json:"data,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

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

func init() {
	SchemeBuilder.Register(&PostgreSQLCluster{}, &PostgreSQLClusterList{})
}
