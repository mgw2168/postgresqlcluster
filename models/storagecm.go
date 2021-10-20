package models

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	scv1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/yaml"
)

const (
	CustomConfigMapName = "pgo-config"
	PATH                = "pgo.yaml"
)

// var StorageSpec crv1.PgStorageSpec

type PgoConfig struct {
	BasicAuth       string                   `json:"BasicAuth"`
	Cluster         ClusterStruct            `json:"Cluster"`
	Pgo             PgoStruct                `json:"Pgo"`
	PrimaryStorage  string                   `json:"PrimaryStorage"`
	WALStorage      string                   `json:"WALStorage"`
	BackupStorage   string                   `json:"BackupStorage"`
	ReplicaStorage  string                   `json:"ReplicaStorage"`
	BackrestStorage string                   `json:"BackrestStorage"`
	PGAdminStorage  string                   `json:"PGAdminStorage"`
	Storage         map[string]StorageStruct `json:"Storage"`
	OpenShift       bool                     `json:"OpenShift"`
}
type StorageStruct struct {
	AccessMode         string
	Size               string
	StorageType        string
	StorageClass       string
	SupplementalGroups string
	MatchLabels        string
}
type PgoStruct struct {
	Audit                          bool
	ConfigMapWorkerCount           *int
	ControllerGroupRefreshInterval *int
	DisableReconcileRBAC           bool
	NamespaceRefreshInterval       *int
	NamespaceWorkerCount           *int
	PGClusterWorkerCount           *int
	PGOImagePrefix                 string
	PGOImageTag                    string
	PGReplicaWorkerCount           *int
	PGTaskWorkerCount              *int
}
type ClusterStruct struct {
	CCPImagePrefix                 string
	CCPImageTag                    string
	Policies                       string
	Metrics                        bool
	Badger                         bool
	Port                           string
	PGBadgerPort                   string
	ExporterPort                   string
	User                           string
	Database                       string
	PasswordAgeDays                string
	PasswordLength                 string
	Replicas                       string
	ServiceType                    v1.ServiceType
	BackrestPort                   int
	BackrestGCSBucket              string
	BackrestGCSEndpoint            string
	BackrestGCSKeyType             string
	BackrestS3Bucket               string
	BackrestS3Endpoint             string
	BackrestS3Region               string
	BackrestS3URIStyle             string
	BackrestS3VerifyTLS            string
	DisableAutofail                bool
	DisableReplicaStartFailReinit  bool
	PodAntiAffinity                string
	PodAntiAffinityPgBackRest      string
	PodAntiAffinityPgBouncer       string
	SyncReplication                bool
	DefaultInstanceResourceMemory  resource.Quantity `json:"DefaultInstanceMemory"`
	DefaultBackrestResourceMemory  resource.Quantity `json:"DefaultBackrestMemory"`
	DefaultPgBouncerResourceMemory resource.Quantity `json:"DefaultPgBouncerMemory"`
	DefaultExporterResourceMemory  resource.Quantity `json:"DefaultExporterMemory"`
	DisableFSGroup                 *bool
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

func getOperatorConfigMap(clientset kubernetes.Interface, namespace string) (*v1.ConfigMap, error) {
	ctx := context.TODO()

	return clientset.CoreV1().ConfigMaps(namespace).Get(ctx, CustomConfigMapName, metav1.GetOptions{})
}

func updateOperatorConfigMap(clientset kubernetes.Interface, namespace string, Cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	ctx := context.TODO()

	return clientset.CoreV1().ConfigMaps(namespace).Update(ctx, Cm, metav1.UpdateOptions{})
}
func DelPod(clientset kubernetes.Interface, namespace string, labels map[string]string) error {
	ctx := context.TODO()
	var options metav1.ListOptions
	if labels != nil {
		options.LabelSelector = fields.Set(labels).String()
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		list, err := clientset.CoreV1().Pods(namespace).List(ctx, options)
		if list == nil && err != nil {
			list = &v1.PodList{}
		}
		for _, p := range list.Items {
			err := clientset.CoreV1().Pods(namespace).Delete(ctx, p.Name, metav1.DeleteOptions{})
			if err != nil {
				log.Fatal(err)
			}
		}
		return err
	})
	if retryErr != nil {
		log.Fatal("Del failed: %v", retryErr)
	}

	return nil
}
func RestartPod(clientset kubernetes.Interface, namespace string) error {
	ctx := context.TODO()
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
	_, err := clientset.AppsV1().Deployments(namespace).Patch(ctx, "postgres-operator", types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (c *PgoConfig) GetStorageSpec(name string) (StorageStruct, error) {
	var err error
	storage := StorageStruct{}

	s, ok := c.Storage[name]
	if !ok {
		err = errors.New("invalid Storage name " + name)
		log.Error(err)
		return storage, err
	}

	storage.StorageClass = s.StorageClass
	storage.AccessMode = s.AccessMode
	storage.Size = s.Size
	storage.StorageType = s.StorageType
	storage.MatchLabels = s.MatchLabels
	storage.SupplementalGroups = s.SupplementalGroups

	if storage.MatchLabels != "" {
		test := strings.Split(storage.MatchLabels, "=")
		if len(test) != 2 {
			err = errors.New("invalid Storage config " + name + " MatchLabels needs to be in key=value format.")
			log.Error(err)
			return storage, err
		}
	}

	return storage, err
}
func (c *PgoConfig) GenStorageSpec(name string, scClass string) (StorageStruct, error) {
	var err error
	storage := StorageStruct{}

	storage.StorageClass = scClass
	storage.AccessMode = "ReadWriteOnce"
	storage.Size = "1Gi"
	storage.StorageType = "dynamic"

	if storage.MatchLabels != "" {
		test := strings.Split(storage.MatchLabels, "=")
		if len(test) != 2 {
			err = errors.New("invalid Storage config " + name + " MatchLabels needs to be in key=value format.")
			log.Error(err)
			return storage, err
		}
	}

	return storage, err
}

func (c *PgoConfig) GetConfig(clientset kubernetes.Interface, namespace string) (*PgoConfig, error) {
	cMap, err := getOperatorConfigMap(clientset, namespace)
	if err != nil {
		log.Errorf("could not get ConfigMap: %s", err.Error())
		return nil, err
	}
	str := cMap.Data[PATH]
	yamlFile := []byte(str)
	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		log.Errorf("Unmarshal: %v", err)
		return nil, err
	}
	return c, err
}

func (c *PgoConfig) UpdateCm(clientset kubernetes.Interface, namespace string, sc *scv1.StorageClass) (*v1.ConfigMap, error) {
	cMap, err := getOperatorConfigMap(clientset, namespace)
	if err != nil {
		log.Errorf("could not get ConfigMap: %s", err.Error())
		return nil, err
	}
	scName := sc.Name
	str := cMap.Data[PATH]
	scyamlFile := []byte(str)
	if err := yaml.Unmarshal(scyamlFile, c); err != nil {
		log.Errorf("Unmarshal: %v", err)
		return nil, err
	}

	n := 0
	for k, v := range c.Storage {
		// StorageSpec, err = c.GetStorageSpec(k)
		if err != nil {
			return nil, err
		}
		if k == scName {
			break
		}
		log.Debug("Configd: %s  %s ", v.StorageClass, k)
		n = n + 1
		if n == len(c.Storage) {
			newsc, err := c.GenStorageSpec(scName, scName)
			if err != nil {
				log.Errorf("AddStorageSpec: %v", err)
			}
			c.Storage[scName] = newsc
			log.Infof("Config Storage class: %q to Configmap pgo-config ", scName)
			if err := RestartPod(clientset, namespace); err != nil {
				log.Error(err)
			}
		}
	}
	//填充cm
	cData, _ := yaml.Marshal(c)
	cMap.Data[PATH] = string(cData)
	cM, err := updateOperatorConfigMap(clientset, namespace, cMap)
	if err != nil {
		log.Error("Config: %v x", err)
	}
	return cM, err
}
