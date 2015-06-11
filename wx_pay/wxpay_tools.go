package wx_pay

//以下接口皆为: 支付工具
//发放代金劵， 需要证书
func SendCoupon(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/mmpaymkttransfers/send_coupon", req)
	
}
//查询代金券批次信息, 不需要证书
func QueryCouponStock(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/mmpaymkttransfers/query_coupon_stock", req)
	
}
//查询代金券信息， 不需要证书
func QueryCoupon(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/promotion/query_coupon", req)
	
}
//----
//红包发放接口， 需要证书
func SendRedPack(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack", req)
	
}
//红包查询接口， 需要证书
func GetHBInfo(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/mmpaymkttransfers/gethbinfo", req)
	
}
//-----
//企业付款(给个人,以openif为准，需网页授权获得，即oauth2web模块的责任)， 需要证书
func Transfers(payClient *WXPAYClient, req map[string]string) (map[string]string, error) {
	
	return payClient.SendXMLByPostMethod("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers", req)
	
}

