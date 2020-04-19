package common

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"
)

//N9eRes N9e http 接口的返回
type N9eRes struct {
	Err string `json:"err"`
}

//HTTPGet 发起一个 http get 请求
func HTTPGet(url string) (body []byte, err error) {
	body = []byte{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ = ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("HTTP Connect Failed, Code is %d, body is %s", resp.StatusCode, string(body))
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// HTTPPost 发起一个 http post 请求
func HTTPPost(url string, data io.Reader, headers map[string]string) (body []byte, err error) {
	body = []byte{}
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ = ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("HTTP Connect Failed, Code is %d, body is %s", resp.StatusCode, string(body))
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	return
}

//N9ePush
func N9ePush(url string, data io.Reader) error {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"
	res, err := HTTPPost(url, data, headers)
	if err != nil {
		return err
	}

	var n9eRes N9eRes
	if err = json.Unmarshal(res, &n9eRes); err != nil {
		return err
	}
	if n9eRes.Err != "" {
		return errors.New(n9eRes.Err)
	}
	return nil
}
