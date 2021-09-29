package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"net/http"
)

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
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
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
	fmt.Println("///res:", string(res), "======end")
	return res, err
}
