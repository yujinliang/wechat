//推广支持
package mp

import (
	
	"log"
	"encoding/json"
	"github.com/yujinliang/wechat/mp/response"
	
)

//二维码
type QrCode struct {
	
	ExpireSeconds int64 `json:"expire_seconds,omitempty"`
	ActionName	  string `json:"action_name"`
	ActionInfo	  struct {
		
		Scene struct {
			
			SceneId int64 `json:"scene_id"`
			
		} `json:"scene"`
		
		
	} `json:"action_info"`
	
}

//二维码
//创建永久二维码
func (wx *WeiXin) CreateQRCode(sceneId int64) (*response.Response, error) {
	
	var code QrCode
	code.ActionName = "QR_SCENE"
	code.ActionInfo.Scene.SceneId = sceneId
	
	data, err := json.Marshal(&code)
	if err != nil {
		
		return nil, err
		
	} 
		
	//just for debug.
	log.Println("CreateQRCode: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/qrcode/create?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return rc, err
	
}
//创建临时二维码
func (wx *WeiXin) CreateQRLimitCode(sceneId int64, expireSeconds int64) (*response.Response, error) {
	
	var code QrCode
	code.ActionName = "QR_LIMIT_SCENE"
	code.ActionInfo.Scene.SceneId = sceneId
	code.ExpireSeconds = expireSeconds
	
	data, err := json.Marshal(&code)
	if err != nil {
		
		return nil, err
		
	} 
		
	//just for debug.
	log.Println("CreateQRLimitCode: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/qrcode/create?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return rc, err
	
}
//二维码图片的url
func (wx *WeiXin) GetQRCodeURL(ticket string) string {
	
	return weixinShowQRScene + "?ticket=" + ticket
	
}
//长链接转为短链接
func (wx *WeiXin) LongURL2ShortURL(longURL string) (string, error) {
	
	var msg struct {
		
		Action 	string `json:"action"`
		LongURL string `json:"long_url"`
	}
	msg.Action  = "long2short"
	msg.LongURL = longURL
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		return "", err
		
	}
	
	//just for debug.
	log.Println("LongURL2ShortURL: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/shorturl?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	
	return rc.ShortURL, nil
	
}
