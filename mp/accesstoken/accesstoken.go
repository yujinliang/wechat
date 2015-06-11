package accesstoken

import (
	
	"time"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"errors"

)

type AccessToken struct {
	
	Token   string
	Expires time.Time //seconds.
	
}
//Attention: just allow to call only one time in weixin.New function.
func CreateAccessToken( weixinhost string, c chan AccessToken, appid string, secret string) {
	
	token := AccessToken{"", time.Now()}
	
	for {
		
		log.Println("old AccessTokenTimeOut: ", token.Expires.Unix(), time.Now().Unix())
		if token.Expires.Unix() <= time.Now().Unix() {
			
			var expires int64
			token.Token, expires = authAccessToken(weixinhost, appid, secret)
			token.Expires = time.Now().Add(time.Second * time.Duration(expires))
			//just for debug.
			log.Println("new AccessTokenTimeOut: ",token.Token, expires, token.Expires.Unix(), time.Now().Unix())
		}
		
		c <- token
		
	}
}

func authAccessToken(weixinhost string, appid string, secret string) (string, int64) {
	
	resp , err := http.Get(weixinhost + "/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret)
	if err != nil {
		
		log.Println("Get access token failed: ", err)
		
	} else {
		
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			
			log.Println("Read Access token failed: ", err)
			
		} else {
			
			var rc struct {
				
				// error fields
				ErrCode int64  `json:"errcode"`
				ErrMsg  string `json:"errmsg"`
				// token fields
				AccessToken string `json:"access_token"`
				ExpiresIn   int64  `json:"expires_in"`
				
			}
			
			if err := json.Unmarshal(body, &rc); err != nil {
				
				log.Println("Parse access token failed: ", err)
				
			} else {
				
				if rc.ErrCode != 0 {
					
					log.Println("wexin return error: ", rc.ErrMsg)
					
				} else {
					
					return rc.AccessToken, rc.ExpiresIn
					
				}
			}
		}
	}
	
	return "", 0
	
}
//获取微信服务器IP地址
func GetWeiXinServerIPList(weixinhost string, c chan AccessToken) ([]string, error) {
	
	token := <- c
	resp , err := http.Get(weixinhost + "/getcallbackip?access_token=" + token.Token)
	if err != nil {
		
		log.Println("Get WeiXin Server IP list failed: ", err)
		return nil, err
		
	} else {
		
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			
			log.Println("Read WeiXin Server IP list failed: ", err)
			return nil, err
			
		} else {
			
			var rc struct {
				
				// error fields
				ErrCode int64  `json:"errcode"`
				ErrMsg  string `json:"errmsg"`
				// token fields
				Ip_list  []string `json:"ip_list,omitempty"`

			}
			
			if err := json.Unmarshal(body, &rc); err != nil {
				
				log.Println("Parse WeiXin Server IP list failed: ", err)
				return nil, err
				
			} else {
				
				if rc.ErrCode != 0 {
					
					log.Println("wexin return error:[%d] ",rc.ErrCode, rc.ErrMsg)
					return nil, err
					
				} else {
					
					return rc.Ip_list, nil
					
				}
			}
		}
	}
	
	return nil, errors.New("GetWeiXinServerIPList Unknow Error!")
	
}