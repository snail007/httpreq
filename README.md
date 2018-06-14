# httpreq

http request for ios &amp; android

## 方法参数说明

所有的方法返回结果都是一个固定的json对象,结构如下:

```json
{
    "StatusCode"   int //返回的http状态码,如果请求失败,那么是0
    "ErrorMessage" string //请求发生错误的时候,这里是错误信息
    "Body"         string //请求返回的内容,如果参数base64body是1,
                    //那么这里的内容是经过base64编码的文本,
                    //使用的时候需要自己base64解码.
    "Headers"      {key:[]string}
}
```

所有方法都有如下几个参数,这里统一说明:

URL: 字符串,请求的URL地址,必须是http://或者https://开头.

jsonHeader: 字符串,需要设置的请求头,json格式的数据,一个json对象,键值都必须是字符串类型.比如:{"User-Agent":"httpreq/1.1"}

timeout: 字符串,内容是数字,请求的超时时间,单位是毫秒,比如:500

base64body: 字符串,返回的内容是否进行base64编码.1:编码,0:不编码,默认:0

tlsJsonConfig:字符串,https请求的时候,需要的一些配置,json格式的数据,一个json对象它的结构如下:
    ```json
    {
        Key             string //pem格式的key文件的文本内容
        Cert            string  //pem格式的crt根证书文件的文本内容
        Cas              []string //字符串数组,值是pem格式的ca根证书文件的文本内容,用于对服务器的证书的检查
        UseSystemCert   string   //是否加载系统的信任证书,用于对服务器的证书的检查,1:加载,0:不加载,默认:1
        CheckServerName string    //是否检查ServerName,1:检查,0:不检查,默认:0
        CheckCert       string //是否检查服务器证书,1:检查,0:不检查,默认:0
    }
    ```
    说明:

    所有的配置都不是必须传递的,根据情况设置.

    1.tls双向认证的时候,如果要检查证书,必须设置Key,Cert,Ca,CheckCert

    2.单向认证的时候,如果要检查证书,必须设置Ca,CheckCert

### PostBody

作用:

Post原始的字符串内容到服务器,需要自己根据内容类型设置合适的头部.

声明如下:

`PostBody(URL, bodyData, jsonHeader, timeout, base64body, tlsJsonConfig string) (resultJSON string)`

bodyData: 字符串,需要发送的数据,数据不会经过任何处理,直接发送给服务器.

### PostJSON

作用:

Post发送JSON字符串内容到服务器,请求时会设置头部:Content-Type: application/json

声明如下:

`PostBody(URL, jsonData, jsonHeader, timeout, base64body, tlsJsonConfig string) (resultJSON string)`

jsonData: 字符串,需要发送的数据,数据不会经过任何处理,直接发送给服务器.

### PostXML

作用:

Post发送XML字符串内容到服务器,请求时会设置头部:Content-Type: text/xml

声明如下:

`PostXML(URL, xmlData, jsonHeader, timeout, base64body, tlsJsonConfig string) (resultJSON string)`
xmlData: 字符串,需要发送的数据,数据不会经过任何处理,直接发送给服务器.

### PostForm

作用:

Post发送经过编码的表单内容到服务器,json格式的formData会被处理为表单数据然后编码,请求时会设置头部:Content-Type: application/x-www-form-urlencoded

声明如下:

`PostXML(URL, formData, jsonHeader, timeout, base64body, tlsJsonConfig string) (resultJSON string)`

formData: 字符串,需要发送的数据,数据会被处理为表单数据然后编码发送给服务器.

### Get

作用:

发送Get请求到服务器.

声明如下:

`Get(URL, jsonParams, jsonHeader, timeout, base64body, tlsJsonConfig string) (resultJSON string)`

jsonParams: 字符串,需要附加到URL后面的参数,json格式的数据,一个json对象,键值都必须是字符串类型.
    比如:{"uid":"123"}
