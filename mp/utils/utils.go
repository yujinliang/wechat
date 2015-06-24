// weixin.go
package utils

import (
	
	"fmt"
	"io"
	"net"
	"sort"
	"hash"
	"time"
	"errors"
	"bytes"
	"strings"
	"net/http"
	"crypto/sha1"
	"crypto/tls"
	"crypto/md5"
	"encoding/xml"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	
)

func InvalidateForWeiXin(token string, r *http.Request, w http.ResponseWriter) {
	
	r.ParseForm()
	var signature string = r.FormValue("signature")
	if len(signature) <= 0 {
		
		http.Error(w, "No Signature!", http.StatusUnauthorized)
		return
		
	}
	
	var timestamp string = r.FormValue("timestamp")
	if len(timestamp) <= 0 {
		
		http.Error(w, "No TimeStamp!", http.StatusUnauthorized)
		return
		
	}
	
	var nonce	  string = r.FormValue("nonce")
	if len(nonce) <= 0 {
		
		http.Error(w, "No Nonce!", http.StatusUnauthorized)
		return
		
	}
	
	var echostr string = r.FormValue("echostr")
	if len(echostr) <= 0 {
		
		http.Error(w, "No Echostr!", http.StatusUnauthorized)
		return
		
	}
	//--
	signatureGen := MakeSignature(token,timestamp,nonce)
	
	if signature == signatureGen {
		
		echostr := strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(w, echostr)
		return 
	}
	
	http.Error(w, "非法请法", http.StatusUnauthorized)
	return
	
}

func MakeSignature(token, timestamp, nonce string) string {
	
	tmpArray := []string{token, timestamp, nonce}
	sort.Strings(tmpArray)
	h := sha1.New()
	io.WriteString(h, strings.Join(tmpArray, ""))
	return fmt.Sprintf("%x", h.Sum(nil))
	
}

//此接口只用于处理微信支付接口xml,其它复杂情况没有考虑.
//没有处理深度，如果有嵌套，只取最内层的key:value.
func XML2Map( stream io.Reader) (map[string]string, error) {
	
	if stream == nil {
		
		return nil, errors.New("Input Parameter is Empty!")
		
	}
	
	m := make(map[string]string)
	d := xml.NewDecoder(stream)
	
	var(
		
		t xml.Token
		err error
		key string
		value bytes.Buffer
		
	)
	
	for {
		
		t, err = d.Token()
		if err != nil {
			
			fmt.Printf("v = %#v\n", m)
			if err != io.EOF {
				
				return nil, err
			
			}
			
			return m, nil
			
		}
		
		switch v := t.(type) {
			
			case xml.StartElement: {
				
				key = v.Name.Local
				value.Reset()
				
			}
			case xml.CharData: {
				
				value.Write(v)
				
			}
			case xml.EndElement: {
				
				if len(key) > 0 {
					
					m[key] = value.String()
					
				} 
				
			}//end of xml.EndElement
			
		}//end of switch
		
	} //end of for.
	
}
//对于value值没有转义，只是加了CDATA保护，即然xml不再解释它，我就先不加转义了.
func Map2XMLBytes(m map[string]string) ([]byte, error) {
	
	if len(m) <= 0 {
		
		return nil, errors.New("Map is Empty!")
		
	}
	
	buf := bytes.NewBufferString("<xml>")
	
	for k, v := range m {
		
		if _, err := buf.WriteString("<" + k + ">"); err != nil {
			
			return nil, errors.New("2xml failed, key:  " + k)
			
		}
		
		if _, err := buf.WriteString("<![CDATA[" + v + "]]>"); err != nil {
			
			return nil, errors.New("2xml failed, value: " + v)
			
		}
		
		if _, err := buf.WriteString("</" + k + ">"); err != nil {
			
			return nil, errors.New("2xml failed, key end</:  " + k)
			
		}
		
	}
	
	if _, err := buf.WriteString("</xml>"); err != nil {
			
			return nil, errors.New("2xml failed, xml end")
			
	}
	
	return buf.Bytes(), nil
	
}
//创建应用证书之httpClient.
// NewTLSHttpClient 创建支持双向证书认证的 http.Client
//引自：https://github.com/chanxuehong/wechat/blob/master/mch/http_client.go
//参考:https://github.com/bigwhite/experiments/blob/master/gohttps/6-dual-verify-certs/client.go
//http://stackoverflow.com/questions/18187136/net-http-ignoring-system-proxy-settings
//http://blog.csdn.net/luciswe/article/details/45890713
//现在没有对服务端身体进行验证即InsecureSkipVerify=true,以后有必要再加。
func NewTLSHttpClient(certFile, keyFile string) (*http.Client, error) {
	
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		
		return nil, err
		
	}
	
	tlsConfig := &tls.Config{
		
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
		
	}

	httpClient := &http.Client{
		
		Transport: &http.Transport{
			
			Dial: (&net.Dialer{
				
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				
			}).Dial,
			
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
			
		},
		
		Timeout: 60 * time.Second,
		
	}
	
	return httpClient, nil
	
}
//引自：https://github.com/chanxuehong/wechat/blob/master/mch/sign.go
//写的非常清晰，我就不再造车轮了，感恩前辈。
// 微信支付签名.
//  parameters: 待签名的参数集合
//  apiKey:     API密钥
//  fn:         func() hash.Hash, 如果 fn == nil 则默认用 md5.New
func SignForWXPay(parameters map[string]string, apiKey string, fn func() hash.Hash) string {
	
	ks := make([]string, 0, len(parameters))
	for k := range parameters {
		
		if k == "sign" {
			continue
		}
		
		ks = append(ks, k)
		
	}
	
	sort.Strings(ks)

	if fn == nil {
		
		fn = md5.New
		
	}
	h := fn()
	signature := make([]byte, h.Size()*2)

	for _, k := range ks {
		
		v := parameters[k]
		if v == "" {
			continue
		}
		h.Write([]byte(k))
		h.Write([]byte{'='})
		h.Write([]byte(v))
		h.Write([]byte{'&'})
	}
	h.Write([]byte("key="))
	h.Write([]byte(apiKey))

	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
	
}

//用于生成随机字符串
//代码摘自: http://blog.csdn.net/luciswe/article/details/45900373	
var NONCE_LETTERS_TOTAL = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
func GetNonceStr(n int) string {
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i:= range b {
		
		b[i] = NONCE_LETTERS_TOTAL[r.Intn(len(NONCE_LETTERS_TOTAL))]
		
	}
	return string(b)
	
}	
//注意当把AuthCodeURL生成的url用作菜单的url时，微信服务器端会把”非法字法的错误信息，如：& 被json.marshal转化为\u0026; 可微信不认这个，报\u0026为非法字符！”
//现在写一个方法将微信认为非法的转义后字符再转化回原字符
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	
	b, err := json.Marshal(v)
	if err != nil {
		
		return nil, err
		
	}
	
	if safeEncoding {
		
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	
	return b, nil
	
}


