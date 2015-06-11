package mp

import (
	
	"io"
	"fmt"
	"time"
	"bytes"
	"errors"
	"sort"
	"strings"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/base64"
	"encoding/xml"
	"github.com/yujinliang/wechat/mp/request"
	
)

const (
	
	// request message types
	MsgTypeText     = "text"
	MsgTypeImage    = "image"
	MsgTypeVoice    = "voice"
	MsgTypeShortVideo = "shortvideo"
	MsgTypeVideo    = "video"
	MsgTypeLocation = "location"
	MsgTypeLink     = "link"
	MsgTypeEvent    = "event"
	// event types
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventScan        = "SCAN"
	EventLocation    = "LOCATION"
	EventClick       = "CLICK"
	EventView        = "VIEW"
	//回复微信基础消息模板
	replyText               = "<xml>%s<MsgType><![CDATA[text]]></MsgType><Content><![CDATA[%s]]></Content></xml>"
	replyImage              = "<xml>%s<MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[%s]]></MediaId></Image></xml>"
	replyVoice              = "<xml>%s<MsgType><![CDATA[voice]]></MsgType><Voice><MediaId><![CDATA[%s]]></MediaId></Voice></xml>"
	replyVideo              = "<xml>%s<MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[%s]]></MediaId><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description></Video></xml>"
	replyMusic              = "<xml>%s<MsgType><![CDATA[music]]></MsgType><Music><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description><MusicUrl><![CDATA[%s]]></MusicUrl><HQMusicUrl><![CDATA[%s]]></HQMusicUrl><ThumbMediaId><![CDATA[%s]]></ThumbMediaId></Music></xml>"
	replyNews               = "<xml>%s<MsgType><![CDATA[news]]></MsgType><ArticleCount>%d</ArticleCount><Articles>%s</Articles></xml>"
	replyHeader             = "<ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime>"
	replyArticle            = "<item><Title><![CDATA[%s]]></Title> <Description><![CDATA[%s]]></Description><PicUrl><![CDATA[%s]]></PicUrl><Url><![CDATA[%s]]></Url></item>"

)

type CDATAText struct {
	
	Text string `xml:",innerxml"`
	
}
type EncryptResponseBody struct {
	
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATAText
	MsgSignature CDATAText
	TimeStamp    string
	Nonce        CDATAText
	
}
// Use to reply music message
type Music struct {
	
	Title        string `json:"title"`
	Description  string `json:"description"`
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
	
}

// Use to reply news message
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PicUrl      string `json:"picurl"`
	Url         string `json:"url"`
}

// Format reply message header
func (wx *WeiXin) replyHeader(originMsg *request.WeiXinRequest) string {
	
	return fmt.Sprintf(replyHeader, originMsg.FromUserName, originMsg.ToUserName, time.Now().Unix())
	
}
//to reply normal weixin msg start
func (wx *WeiXin) ReplyText(text string, originMsg *request.WeiXinRequest) string {
	
	return fmt.Sprintf(replyText, wx.replyHeader(originMsg),text)

}
func (wx *WeiXin) ReplyImage(mediaId string, originMsg *request.WeiXinRequest) string {
	
	return fmt.Sprintf(replyImage, wx.replyHeader(originMsg), mediaId)
	
}
func (wx *WeiXin) ReplyVoice(mediaId string, originMsg *request.WeiXinRequest) string {
	
	return fmt.Sprintf(replyVoice, wx.replyHeader(originMsg), mediaId)
	
}
func (wx *WeiXin) ReplyVideo(mediaId string, title string, description string, originMsg *request.WeiXinRequest) string {
	
	return fmt.Sprintf(replyVideo, wx.replyHeader(originMsg), mediaId, title, description)
	
}
func (wx *WeiXin) ReplyMusic(m *Music, originMsg *request.WeiXinRequest) string {
	
	return fmt.Sprintf(replyMusic, wx.replyHeader(originMsg), m.Title, m.Description, m.MusicUrl, m.HQMusicUrl, m.ThumbMediaId)
	
}
func (wx *WeiXin) ReplyNews(articles []Article, originMsg *request.WeiXinRequest) string {
	
	var articleItems string
	for _, item := range articles {
		
		articleItems += fmt.Sprintf(replyArticle,item.Title, item.Description, item.PicUrl, item.Url)
		
	}
	return fmt.Sprintf(replyNews, wx.replyHeader(originMsg), len(articles), articleItems)
	
}
func value2CDATA(v string) CDATAText {
	
	return CDATAText{"<![CDATA[" + v + "]]>"}
	
}
func (wx *WeiXin) MakeEncryptResponse(replyMsg []byte, timestamp, nonce string)([]byte, error) {
	
	encryptResponse := &EncryptResponseBody{}
	
	encryptXMLData, err := wx.makeEncryptReplyXMLData(replyMsg)
	if err != nil {
		
		return nil, err
		
	}
	encryptResponse.Encrypt = value2CDATA(encryptXMLData)
	encryptResponse.MsgSignature = value2CDATA(makeMsgSignature(wx.token, timestamp, nonce, encryptXMLData))
	encryptResponse.TimeStamp = timestamp
	encryptResponse.Nonce = value2CDATA(nonce)
	
	return xml.MarshalIndent(encryptResponse, " ", " ")
	
}
//注意： 参数：replyMsg 为上方Reply*函数的返回值
func (wx *WeiXin) makeEncryptReplyXMLData(replyMsg []byte) (string, error) {
	
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(replyMsg)))
	if err != nil {
		
		return "", err
		
	}
	
	msgLenght := buf.Bytes()
	randomBytes := []byte("wbcdvfghujklmnzp")
	plainData := bytes.Join([][]byte{randomBytes, msgLenght, replyMsg, []byte(wx.appId)}, nil)
	
	var aesKey []byte
	if wx.useCurrentAESKey {
		
		aesKey = wx.currentAesKey
		
	} else if len(wx.lastAesKey) > 0 {
		
		aesKey = wx.lastAesKey
		
	} else {
		
		return "", errors.New("No availabal encodingAESKey")
		
	}
	
	cipherData, err := aesEncrypt(plainData, aesKey)
	if err != nil {
		
		return "", err
		
	}
	
	return base64.StdEncoding.EncodeToString(cipherData), nil
	
}

func aesEncrypt(plainData []byte, aesKey []byte) ([]byte, error) {
	
	k := len(aesKey)
	if len(plainData)%k != 0 {
		
		plainData = PKCS7Pad(plainData, k)
		
	}
	fmt.Printf("aesEncrypt: after padding, plainData length = %d\n", len(plainData))

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		
		return nil, err
		
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		
		return nil, err
		
	}

	cipherData := make([]byte, len(plainData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, plainData)

	return cipherData, nil
	
}

func aesDecrypt(cipherData []byte, aesKey []byte) ([]byte, error) {
	
	k := len(aesKey) //PKCS#7
	if len(cipherData)%k != 0 {
		
		return nil, errors.New("crypto/cipher: ciphertext size is not multiple of aes key length")
		
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		
		return nil, err
		
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		
		return nil, err
		
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(plainData, cipherData)
	return plainData, nil
	
}

//from github.com/vgorin/cryptogo
func PKCS7Pad(message []byte, blocksize int) (padded []byte) {
	
	// block size must be bigger or equal 2
	if blocksize < 1<<1 {
		
		//panic("block size is too small (minimum is 2 bytes)")
		return nil
		
	}
	// block size up to 255 requires 1 byte padding
	if blocksize < 1<<8 {
		
		// calculate padding length
		padlen := PadLength(len(message), blocksize)

		// define PKCS7 padding block
		padding := bytes.Repeat([]byte{byte(padlen)}, padlen)

		// apply padding
		padded = append(message, padding...)
		return padded
		
	}
	// block size bigger or equal 256 is not currently supported
	//panic("unsupported block size")
	return nil
	
}
// PadLength calculates padding length, from github.com/vgorin/cryptogo
func PadLength(slice_length, blocksize int) (padlen int) {
	
	padlen = blocksize - slice_length%blocksize
	if padlen == 0 {
		
		padlen = blocksize
		
	}
	return padlen
	
}

func makeMsgSignature(token, timestamp, nonce, msg_encrypt string) string {
	
	sl := []string{token, timestamp, nonce, msg_encrypt}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
	
}
