package utils

import (
	"bytes"
	"github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
)

func HttpPost(url string, body []byte) ([]byte, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		seelog.Error("http new request error:%v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	var client = http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		seelog.Error("http client do error:%v", err)
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
