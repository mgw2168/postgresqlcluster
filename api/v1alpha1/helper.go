package v1alpha1

func (in *PostgreSQLCluster) HasValidS3Conf() bool {
	if in.Spec.BackrestS3Key != "" && in.Spec.BackrestS3KeySecret != "" && in.Spec.BackrestS3Endpoint != "" &&
		in.Spec.BackrestS3Region != "" && in.Spec.BackrestS3Bucket != "" {
		return true
	}
	return false
}
