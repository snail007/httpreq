package httpreq

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	version = "1.1"
)

type tlsConfig struct {
	Key             string
	Cert            string
	Cas             []string
	UseSystemCert   string
	CheckServerName string
	CheckCert       string
}
type result struct {
	StatusCode   int
	ErrorMessage string
	Body         string
	Headers      map[string][]string
}

func Version() string {
	return version
}

func Get(URL, jsonParams, jsonHeader, timeout, base64body, tlsJsonConfig string) (result string) {
	paramsString := ""
	jsonDataMap := map[string]string{}
	if e := json.Unmarshal([]byte(jsonParams), &jsonDataMap); e == nil {
		postParams := []string{}
		for k, v := range jsonDataMap {
			postParams = append(postParams, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
		paramsString = strings.Join(postParams, "&")
	}
	if paramsString != "" {
		if strings.Contains(URL, "?") {
			URL += "&"
		} else {
			URL += "?"
		}
		URL += paramsString
	}
	return request("GET", URL, "", jsonHeader, timeout, base64body, tlsJsonConfig)
}
func PostBody(URL, bodyData, jsonHeader, timeout, base64body, tlsJsonConfig string) (result string) {
	return request("POST", URL, bodyData, jsonHeader, timeout, base64body, tlsJsonConfig)
}
func PostJSON(URL, jsonData, jsonHeader, timeout, base64body, tlsJsonConfig string) (result string) {
	return postXXX(URL, "application/json", jsonData, jsonHeader, timeout, "0", base64body, tlsJsonConfig)
}
func PostXML(URL, xmlData, jsonHeader, timeout, base64body, tlsJsonConfig string) (result string) {
	return postXXX(URL, "text/xml", xmlData, jsonHeader, timeout, "0", base64body, tlsJsonConfig)
}
func PostForm(URL, formData, jsonHeader, timeout, base64body, tlsJsonConfig string) (result string) {
	return postXXX(URL, "application/x-www-form-urlencoded", formData, jsonHeader, timeout, "1", base64body, tlsJsonConfig)
}
func postXXX(URL, contentType, jsonData, jsonHeader, timeout, isForm, base64body, tlsJsonConfig string) (result string) {
	if jsonHeader == "" {
		jsonHeader = `{"Content-Type":"` + contentType + `"}`
	} else {
		var h map[string]string
		json.Unmarshal([]byte(jsonHeader), &h)
		h["Content-Type"] = contentType
		b, _ := json.Marshal(h)
		jsonHeader = string(b)
	}
	postParamsString := ""
	if jsonData != "" {
		if isForm == "1" {
			jsonDataMap := map[string]string{}
			if e := json.Unmarshal([]byte(jsonData), &jsonDataMap); e == nil {
				postParams := []string{}
				for k, v := range jsonDataMap {
					postParams = append(postParams, url.QueryEscape(k)+"="+url.QueryEscape(v))
				}
				postParamsString = strings.Join(postParams, "&")
			}
		} else {
			postParamsString = jsonData
		}
	}
	return request("POST", URL, postParamsString, jsonHeader, timeout, base64body, tlsJsonConfig)
}
func request(method, URL, paramsString, jsonHeader, timeout, base64body, tlsJsonConfig string) (ret string) {
	ret0 := result{
		Headers: map[string][]string{},
	}
	defer func() {
		if e := recover(); e != nil {
			ret0.ErrorMessage = fmt.Sprintf("%s", e)
		}
		b, _ := json.Marshal(ret0)
		ret = string(b)
	}()
	timeout0, _ := strconv.Atoi(timeout)

	//configuration for request
	var tr *http.Transport
	var client *http.Client
	if strings.Contains(URL, "https://") {
		tlsConfig0 := tlsConfig{}
		if tlsJsonConfig != "" {
			err := json.Unmarshal([]byte(tlsJsonConfig), &tlsConfig0)
			if err != nil {
				ret0.ErrorMessage = err.Error()
				return
			}
		}
		u, _ := url.Parse(URL)
		serverName := u.Hostname()
		cas := [][]byte{}
		for _, ca := range tlsConfig0.Cas {
			cas = append(cas, []byte(ca))
		}
		conf, err := getRequestTlsConfig([]byte(tlsConfig0.Cert), []byte(tlsConfig0.Key), cas, serverName, tlsConfig0.UseSystemCert != "0", tlsConfig0.CheckServerName == "1", tlsConfig0.CheckCert == "1")
		if err != nil {
			ret0.ErrorMessage = err.Error()
			return
		}
		tr = &http.Transport{TLSClientConfig: conf}
		client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout0), Transport: tr}
	} else {
		tr = &http.Transport{}
		client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout0), Transport: tr}
	}
	defer tr.CloseIdleConnections()
	var bodyReader io.Reader
	if strings.ToLower(method) == "post" {
		bodyReader = bytes.NewBuffer([]byte(paramsString))
	}

	req, err := http.NewRequest(strings.ToUpper(method), URL, bodyReader)
	if err != nil {
		return
	}
	//set headers
	req.Header.Set("User-Agent", "httpreq/"+version)
	if jsonHeader != "" {
		jsonHeaderMap := map[string]string{}
		if e := json.Unmarshal([]byte(jsonHeader), &jsonHeaderMap); e == nil {
			for k, v := range jsonHeaderMap {
				req.Header.Set(k, v)
			}
		}
	}

	//request
	resp, err := client.Do(req)

	if err != nil {
		ret0.ErrorMessage = err.Error()
		return
	}
	defer resp.Body.Close()

	//status code
	ret0.StatusCode = resp.StatusCode
	//headers
	for k, v := range resp.Header {
		ret0.Headers[k] = v
	}
	//body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ret0.ErrorMessage = err.Error()
		return
	}
	if base64body == "1" {
		ret0.Body = base64.StdEncoding.EncodeToString(b)
	} else {
		ret0.Body = string(b)
	}
	return
}

func getRequestTlsConfig(certBytes, keyBytes []byte, caCertBytes [][]byte, serverName string, useSystemCert, checkServerName, checkCert bool) (conf *tls.Config, err error) {
	conf = &tls.Config{ServerName: serverName}
	var serverCertPool *x509.CertPool

	if useSystemCert {
		serverCertPool, err = x509.SystemCertPool()
		if err != nil {
			return
		}
	} else {
		serverCertPool = x509.NewCertPool()
	}

	if caCertBytes != nil && len(caCertBytes) > 0 && len(caCertBytes[0]) > 0 {
		for _, ca := range caCertBytes {
			ok := serverCertPool.AppendCertsFromPEM(ca)
			if !ok {
				err = errors.New("failed to parse root certificate")
				return
			}
		}
	}
	conf.RootCAs = serverCertPool

	if certBytes != nil && keyBytes != nil && len(certBytes) > 0 && len(keyBytes) > 0 {
		var cert tls.Certificate
		cert, err = tls.X509KeyPair(certBytes, keyBytes)
		if err != nil {
			return
		}
		conf.Certificates = []tls.Certificate{cert}
	}

	conf.InsecureSkipVerify = true
	conf.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if checkCert {
			opts := x509.VerifyOptions{
				Roots: serverCertPool,
			}
			for _, rawCert := range rawCerts {
				cert, _ := x509.ParseCertificate(rawCert)
				_, err := cert.Verify(opts)
				if err != nil {
					return err
				}
			}
		}
		if checkServerName {
			for _, rawCert := range rawCerts {
				cert, _ := x509.ParseCertificate(rawCert)
				err := cert.VerifyHostname(serverName)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
	return
}
