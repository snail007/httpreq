package httpreq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPost(t *testing.T) {
	caStringB, _ := ioutil.ReadFile("ca.crt")
	certStringB, _ := ioutil.ReadFile("chrome.crt")
	keyStringB, _ := ioutil.ReadFile("chrome.key")
	conf := tlsConfig{
		CheckCert:       "1",
		CheckServerName: "0",
		Ca:              string(caStringB),
		Cert:            string(certStringB),
		Key:             string(keyStringB),
	}
	c, _ := json.Marshal(conf)
	headers := map[string]string{
		"User-Agent": "httpreq/1.1",
		//"Content-Type": "application/x-www-form-urlencoded",
	}
	h, _ := json.Marshal(headers)

	body := map[string]string{
		"uid": "123",
	}
	b, _ := json.Marshal(body)

	res := PostForm("https://127.0.0.1:8080/", string(b), string(h), "1000", "0", string(c))
	r := result{}
	json.Unmarshal([]byte(res), &r)
	if r.StatusCode == 0 {
		t.Fatal(r.ErrorMessage)
	} else {
		fmt.Printf("StatusCode:%d,Body:%s\n", r.StatusCode, r.Body)
	}
}
