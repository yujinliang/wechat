package wx_pay

import (
	
	"fmt"
	"bytes"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/xml"
	"github.com/yujinliang/wechat/mp/utils"
	
)

const(
	
	APPID 			= "wx426b3015555a46be" //APPID：绑定支付的APPID（必须配置）
	APPSECRET 		= "01c6d59a3f9024db6336662ac95c8e74" //APPSECRET：公众帐号secert（仅JSAPI支付的时候需要配置）
	MCHID 			= "1225312702" //商户号（必须配置）
	KEY	  			= "e10adc3949ba59abbe56e057f20f883e" //商户支付密钥(API密钥)，参考开户邮件设置（必须配置
	//证书路径设置.
	SSLCERT_PATH 	= "./cert/apiclient_cert.pem"
	SSLKEY_PATH	  	= "./cert/apiclient_key.pem"
	
	//status.
	RETURN_CODE_SUCCESS	 = "SUCCESS"
	RETURN_CODE_FAIL	 = "FAIL"
	RESULT_CODE_SUCCESS = "SUCCESS"
	RESULT_CODE_FAIL	 = "FAIL"
	
)

//参考：https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/07.1.md
//https://github.com/chanxuehong/wechat/blob/master/mch/error.go
type WXPayError struct {
	
	XMLName	   xml.Name	`xml:"xml"`
	ReturnCode string 	`xml:"return_code"`
	ReturnMsg  string 	`xml:"return_msg,omitempty"`
	
}
func (e *WXPayError) Error() string {
	
	return fmt.Sprintf("return_code: %q, return_msg: %q", e.ReturnCode, e.ReturnMsg)
}
type WXPAYClient struct {
	
	curl *http.Client
	
}

func NewWXPAYClient(httpClient *http.Client) *WXPAYClient {
	
	if httpClient == nil {
		
		//默认应用无证书httpClient
		httpClient = http.DefaultClient
		
	}
	
	return &WXPAYClient {
		
		curl : httpClient,

	}
}
//发包接口
//参考:https://github.com/chanxuehong/wechat/blob/master/mch/client.go
//微信支付接口说明: 此处只检查到return_code, 而：result_code是业务逻辑相关的，应放到业务代码中检查.
func (pay *WXPAYClient) SendXMLByPostMethod(url string, req map[string]string) (map[string]string, error) {
	
	xmlBytes, err := utils.Map2XMLBytes(req)
	if err != nil {
		
		return nil, err
		
	}
	
	rc, err := pay.curl.Post(url, "text/xml; charset=utf-8", bytes.NewBuffer(xmlBytes))
	if err != nil {
		
		return nil, err
		
	}
	defer rc.Body.Close()
	
	if rc.StatusCode != http.StatusOK {
		
		return nil, fmt.Errorf("http status: %s", rc.Status)
		
	}
	
	xmlMap, err := utils.XML2Map(rc.Body)
	if err != nil {
		
		return nil, err
		
	}
	
	//check weixin pay, return code.
	returnCode, ok := xmlMap["return_code"]
	if !ok {
		
		return xmlMap, fmt.Errorf("no return_code") 
		
	}
	if returnCode != RETURN_CODE_SUCCESS {
		
		return xmlMap, fmt.Errorf("return_code: %s, return_msg: %s", returnCode, xmlMap["return_msg"])
		
	}
	
	//认证签名
	signatureOrigin, ok := xmlMap["sign"]
	if !ok {
		
		return xmlMap, errors.New("no sign")
	}
	signatureNow := utils.SignForWXPay(xmlMap, KEY, nil)
	if signatureOrigin != signatureNow {
		
		return xmlMap, fmt.Errorf("signature mismatch, origin: %s , now: %s", signatureOrigin, signatureNow)
		
	}
	
	return xmlMap, nil
	
}
//统一下单，不需要证书.
func UnifiedOrder(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/pay/unifiedorder", req)
	
}
//查询订单，不需要证书.
func OrderQuery(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/pay/orderquery", req)
	
}
//关闭订单，不需要证书.
func CloseOrder(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/pay/closeorder", req)
	
}
//申请退款，需要证书.
//请用utils.NewTLSHttpClient创建双向证书httpClient.
func Refund(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/secapi/pay/refund", req)
	
}
//查询退款,不需要证书.
func RefundQuery(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/pay/refundquery", req)
	
}
//下载对账单，不需要证书.
func DownloadBill(payClient *WXPAYClient, req map[string]string) ([]byte, error) {
	
	xmlData, err := utils.Map2XMLBytes(req)
	if err != nil {
		
		return nil, err
		
	}
	
	rc, err := payClient.curl.Post("https://api.mch.weixin.qq.com/pay/downloadbill", "text/xml; charset=utf-8", bytes.NewBuffer(xmlData))
	if err != nil {
		
		return nil, err
		
	}
	defer rc.Body.Close()
	
	if rc.StatusCode != http.StatusOK {
		
		return nil, fmt.Errorf("http status: %s", rc.Status)
		
	}
	
	respContent, err := ioutil.ReadAll(rc.Body)
	if err != nil {
		
		return nil, err
		
	}
	
	var result WXPayError
	if err := xml.Unmarshal(respContent, &result); err == nil {
		
		return nil, &result
		
	}
	
	return respContent, nil
	
}
//测速上报，不需要证书.
func Report(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/payitil/report", req)
	
}
//长链接转化为短链接.
func ToShortURL(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/tools/shorturl", req)
	
}
//主动扫支付（扫码支付) 与公众号支付(js) 接口相同， 都是以上接口!!!
//公众号支付(js), 支付结果通用通知 : http://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_7
//支付完成后，微信会把相关支付结果和用户信息发送给商户，商户需要接收处理，并返回应答。
//该链接是通过【统一下单API】中提交的参数notify_url设置，如果链接无法访问，商户将无法接收到微信通知。
//回调函数原型声明
type PayNotifyCallback func(map[string]string, *http.ResponseWriter) error
//支付结果通用通知，到达时， 调用此接口，注册处理回调函数.
func Notify(r *http.Request, w *http.ResponseWriter, callback PayNotifyCallback) error {
	
	//1. 从http request解析出parameters, 汇总进map.
	xmlMap, err := utils.XML2Map(r.Body)
	if err != nil {
		
		return err
		
	}
	//2. call callback, if no empty; 并将处理结果以下面结构发回微信服务器端.
	//<xml>
    //<return_code><![CDATA[SUCCESS]]></return_code>
    //<return_msg><![CDATA[OK]]></return_msg>
	//</xml>
	return callback(xmlMap, w)
	
}
//----
//只有被动扫支付有以下两个接口不同而已！！！
//被动扫支付（刷卡支付)
func Micropay(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/pay/micropay", req)
	
}
//被动扫支付（刷卡支付）- 撤销订单
func ReverseOrder(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/secapi/pay/reverse", req)
	
}
//以上接口皆为：支付方式

