//摘自: https://github.com/chanxuehong/wechat/blob/master/mch/native_url.go
//同时参考：微信支付php sdk.
package wx_pay

import (
	
	"net/url"
	"github.com/yujinliang/wechat/mp/utils"
	
)
//参考： http://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=6_4
// 扫码原生支付模式1的地址
func NativeURL1(appId, mchId, productId, timestamp, nonceStr, apiKey string) string {
	
	m := make(map[string]string, 5)
	m["appid"] = appId
	m["mch_id"] = mchId
	m["product_id"] = productId
	m["time_stamp"] = timestamp
	m["nonce_str"] = nonceStr

	signature := utils.SignForWXPay(m, apiKey, nil)

	return "weixin://wxpay/bizpayurl?sign=" + signature +
		"&appid=" + url.QueryEscape(appId) +
		"&mch_id=" + url.QueryEscape(mchId) +
		"&product_id=" + url.QueryEscape(productId) +
		"&time_stamp=" + url.QueryEscape(timestamp) +
		"&nonce_str=" + url.QueryEscape(nonceStr)
		
}
//参考：http://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=6_5
// 扫码原生支付模式2的地址
func NativeURL2(codeURL string) string {
	
	return "weixin://wxpay/bizpayurl?sr=" + url.QueryEscape(codeURL)
	
}
//以上两接口，用于生成支付url, 将此url生成二维码，以供微信客户端扫一扫
//--------
//根据: http://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7
//解析“统一下单“接口返回的结果，用于构建jsapi参数，发起“支付预支付订单调用”
func GetJsApiParameters(UnifiedOrderResultMap map[string]string) string {
	
	//参数:Appid, timeStamp, nonceStr, prepay_id, MD5, sign.
	//json_encode to string.
	return ""
	
}

