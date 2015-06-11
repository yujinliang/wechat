// WeiXin.go
package mp

import (
	
	"log"
	"net/http"
	"net/url"
	"encoding/xml"
	"encoding/base64"
	"github.com/yujinliang/wechat/mp/request"
	"github.com/yujinliang/wechat/mp/accesstoken"
	"github.com/yujinliang/wechat/mp/utils"
	
)

const (

	// media types
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeThumb = "thumb"

	// environment constants
	weixinRootHost  = "https://api.weixin.qq.com"
	weixinHost      = "https://api.weixin.qq.com/cgi-bin"
	weixinQRScene   = "https://api.weixin.qq.com/cgi-bin/qrcode"
	weixinShowQRScene = "https://mp.weixin.qq.com/cgi-bin/showqrcode"
	weixinFileURL   = "api.weixin.qq.com/cgi-bin/media"
	weixinKFFileURL = "http://api.weixin.qq.com/customservice/kfaccount/uploadheadimg"
	retryNum       = 3

)

//http request handler.
type RequestHandlerFunc func(wx *WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string)

type WeiXin struct {
	
	token 		string
	appId 		string
	appSecret 	string
	useCurrentAESKey		bool
	currentEncodingAESKey 	string
	lastEncodingAESKey 		string
	currentAesKey			[]byte
	lastAesKey				[]byte
	
	//for http request route map
	httpRoutes	map[string] RequestHandlerFunc
	//for send custom message.
	accessTokenSupplierChan chan accesstoken.AccessToken //no buffer channel , synchronous mode.
	
}

// 安全模式, 微信服务器推送过来的 http body
type EncryptedRequestHttpBody struct {
	
	XMLName struct{} `xml:"xml" json:"-"`

	ToUserName   string `xml:"ToUserName" json:"ToUserName"`
	EncryptedMsg string `xml:"Encrypt"    json:"Encrypt"`
	
}

func encodingAESKey2AESKey(encodingKey string) []byte {
	
	data, _ := base64.StdEncoding.DecodeString(encodingKey + "=")
	return data
	
}
func New(token string, appid string, secret string, currentEncodingAESKey string, lastEncodingAESKey string) *WeiXin {
	
	wx := &WeiXin{}
	wx.token = token
	wx.appId = appid
	wx.appSecret = secret
	wx.currentEncodingAESKey = currentEncodingAESKey
	wx.lastEncodingAESKey    = lastEncodingAESKey
	wx.useCurrentAESKey = true
	
	if len(wx.currentEncodingAESKey) > 0 {
		
		wx.currentAesKey = encodingAESKey2AESKey(wx.currentEncodingAESKey)
	
	} else {
		
		//必须提供当前设定的encodingAESKey
		return nil
		
	}
	//如果能提供上一次设定EncodingAESKey, 则当用当前设定的encodingAESKey解密失败时， 则转而用上一次设定过的encodingAESKey,
	//何故如此： 在我们重新设定encodingAESKey之后， 微信集群不能实时同步更新， 很可能一段时间内部分服务器还在使用旧的encodingAESKey,
	//故此我方才如此设计，以使系统强大，健壮。
	if len(wx.lastEncodingAESKey) > 0 {
		
		wx.lastAesKey = encodingAESKey2AESKey(wx.lastEncodingAESKey)
		
	}
	wx.httpRoutes = make(map[string] RequestHandlerFunc) //key: msg type , value : msg handler.
	
	//get access token, run in a thread.
	if len(appid) > 0 && len(secret) > 0 {
		
		wx.accessTokenSupplierChan = make(chan accesstoken.AccessToken)
		go accesstoken.CreateAccessToken(weixinHost, wx.accessTokenSupplierChan, appid, secret)
		
	}
	
	return wx
	
}
func (wx *WeiXin) GetAppId() string {
	
	return wx.appId
	
}
func (wx *WeiXin) GetAppSecret() string {
	
	return wx.appSecret
	
}
func (wx *WeiXin) GetToken() string {
	
	return wx.token
	
}
// register request handler
func (wx *WeiXin) HandleFunc(pattern string, handler RequestHandlerFunc) {
	
	if len(pattern) > 0 && handler != nil {
		
		wx.httpRoutes[pattern] = handler
	}
	
}
func GenHttpRouteKey(msgType string, msgEvent string) string {
	
	return  msgType + msgEvent
	
}

//route http request to specified handler
func (wx *WeiXin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
		
	case "POST": {
		
	//1. 解密
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		
		return
		
	}
	switch encryptType := queryValues.Get("encrypt_type"); encryptType {
		
		case "aes": {
			
			signature := queryValues.Get("signature") //just read it , no check.
			msgSignature1 := queryValues.Get("msg_signature")
			if len(msgSignature1) <= 0 {
				
				return
				
			}
			timestamp := queryValues.Get("timestamp")
			if len(timestamp) <= 0 {
				
				return
				
			}
			nonce := queryValues.Get("nonce")
			if len(nonce) <= 0 {
				
				return
				
			}
			var encryptedMsg EncryptedRequestHttpBody
			if err := xml.NewDecoder(r.Body).Decode(&encryptedMsg); err != nil {
				
				return
				
			}
			
			//验证tousername 是否为当前公众号
			//验证签名
			msgSignature2 := makeMsgSignature(wx.token, timestamp, nonce, encryptedMsg.EncryptedMsg)
			if msgSignature1 != msgSignature2 {
				
				return
				
			}
			//开发解密
			cipherData, err := base64.StdEncoding.DecodeString(encryptedMsg.EncryptedMsg)
			if err != nil {
				
				return
				
			}
				
			plainData, err := aesDecrypt(cipherData, wx.currentAesKey)
			if err != nil {
				
				//try last aesKey to decrypt.
				if len(wx.lastAesKey) > 0 {
						
					plainData, err = aesDecrypt(cipherData, wx.lastAesKey)
					if err != nil {
							
						//两个密钥都失败， 则放弃, 并切换回当前密钥,以供处理下次请求时用。
						wx.useCurrentAESKey = true
						return
							
					}
						
					wx.useCurrentAESKey = false
						
				} else {
						//没有提供上一次的密钥， 则放弃。
						wx.useCurrentAESKey = true
						return
						
				}
					
			} else {
				
				wx.useCurrentAESKey = true
				
			}
			
			mixedMsg := &request.WeiXinRequest{}
			if rc := mixedMsg.TryUnpackWeiXinRequestForEncrypted(wx.appId, wx.token,w,plainData,timestamp,nonce,signature); rc != true {
				
				return
				
			}
			
			if len(mixedMsg.MsgType) > 0 {
			
				log.Println("Weixin receive encrypted message:", mixedMsg.MsgType)
			
				if handler, ok := wx.httpRoutes[GenHttpRouteKey(mixedMsg.MsgType, mixedMsg.Event)]; ok == true {
				
					handler(wx, w, mixedMsg, timestamp, nonce)
				
				}
				
			}
			
		}//end of aes.
		case "", "raw": {
			//2. 以下为处理明文
			msg := &request.WeiXinRequest{}
			if rc := msg.TryUnpackWeiXinRequest(wx.token, w, r); rc == true {
		
				if len(msg.MsgType) > 0 {
			
					log.Println("Weixin receive message:", msg.MsgType)
			
					if handler, ok := wx.httpRoutes[GenHttpRouteKey(msg.MsgType, msg.Event)]; ok == true {
				
						timestamp := queryValues.Get("timestamp")
						if len(timestamp) <= 0 {
				
							return
				
						}
						nonce := queryValues.Get("nonce")
						if len(nonce) <= 0 {
				
							return
				
						}
						//call callback function.
						handler(wx, w, msg, timestamp, nonce)
				
					}
					
				}
				
			}
	
		}//end of raw
		default: {
			//未知加密类型，直接忽略，不再处理。
			return
			
		}
	
	}//end of switch
	
	}//end of post.
	
	case "GET": {
		
		utils.InvalidateForWeiXin(wx.token, r, w)
		
	}//end of get.
	
	} //end of method
}

