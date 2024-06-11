#golang go语言 微信公众号 SDK
    

#调用方法

* 请参看 httpserver.go 便知， 此为调用sdk的例子，可直接编译之，生成可执行文件。
* 如果编译httpserver_httprouter.go 则应首先在$GOPATH/src目录下执行go get github.com/julienschmidt/httprouter
    因为这个测试入口调用了httprouter第三方http路由库，用于分流不同域名的请求至相应系统.
* 如果编译httpserver_httprouter_negroni.go则应首先在$GOPATH/src目录下执行go get                                            github.com/goincremental/negroni-sessions , go get github.com/julienschmidt/httprouter和 go get github.com/codegangsta/negroni 此为加入negroni web server      middleware及seesion支持的版本.

#测试环境

* mac osx yosemite 10.10.3
* go version go1.4.2 darwin/amd64
* http隧道至本地开发机：./ngrok -config ngrok.cfg  -subdomain webapp.jinliangyu_weinxin_dev 8080  另一个域名：./ngrok -config ngrok.cfg  -subdomain wechat.jinliangyu_weinxin_dev 8080

#测试说明

* 微信首次接入验证通过.
* 被动消息收发测试通过(包括密文模式)
* 主动客服接口-发消息（https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=ACCESS_TOKEN）测试通过.
* 创建菜单测试通过.
* 上传下载临时素材测试通过.
* 创建临时，永久二维码，获取二维码url测试通过.
* oauth2基础测试通过，可以获取到:code,accesstoken, 最终可以获取到userinfo,如:openid,nickname等. 
* 高级群发接口:预览接口通过； 基于分组群发接口流程通过，只是没有权限所以看不到实际效果.
* oauth2生成的授权url 与 菜单共用时报错已修复，因为golang json包marshal时自动对一些字符如:&做了转义，但微信不认。
 其它接口测试进行中.


#参考代码出处列表

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
    Go语言编程，许式伟，吕桂华著
    Go语言程序设计，英Mark Summerfield著
    Go Web 编程， 谢孟军著
    搜索引擎，许多blog，未能详尽者，在此一并感谢!


#共同学习

    @ Wechat - 2015
    作者: 于金良
    weixin: lingshanxingzhe-pure
    邮箱: 285779289@qq.com
    csdn: https://blog.csdn.net/htyu_0203_39
    zhihu: https://www.zhihu.com/people/yujinliang-pure
    心声：一入江湖无踪影，归来依旧少年郎！
    
    * 早年写的一个golang学习随笔：
    https://note.youdao.com/s/M64kuqqT
