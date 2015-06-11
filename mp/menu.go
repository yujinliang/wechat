package mp

import (
	
	"log"
	"encoding/json"
	"github.com/yujinliang/wechat/mp/response"
	
)

const (
	
	// button types
	ButtonTypeClick = "click"
	ButtonTypeView  = "view"
	
)
//自定义菜单
type Menu struct {
	
	Buttons []MenuButton `json:"button,omitempty"`
	
}

type MenuButton struct {
	
	Name       string       `json:"name"`
	Type       string       `json:"type,omitempty"`
	Key        string       `json:"key,omitempty"`
	Url        string       `json:"url,omitempty"`
	MediaId	   string		 `json:"media_id,omitempty"`
	SubButtons []MenuButton `json:"sub_button,omitempty"`
	
}

//自定义菜单
func (wx *WeiXin) CreateMenu(menu *Menu) error {
	
	data, err := json.Marshal(menu)
	if err != nil {
		
		return err
		
	} 
		
	//just for debug.
	log.Println("CreateMenu: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/menu/create?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
func (wx *WeiXin) GetMenu() (*Menu, error) {
	
	reply, err := response.SendGetRequestWithByteResult(weixinHost+"/menu/get?access_token=", retryNum, wx.accessTokenSupplierChan)
	if err != nil {
		
		return nil, err
		
	}
	
	var result struct {
		
		MenuCore *Menu `json:"menu"`
		
	}
	
	if err:= json.Unmarshal(reply, &result); err != nil {
		
		return nil, err
		
	}
	
	return result.MenuCore, nil
	
}
func (wx *WeiXin) DeleteMenu() error {
	
	_, err := response.SendGetRequestWithByteResult(weixinHost + "/menu/delete?access_token=", retryNum, wx.accessTokenSupplierChan)
	return err
	
}
