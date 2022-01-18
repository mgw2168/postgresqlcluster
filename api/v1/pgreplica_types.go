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

package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PgreplicaSpec defines the desired state of Pgreplica
type PgreplicaSpec struct {
	Name           string        `json:"name"`
	ClusterName    string        `json:"clustername"`
	ReplicaStorage PgStorageSpec `json:"replicastorage"`
	// ServiceType references the type of Service that should be used when
	// deploying PostgreSQL instances
	ServiceType v1.ServiceType    `json:"serviceType"`
	Status      string            `json:"status"`
	UserLabels  map[string]string `json:"userlabels"`
	// NodeAffinity is an optional structure that dictates how an instance should
	// be deployed in an environment
	NodeAffinity *v1.NodeAffinity `json:"nodeAffinity"`
	// Tolerations are an optional list of Pod toleration rules that are applied
	// to the PostgreSQL instance.
	Tolerations []v1.Toleration `json:"tolerations"`
}

// PgreplicaStatus defines the observed state of Pgreplica
type PgreplicaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	State   PgreplicaState `json:"state,omitempty"`
	Message string         `json:"message,omitempty"`
}

// PgreplicaState ...
// swagger:ignore
type PgreplicaState string

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pgreplica is the Schema for the pgreplicas API
type Pgreplica struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PgreplicaSpec   `json:"spec,omitempty"`
	Status PgreplicaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PgreplicaList contains a list of Pgreplica
type PgreplicaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pgreplica `json:"items"`
}

func (in *PgreplicaList) IsReplicaExisted(replicaName string) bool {
	for _, r := range in.Items {
		if r.GetName() == replicaName {
			return true
		}
	}
	return false
}

func init() {
	SchemeBuilder.Register(&Pgreplica{}, &PgreplicaList{})
}
