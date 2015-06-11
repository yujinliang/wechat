// httpserver.go
package main

import (
	
	"net/http"
	"log"
	"mp"
	"mp/request"
	"mp/oauth2web"
	
)

const (
	
	token = "weixin_@#$_x"
	
)

func HandleVoiceMsg(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("语音消息!", r)
	w.Write([]byte(replyText))
	
	createMenu(wx)
	
}
func HandleTextMsg(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("文本消息!", r)
	//data, _ := wx.MakeEncryptResponse([]byte(replyText), timestamp, nonce)
	w.Write([]byte(replyText))
	
	//send custom message
	//wx.PostText(r.FromUserName, "我是客服消息， 你好！", "")
	
}
func HandleImgeMsg(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("图片消息!", r)
	w.Write([]byte(replyText))
}
func HandleVideoMsg(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("视频消息!", r)
	w.Write([]byte(replyText))
}
func HandleLocationMsg(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("位置消息!", r)
	w.Write([]byte(replyText))
}
func HandleLinkMsg(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("link消息!", r)
	w.Write([]byte(replyText))
}
func HandleSubscribeEvent(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("订阅事件!", r)
	w.Write([]byte(replyText))
}
func HandleUnSubscribeEvent(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("取消订阅事件!", r)
	w.Write([]byte(replyText))
}
func HandleScanEvent(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("扫二维码事件！", r)
	w.Write([]byte(replyText))
}
func HandleLocationEvent(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("位置事件!", r)
	w.Write([]byte(replyText))
}
func HandleMenuClickEvent(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("菜单点击事件!", r)
	w.Write([]byte(replyText))
}
func HandleMenuViewEvent(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	
	replyText := wx.ReplyText("打开网页事件!", r)
	w.Write([]byte(replyText))
}
//test menu
func createMenu(wx *mp.WeiXin) {
	
	menu := &mp.Menu{make([]mp.MenuButton,3)}
	menu.Buttons[0].Name = "我要打七"
	menu.Buttons[0].Type = mp.ButtonTypeView
	menu.Buttons[0].Url  = "https://mp.weixin.qq.com"
	menu.Buttons[1].Name = "结缘法宝"
	menu.Buttons[1].Type = mp.ButtonTypeView
	menu.Buttons[1].Url  = "https://mp.weixin.qq.com"
	menu.Buttons[2].Name = "打七论坛"
	menu.Buttons[2].SubButtons = make([]mp.MenuButton,1)
	menu.Buttons[2].SubButtons[0].Name = "分享"
	menu.Buttons[2].SubButtons[0].Type = mp.ButtonTypeClick
	menu.Buttons[2].SubButtons[0].Key  = "TestKey001"
	
	err := wx.CreateMenu(menu)
	
	if err != nil {
		
		log.Println(err)
		
	}
}
func testOAuth2(wx *mp.WeiXin, AppId, AppSecret, RedirectURL string, Scope ...string) {
	
	oauthConfig := oauth2web.NewOAuth2Config(AppId, AppSecret, RedirectURL, Scope...)
	oClient := &oauth2web.Client{OAuth2Config:oauthConfig}
	oClient.CheckAccessTokenValid()
	
}
func main() {
	
	//http.HandleFunc("/", InvalidateForWeiXin)
	//err := http.ListenAndServe(":8080", nil)
	wx := mp.New(token, "wx84c5988f9d70106d", "c888cc4fe859e212b2fa8642a3056861" ,"yujinliangyujinliangyujinliangyujinliang123", "")
	wx.HandleFunc(mp.MsgTypeText,  HandleTextMsg)
	wx.HandleFunc(mp.MsgTypeVoice, HandleVoiceMsg)
	wx.HandleFunc(mp.MsgTypeImage, HandleImgeMsg)
	wx.HandleFunc(mp.MsgTypeVideo, HandleVideoMsg)
	wx.HandleFunc(mp.MsgTypeLocation, HandleLocationMsg)
	wx.HandleFunc(mp.MsgTypeLink, HandleLinkMsg)
	//event.
	wx.HandleFunc(mp.GenHttpRouteKey(mp.MsgTypeEvent, mp.EventSubscribe), HandleSubscribeEvent)
	wx.HandleFunc(mp.GenHttpRouteKey(mp.MsgTypeEvent, mp.EventUnsubscribe), HandleUnSubscribeEvent)
	wx.HandleFunc(mp.GenHttpRouteKey(mp.MsgTypeEvent, mp.EventScan), HandleScanEvent)
	wx.HandleFunc(mp.GenHttpRouteKey(mp.MsgTypeEvent, mp.EventLocation), HandleLocationEvent)
	wx.HandleFunc(mp.GenHttpRouteKey(mp.MsgTypeEvent, mp.EventClick), HandleMenuClickEvent)
	wx.HandleFunc(mp.GenHttpRouteKey(mp.MsgTypeEvent, mp.EventView), HandleMenuViewEvent)
	
	http.Handle("/",wx)
	err := http.ListenAndServe(":8080", nil)
	
	if err != nil {
		
		log.Fatal("ListenAndServe: ", err)
		
	}
}
