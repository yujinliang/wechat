// httpserver.go
package main

import (
	
	"fmt"
	"net/http"
	"log"
	"github.com/yujinliang/wechat/mp"
	"github.com/yujinliang/wechat/mp/request"
	"github.com/yujinliang/wechat/mp/oauth2web"
	"github.com/julienschmidt/httprouter"
	
)

const (
	
	token = "weixin_@#$_x"
	appId = "wx84c5988f9d70106d"
	appSecret = "c888cc4fe859e212b2fa8642a3056861"
	encodingAESKey = "yujinliangyujinliangyujinliangyujinliang123"
	
)

func HandleVoiceMsg(wx *mp.WeiXin, w http.ResponseWriter, r *request.WeiXinRequest, timestamp, nonce string) {
	

	
	//test menu.
	//createMenu(wx)
	
	//test oauth2
	oauthConfig := oauth2web.NewOAuth2Config(wx.GetAppId(), wx.GetAppSecret(), "http://webapp.jinliangyu_weinxin_dev.tunnel.mobi/showuserinfo", "snsapi_userinfo")
	oauthUrl := oauthConfig.AuthCodeURL("testOauth2")
	//oClient := &oauth2web.Client{OAuth2Config:oauthConfig}
	//oClient.CheckAccessTokenValid()
	
	replyText := wx.ReplyText(oauthUrl, r)
	w.Write([]byte(replyText))

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
//-----------
//切换至不同的域名 start
type HostSwitch map[string]http.Handler
// Implement the ServerHTTP method on our new type
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Check if a http.Handler is registered for the given host.
    // If yes, use it to handle the request.
    if handler := hs[r.Host]; handler != nil {
        handler.ServeHTTP(w, r)
    } else {
        // Handle host names for wich no handler is registered
        http.Error(w, r.Host, 403) // Or Redirect?
    }
}
//切换至不同的域名 end
//微网站 start
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	
}
func ShowUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
	r.ParseForm()
	code  := r.FormValue("code")
	state := r.FormValue("state")
	fmt.Fprintf(w, "code: %s, state: %s\n", code, state)
	//get access token.
	oauthConfig := oauth2web.NewOAuth2Config(appId, appSecret, "http://webapp.jinliangyu_weinxin_dev.tunnel.mobi/showuserinfo", "snsapi_userinfo")
	oClient := &oauth2web.Client{OAuth2Config:oauthConfig}
	oClient.ExchangeOAuth2AccessTokenByCode(code)
	info, _ := oClient.UserInfo("zh_CN")
	fmt.Fprintf(w, "openid:%s, nickname:%s, sex:%s, city:%s, province:%s, country:%s,UnionId:%s, HeadImageURL:%s, Privilege:%v", info.OpenId,info.Nickname, info.Sex, info.City, info.Province, info.Country, info.UnionId, info.HeadImageURL,info.Privilege )
	
}

//微网站 end

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
	
	//mux bind.
	router := httprouter.New()
	router.GET("/hello/:name", Hello)
	router.GET("/showuserinfo", ShowUserInfo)
	
	//mux chain.
	muxchain := make(HostSwitch)
	muxchain["wechat.jinliangyu_weinxin_dev.tunnel.mobi"] = wx
	muxchain["webapp.jinliangyu_weinxin_dev.tunnel.mobi"] = router
	
	//http.Handle("/wechat",wx)
	err := http.ListenAndServe(":8080", muxchain)
	
	if err != nil {
		
		log.Fatal("ListenAndServe: ", err)
		
	}
}
