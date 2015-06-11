package mp

import (
	
	"log"
	"errors"
	"encoding/json"
	"mp/response"
	
)
//模板消息
type TemplateMessage struct {
	
	ToUser     string `json:"touser"`             // 必须, 接受者OpenID
	TemplateId string `json:"template_id"`        // 必须, 模版ID
	URL        string `json:"url,omitempty"`      // 可选, 用户点击后跳转的URL，该URL必须处于开发者在公众平台网站中设置的域中
	TopColor   string `json:"topcolor,omitempty"` // 可选, 整个消息的颜色, 可以不设置

	// 必须, JSON 格式的 []byte, 满足特定的模板需求
	RawJSONData json.RawMessage `json:"data"`
	
}

//模板消息
//1. 设定行业
func (wx *WeiXin) SetIndustry(industryId ...int64) error {
	
	if len(industryId) < 2 {
		
		return errors.New("industryId 的个数不能小于2")
		
	}
	
	var msg struct {
		
		IndustryId1 int64 `json:"industry_id1"`
		IndustryId2 int64 `json:"industry_id2"`
		
	}
	msg.IndustryId1 = industryId[0]
	msg.IndustryId2 = industryId[1]
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		log.Println("SetIndustry marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("SetIndustry: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/template/api_set_industry?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
// 从行业模板库选择模板添加到账号后台, 并返回模板id.
//  templateIdShort: 模板库中模板的编号，有“TM**”和“OPENTMTM**”等形式.
func (wx *WeiXin) AddTemplate2MyMP(templateIdShort string) (string, error) {
	
	var msg struct {
		
		TemplateIdShort string `json:"template_id_short"`
		
	}
	msg.TemplateIdShort = templateIdShort
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		log.Println("AddTemplate2MyMP marshal failed: ", err)
		return "", err
		
	}
	
	//just for debug.
	log.Println("AddTemplate2MyMP: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/template/api_add_template?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	
	return rc.TemplateId, nil
	
}
//发送模板消息
func (wx *WeiXin) SendTemplateMessage(msg *TemplateMessage) (string, error) {
	
	data, err := json.Marshal(msg)
	if err != nil {
		
		log.Println("SendTemplateMessage marshal failed: ", err)
		return "", err
		
	}
	
	//just for debug.
	log.Println("SendTemplateMessage: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/message/template/send?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	
	return rc.MsgId, nil
	
}
