//皆引自https://github.com/chanxuehong/wechat/blob/master/mp/user/oauth2/oauth2.go
//感谢前辈
package oauth2web

import (
	
	"bytes"
	"errors"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	
)

type OAuth2Error struct {
	
	// StructField 固定这个顺序, RETRY 依赖这个顺序
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	
}

func (e *OAuth2Error) Error() string {
	
	return fmt.Sprintf("errcode: %d, errmsg: %s", e.ErrCode, e.ErrMsg)
	
}
type OAuth2Config struct {
	
	AppId     string
	AppSecret string
	// 目前有 snsapi_base, snsapi_userinfo.
	// 应用授权作用域，多个作用域用逗号（,）分隔;
	Scope	  string
	// 用户授权后跳转的目的地址
	// 用户授权后跳转到 RedirectURL?code=CODE&state=STATE
	// 用户禁止授权跳转到 RedirectURL?state=STATE
	RedirectURL string
	
}
const (
	
	Language_oauth2_zh_CN = "zh_CN" // 简体中文
	Language_oauth2_zh_TW = "zh_TW" // 繁体中文
	Language_oauth2_en    = "en"    // 英文
	
)

const (
	
	SexUnknown_oauth2 = 0 // 未知
	SexMale_oauth2    = 1 // 男性
	SexFemale_oauth2  = 2 // 女性
	
)

const ErrCodeOK = 0

type OAuth2UserInfo struct {
	
	OpenId   string `json:"openid"`   // 用户的唯一标识
	Nickname string `json:"nickname"` // 用户昵称
	Sex      int    `json:"sex"`      // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City     string `json:"city"`     // 普通用户个人资料填写的城市
	Province string `json:"province"` // 用户个人资料填写的省份
	Country  string `json:"country"`  // 国家，如中国为CN

	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
	// 用户没有头像时该项为空
	HeadImageURL string `json:"headimgurl,omitempty"`

	// 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	Privilege []string `json:"privilege"`

	// 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
	UnionId string `json:"unionid"`
	
}

func NewOAuth2Config(AppId, AppSecret, RedirectURL string, Scope ...string) *OAuth2Config {
	
	return &OAuth2Config{
		AppId:       AppId,
		AppSecret:   AppSecret,
		Scope:       strings.Join(Scope, ","),
		RedirectURL: RedirectURL,
	}
	
}
// 构造请求用户授权获取code的地址.
//  appId:       公众号的唯一标识
//  redirectURL: 授权后重定向的回调链接地址
//               如果用户同意授权，页面将跳转至 redirect_uri/?code=CODE&state=STATE。
//               若用户禁止授权，则重定向后不会带上code参数，仅会带上state参数redirect_uri?state=STATE
//  scope:       应用授权作用域，
//               snsapi_base （不弹出授权页面，直接跳转，只能获取用户openid），
//               snsapi_userinfo （弹出授权页面，可通过openid拿到昵称、性别、所在地。
//               并且，即使在未关注的情况下，只要用户授权，也能获取其信息）
//  state:       重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
func authCodeURL(appId, redirectURL, scope, state string) string {
	
	return "https://open.weixin.qq.com/connect/oauth2/authorize" +
		"?appid=" + url.QueryEscape(appId) +
		"&redirect_uri=" + url.QueryEscape(redirectURL) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
		
}
// 请求用户授权获取code的地址.
func (cfg *OAuth2Config) AuthCodeURL(state string) string {
	
	return authCodeURL(cfg.AppId, cfg.RedirectURL, cfg.Scope, state)
	
}
//注意当把AuthCodeURL生成的url用作菜单的url时，微信服务器端会把”非法字法的错误信息，如：& 被json.marshal转化为\u0026; 可微信不认这个，报\u0026为非法字符！”
//现在写一个方法将微信认为非法的转义后字符再转化回原字符
func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	
	b, err := json.Marshal(v)
	if err != nil {
		
		return nil, err
		
	}
	
	if safeEncoding {
		
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	
	return b, nil
	
}
// 用户相关的 oauth2 token 信息
type OAuth2Token struct {
	
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64 // 过期时间, unixtime, 分布式系统要求时间同步, 建议使用 NTP

	OpenId  string
	UnionId string   // UnionID机制
	Scopes  []string // 用户授权的作用域
	
}

// 判断授权的 OAuth2Token.AccessToken 是否过期, 过期返回 true, 否则返回 false
func (token *OAuth2Token) accessTokenExpired() bool {
	
	return time.Now().Unix() >= token.ExpiresAt
	
}

type Client struct {
	
	*OAuth2Config
	*OAuth2Token // 程序会自动更新最新的 OAuth2Token 到这个字段, 如有必要该字段可以保存起来

	HttpClient *http.Client // 如果 httpClient == nil 则默认用 http.DefaultClient
	
}

func (clt *Client) httpClient() *http.Client {
	
	if clt.HttpClient != nil {
		
		return clt.HttpClient
		
	}
	
	return http.DefaultClient
	
}

// 通过code换取网页授权access_token.
//  NOTE:
//  1. Client 需要指定 OAuth2Config
//  2. 如果指定了 OAuth2Token, 则会更新这个 OAuth2Token, 同时返回的也是指定的 OAuth2Token;
//     否则会重新分配一个 OAuth2Token.
func (clt *Client) ExchangeOAuth2AccessTokenByCode(code string) (token *OAuth2Token, err error) {
	
	if clt.OAuth2Config == nil {
		
		err = errors.New("没有提供 OAuth2Config")
		return
		
	}

	tk := clt.OAuth2Token
	if tk == nil {
		
		tk = new(OAuth2Token)
		
	}

	_url := "https://api.weixin.qq.com/sns/oauth2/access_token" +
		"?appid=" + url.QueryEscape(clt.AppId) +
		"&secret=" + url.QueryEscape(clt.AppSecret) +
		"&code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"
		
	if err = clt.updateToken(tk, _url); err != nil {
		
		return
		
	}

	clt.OAuth2Token = tk
	token = tk
	return
	
}

// 刷新access_token（如果需要）.
//  NOTE: Client 需要指定 OAuth2Config, OAuth2Token
func (clt *Client) TokenRefresh() (token *OAuth2Token, err error) {
	
	if clt.OAuth2Config == nil {
		
		err = errors.New("没有提供 OAuth2Config")
		return
		
	}
	if clt.OAuth2Token == nil {
		
		err = errors.New("没有提供 OAuth2Token")
		return
		
	}
	if clt.RefreshToken == "" {
		
		err = errors.New("没有有效的 RefreshToken")
		return
		
	}

	_url := "https://api.weixin.qq.com/sns/oauth2/refresh_token" +
		"?appid=" + url.QueryEscape(clt.AppId) +
		"&grant_type=refresh_token&refresh_token=" + url.QueryEscape(clt.RefreshToken)
		
	if err = clt.updateToken(clt.OAuth2Token, _url); err != nil {
		
		return
		
	}

	token = clt.OAuth2Token
	return
	
}

// 检验授权凭证（access_token）是否有效.
//  NOTE:
//  1. Client 需要指定 OAuth2Token
//  2. 先判断 err 然后再判断 valid
func (clt *Client) CheckAccessTokenValid() (valid bool, err error) {
	
	if clt.OAuth2Token == nil {
		
		err = errors.New("没有提供 OAuth2Token")
		return
		
	}
	if clt.AccessToken == "" {
		
		err = errors.New("没有有效的 AccessToken")
		return
		
	}
	if clt.OpenId == "" {
		
		err = errors.New("没有有效的 OpenId")
		return
		
	}

	_url := "https://api.weixin.qq.com/sns/auth?access_token=" + url.QueryEscape(clt.AccessToken) +
		"&openid=" + url.QueryEscape(clt.OpenId)
		
	httpResp, err := clt.httpClient().Get(_url)
	if err != nil {
		
		return
		
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
		
	}

	var result OAuth2Error
	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		
		return
		
	}

	switch result.ErrCode {
		
	case ErrCodeOK: {
		
		valid = true
		return
		
	}
	case 40001: {
		
		return
		
	}
	default: {
		
		err = &result
		return
		
	}
	
	}
	
}

// 从服务器获取新的 token 更新 tk
func (clt *Client) updateToken(tk *OAuth2Token, url string) (err error) {
	
	if tk == nil {
		
		return errors.New("nil OAuth2Token")
		
	}

	httpResp, err := clt.httpClient().Get(url)
	if err != nil {
		
		return
		
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		
		return fmt.Errorf("http.Status: %s", httpResp.Status)
		
	}

	var result struct {
		
		OAuth2Error
		AccessToken  string `json:"access_token"`  // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
		RefreshToken string `json:"refresh_token"` // 用户刷新access_token
		ExpiresIn    int64  `json:"expires_in"`    // access_token接口调用凭证超时时间，单位（秒）
		OpenId       string `json:"openid"`        // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
		UnionId      string `json:"unionid"`       // UnionID机制
		Scope        string `json:"scope"`         // 用户授权的作用域，使用逗号（,）分隔
		
	}

	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		
		return
		
	}

	if result.ErrCode != ErrCodeOK {
		
		return &result.OAuth2Error
		
	}

	// 由于网络的延时, 分布式服务器之间的时间可能不是绝对同步, access_token 过期时间留了一个缓冲区;
	switch {
		
	case result.ExpiresIn > 31556952: {// 60*60*24*365.2425
	
		err = errors.New("expires_in too large: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
		
	}
	case result.ExpiresIn > 60*60: {
		
		result.ExpiresIn -= 60 * 20
		
	}
	case result.ExpiresIn > 60*30: {
		
		result.ExpiresIn -= 60 * 10
		
	}
	case result.ExpiresIn > 60*15: {
		
		result.ExpiresIn -= 60 * 5
		
	}
	case result.ExpiresIn > 60*5: {
		
		result.ExpiresIn -= 60
		
	}
	case result.ExpiresIn > 60: {
		
		result.ExpiresIn -= 20
		
	}
	default: {
		
		err = errors.New("expires_in too small: " + strconv.FormatInt(result.ExpiresIn, 10))
		return
		
	}
		
	}//end of switch.

	tk.AccessToken = result.AccessToken
	
	if result.RefreshToken != "" {
		
		tk.RefreshToken = result.RefreshToken
		
	}
	tk.ExpiresAt = time.Now().Unix() + result.ExpiresIn

	tk.OpenId = result.OpenId
	tk.UnionId = result.UnionId

	strs := strings.Split(result.Scope, ",")
	tk.Scopes = make([]string, 0, len(strs))
	for _, str := range strs {
		
		str = strings.TrimSpace(str)
		if str == "" {
			
			continue
			
		}
		tk.Scopes = append(tk.Scopes, str)
		
	}

	return
}

// 获取用户信息(需scope为 snsapi_userinfo).
//  NOTE:
//  1. Client 需要指定 OAuth2Config, OAuth2Token
//  2. lang 可能的取值是 zh_CN, zh_TW, en, 如果留空 "" 则默认为 zh_CN.
func (clt *Client) UserInfo(lang string) (info *OAuth2UserInfo, err error) {
	
	switch lang {
		
	case "": {
		
		lang = Language_oauth2_en
		
	}
	case Language_oauth2_zh_CN, Language_oauth2_zh_TW, Language_oauth2_en: {
		
		//do nothing.
	}
	default: {
		
		err = errors.New("错误的 lang 参数")
		return
		
	}
	
	}

	if clt.OAuth2Config == nil { // clt.TokenRefresh() 需要
	
		err = errors.New("没有提供 OAuth2Config")
		return
		
	}
	if clt.OAuth2Token == nil {
		
		err = errors.New("没有提供 OAuth2Token")
		return
		
	}

	if clt.accessTokenExpired() {
		
		if _, err = clt.TokenRefresh(); err != nil {
			
			return
			
		}
		
	}

	if clt.AccessToken == "" {
		
		err = errors.New("没有有效的 AccessToken")
		return
		
	}
	if clt.OpenId == "" {
		
		err = errors.New("没有有效的 OpenId")
		return
		
	}

	_url := "https://api.weixin.qq.com/sns/userinfo" +
		"?access_token=" + url.QueryEscape(clt.AccessToken) +
		"&openid=" + url.QueryEscape(clt.OpenId) +
		"&lang=" + url.QueryEscape(lang)
		
	httpResp, err := clt.httpClient().Get(_url)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		
		err = fmt.Errorf("http.Status: %s", httpResp.Status)
		return
		
	}

	var result struct {
		
		OAuth2Error
		OAuth2UserInfo
		
	}

	if err = json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		
		return
		
	}

	if result.ErrCode != ErrCodeOK {
		
		err = &result.OAuth2Error
		return
		
	}
	info = &result.OAuth2UserInfo
	return
	
}

