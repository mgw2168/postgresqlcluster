package cluster

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/kubesphere/api/v1alpha1"
	"github.com/kubesphere/k8sclient"
	"github.com/kubesphere/pkg"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"regexp"
)

const authorizedKeys = `c3NoLWVkMjU1MTkgQUFBQUMzTnphQzFsWkRJMU5URTVBQUFBSUc5TmFCRW9lM1U4MVp2NjZndHQyMmw5ZHNlTkcvdmkrcEtIekx1NjhrQk0K`
const config = `SG9zdCAqCglTdHJpY3RIb3N0S2V5Q2hlY2tpbmcgbm8KCUlkZW50aXR5RmlsZSAvdG1wL2lkX2VkMjU1MTkKCVBvcnQgMjAyMgoJVXNlciBwZ2JhY2tyZXN0Cg==`
const idED25519 = `LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUFNd0FBQUF0egpjMmd0WldReU5UVXhPUUFBQUNCdlRXZ1JLSHQxUE5XYit1b0xiZHRwZlhiSGpSdjc0dnFTaDh5N3V2SkFUQUFBCkFJaDBCZkdpZEFYeG9nQUFBQXR6YzJndFpXUXlOVFV4T1FBQUFDQnZUV2dSS0h0MVBOV2IrdW9MYmR0cGZYYkgKalJ2NzR2cVNoOHk3dXZKQVRBQUFBRUIxMGlqVVYvc2JpZ2k0ekdGbU9PQmkzSGhGelZPZzljMjZNR0o4QjAwRgpZRzlOYUJFb2UzVTgxWnY2Nmd0dDIybDlkc2VORy92aStwS0h6THU2OGtCTUFBQUFBQUVDQXdRRgotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0K`
const sshHostED25519Key = `LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUFNd0FBQUF0egpjMmd0WldReU5UVXhPUUFBQUNCdlRXZ1JLSHQxUE5XYit1b0xiZHRwZlhiSGpSdjc0dnFTaDh5N3V2SkFUQUFBCkFJaDBCZkdpZEFYeG9nQUFBQXR6YzJndFpXUXlOVFV4T1FBQUFDQnZUV2dSS0h0MVBOV2IrdW9MYmR0cGZYYkgKalJ2NzR2cVNoOHk3dXZKQVRBQUFBRUIxMGlqVVYvc2JpZ2k0ekdGbU9PQmkzSGhGelZPZzljMjZNR0o4QjAwRgpZRzlOYUJFb2UzVTgxWnY2Nmd0dDIybDlkc2VORy92aStwS0h6THU2OGtCTUFBQUFBQUVDQXdRRgotLS0tLUVORCBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0K`
const sshdConfig = `IwkkT3BlbkJTRDogc3NoZF9jb25maWcsdiAxLjEwMCAyMDE2LzA4LzE1IDEyOjMyOjA0IG5hZGR5IEV4cCAkCgojIFRoaXMgaXMgdGhlIHNzaGQgc2VydmVyIHN5c3RlbS13aWRlIGNvbmZpZ3VyYXRpb24gZmlsZS4gIFNlZQojIHNzaGRfY29uZmlnKDUpIGZvciBtb3JlIGluZm9ybWF0aW9uLgoKIyBUaGlzIHNzaGQgd2FzIGNvbXBpbGVkIHdpdGggUEFUSD0vdXNyL2xvY2FsL2JpbjovdXNyL2JpbgoKIyBUaGUgc3RyYXRlZ3kgdXNlZCBmb3Igb3B0aW9ucyBpbiB0aGUgZGVmYXVsdCBzc2hkX2NvbmZpZyBzaGlwcGVkIHdpdGgKIyBPcGVuU1NIIGlzIHRvIHNwZWNpZnkgb3B0aW9ucyB3aXRoIHRoZWlyIGRlZmF1bHQgdmFsdWUgd2hlcmUKIyBwb3NzaWJsZSwgYnV0IGxlYXZlIHRoZW0gY29tbWVudGVkLiAgVW5jb21tZW50ZWQgb3B0aW9ucyBvdmVycmlkZSB0aGUKIyBkZWZhdWx0IHZhbHVlLgoKIyBJZiB5b3Ugd2FudCB0byBjaGFuZ2UgdGhlIHBvcnQgb24gYSBTRUxpbnV4IHN5c3RlbSwgeW91IGhhdmUgdG8gdGVsbAojIFNFTGludXggYWJvdXQgdGhpcyBjaGFuZ2UuCiMgc2VtYW5hZ2UgcG9ydCAtYSAtdCBzc2hfcG9ydF90IC1wIHRjcCAjUE9SVE5VTUJFUgojClBvcnQgMjAyMgojQWRkcmVzc0ZhbWlseSBhbnkKI0xpc3RlbkFkZHJlc3MgMC4wLjAuMAojTGlzdGVuQWRkcmVzcyA6OgoKSG9zdEtleSAvc3NoZC9zc2hfaG9zdF9lZDI1NTE5X2tleQoKIyBDaXBoZXJzIGFuZCBrZXlpbmcKI1Jla2V5TGltaXQgZGVmYXVsdCBub25lCgojIExvZ2dpbmcKI1N5c2xvZ0ZhY2lsaXR5IEFVVEgKU3lzbG9nRmFjaWxpdHkgQVVUSFBSSVYKI0xvZ0xldmVsIElORk8KCiMgQXV0aGVudGljYXRpb246CgojTG9naW5HcmFjZVRpbWUgMm0KUGVybWl0Um9vdExvZ2luIG5vClN0cmljdE1vZGVzIG5vCiNNYXhBdXRoVHJpZXMgNgojTWF4U2Vzc2lvbnMgMTAKClB1YmtleUF1dGhlbnRpY2F0aW9uIHllcwoKIyBUaGUgZGVmYXVsdCBpcyB0byBjaGVjayBib3RoIC5zc2gvYXV0aG9yaXplZF9rZXlzIGFuZCAuc3NoL2F1dGhvcml6ZWRfa2V5czIKIyBidXQgdGhpcyBpcyBvdmVycmlkZGVuIHNvIGluc3RhbGxhdGlvbnMgd2lsbCBvbmx5IGNoZWNrIC5zc2gvYXV0aG9yaXplZF9rZXlzCiNBdXRob3JpemVkS2V5c0ZpbGUJL3BnY29uZi9hdXRob3JpemVkX2tleXMKQXV0aG9yaXplZEtleXNGaWxlCS9zc2hkL2F1dGhvcml6ZWRfa2V5cwoKI0F1dGhvcml6ZWRQcmluY2lwYWxzRmlsZSBub25lCgojQXV0aG9yaXplZEtleXNDb21tYW5kIG5vbmUKI0F1dGhvcml6ZWRLZXlzQ29tbWFuZFVzZXIgbm9ib2R5CgojIEZvciB0aGlzIHRvIHdvcmsgeW91IHdpbGwgYWxzbyBuZWVkIGhvc3Qga2V5cyBpbiAvZXRjL3NzaC9zc2hfa25vd25faG9zdHMKI0hvc3RiYXNlZEF1dGhlbnRpY2F0aW9uIG5vCiMgQ2hhbmdlIHRvIHllcyBpZiB5b3UgZG9uJ3QgdHJ1c3Qgfi8uc3NoL2tub3duX2hvc3RzIGZvcgojIEhvc3RiYXNlZEF1dGhlbnRpY2F0aW9uCiNJZ25vcmVVc2VyS25vd25Ib3N0cyBubwojIERvbid0IHJlYWQgdGhlIHVzZXIncyB+Ly5yaG9zdHMgYW5kIH4vLnNob3N0cyBmaWxlcwojSWdub3JlUmhvc3RzIHllcwoKIyBUbyBkaXNhYmxlIHR1bm5lbGVkIGNsZWFyIHRleHQgcGFzc3dvcmRzLCBjaGFuZ2UgdG8gbm8gaGVyZSEKI1Bhc3N3b3JkQXV0aGVudGljYXRpb24geWVzCiNQZXJtaXRFbXB0eVBhc3N3b3JkcyBubwpQYXNzd29yZEF1dGhlbnRpY2F0aW9uIG5vCgojIENoYW5nZSB0byBubyB0byBkaXNhYmxlIHMva2V5IHBhc3N3b3JkcwpDaGFsbGVuZ2VSZXNwb25zZUF1dGhlbnRpY2F0aW9uIHllcwojQ2hhbGxlbmdlUmVzcG9uc2VBdXRoZW50aWNhdGlvbiBubwoKIyBLZXJiZXJvcyBvcHRpb25zCiNLZXJiZXJvc0F1dGhlbnRpY2F0aW9uIG5vCiNLZXJiZXJvc09yTG9jYWxQYXNzd2QgeWVzCiNLZXJiZXJvc1RpY2tldENsZWFudXAgeWVzCiNLZXJiZXJvc0dldEFGU1Rva2VuIG5vCiNLZXJiZXJvc1VzZUt1c2Vyb2sgeWVzCgojIEdTU0FQSSBvcHRpb25zCiNHU1NBUElBdXRoZW50aWNhdGlvbiB5ZXMKI0dTU0FQSUNsZWFudXBDcmVkZW50aWFscyBubwojR1NTQVBJU3RyaWN0QWNjZXB0b3JDaGVjayB5ZXMKI0dTU0FQSUtleUV4Y2hhbmdlIG5vCiNHU1NBUElFbmFibGVrNXVzZXJzIG5vCgojIFRoaXMgaXMgc2V0IGV4cGxpY2l0bHkgdG8gKm5vKiBhcyB3ZSBhcmUgb25seSB1c2luZyBwdWJrZXkgYXV0aGVudGljYXRpb24gYW5kCiMgYmVjYXVzZSBlYWNoIGNvbnRhaW5lciBpcyBpc29sYXRlZCB0byBvbmx5IGFuIHVucHJpdmlsZWdlZCB1c2VyLgpVc2VQQU0gbm8KCiNBbGxvd0FnZW50Rm9yd2FyZGluZyB5ZXMKI0FsbG93VGNwRm9yd2FyZGluZyB5ZXMKI0dhdGV3YXlQb3J0cyBubwpYMTFGb3J3YXJkaW5nIHllcwojWDExRGlzcGxheU9mZnNldCAxMAojWDExVXNlTG9jYWxob3N0IHllcwojUGVybWl0VFRZIHllcwojUHJpbnRNb3RkIHllcwojUHJpbnRMYXN0TG9nIHllcwojVENQS2VlcEFsaXZlIHllcwojVXNlTG9naW4gbm8KI1Blcm1pdFVzZXJFbnZpcm9ubWVudCBubwojQ29tcHJlc3Npb24gZGVsYXllZAojQ2xpZW50QWxpdmVJbnRlcnZhbCAwCiNDbGllbnRBbGl2ZUNvdW50TWF4IDMKI1Nob3dQYXRjaExldmVsIG5vCiNVc2VETlMgeWVzCiNQaWRGaWxlIC92YXIvcnVuL3NzaGQucGlkCiNNYXhTdGFydHVwcyAxMDozMDoxMDAKI1Blcm1pdFR1bm5lbCBubwojQ2hyb290RGlyZWN0b3J5IG5vbmUKI1ZlcnNpb25BZGRlbmR1bSBub25lCgojIG5vIGRlZmF1bHQgYmFubmVyIHBhdGgKI0Jhbm5lciBub25lCgojIEFjY2VwdCBsb2NhbGUtcmVsYXRlZCBlbnZpcm9ubWVudCB2YXJpYWJsZXMKQWNjZXB0RW52IExBTkcgTENfQ1RZUEUgTENfTlVNRVJJQyBMQ19USU1FIExDX0NPTExBVEUgTENfTU9ORVRBUlkgTENfTUVTU0FHRVMKQWNjZXB0RW52IExDX1BBUEVSIExDX05BTUUgTENfQUREUkVTUyBMQ19URUxFUEhPTkUgTENfTUVBU1VSRU1FTlQKQWNjZXB0RW52IExDX0lERU5USUZJQ0FUSU9OIExDX0FMTCBMQU5HVUFHRQpBY2NlcHRFbnYgWE1PRElGSUVSUwoKIyBvdmVycmlkZSBkZWZhdWx0IG9mIG5vIHN1YnN5c3RlbXMKU3Vic3lzdGVtCXNmdHAJL3Vzci9saWJleGVjL29wZW5zc2gvc2Z0cC1zZXJ2ZXIKCiMgRXhhbXBsZSBvZiBvdmVycmlkaW5nIHNldHRpbmdzIG9uIGEgcGVyLXVzZXIgYmFzaXMKI01hdGNoIFVzZXIgYW5vbmN2cwojCVgxMUZvcndhcmRpbmcgbm8KIwlBbGxvd1RjcEZvcndhcmRpbmcgbm8KIwlQZXJtaXRUVFkgbm8KIwlGb3JjZUNvbW1hbmQgY3ZzIHNlcnZlcgoKIyBlbnN1cmUgbnNzX3dyYXBwZXIgZW52IHZhcnMgYXJlIHNldCB3aGVuIGV4ZWN1dGluZyBjb21tYW5kcyBhcyBuZWVkZWQgZm9yIE9wZW5TaGlmdCBjb21wYXRpYmlsaXR5CkZvcmNlQ29tbWFuZCBOU1NfV1JBUFBFUl9TVUJESVI9c3NoIC4gL29wdC9yYWRvbmRiL2Jpbi9uc3Nfd3JhcHBlcl9lbnYuc2ggJiYgJFNTSF9PUklHSU5BTF9DT01NQU5ECg==`

const restoreFromPattern = `dmp-mds-(?s:(.*?))-to-%s`
const managedSecretNamePattern = `%s-backrest-repo-config`
const managedPVCNamePattern = `%s-pgbr-repo`
const managedDeploymentNamePattern = `%s-backrest-shared-repo`

func GetRestoreFromName(pg *v1alpha1.PostgreSQLCluster) string {
	reg := regexp.MustCompile(fmt.Sprintf(restoreFromPattern, pg.Name))
	result := reg.FindAllStringSubmatch(pg.Spec.RestoreFrom, -1)
	if result == nil || len(result[0]) < 2 {
		return pg.Spec.RestoreFrom
	}

	return shaName(pg.Spec.RestoreFrom)
}

func shaName(s string) string {
	return "dmp-mds-" + pkg.Sha1Str(s)
}

func SecretFromPgcluster(pg *v1alpha1.PostgreSQLCluster, repoPath string) *corev1.Secret {
	s := &corev1.Secret{}
	s.Name = fmt.Sprintf(managedSecretNamePattern, shaName(pg.Spec.RestoreFrom))
	s.Namespace = pg.Namespace

	data := make(map[string][]byte)

	data["authorized_keys"], _ = base64.StdEncoding.DecodeString(authorizedKeys)
	data["config"], _ = base64.StdEncoding.DecodeString(config)
	data["id_ed25519"], _ = base64.StdEncoding.DecodeString(idED25519)
	data["ssh_host_ed25519_key"], _ = base64.StdEncoding.DecodeString(sshHostED25519Key)
	data["sshd_config"], _ = base64.StdEncoding.DecodeString(sshdConfig)
	data["aws-s3-key"], _ = base64.StdEncoding.DecodeString(pg.Spec.BackrestS3Key)
	data["aws-s3-key-secret"], _ = base64.StdEncoding.DecodeString(pg.Spec.BackrestS3KeySecret)

	s.Data = data

	annotations := make(map[string]string)
	annotations["pg-port"] = "5432"
	annotations["repo-path"] = fmt.Sprintf(`/%s`, repoPath)
	annotations["s3-bucket"] = pg.Spec.BackrestS3Bucket
	annotations["s3-endpoint"] = pg.Spec.BackrestS3Endpoint
	annotations["s3-region"] = pg.Spec.BackrestS3Region
	annotations["s3-uri-style"] = pg.Spec.BackrestS3URIStyle
	annotations["s3-verify-tls"] = "false"
	annotations["sshd-port"] = "2022"
	s.Annotations = annotations

	return s
}

func PVCFromPgcluster(pg *v1alpha1.PostgreSQLCluster) *corev1.PersistentVolumeClaim {
	p := &corev1.PersistentVolumeClaim{}

	p.Name = fmt.Sprintf(managedPVCNamePattern, shaName(pg.Spec.RestoreFrom))
	p.Namespace = pg.Namespace
	p.Spec.AccessModes = append(p.Spec.AccessModes, corev1.ReadWriteOnce)
	p.Spec.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceStorage: resource.MustParse("1Gi"),
		},
	}

	storageClass := pg.Spec.StorageConfig
	p.Spec.StorageClassName = &storageClass
	volumeMode := corev1.PersistentVolumeFilesystem
	p.Spec.VolumeMode = &volumeMode

	return p
}

func CreateManagedResource(pg *v1alpha1.PostgreSQLCluster) error {
	reg := regexp.MustCompile(fmt.Sprintf(restoreFromPattern, pg.Name))
	result := reg.FindAllStringSubmatch(pg.Spec.RestoreFrom, -1)
	if result == nil || len(result[0]) < 2 {
		return nil
	}

	secret := SecretFromPgcluster(pg, result[0][1])
	pvc := PVCFromPgcluster(pg)

	k8s := k8sclient.GetKubernetesClient()

	_, err := k8s.CoreV1().Secrets(pg.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		klog.Errorf("create secret %s error: %s", secret.Name, err)
		return err
	}
	_, err = k8s.CoreV1().PersistentVolumeClaims(pg.Namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		klog.Error("create pvc %s error: %s", pvc.Name, err)
		return err
	}

	return nil
}

func DeleteManagedResource(pg *v1alpha1.PostgreSQLCluster) {
	reg := regexp.MustCompile(fmt.Sprintf(restoreFromPattern, pg.Name))
	result := reg.FindAllStringSubmatch(pg.Spec.RestoreFrom, -1)
	if result == nil || len(result[0]) < 2 {
		return
	}

	name := shaName(pg.Spec.RestoreFrom)

	k8s := k8sclient.GetKubernetesClient()
	secretName := fmt.Sprintf(managedSecretNamePattern, name)
	pvcName := fmt.Sprintf(managedPVCNamePattern, name)
	deployName := fmt.Sprintf(managedDeploymentNamePattern, name)

	err := k8s.AppsV1().Deployments(pg.Namespace).Delete(context.TODO(), deployName, metav1.DeleteOptions{})
	if err != nil {
		klog.Warningf("delete managed deployment %s error: %s", deployName, err)
	}
	err = k8s.CoreV1().PersistentVolumeClaims(pg.Namespace).Delete(context.TODO(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		klog.Warningf("delete managed pvc %s error: %s", pvcName, err)
	}
	err = k8s.CoreV1().Secrets(pg.Namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		klog.Warningf("delete managed secret %s error: %s", pvcName, err)
	}

	return
}
