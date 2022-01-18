package controllers

import (
	"context"
	v1 "github.com/kubesphere/api/v1"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/models/cluster"
	"github.com/kubesphere/models/user"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r PostgreSQLClusterReconciler) compareWithPGCluster(ctx context.Context, desiredPGC *v1alpha1.PostgreSQLCluster) (bool, error) {
	currentPGC := v1.Pgcluster{}
	// make sure pgcluster CR is existed
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: desiredPGC.Namespace,
		Name:      desiredPGC.Name,
	}, &currentPGC); err != nil {
		return false, err
	}

	// list all replica for current cluster
	currentPGReplicas := v1.PgreplicaList{}
	listOpts := client.ListOptions{
		Namespace:     currentPGC.GetNamespace(),
		LabelSelector: labels.SelectorFromSet(map[string]string{"pg-cluster": currentPGC.GetName()}),
	}
	if err := r.Client.List(ctx, &currentPGReplicas, &listOpts); err != nil {
		return false, client.IgnoreNotFound(err)
	}

	if currentPGC.Spec.PrimaryStorage.Size != "" && currentPGC.Spec.PrimaryStorage.Size != desiredPGC.Spec.PVCSize {
		if err := cluster.UpdatePgCluster(desiredPGC, true); err != nil {
			klog.Errorf("update pvc error: %s", err)
			return false, err
		}
		klog.Infof("pvc updated to:%s", desiredPGC.Spec.PVCSize)
	}

	// update cpu/mem
	if !currentPGC.Spec.Resources.Cpu().Equal(resource.MustParse(desiredPGC.Spec.CPURequest)) ||
		!currentPGC.Spec.Limits.Cpu().Equal(resource.MustParse(desiredPGC.Spec.CPULimit)) ||
		!currentPGC.Spec.Resources.Memory().Equal(resource.MustParse(desiredPGC.Spec.MemoryRequest)) ||
		!currentPGC.Spec.Limits.Memory().Equal(resource.MustParse(desiredPGC.Spec.MemoryLimit)) {
		if err := cluster.UpdatePgCluster(desiredPGC, false); err != nil {
			klog.Errorf("update cpu and memory error: %s", err.Error())
			return false, err
		}
		klog.Infof("cpu updated to:%s/%s, memory update to:%s:%s",
			desiredPGC.Spec.CPURequest, desiredPGC.Spec.CPULimit,
			desiredPGC.Spec.MemoryRequest, desiredPGC.Spec.MemoryLimit)
	}

	// scale out
	if len(currentPGReplicas.Items) < desiredPGC.Spec.ReplicaCount {
		replicaCount := desiredPGC.Spec.ReplicaCount - len(currentPGReplicas.Items)
		if err := cluster.ScaleUpPgCluster(desiredPGC, replicaCount); err != nil {
			klog.Errorf("scale up error: %s", err.Error())
			return false, err
		}
		klog.Infof("replica updated to:%d", desiredPGC.Spec.ReplicaCount)
	}

	// ReplicaName the replica will to be deleted
	if desiredPGC.Spec.ReplicaName != "" {
		if !currentPGReplicas.IsReplicaExisted(desiredPGC.Spec.ReplicaName) {
			desiredPGC.Spec.ReplicaName = ""
			return true, nil
		}
		if err := cluster.ScaleDownPgCluster(desiredPGC); err != nil {
			klog.Errorf("scale down error: %s", err.Error())
			return false, err
		}
		klog.Infof("replica %s deleted", desiredPGC.Spec.ReplicaName)
		desiredPGC.Spec.ReplicaName = ""
		return true, nil
	}

	// restart
	if desiredPGC.Spec.Restart {
		if err := cluster.RestartCluster(desiredPGC); err != nil {
			klog.Errorf("restart cluster error: %s", err.Error())
			return false, err
		}
		klog.Infof("cluster restarted")
		desiredPGC.Spec.Restart = false
		return true, nil
	}

	return false, nil
}

func (r PostgreSQLClusterReconciler) updateUserList(desiredPGC *v1alpha1.PostgreSQLCluster) error {
	currentUserList, err := user.GetExistedUserList(desiredPGC)
	if err != nil {
		return err
	}

	// delete users which are not in Spec.Users/Spec.Username
	for _, u := range currentUserList.Results {
		if u.Username == "postgres" || u.Username == "ccp_monitoring" || u.Username == "primaryuser" {
			continue
		}
		if !desiredPGC.IsUserExisted(u.Username) {
			err = user.DeletePgUser(desiredPGC, u.Username)
			if err != nil {
				klog.Errorf("delete user error: %s", err.Error())
				return err
			}
			klog.Infof("user:%s not existed in spec, be deleted", u.Username)
			continue
		}
	}

	userListCopy := desiredPGC.Spec.Users
	userListCopy = append(userListCopy, v1alpha1.User{UserName: desiredPGC.Spec.Username, Password: desiredPGC.Spec.Password})
	// update users in Spec.Users/Spec.Username into PGC
	for _, u := range userListCopy {
		if u.UserName == "postgres" || u.UserName == "ccp_monitoring" || u.UserName == "primaryuser" {
			continue
		}

		index, existed := currentUserList.IsUserExisted(u)
		// create user
		if !existed {
			err = user.CreatePgUser(desiredPGC, u.UserName, u.Password)
			if err != nil {
				klog.Errorf("create user error: %s", err.Error())
				return err
			}
			klog.Infof("user:%s not existed in cluster, be created", u.UserName)
			continue
		}

		needUpdate := currentUserList.Results[index].NeedUpdate(u)
		// update user password
		if needUpdate {
			err = user.UpdatePgUser(desiredPGC, u.UserName, u.Password)
			if err != nil {
				klog.Errorf("update password error: %s", err.Error())
				return err
			}
			klog.Infof("user:%s password is updated", u.UserName)
			continue
		}
	}

	return nil
}

func (r *PostgreSQLClusterReconciler) fetchStateThenReturn(ctx context.Context, fakePGC *v1alpha1.PostgreSQLCluster) (ctrl.Result, error) {
	pgInstance := v1.Pgcluster{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: fakePGC.Namespace,
		Name:      fakePGC.Name,
	}, &pgInstance)
	if err != nil {
		klog.Errorf("get pgcluster resource error: %s", err)
		return ctrl.Result{}, err
	}

	if string(pgInstance.Status.State) != "" {
		fakePGC.Status.State = string(pgInstance.Status.State)
	}

	if err := r.Client.Status().Update(ctx, fakePGC); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: checkTime}, nil
}

func (r *PostgreSQLClusterReconciler) updateSpecThenReturn(ctx context.Context, fakePGC *v1alpha1.PostgreSQLCluster) (ctrl.Result, error) {
	pgInstance := v1.Pgcluster{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Namespace: fakePGC.Namespace,
		Name:      fakePGC.Name,
	}, &pgInstance)
	if err != nil {
		klog.Errorf("get pgcluster resource error: %s", err)
	}

	if string(pgInstance.Status.State) != "" {
		fakePGC.Status.State = string(pgInstance.Status.State)
	}

	if err := r.Update(ctx, fakePGC); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: checkTime}, nil
}
