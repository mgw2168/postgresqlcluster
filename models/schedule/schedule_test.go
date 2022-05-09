package schedule

import (
	"github.com/kubesphere/api/v1alpha1"
	"testing"
)

func TestCreateSchedule(t *testing.T) {
	pgcluster := &v1alpha1.PostgreSQLCluster{}

	pgcluster.Name = "radondb-ukcv10"
	pgcluster.Namespace = "dev"
	pgcluster.Spec.FullBackupSchedule = "*/5 * * * *"

	err := CreateSchedule(pgcluster, "full")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}

func TestDeleteSchedule(t *testing.T) {
	pgcluster := &v1alpha1.PostgreSQLCluster{}

	pgcluster.Name = "radondb-ukcv10"
	pgcluster.Namespace = "dev"
	pgcluster.Spec.FullBackupSchedule = ""

	err := DeleteSchedule(pgcluster, "full")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}
