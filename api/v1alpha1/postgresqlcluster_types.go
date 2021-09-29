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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PostgreSQLClusterSpec defines the desired state of PostgreSQLCluster
type PostgreSQLClusterSpec struct {
	// ******** create cluster params ********
	ClientVersion   string `json:"clientVersion"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	SyncReplication *bool  `json:"syncReplication"`
	CCPImage        string `json:"ccpImage"`
	CCPImageTag     string `json:"ccpImageTag"`
	ReplicaCount    int    `json:"replicaCount"`
	CPULimit        string `json:"cpuLimit"`
	CPURequest      string `json:"cpuRequest"`
	MemoryLimit     string `json:"memoryLimit"`
	MemoryRequest   string `json:"memoryRequest"`
	Database        string `json:"database"`
	Username        string `json:"username"`
	Password        string `json:"password"`

	// ******** update cluster

	// ******** delete cluster
	ClusterName   string `json:"clusterName"`
	Selector      string `json:"selector"`
	AllFlag       bool   `json:"allFlag"`
	DeleteBackups bool   `json:"deleteBackups"`
	DeleteData    bool   `json:"deleteData"`
}

// PostgreSQLClusterStatus defines the observed state of PostgreSQLCluster
type PostgreSQLClusterStatus struct {
	//Condition []string `json:"condition"`
	State   string `json:"state,omitempty"`
	Version string `json:"version,omitempty"`
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
