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
	// ******* create cluster
	Name           string `json:"name,omitempty"`
	ClusterName    string `json:"clustername,omitempty"`
	Policies       string `json:"policies,omitempty"`
	CCPImage       string `json:"ccpimage,omitempty"`
	CCPImageTag    string `json:"ccpimagetag,omitempty"`
	CCPImagePrefix string `json:"ccpimageprefix,omitempty"`
	PGOImagePrefix string `json:"pgoimageprefix,omitempty"`
	Port           string `json:"port,omitempty"`

	DisableAutofail bool `json:"disableAutofail,omitempty"`
	// PGBadger, if set to true, enables the pgBadger sidecar
	PGBadger     bool   `json:"pgBadger,omitempty"`
	PGBadgerPort string `json:"pgbadgerport,omitempty"`
	// Exporter, if set to true, enables the exporter sidecar
	Exporter     bool   `json:"exporter,omitempty"`
	ExporterPort string `json:"exporterport,omitempty"`

	Limits v1.ResourceList `json:"limits,omitempty"`
	// BackrestResources, if specified, contains the container request resources
	// for the pgBackRest Deployment for this PostgreSQL cluster
	BackrestResources v1.ResourceList `json:"backrestResources,omitempty"`
	// BackrestLimits, if specified, contains the container resource limits
	// for the pgBackRest Deployment for this PostgreSQL cluster
	BackrestLimits v1.ResourceList `json:"backrestLimits,omitempty"`
	// ExporterResources, if specified, contains the container request resources
	// for the RadonDB Postgres Exporter Deployment for this PostgreSQL cluster
	ExporterResources v1.ResourceList `json:"exporterResources,omitempty"`
	// ExporterLimits, if specified, contains the container resource limits
	// for the RadonDB Postgres Exporter Deployment for this PostgreSQL cluster
	ExporterLimits v1.ResourceList `json:"exporterLimits"`

	// PgBouncer contains all of the settings to properly maintain a pgBouncer
	// implementation
	//PgBouncer           PgBouncerSpec         `json:"pgBouncer"`
	User                string                `json:"user,omitempty"`
	PasswordType        string                `json:"passwordType,omitempty"`
	Database            string                `json:"database,omitempty"`
	Replicas            string                `json:"replicas,omitempty"`
	Status              string                `json:"status,omitempty"`
	CustomConfig        string                `json:"customconfig,omitempty"`
	UserLabels          map[string]string     `json:"userlabels,omitempty"`
	//NodeAffinity        NodeAffinitySpec      `json:"nodeAffinity"`
	//PodAntiAffinity     PodAntiAffinitySpec   `json:"podAntiAffinity"`
	SyncReplication     *bool                 `json:"syncReplication,omitempty,omitempty"`
	BackrestConfig      []v1.VolumeProjection `json:"backrestConfig,omitempty,omitempty"`
	BackrestGCSBucket   string                `json:"backrestGCSBucket,omitempty,omitempty"`
	BackrestGCSEndpoint string                `json:"backrestGCSEndpoint,omitempty,omitempty"`
	BackrestGCSKeyType  string                `json:"backrestGCSKeyType,omitempty,omitempty"`
	BackrestS3Bucket    string                `json:"backrestS3Bucket,omitempty,omitempty"`
	BackrestS3Region    string                `json:"backrestS3Region,omitempty,omitempty"`
	BackrestS3Endpoint  string                `json:"backrestS3Endpoint,omitempty,omitempty"`
	BackrestS3URIStyle  string                `json:"backrestS3URIStyle,omitempty,omitempty"`
	BackrestS3VerifyTLS string                `json:"backrestS3VerifyTLS,omitempty,omitempty"`
	BackrestRepoPath    string                `json:"backrestRepoPath,omitempty,omitempty"`

	TLSOnly              bool                     `json:"tlsOnly,omitempty"`
	Standby              bool                     `json:"standby,omitempty"`
	Shutdown             bool                     `json:"shutdown,omitempty"`
	//PGDataSource         PGDataSourceSpec         `json:"pgDataSource,omitempty"`

	// Annotations contains a set of Deployment (and by association, Pod)
	// annotations that are propagated to all managed Deployments
	//Annotations ClusterAnnotations `json:"annotations,omitempty"`

	// ServiceType references the type of Service that should be used when
	// deploying PostgreSQL instances
	ServiceType v1.ServiceType `json:"serviceType,omitempty"`

	// Tolerations are an optional list of Pod toleration rules that are applied
	// to the PostgreSQL instance.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// ******* update cluster
	
}

// PostgreSQLClusterStatus defines the observed state of PostgreSQLCluster
type PostgreSQLClusterStatus struct {
	Condition []string `json:"condition"`
	State     string   `json:"state"`
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

// PostgreSQLClusterList contains a list of PostgreSQLCluster
type PostgreSQLClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgreSQLCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PostgreSQLCluster{}, &PostgreSQLClusterList{})
}
