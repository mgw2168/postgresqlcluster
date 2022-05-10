package cluster

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestSecretFromPgcluster(t *testing.T) {
	//pgcluster := &v1alpha1.PostgreSQLCluster{}
	//pgcluster.Name = "radondb-ukcv10"
	//pgcluster.Namespace = "dev"
	//pgcluster.Spec.RestoreFrom = "dmp-mds-radondb-ukcv09-to-radondb-ukcv10"
	//pgcluster.Spec.BackrestS3Key = "b3BlbnBpdHJpeG1pbmlvYWNjZXNza2V5"
	//pgcluster.Spec.BackrestS3KeySecret = "b3BlbnBpdHJpeG1pbmlvc2VjcmV0a2V5"
	//
	//s := SecretFromPgcluster(pgcluster)
	//
	//t.Log(s)

	reg := regexp.MustCompile(fmt.Sprintf(restoreFromPattern, "radondb-ukcv10"))
	result := reg.FindAllStringSubmatch("dmp-mds-radondb-ukcv09-to-radondb-ukcv10", -1)
	t.Log(result[0][1])

	t.Log(time.Now().Format("20060102-150405"))
}
