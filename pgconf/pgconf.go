package pgconf

import (
	"context"
	"fmt"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/k8sclient"
	"github.com/kubesphere/models/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"reflect"
	"sigs.k8s.io/yaml"
	"time"
)

const (
	PGHAConfigMapName = "%s-pgha-config"
	PGHADCSConfigName = "%s-dcs-config"
)

var postMasterItem = map[string]string{
	"shared_buffers":            "shared_buffers",
	"max_wal_senders":           "max_wal_senders",
	"max_connections":           "max_connections",
	"wal_buffers":               "wal_buffers",
	"max_replication_slots":     "max_replication_slots",
	"max_prepared_transactions": "max_prepared_transactions",
	"wal_level":                 "wal_level",
	"port":                      "port",
}

// DCSConfig represents the cluster-wide configuration that is stored in the Distributed
// Configuration Store (DCS).
type DCSConfig struct {
	LoopWait              int                `json:"loop_wait,omitempty"`
	TTL                   int                `json:"ttl,omitempty"`
	RetryTimeout          int                `json:"retry_timeout,omitempty"`
	MaximumLagOnFailover  int                `json:"maximum_lag_on_failover,omitempty"`
	MasterStartTimeout    int                `json:"master_start_timeout,omitempty"`
	SynchronousMode       bool               `json:"synchronous_mode,omitempty"`
	SynchronousModeStrict bool               `json:"synchronous_mode_strict,omitempty"`
	PostgreSQL            *PostgresDCS       `json:"postgresql,omitempty"`
	StandbyCluster        *StandbyDCS        `json:"standby_cluster,omitempty"`
	Slots                 map[string]SlotDCS `json:"slots,omitempty"`
}

// PostgresDCS represents the PostgreSQL settings that can be applied cluster-wide to a
// PostgreSQL cluster via the DCS.
type PostgresDCS struct {
	UsePGRewind  bool                   `json:"use_pg_rewind,omitempty"`
	UseSlots     bool                   `json:"use_slots,omitempty"`
	RecoveryConf map[string]interface{} `json:"recovery_conf,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// StandbyDCS represents standby cluster settings that can be applied cluster-wide via the DCS.
type StandbyDCS struct {
	Host                  string                 `json:"host,omitempty"`
	Port                  int                    `json:"port,omitempty"`
	PrimarySlotName       map[string]interface{} `json:"primary_slot_name,omitempty"`
	CreateReplicaMethods  []string               `json:"create_replica_methods,omitempty"`
	RestoreCommand        string                 `json:"restore_command,omitempty"`
	ArchiveCleanupCommand string                 `json:"archive_cleanup_command,omitempty"`
	RecoveryMinApplyDelay int                    `json:"recovery_min_apply_delay,omitempty"`
}

// SlotDCS represents slot settings that can be applied cluster-wide via the DCS.
type SlotDCS struct {
	Type     string `json:"type,omitempty"`
	Database string `json:"database,omitempty"`
	Plugin   string `json:"plugin,omitempty"`
}

func MergeConfig(newObj *v1alpha1.PostgreSQLCluster) error {
	configMapName := fmt.Sprintf(PGHAConfigMapName, newObj.GetName())
	dcsConfigName := fmt.Sprintf(PGHADCSConfigName, newObj.GetName())
	client := k8sclient.GetKubernetesClient()

	var updateItems []string

	clusterConfig, err := client.CoreV1().
		ConfigMaps(newObj.GetNamespace()).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	dcsConf := &DCSConfig{}
	if err = yaml.Unmarshal([]byte(clusterConfig.Data[dcsConfigName]), dcsConf); err != nil {
		return err
	}

	newParams := &PostgresDCS{}
	if err = yaml.Unmarshal([]byte(newObj.Spec.ClusterConfig), newParams); err != nil {
		return err
	}

	for k, v := range newParams.Parameters {
		if p, ok := dcsConf.PostgreSQL.Parameters[k]; ok {
			if !reflect.DeepEqual(v, p) && reflect.TypeOf(v) == reflect.TypeOf(p) {
				dcsConf.PostgreSQL.Parameters[k] = newParams.Parameters[k]
				klog.Infof("update existed config item:%s from %v to %v", k, p, v)
				updateItems = append(updateItems, k)
			}
		} else {
			dcsConf.PostgreSQL.Parameters[k] = newParams.Parameters[k]
			klog.Infof("update a new config item:%s to %v", k, p, v)
			updateItems = append(updateItems, k)
		}
	}

	// TODO support modify recovery_conf

	content, err := yaml.Marshal(dcsConf)
	if err != nil {
		klog.Errorf("unable to marshal dcsconf, error: %s", err)
		return err
	}

	jsonOpBytes, err := NewJSONPatch().Replace("data", dcsConfigName)(string(content)).Bytes()
	if err != nil {
		return err
	}

	_, err = client.CoreV1().ConfigMaps(clusterConfig.GetNamespace()).Patch(context.TODO(), clusterConfig.GetName(), types.JSONPatchType, jsonOpBytes, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	if hasPostMasterItem(updateItems) {
		// waiting pgha-config be applied by pgo
		klog.Info("to make some config items applied, cluster will restart after 3s")
		time.Sleep(3 * time.Second)

		err = cluster.RestartCluster(newObj)
		if err != nil {
			klog.Errorf("restart cluster error: %s", err)
			return err
		}
	}

	return nil
}

func hasPostMasterItem(updateItems []string) bool {
	for _, v := range updateItems {
		if _, ok := postMasterItem[v]; ok {
			return true
		}
	}
	return false
}
