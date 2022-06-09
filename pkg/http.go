package pkg

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/kubesphere/k8sclient"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"net/http"
	"os"
)

var certContent, certKeyContent []byte
var tlsConfig *tls.Config
var cert tls.Certificate

const PgoTLSSecretName = "pgo.tls"
const disableTLSEnvName = "DISABLE_TLS"

func init() {
	if os.Getenv(disableTLSEnvName) != "" {
		return
	}

	var err error

	k8s := k8sclient.GetKubernetesClient()

	s, err := k8s.CoreV1().Secrets(PgoNamespace).Get(context.TODO(), PgoTLSSecretName, metav1.GetOptions{})
	if err != nil {
		klog.Fatal("unable get tls secret")
	}

	certContent = s.Data["tls.crt"]
	certKeyContent = s.Data["tls.key"]

	cert, err = tls.X509KeyPair(certContent, certKeyContent)
	if err != nil {
		klog.Fatal("unable to load cert")
	}

	tlsConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
}

func Call(method, route string, body interface{}) (resByte []byte, err error) {
	var req []byte
	if body != nil {
		req, err = json.Marshal(body)
		if err != nil {
			return
		}
	}
	resByte, err = DoRequest(method, route, req)
	if err != nil {
		klog.Errorf("call : post json url %s,error %s", route, err)
		return
	}
	return resByte, err
}

func DoRequest(method, url string, body []byte) ([]byte, error) {
	if os.Getenv(disableTLSEnvName) != "" {
		return InsecureRequest(method, url, body)
	}

	return SecureRequest(method, url, body)
}

func InsecureRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, "http://"+url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.SetBasicAuth(UserName, PassWord)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	return res, err
}

func SecureRequest(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, "https://"+url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.SetBasicAuth(UserName, PassWord)

	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	return res, err
}
