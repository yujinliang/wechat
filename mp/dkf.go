package mp

import (
	
	"encoding/json"
	"log"
	"io"
	"bytes"
	"fmt"
	"time"
	"errors"
	"github.com/yujinliang/wechat/mp/response"
	"github.com/yujinliang/wechat/mp/request"
	
)

const (
	
	transferCustomerService = "<xml>" + replyHeader + "<MsgType><![CDATA[transfer_customer_service]]></MsgType></xml>"
	
)
//客服信息
type KF_Info struct {
	
	KF_account string `json:"kf_account"`
	KF_nick    string `json:"kf_nick"`
	KF_id	   string `json:"kf_id"`
	KF_headimgurl string `json:"kf_headimgurl"`
	
}
type KF_List struct {
	
	response.Response
	KF_list []KF_Info `json:"kf_list"`
	
}
type KF_Online_Info struct {
	
	KF_account   string `json:"kf_account"`
	KF_id	     string `json:"kf_id"`
	status	     int8	  `json:"status"`
	AutoAccept   int	  `json:"auto_accept"`
	AcceptedCase int  `json:"accepted_case"`
	
}
type Online_KF_List struct {
	
	response.Response
	online_kf_list []KF_Online_Info `json:"kf_online_list"`
	
}
//发送到多客服
func (wx *WeiXin) TransferCustomerService(kfId string, originMsg *request.WeiXinRequest) string {
	
	return fmt.Sprintf(transferCustomerService, kfId, originMsg.FromUserName, time.Now().Unix())
	
}
//客服帐号管理
func (wx *WeiXin) AddKFAccount(kf_account string, nickname string, password string) error {
	
	var KF_info_msg struct {
		
		KF_account string `json:"kf_account"`
		Nickname   string `json:"nickname"`
		Password   string `json:"password"`
	}
	
	KF_info_msg.KF_account = kf_account
	KF_info_msg.Nickname	= nickname
	KF_info_msg.Password	= password
	
	data, err := json.Marshal(&KF_info_msg)
	if err != nil {
		
		log.Println("AddKFAccount marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("AddKFAccount: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/customservice/kfaccount/add?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
func (wx *WeiXin) ModifyKFAccount(kf_account string, nickname string, password string) error {
	
	var KF_info_msg struct {
		
		KF_account string `json:"kf_account"`
		Nickname   string `json:"nickname"`
		Password   string `json:"password"`
	}
	
	KF_info_msg.KF_account = kf_account
	KF_info_msg.Nickname	= nickname
	KF_info_msg.Password	= password
	
	data, err := json.Marshal(&KF_info_msg)
	if err != nil {
		
		log.Println("ModifyKFAccount marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("ModifyKFAccount: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/customservice/kfaccount/update?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
func (wx *WeiXin) DeleteKFAccount(kf_account string, nickname string, password string) error {
	
	var KF_info_msg struct {
		
		KF_account string `json:"kf_account"`
		Nickname   string `json:"nickname"`
		Password   string `json:"password"`
	}
	
	KF_info_msg.KF_account = kf_account
	KF_info_msg.Nickname	= nickname
	KF_info_msg.Password	= password
	
	data, err := json.Marshal(&KF_info_msg)
	if err != nil {
		
		log.Println("DeleteKFAccount marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("DeleteKFAccount: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/customservice/kfaccount/del?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
//设置客服帐号的头像
func (wx *WeiXin) UploadKFHeadImg(kf_account string, filename string, reader io.Reader) (string, error) {
	
	url := weixinKFFileURL + "?kf_account=" + kf_account + "&access_token="
	return response.UploadMedia(url, retryNum, filename, nil, wx.accessTokenSupplierChan, reader)
	
}
//获取客服列表
func (wx *WeiXin) GetKFList() (*KF_List, error) {
	
	url := weixinHost + "/customservice/getkflist?access_token="
	rc, err := response.SendGetRequestWithByteResult(url, retryNum, wx.accessTokenSupplierChan)
	if err != nil {
		
		return nil, err
		
	}
	
	if rc != nil {
		
		if bytes.Contains(rc, []byte("kf_list")) {
			
			var kf_list KF_List
			if err := json.Unmarshal(rc, &kf_list); err != nil {
			
				return nil, err
			
			} else {
			
				return &kf_list, nil
			
			}
		
		} 

	}
	
	return nil, errors.New("GetKFList Unknow Error!")
	
}
//获取在线客服列表
func (wx *WeiXin) GetOnlineKFList() (*Online_KF_List, error) {
	
	url := weixinHost + "/customservice/getonlinekflist?access_token="
	rc, err := response.SendGetRequestWithByteResult(url, retryNum, wx.accessTokenSupplierChan)
	if err != nil {
		
		return nil, err
		
	}
	
	if rc != nil {
		
		if bytes.Contains(rc, []byte("kf_online_list")) {
			
			var kf_list Online_KF_List
			if err := json.Unmarshal(rc, &kf_list); err != nil {
			
				return nil, err
			
			} else {
			
				return &kf_list, nil
			
			}
		
		} 

	}
	
	return nil, errors.New("GetOnlineKFList Unknow Error!")
	
}
//客服帐号管理 end

