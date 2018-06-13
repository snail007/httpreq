package httpreq
import(
	"encoding/json"
	"net/url"
	"strings"
	"net/http"
	"crypto/tls"
	"time"
	"bytes"
	"io/ioutil"
	"crypto/x509"
	"errors"
	"strconv"
)
type tlsConfig struct{
	Key string
	Cert string
	Ca string
	CheckServerName string
	CheckCert string
}
type result struct{
	StatusCode int
	ErrorMessage string
	Body string
	Headers map[string][]string
}
func Post(URL  ,jsonData  ,jsonHeader ,timeout ,tlsJsonConfig string )(ret string){
	timeout0,_:=strconv.Atoi(timeout)
	ret0:=result{
		Headers:map[string][]string{},
	}
	defer func(){
		b,_:=json.Marshal(ret0)
		ret=string(b)
	}()

	//convert post data to string
	postParamsString := ""
	if jsonData!=""{
		jsonDataMap:=map[string]string{}
		if e:=json.Unmarshal([]byte(jsonData),&jsonDataMap);e==nil {
			postParams := []string{}
			for k, v := range jsonDataMap {
				postParams = append(postParams, url.QueryEscape(k)+"="+url.QueryEscape(v))
			}
			postParamsString = strings.Join(postParams, "&")
		}
	}
	//configuration for request
	var tr *http.Transport
	var client *http.Client
	if strings.Contains(URL, "https://") {
		tlsConfig0:=tlsConfig{}
		if tlsJsonConfig!=""{
			err:=json.Unmarshal([]byte(tlsJsonConfig),&tlsConfig0)
			if err!=nil{
				ret0.ErrorMessage=err.Error()
				return
			}
		}
		u,_:=url.Parse(URL)
		serverName:=u.Hostname()
		conf,err:=getRequestTlsConfig([]byte(tlsConfig0.Cert),[]byte(tlsConfig0.Key),[]byte(tlsConfig0.Ca),serverName,tlsConfig0.CheckServerName=="1",tlsConfig0.CheckCert=="1")
		if err!=nil{
			ret0.ErrorMessage=err.Error()
			return
		}
		tr = &http.Transport{TLSClientConfig: conf}
		client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout0), Transport: tr}
	} else {
		tr = &http.Transport{}
		client = &http.Client{Timeout: time.Millisecond * time.Duration(timeout0), Transport: tr}
	}
	defer tr.CloseIdleConnections()

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer([]byte(postParamsString)))
	if err != nil {
		return
	}
	//set headers
	if jsonHeader!=""{
		jsonHeaderMap:=map[string]string{}
		if e:=json.Unmarshal([]byte(jsonHeader),&jsonHeaderMap);e==nil {
			
			for k, v := range jsonHeaderMap {
				req.Header.Set(k, v)
			}
			
		}
	}

	//request
	resp, err := client.Do(req)

	if err != nil {
		ret0.ErrorMessage=err.Error()
		return
	}
	defer resp.Body.Close()

	//status code
	ret0.StatusCode=resp.StatusCode
	//headers
	for k,v:=range resp.Header{
		ret0.Headers[k]=v
	}
	//body
	b, err := ioutil.ReadAll(resp.Body)
	if err!=nil{
		ret0.ErrorMessage=err.Error()
		return
	}
	ret0.Body=string(b)
	return 
}
func Get(reqJSON string)(result string){

	return ""
	}

func getRequestTlsConfig(certBytes, keyBytes, caCertBytes []byte,serverName string,checkServerName,checkCert bool) (conf *tls.Config, err error) {
	conf = &tls.Config{ServerName:serverName,}
	
	serverCertPool := x509.NewCertPool()
	if caCertBytes!=nil&&len(caCertBytes)>0{
		ok := serverCertPool.AppendCertsFromPEM(caCertBytes)
		if !ok {
			err = errors.New("failed to parse root certificate")
			return
		}
	}
	conf.RootCAs=serverCertPool 

	if certBytes!=nil&&keyBytes!=nil&&len(certBytes)>0&&len(keyBytes)>0{
		var cert tls.Certificate
		cert, err = tls.X509KeyPair(certBytes, keyBytes)
		if err != nil {
			return
		}
		conf.Certificates=[]tls.Certificate{cert}
	}

	conf.InsecureSkipVerify = true
	conf.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if checkCert{
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
		if checkServerName{ 
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
	