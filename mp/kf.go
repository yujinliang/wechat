//注意其与autoreply.go共用一些相同的数据类定义,如: Article, Music, 以及一些常量.
package mp

import (
	
	"log"
	"encoding/json"
	"github.com/yujinliang/wechat/mp/response"
	"github.com/yujinliang/wechat/mp/accesstoken"
	
)
//用指定客服帐号发消息给客户.
type Customservice struct {
			
	KF_account string `json:kf_account,omitempty`
					
}

//客服接口 - 发送消息
func (wx *WeiXin) PostText(touser string, text string, kf_account string) error {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Text    struct {
			
			Content string `json:"content"`
		} `json:"text"`
		
		Customservice *Customservice `json:"customservice,omitempty"`
		
	}
		
		msg.ToUser = touser
		msg.MsgType = "text"
		msg.Text.Content = text
		if len(kf_account) > 0 {
			
			msg.Customservice = &Customservice{KF_account:kf_account}
			
		} 
		return postMessage(wx.accessTokenSupplierChan, &msg)
	
}
func (wx *WeiXin) PostImage(touser string, mediaId string, kf_account string) error {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Image	struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"image"`
		
		Customservice *Customservice `json:"customservice,omitempty"`
		
	}
	
	msg.ToUser = touser
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	if len(kf_account) > 0 {
			
		msg.Customservice = &Customservice{KF_account:kf_account}
			
	} 
	return postMessage(wx.accessTokenSupplierChan, &msg)
	
}
func (wx *WeiXin) PostVoice(touser string, mediaId string, kf_account string) error {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Voice   struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"voice"`
		
		Customservice *Customservice `json:"customservice,omitempty"`
		
	}
	msg.ToUser = touser
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	if len(kf_account) > 0 {
			
		msg.Customservice = &Customservice{KF_account:kf_account}
			
	} 
	return postMessage(wx.accessTokenSupplierChan, &msg)
	
}
func (wx *WeiXin) PostVideo(touser string, mediaId string, thumbMediaId string, title string, description string, kf_account string) error {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Video   struct {
			
			MediaId      string `json:"media_id"`
			ThumbMediaId string `json:"thumb_media_id"`
			Title        string `json:"title"`
			Description  string `json:"description"`
			
		} `json:"video"`
		
		Customservice *Customservice `json:"customservice,omitempty"`
		
	}
	msg.ToUser = touser
	msg.MsgType = "video"
	msg.Video.MediaId = mediaId
	msg.Video.ThumbMediaId = thumbMediaId
	msg.Video.Title = title
	msg.Video.Description = description
	if len(kf_account) > 0 {
			
		msg.Customservice = &Customservice{KF_account:kf_account}
			
	} 
	return postMessage(wx.accessTokenSupplierChan, &msg)
	
}
func (wx *WeiXin) PostMusic(touser string, music *Music, kf_account string) error {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		Music   *Music `json:"music"`
		
		Customservice *Customservice `json:"customservice,omitempty"`
		
	}
	msg.ToUser = touser
	msg.MsgType = "music"
	msg.Music = music
	if len(kf_account) > 0 {
			
		msg.Customservice = &Customservice{KF_account:kf_account}
			
	} 
	return postMessage(wx.accessTokenSupplierChan, &msg)
	
}
func (wx *WeiXin) PostNews(touser string, articles []Article, kf_account string) error {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		News    struct {
			
			Articles []Article `json:"articles"`
			
		} `json:"news"`
		
		Customservice *Customservice `json:"customservice,omitempty"`
		
	}
	msg.ToUser = touser
	msg.MsgType = "news"
	msg.News.Articles = articles
	if len(kf_account) > 0 {
			
		msg.Customservice = &Customservice{KF_account:kf_account}
			
	} 
	return postMessage(wx.accessTokenSupplierChan, &msg)
	
}
func postMessage(c chan accesstoken.AccessToken, msg interface{}) error {
	
	data, err := json.Marshal(msg)
	if err != nil {
		
		log.Println("postMessage marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("postMessage: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/message/custom/send?access_token=", retryNum, c, data)
	return err
	
}
//客服接口 - 发送消息 end
