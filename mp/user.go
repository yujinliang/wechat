package mp

import (
	
	"errors"
	"log"
	"encoding/json"
	"mp/response"
	
)

//用户分组
type UserGroup struct {
	
	Id        int64  `json:"id"`    // 分组id, 由微信分配
	Name      string `json:"name"`  // 分组名字, UTF8编码
	UserCount int    `json:"count"` // 分组内用户数量
	
}
//用户基本信息， 非OAuth
type UserInfo struct {
	
	Subscribe  int8 `json:"subscribe"` //0：未关注该公众号，拉取不到其余信息； 1：为已关注.
	OpenId   string `json:"openid"`   // 用户的标识，对当前公众号唯一
	Nickname string `json:"nickname"` // 用户的昵称
	Sex      int    `json:"sex"`      // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Language string `json:"language"` // 用户的语言，zh_CN，zh_TW，en
	City     string `json:"city"`     // 用户所在城市
	Province string `json:"province"` // 用户所在省份
	Country  string `json:"country"`  // 用户所在国家

	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），
	// 用户没有头像时该项为空
	HeadImageURL string `json:"headimgurl,omitempty"`

	// 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
	SubscribeTime int64 `json:"subscribe_time"`

	// 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	UnionId string `json:"unionid,omitempty"`

	// 备注名
	Remark string `json:"remark,omitempty"`
	
	// 组id
	Groupid string `json:"groupid,omitempty"`
	
}
//获取关注者列表返回的数据结构
type UserListResult struct {
	
	TotalCount int `json:"total"` // 关注该公众账号的总用户数
	GotCount   int `json:"count"` // 拉取的OPENID个数，最大值为10000

	Data struct {
		
		OpenId []string `json:"openid,omitempty"`
		
	} `json:"data"` // 列表数据，OPENID的列表

	// 拉取列表的后一个用户的OPENID, 如果 next_openid == "" 则表示没有了用户数据
	NextOpenId string `json:"next_openid"`
	
}

//用户管理
//1. 创建分组
// 一个公众账号，最多支持创建100个分组。
func (wx *WeiXin) CreateUserGroup(groupName string) (*UserGroup, error) {
	
	if len(groupName) <= 0 {
		
		return nil, errors.New("Empty group name!")
		
	}
	
	var msg struct {
		
		Group struct {
			
			Name string `json:"name"`
			
		} `json:"group"`
		
	}
	msg.Group.Name = groupName
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		return nil, err
		
	}
	
	//just for debug.
	log.Println("CreateUserGroup: ", string(data))
	
	rc, err := response.SendPostRequestWithByteResult(weixinHost + "/groups/create?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return nil, err
		
	}
	
	var result struct {
		
		UserGroup `json:"group"`
		
	}
	if err := json.Unmarshal(rc, &result); err != nil {
		
		return nil, err
		
	}
	
	result.UserGroup.UserCount = 0
	return &result.UserGroup, nil
	
}
//2. 获取分组列表
func (wx *WeiXin) GetGroupList() ([]UserGroup, error) {
	
	rc, err := response.SendGetRequestWithByteResult(weixinHost + "/groups/get?access_token=", retryNum, wx.accessTokenSupplierChan)
	if err != nil {
		
		return nil, err
		
	}
	
	var result struct {
		
		GroupList []UserGroup `json:"groups"`
		
	}
	
	result.GroupList = 	make([]UserGroup, 0, 16)
	
	if err := json.Unmarshal(rc, &result); err != nil {
		
		return nil, err
		
	}
	
	return result.GroupList, nil
	
}
//3. 查询用户所在分组id
func (wx *WeiXin) VerifyUserInWhichGroup(openId string) (string, error) {
	
	var msg struct {
		
		OpenId string `json:"openid"`
		
	}
	msg.OpenId = openId
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		return "", err
		
	}
	//just for debug.
	log.Println("VerifyUserInWhichGroup: ", string(data))
	
	rc, err := response.SendPostRequestWithByteResult(weixinHost + "/groups/getid?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	
	var result struct {
		
		GroupId string `json:"groupid"`
		
	}
	if err := json.Unmarshal(rc, &result); err != nil {
		
		return "", err
		
	}
	
	return result.GroupId, nil
	
}
//4. 修改分组名
func (wx *WeiXin) ModifyGroupName(groupId string, newName string) error {
	
	if len(groupId) <=0 || len(newName) <= 0 {
		
		return errors.New("Empty groupid, or new name!")
		
	}
	
	var msg struct {
		
		Group struct {
			
			Id 	 string `json:"id"`
			Name string `json:"name"`
			
		} `json:"group"`
	}
	msg.Group.Id = groupId
	msg.Group.Name = newName
	data, err := json.Marshal(&msg) 
	if err != nil {
			
		return err
			
	}
		
	//just for debug.
	log.Println("ModifyGroupName: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/groups/update?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
//5. 移动用户分组
func (wx *WeiXin) MoveUser2Group(openId string, groupId string) error {
	
	if len(groupId) <=0 || len(openId) <= 0 {
		
		return errors.New("Empty groupid, or openid!")
		
	}
	
	var msg struct {
		
		OpenId    string `json:"openid"`
		ToGroupId string `json:"to_groupid"`
		
	}
	msg.OpenId = openId
	msg.ToGroupId = groupId

	data, err := json.Marshal(&msg) 
	if err != nil {
			
		return err
			
	}
		
	//just for debug.
	log.Println("MoveUser2Group: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/groups/members/update?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
//6. 批量移动用户分组
//openid个数<=50
func (wx *WeiXin) MoveUsers2Group(openId[] string, groupId string) error {
	
	if len(groupId) <=0 || len(openId) <= 0 {
		
		return errors.New("Empty groupid, or openid!")
		
	}
	
	var msg struct {
		
		OpenIdList    []string `json:"openid_list"`
		ToGroupId 	   string `json:"to_groupid"`
		
	}
	msg.OpenIdList = openId
	msg.ToGroupId = groupId

	data, err := json.Marshal(&msg) 
	if err != nil {
			
		return err
			
	}
		
	//just for debug.
	log.Println("MoveUsers2Group: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/groups/members/batchupdate?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
//7.删除分组
func (wx *WeiXin) DeleteGroup(groupId string) error {
	
	if len(groupId) <=0 {
		
		return errors.New("Empty groupid!")
		
	}
	
	var msg struct {
		
		Group struct {
			
			Id string `json:"id"`
			
		} `json:"group"`
		
	}
	msg.Group.Id = groupId

	data, err := json.Marshal(&msg) 
	if err != nil {
			
		return err
			
	}
		
	//just for debug.
	log.Println("DeleteGroup: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/groups/delete?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
//8.设置用户备注名
func (wx *WeiXin) UpdateRemark(openId string, remark string) error {
	
	if len(remark) <=0 || len(openId) <= 0 {
		
		return errors.New("Empty remark, or openid!")
		
	}
	
	var msg struct {
		
		OpenId     string `json:"openid"`
		Remark 	   string `json:"remark"`
		
	}
	msg.OpenId = openId
	msg.Remark = remark

	data, err := json.Marshal(&msg) 
	if err != nil {
			
		return err
			
	}
		
	//just for debug.
	log.Println("UpdateRemark: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/user/info/updateremark?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	return err
	
}
//9.获取用户基本信息
//lang = ["zh_CN", "zh_TW", "en"]
func (wx *WeiXin) GetUserInfo(openId string, lang string) (*UserInfo, error) {
	
	if len(openId) <= 0 || len(lang) <= 0 {
		
		return nil, errors.New("Empty openid or lang!")
		
	}
	
	switch lang {
		
		case "": {
			
			lang = "en"
			
		}
		case "zh_CN","zh_TW","en": {
			
			//do nothing.
		}
		default: {
			
			return nil, errors.New("Invalid lang: " + lang)
			
		}
		
	}
	
	rc, err := response.SendGetRequestWithByteResult( weixinHost + "/user/info?openid=" + openId + "&lang=" + lang + "&access_token=", retryNum, wx.accessTokenSupplierChan)
	if err != nil {
		
		return nil, err
		
	}
	
	var user UserInfo
	if err = json.Unmarshal(rc, &user); err != nil {
		
		return nil, err
		
	}
	
	return &user, nil
	
}
func (wx *WeiXin) GetUserList(nextOpenId string) (*UserListResult, error) {
	
	var url string
	if len(nextOpenId) > 0 {
		
		url = weixinHost + "/user/get?next_openid=" + nextOpenId + "&access_token="
		
	} else {
		
		url = weixinHost + "/user/get?access_token="
		
	}
	
	rc, err := response.SendGetRequestWithByteResult(url, retryNum, wx.accessTokenSupplierChan)
	if err != nil {
		
		return nil, err
		
	}
	
	var list UserListResult
	if err = json.Unmarshal(rc, &list); err != nil {
		
		return nil, err
		
	}
	
	return &list, nil
	
}

