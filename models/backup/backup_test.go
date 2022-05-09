package backup

import (
	"encoding/json"
	"github.com/kubesphere/api/v1alpha1"
	"testing"
)

func TestShowBackup(t *testing.T) {
	pgcluster := &v1alpha1.PostgreSQLCluster{}
	pgcluster.Name = "radondb-ukcv10"
	pgcluster.Namespace = "dev"

	err := ShowBackup(pgcluster)
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := json.Marshal(pgcluster.Status.Backups)

	t.Log(string(bytes))

}

func TestPerformBackup(t *testing.T) {
	pgcluster := &v1alpha1.PostgreSQLCluster{}
	pgcluster.Name = "radondb-ukcv10"
	pgcluster.Namespace = "dev"
	pgcluster.Spec.PerformBackup = "1652081430|full"

	err := PerformBackup(pgcluster)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestDeleteBackup(t *testing.T) {
	pgcluster := &v1alpha1.PostgreSQLCluster{}
	pgcluster.Name = "radondb-ukcv10"
	pgcluster.Namespace = "dev"
	pgcluster.Spec.BackupToDelete = "1652081430|20220509-114613F"

	err := DeleteBackup(pgcluster)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}
