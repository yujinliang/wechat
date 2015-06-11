package mp

import (
	
	"log"
	"io"
	"bytes"
	"errors"
	"encoding/json"
	"mp/response"
	
)

//永久图文素材
type PermanentNews struct {
	
	Title 		  		string `json:"title"`
	ThumbMediaId  		string `json:"thumb_media_id"`
	Author		  		string `json:"author,omitempty"`
	Digest		  		string `json:"digest,omitempty"`
	ShowCoverPic  		int8   `json:"show_cover_pic"`
	Content 	  		string `json:"content"`
	ContentSourceUrl 	string `json:"content_source_url,omitempty"`
	
}
//用于下载永久图文和视频素材
type PermanentNews_Video struct {
	
	News_Item []PermanentNews `json:"news_item,omitempty"`
	//for video
	Title	    string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	DownUrl		string	`json:"down_url,omitempty"`
	
}
//素材数
type MaterialsCount struct {
	
	VoiceCount uint32 `json:"voice_count"`
	VideoCount uint32 `json:"video_count"`
	ImageCount uint32 `json:"image_count"`
	NewsCount  uint32 `json:"news_count"`
	
}
//素材列表
type MaterialList struct {
	
	TotalCount uint `json:"total_count"`
	ItemCount  uint `json:"item_count"`
	Item	   []struct {
		
		MediaId    string `json:"media_id"`
		UpdateTime int64  `json:"update_time"`
		Name	   string `json:"name,omitempty"`
		Content	   struct {
			
			NewsItem []struct {
				
				Title 		   		string `json:"title,omitempty"`
				ThumbMediaId  		string `json:"thumb_media_id,omitempty"`
				ShowCoverPic 		int8   `json:"show_cover_pic"`
				Author				string `json:"author,omitempty"`
				Digest				string `json:"digest,omitempty"`
				Content				string `json:"content"`
				Url 				string `json:"url,omitempty"`
				ContentSourceUrl 	string `json:"content_source_url,omitempty"`
				
			} `json:"news_item,omitempty"`
			
		} `json:"content,omitempty"`
		
	} `json:"item,omitempty"`
	
}
//素村管理
//1上传临时素材， 3 days
//参数filename必须形如: xxx.jpg等， 扩展名必须为小写.
//os.Open打开文件时，扩展名必须为小写，否则找到不文件.
func (wx *WeiXin) UploadTmpMedia(mediaType string, filename string, reader io.Reader) (string, error) {
	
	url := "https://" + weixinFileURL + "/upload?type=" + mediaType + "&access_token="
	return response.UploadMedia(url, retryNum, filename, nil, wx.accessTokenSupplierChan, reader)
	
}
//2.获取临时素材
//os.Create中须指定形如: xxx.jpg的文件.
func (wx *WeiXin) DownloadTmpMedia(mediaId string, mediaType string, writer io.Writer) ([]byte, error) {
	
	scheme := "https://"
	if mediaType == MediaTypeVideo {
		
		scheme = "http://"
		
	}
	url := scheme + weixinFileURL + "/get?media_id=" + mediaId + "&access_token="
	jsonBytes, err := response.DownloadMedia(url, retryNum, wx.accessTokenSupplierChan, writer)
	
	return jsonBytes, err
	
}
//3. 获取永久素材
//注意： 图文， 视频， 这两种类型只会返回json串， 其它类型素材则以文件流形式直接写入io.Writer.
func (wx *WeiXin) DownloadPermanentMaterial(mediaId string, writer io.Writer) (*PermanentNews_Video ,error) {
	
	var msg struct {
		
		MediaId string `json:"media_id"`
		
	}
	msg.MediaId = mediaId
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		log.Println("DownloadPermanentMaterial marshal failed: ", err)
		return nil, err
		
	}
	
	url := weixinHost + "/material/get_material?access_token="
	jsonBytes, err := response.DownloadMediaByPost(url, retryNum, wx.accessTokenSupplierChan, writer, data)
	
	if err == nil && len(jsonBytes) > 0 {
		
		var result PermanentNews_Video
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
		
			return nil, err
		
		}
	
		return &result, nil
	
	}
	
	return nil, err
	
}
//4.新增永久图文素村
func (wx *WeiXin) UploadPermanentNews(news []PermanentNews) (string, error) {
	
	var newsMsg struct {
		
		Articles []PermanentNews `json:"articles"`
	}
	newsMsg.Articles = news
	
	data, err := json.Marshal(&newsMsg)
	if err != nil {
		
		log.Println("UploadPermanentNews marshal failed: ", err)
		return "", err
		
	}
	
	//just for debug.
	log.Println("UploadPermanentNews: ", string(data))
	
	rc, err := response.SendPostRequest(weixinHost + "/material/add_news?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return "", err
		
	}
	return rc.MediaId, nil
	
}
//5.新增永久视频素材
func (wx *WeiXin) UploadPermanentVideoMaterial(filename string,title string, introduction string, reader io.Reader) (string, error) {
	
	var desc struct {
		
		Title string `json:"title"`
		Introduction string `json:"introduction"`
	}
	desc.Title = title
	desc.Introduction = introduction
	
	descBytes, err := json.Marshal(&desc)
	if err != nil {
		
		return "", err
		
	}
	
	url := weixinHost + "/material/add_material?type=video&access_token="
	return response.UploadMedia(url, retryNum, filename, descBytes, wx.accessTokenSupplierChan, reader)
	
}
//6.上传非视频类永久素材
func (wx *WeiXin) UploadPermanentNoVideoMaterial(mediaType string, filename string, reader io.Reader) (string, error) {
	
	url := weixinHost + "/material/add_material?type=" + mediaType + "&access_token="
	return response.UploadMedia(url, retryNum, filename, nil, wx.accessTokenSupplierChan, reader)
	
}
//7.删除永久素材
func (wx *WeiXin) DeletePermanentMaterial(permanent_mediaId string) error {
	
	var deletePermanentMaterialMsg struct {
		
		MediaId string `json:"media_id"`
		
	}
	deletePermanentMaterialMsg.MediaId = permanent_mediaId
	
	data, err := json.Marshal(&deletePermanentMaterialMsg)
	if err != nil {
		
		log.Println("DeletePermanentMaterial marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("DeletePermanentMaterial: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/material/del_material?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return err
		
	}
	
	return nil
	
}
//7.修改永久图文素材
func (wx *WeiXin) ModifyPermanentNews(mediaId string, index int, articles *PermanentNews) error {
	
	var msg struct {
		
		MediaId  string `json:"media_id"`
		Index	 int	 `json:"index"`
		Articles *PermanentNews `json:"articles,omitempty"`
		
	}
	msg.MediaId  = mediaId
	msg.Index    = index
	msg.Articles = articles
	
	data, err := json.Marshal(&msg)
	if err != nil {
		
		log.Println("ModifyPermanentNews marshal failed: ", err)
		return err
		
	}
	
	//just for debug.
	log.Println("ModifyPermanentNews: ", string(data))
	
	_, err = response.SendPostRequest(weixinHost + "/material/update_news?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return err
		
	}
	
	return nil
	
}
//素材总数
func (wx *WeiXin) GetMaterialCount() (*MaterialsCount, error) {
	
	url := weixinHost + "/material/get_materialcount?access_token="
	rc, err := response.SendGetRequestWithByteResult(url, retryNum, wx.accessTokenSupplierChan)
	if err != nil {
		
		return nil, err
		
	}
	
	if rc != nil {
		
		if bytes.Contains(rc, []byte("_count")) {
			
			var c MaterialsCount
			if err := json.Unmarshal(rc, &c); err != nil {
			
				return nil, err
			
			} else {
			
				return &c, nil
			
			}
		
		} 

	}
	
	return nil, errors.New("getMaterialCount Unknow Error!")
}
//获取素材列表
func (wx *WeiXin) GetMaterialsList(mediaType string, offset , count int) (*MaterialList, error) {
	
	var msg struct {
		
		Type   string `json:"type"`
		Offset int	  `json:"offset"` 
		Count  int	  `json:"count"`
	}
	msg.Type   = mediaType
	msg.Offset = offset
	msg.Count  = count
	
		data, err := json.Marshal(&msg)
	if err != nil {
		
		return nil, err
		
	}
	
	//just for debug.
	log.Println("GetMaterialsList: ", string(data))
	
	rc, err := response.SendPostRequestWithByteResult(weixinHost + "/material/batchget_material?access_token=", retryNum, wx.accessTokenSupplierChan, data)
	if err != nil {
		
		return nil, err
		
	}
	
	var result MaterialList
	if err := json.Unmarshal(rc, &result); err != nil {
		
		return nil, err
		
	}
	
	return &result, nil
	
}