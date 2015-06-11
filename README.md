# golang go语言 微信公众号 SDK
    >>本人奔四了，老菜鸟程序员，用的语言较杂，皆不精，然对技术摸索多年未成系统，虽用C++开发多年，却始终不得其要，心力憔悴，
    几欲放弃，无奈饭碗不敢丢弃，苦寻大道至简之语言，绝望之际，遂意退出IT圈，终遇go语言，重燃希望之火，正努力学习中，
    故写的比较匆忙，一者学习go语言开发，二者学习微信公众号开发，主要是为了学习，实现功能，许多欠妥之处，以后慢慢改，
    参考使用了许多前辈的代码，感谢你们,一起学习一起进步！


##调用方法

* 请参看 httpserver.go 便知， 此为调用sdk的例子，可直接编译之，生成可执行文件。
* 如果编译httpserver_httprouter.go 则应首先在$GOPATH/src目录下执行go get github.com/julienschmidt/httprouter
    因为这个测试入口调用了httprouter第三方http路由库，用于分流不同域名的请求至相应系统.
* 如果编译httpserver_httprouter_negroni.go则应首先在$GOPATH/src目录下执行go get                                            github.com/goincremental/negroni-sessions , go get github.com/julienschmidt/httprouter和 go get github.com/codegangsta/negroni 此为加入negroni web server      middleware及seesion支持的版本.

##测试环境

* mac osx yosemite 10.10.3
* go version go1.4.2 darwin/amd64
* http隧道至本地开发机：./ngrok -config ngrok.cfg  -subdomain webapp.jinliangyu_weinxin_dev 8080  另一个域名：./ngrok -config ngrok.cfg  -subdomain wechat.jinliangyu_weinxin_dev 8080

##测试说明

* 微信首次接入验证通过.
* 被动消息收发测试通过(包括密文模式)
* 主动客服接口-发消息（https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=ACCESS_TOKEN）测试通过.
* 创建菜单测试通过.
* 上传下载临时素材测试通过.
* 创建临时，永久二维码，获取二维码url测试通过.
* oauth2基础测试通过，可以获取到:code,accesstoken, 最终可以获取到userinfo,如:openid,nickname等. 
 其它接口测试进行中.


##参考代码出处列表

    http://blog.csdn.net/luciswe/article/details/45890713
    http://blog.csdn.net/luciswe/article/details/45913053
    http://tonybai.com/2015/04/30/go-and-https/
    https://github.com/bigwhite/experiments/blob/master/gohttps/6-dual-verify-certs/client.go
    https://github.com/sidbusy/weixinmp   //清晰简单
    https://github.com/wizjin/weixin //复杂一些， 因为重写了http路由
    https://github.com/chanxuehong/wechat  //最复杂，更新最快。
    https://github.com/leenanxi/wechat2 //第2复杂， 更新不如上一个快。
    https://github.com/bigwhite/experiments/tree/master/wechat_examples
    https://github.com/k4shifz/Go-WeChat
    http://stackoverflow.com/questions/12122159/golang-how-to-do-a-https-request-with-bad-certificate
    http://www.peterbe.com/plog/my-favorite-go-multiplexer
    搜索引擎，许多blog，不能详尽者，在此一并感谢!


