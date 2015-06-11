package mp

import (
	
	"encoding/json"
	"log"
	"mp/response"
	
)

//高级群发接口 - 图文素材
type MPNews struct {
	
	Title 		  		string `json:"title"`
	ThumbMediaId  		string `json:"thumb_media_id"`
	Author		  		string `json:"author,omitempty"`
	Digest		  		string `json:"digest,omitempty"`
	ShowCoverPic  		int8   `json:"show_cover_pic"`
	Content 	  		string `json:"content"`
	ContentSourceUrl 	string `json:"content_source_url,omitempty"`
	
}
//上传图文素材,为高级群发接口
func (wx *WeiXin) UploadNews(news []MPNews) (string, error) {
	
	var newsMsg struct {
		
		Articles []MPNews `json:"articles"`
	}
	newsMsg.Articles = news
	
	data, err := json.Marshal(&newsMsg)
	if err != nil {
		
		log.Println("UploadNews marshal failed: ", err)
		return "", err
		
	}
	
	//just for debug.
	log.Println("UploadNews: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/media/uploadnews?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	return rc.MediaId, nil
	
}
//根据分组群发消息
func (wx *WeiXin) SendTextByGroupID(groupId string, content string, is_to_all bool) (msgid string, err error) {
	
	var msg struct {
		
		Filter struct {
			
			GroupId   string `json:"group_id"`
			Is_to_all bool	  `json:"is_to_all"`
			
		} `json:"filter"`
		
		MsgType string `json:"msgtype"`
		
		Text struct {
			
			Content string `json:"content"`
			
		} `json:"text"`
		
	}
	msg.MsgType = "text"
	msg.Filter.GroupId = groupId
	msg.Filter.Is_to_all = is_to_all
	msg.Text.Content = content

	return wx.SendAll(weixinHost + "/message/mass/sendall?access_token=", msg)
	
}

func (wx *WeiXin) SendImageByGroupID(groupId string, mediaId string, is_to_all bool) (msgid string, err error) {
	
	var msg struct {
		
		Filter struct {
			
			GroupId   string `json:"group_id"`
			Is_to_all bool	  `json:"is_to_all"`
			
		} `json:"filter"`
		
		MsgType string `json:"msgtype"`
		
		Image struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"image"`
		
	}
	msg.Filter.GroupId = groupId
	msg.Filter.Is_to_all = is_to_all
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/sendall?access_token=", msg)
	
}

func (wx *WeiXin) SendVoiceByGroupID(groupId string, mediaId string, is_to_all bool) (msgid string, err error) {
	
	var msg struct {
		
		Filter struct {
			
			GroupId   string `json:"group_id"`
			Is_to_all bool	  `json:"is_to_all"`
			
		} `json:"filter"`
		
		MsgType string `json:"msgtype"`
		
		Voice struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"voice"`
		
	}
	msg.Filter.GroupId = groupId
	msg.Filter.Is_to_all = is_to_all
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/sendall?access_token=", msg)
	
}

func (wx *WeiXin) SendVideoByGroupID(groupId string, mediaId string, is_to_all bool) (msgid string, err error) {
	
	var msg struct {
		
		Filter struct {
			
			GroupId   string `json:"group_id"`
			Is_to_all bool	  `json:"is_to_all"`
			
		} `json:"filter"`
		
		MsgType string `json:"msgtype"`
		
		Video struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"mpvideo"`
		
	}
	msg.Filter.GroupId = groupId
	msg.Filter.Is_to_all = is_to_all
	msg.MsgType = "mpvideo"
	msg.Video.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/sendall?access_token=", msg)
	
}

func (wx *WeiXin) SendNewsByGroupID(groupId string, mediaId string, is_to_all bool) (msgid string, err error) {
	
	var msg struct {
		
		Filter struct {
			
			GroupId   string `json:"group_id"`
			Is_to_all bool	  `json:"is_to_all"`
			
		} `json:"filter"`
		
		MsgType string `json:"msgtype"`
		
		News struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"mpnews"`
		
	}
	msg.Filter.GroupId = groupId
	msg.Filter.Is_to_all = is_to_all
	msg.MsgType = "mpnews"
	msg.News.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/sendall?access_token=", msg)
	
}
func (wx *WeiXin) SendWXCardByGroupID(groupId string, cardId string, is_to_all bool) (msgid string, err error) {
	
	var msg struct {
		
		Filter struct {
			
			GroupId   string `json:"group_id"`
			Is_to_all bool	  `json:"is_to_all"`
			
		} `json:"filter"`
		
		MsgType string `json:"msgtype"`
		
		WXCard struct {
			
			CardID string `json:"card_id"`
			
		} `json:"wxcard"`
		
	}
	msg.Filter.GroupId = groupId
	msg.Filter.Is_to_all = is_to_all
	msg.MsgType = "wxcard"
	msg.WXCard.CardID = cardId
	
	return wx.SendAll(weixinHost + "/message/mass/sendall?access_token=", msg)
	
}
//基础发送接口
func (wx *WeiXin) SendAll(url string, msg interface{}) (msgid string, err error) {
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		log.Println("SendAll marshal failed: ", err)
		return "", err
		
	}
	
	//just for debug.
	log.Println("SendAll: ", string(data))
	
	rc, err := response.SendPostRequest(url, retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	return rc.MassMsgId, nil
	
}
func (wx *WeiXin) UploadVideoForMass(mediaId string, title string, description string) (string, error) {
	
	var msg struct {
		
		MediaId string `json:"media_id"`
		Title	string	`json:"title"`
		Description string `json:"description"`
		
	}
	msg.MediaId = mediaId
	msg.Title   = title
	msg.Description = description
	
	data, err := json.Marshal(msg)
	if err != nil {
		
		return "", err
		
	}
	
		//just for debug.
	log.Println("UploadVideoForMass: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/media/uploadvideo?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	
	return rc.MediaId, nil
	
}
//根据OpenID列表群发[订阅号不可用， 服务号认证后可用]
func (wx *WeiXin) SendNewsByOpenIDs(toUser []string, mediaId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  []string `json:"touser,omitempty"` // 长度不能超过 ToUserCountLimit
		MsgType string   `json:"msgtype"`
		
		News struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"mpnews"`
		
	}
	msg.MsgType = "mpnews"
	msg.ToUser  = toUser
	msg.News.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/send?access_token=", msg)
	
}
func (wx *WeiXin) SendTextByOpenIDs(toUser []string, content string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  []string `json:"touser,omitempty"` // 长度不能超过 ToUserCountLimit
		MsgType string   `json:"msgtype"`
		
		Text struct {
			
			Content string `json:"content"`
			
		} `json:"text"`
		
	}
	msg.MsgType = "text"
	msg.ToUser  = toUser
	msg.Text.Content = content
	
	return wx.SendAll(weixinHost + "/message/mass/send?access_token=", msg)
	
}
func (wx *WeiXin) SendImageByOpenIDs(toUser []string, mediaId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  []string `json:"touser,omitempty"` // 长度不能超过 ToUserCountLimit
		MsgType string   `json:"msgtype"`
		
		Image struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"image"`
		
	}
	msg.MsgType = "image"
	msg.ToUser  = toUser
	msg.Image.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/send?access_token=", msg)
	
}
func (wx *WeiXin) SendVoiceByOpenIDs(toUser []string, mediaId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  []string `json:"touser,omitempty"` // 长度不能超过 ToUserCountLimit
		MsgType string   `json:"msgtype"`
		
		Voice struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"voice"`
		
	}
	msg.MsgType = "voice"
	msg.ToUser  = toUser
	msg.Voice.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/send?access_token=", msg)
	
}
func (wx *WeiXin) SendVideoByOpenIDs(toUser []string, mediaId string, title string, description string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  []string `json:"touser,omitempty"` // 长度不能超过 ToUserCountLimit
		MsgType string   `json:"msgtype"`
		
		Video struct {
			
			MediaId     string `json:"media_id"`
			Title       string `json:"title,omitempty"`
			Description string `json:"description,omitempty"`
			
		} `json:"video"`
		
	}
	msg.MsgType = "video"
	msg.ToUser  = toUser
	msg.Video.MediaId = mediaId
	msg.Video.Title   = title
	msg.Video.Description = description
	
	return wx.SendAll(weixinHost + "/message/mass/send?access_token=", msg)
	
}
func (wx *WeiXin) SendWXCardByOpenIDs(toUser []string, cardId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  []string `json:"touser,omitempty"` // 长度不能超过 ToUserCountLimit
		MsgType string   `json:"msgtype"`
		
		WXCard struct {
			
			CardID     string `json:"card_id"`
			
		} `json:"wxcard"`
		
	}

	msg.ToUser  = toUser
	msg.WXCard.CardID = cardId
	msg.MsgType = "wxcard"
	
	return wx.SendAll(weixinHost + "/message/mass/send?access_token=", msg)
	
}
//删除群发
func (wx *WeiXin) DeleteMassMessage(msgId string) (err error) {
	
	var msg struct {
		
		MassMsg_id string `json:"msg_id"`	
		
	}

	msg.MassMsg_id = msgId
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		log.Println("DeleteMassMessage marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("DeleteMassMessage: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/message/mass/delete?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return err
		
	}
	return nil
	
}
//查询群发消息发送状态
func (wx *WeiXin) QueryMassMsgStatus(msgId string) (*response.Response, error) {
	
	var msg struct {
		
		MassMsg_id string `json:"msg_id"`	
		
	}

	msg.MassMsg_id = msgId
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		log.Println("QueryMassMsgStatus marshal failed: ", err)
		return nil, err
		
	}
	
	//just for debug.
	log.Println("QueryMassMsgStatus: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/message/mass/get?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return nil, err
		
	}
	return rc, nil
	
}
//预览消息发送接口
func (wx *WeiXin) PreviewText(toUser string, content string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		
		Text struct {
			
			Content string `json:"content"`
			
		} `json:"text"`
		
	}

	msg.MsgType = "text"
	msg.ToUser = toUser
	msg.Text.Content = content
	
	return wx.SendAll(weixinHost + "/message/mass/preview?access_token=", msg)
	
}
func (wx *WeiXin) PreviewImage(toUser string, mediaId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		
		Image struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"image"`
		
	}

	msg.MsgType = "image"
	msg.ToUser = toUser
	msg.Image.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/preview?access_token=", msg)
	
}
func (wx *WeiXin) PreviewVoice(toUser string, mediaId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		
		Voice struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"voice"`
		
	}

	msg.MsgType = "voice"
	msg.ToUser = toUser
	msg.Voice.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/preview?access_token=", msg)
	
}
func (wx *WeiXin) PreviewVideo(toUser string, mediaId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		
		Video struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"mpvideo"`
		
	}

	msg.MsgType = "mpvideo"
	msg.ToUser = toUser
	msg.Video.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/preview?access_token=", msg)
	
}
func (wx *WeiXin) PreviewNews(toUser string, mediaId string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		
		News struct {
			
			MediaId string `json:"media_id"`
			
		} `json:"mpnews"`
		
	}

	msg.MsgType = "mpnews"
	msg.ToUser = toUser
	msg.News.MediaId = mediaId
	
	return wx.SendAll(weixinHost + "/message/mass/preview?access_token=", msg)
	
}
func (wx *WeiXin) PreviewWXCard(toUser string, cardId string, cardExt string) (msgid string, err error) {
	
	var msg struct {
		
		ToUser  string `json:"touser"`
		MsgType string `json:"msgtype"`
		
		WXCard struct {
			
			CardId  string `json:"card_id"`
			CardExt string `json:"card_ext"`
			
		} `json:"wxcard"`
		
	}

	msg.MsgType = "wxcard"
	msg.ToUser = toUser
	msg.WXCard.CardId  = cardId
	msg.WXCard.CardExt = cardExt
	
	return wx.SendAll(weixinHost + "/message/mass/preview?access_token=", msg)
	
}