package models

import "github.com/kubesphere/api/v1alpha1"

func InSlice(pgCluster *v1alpha1.PostgreSQLCluster, username string) bool {
	for _, u := range pgCluster.Spec.Users {
		if u.UserName == username {
			return true
		}
	}
	return false
}
