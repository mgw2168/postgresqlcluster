package models

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/core/v1"
	scv1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
	"strings"
	"sync"
	"time"
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
				klog.Fatal(err)
			}
		}
		return err
	})
	if retryErr != nil {
		klog.Fatalf("Del failed: %v", retryErr)
	}

	return nil
}
func RestartPod(clientset kubernetes.Interface, namespace string) error {
	ctx := context.TODO()
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
	_, err := clientset.AppsV1().Deployments(namespace).Patch(ctx, "postgres-operator", types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		klog.Error(err)
	}
	return nil
}

func (c *PgoConfig) GetStorageSpec(name string) (StorageStruct, error) {
	var err error
	storage := StorageStruct{}

	s, ok := c.Storage[name]
	if !ok {
		err = errors.New("invalid Storage name " + name)
		klog.Error(err)
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
			klog.Error(err)
			return storage, err
		}
	}

	return storage, err
}
func (c *PgoConfig) GenStorageSpec(scName string) (StorageStruct, error) {
	var err error
	storage := StorageStruct{}

	storage.StorageClass = scName
	storage.AccessMode = "ReadWriteOnce"
	storage.Size = "1Gi"
	storage.StorageType = "dynamic"

	if storage.MatchLabels != "" {
		test := strings.Split(storage.MatchLabels, "=")
		if len(test) != 2 {
			err = errors.New("invalid Storage config " + scName + " MatchLabels needs to be in key=value format.")
			klog.Error(err)
			return storage, err
		}
	}

	return storage, err
}

func (c *PgoConfig) GetConfig(clientset kubernetes.Interface, namespace string) (*PgoConfig, error) {
	cMap, err := getOperatorConfigMap(clientset, namespace)
	if err != nil {
		klog.Errorf("could not get ConfigMap: %s", err.Error())
		return nil, err
	}
	str := cMap.Data[PATH]
	yamlFile := []byte(str)
	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		klog.Errorf("Unmarshal: %v", err)
		return nil, err
	}
	return c, err
}

func (c *PgoConfig) UpdateCm(clientset kubernetes.Interface, namespace string, sc *scv1.StorageClass) (*v1.ConfigMap, error) {
	cMap, err := getOperatorConfigMap(clientset, namespace)
	if err != nil {
		klog.Errorf("could not get ConfigMap: %s", err.Error())
		return nil, err
	}
	scName := sc.Name
	str := cMap.Data[PATH]
	scyamlFile := []byte(str)
	if err := yaml.Unmarshal(scyamlFile, c); err != nil {
		klog.Errorf("Unmarshal: %v", err)
		return nil, err
	}
	if _, ok := c.Storage[scName]; ok {
		return nil, nil
	}
	newsc, err := c.GenStorageSpec(scName)
	if err != nil {
		klog.Errorf("AddStorageSpec: %v", err)
		return nil, nil
	}
	mx := &sync.Mutex{}
	mx.Lock()
	c.Storage[scName] = newsc
	mx.Unlock()
	klog.Infof("Config Storage class: %s to Configmap pgo-config", scName)

	//填充cm
	cData, _ := yaml.Marshal(c)
	cMap.Data[PATH] = string(cData)
	cM, err := updateOperatorConfigMap(clientset, namespace, cMap)
	if err != nil {
		klog.Errorf("Config: %v x", err)
	}

	klog.Infof("Config Storage class: %s to Configmap pgo-config ", scName)
	if err := RestartPod(clientset, namespace); err != nil {
		klog.Error(err)
	}

	return cM, err
}

func (c *PgoConfig) UpdateCmInformer(clientset kubernetes.Interface, namespace string, cm *v1.ConfigMap) (*v1.ConfigMap, error) {
	sc, err := clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorf("list storage error: %s", err)
	}
	str := cm.Data[PATH]
	scyamlFile := []byte(str)
	if err := yaml.Unmarshal(scyamlFile, c); err != nil {
		klog.Errorf("Unmarshal: %v", err)
		return nil, err
	}

	mx := &sync.Mutex{}
	for i := range sc.Items {
		if _, ok := c.Storage[sc.Items[i].Name]; ok {
			break
		}
		newSc, err := c.GenStorageSpec(sc.Items[i].Name)
		if err != nil {
			klog.Errorf("AddStorageSpec: %v", err)
		}

		mx.Lock()
		c.Storage[sc.Items[i].Name] = newSc
		mx.Unlock()
		klog.Infof("Config Storage class: %s to Configmap pgo-config ", sc.Items[i].Name)
	}

	if err := RestartPod(clientset, namespace); err != nil {
		klog.Error(err)
	}

	//填充cm
	cData, _ := yaml.Marshal(c)
	cm.Data[PATH] = string(cData)
	cM, err := updateOperatorConfigMap(clientset, namespace, cm)
	if err != nil {
		klog.Errorf("Config: %v x", err)
	}
	return cM, err
}
