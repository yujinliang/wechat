package request

import (
	
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"io"
	"bytes"
	"net/http"
	"sort"
	"encoding/binary"
	"strings"
	"errors"
	
)

// Common weixin request header
type RequestCommonFields struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
}

// Weixin request
type WeiXinRequest struct {
	RequestCommonFields
	MsgId int64 `xml:"MsgId" json:"MsgId"`
	MsgID int64 `xml:"MsgID" json:"MsgID"`
	Content      string
	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string
	LocationX    float32 `xml:"Location_X"`
	LocationY    float32 `xml:"Location_Y"`
	Scale        float32
	Label        string
	Title        string
	Description  string
	Url          string
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float32
	Longitude    float32
	Precision    float32
	Recognition  string
	Status	      string
	TotalCount   int     `xml:"TotalCount"  json:"TotalCount"`
	FilterCount  int     `xml:"FilterCount" json:"FilterCount"`
	SentCount    int     `xml:"SentCount"   json:"SentCount"`
	ErrorCount   int     `xml:"ErrorCount"  json:"ErrorCount"`
	KfAccount	  string
	FromKfAccount string
	ToKfAccount	  string
	ScanCodeInfo struct {
		
		ScanType   string `xml:"ScanType"`
		ScanResult string `xml:"ScanResult"`
		
	} `xml:"ScanCodeInfo"`
	SendLocationInfo struct {
		
		LocationX float64 `xml:"Location_X"`
		LocationY float64 `xml:"Location_Y"`
		Scale     int     `xml:"Scale"`
		Label     string  `xml:"Label"`
		Poiname   string  `xml:"Poiname"`
		
	} `xml:"SendLocationInfo"`
	
}

func (r *WeiXinRequest)	 unpackRequest(req *http.Request) error {
	
	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		
		return err
		
	}
	defer req.Body.Close()
	
	if err := xml.Unmarshal(raw, r); err != nil {
		
		return err
		
	}
	
	return nil
	
}
func (r *WeiXinRequest)	 unpackRequestForEncrypted(plainData []byte, appId string) error {
	
	//read length.
	buf := bytes.NewBuffer(plainData[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	
	//appId validation.
	appIDStart := 20 + length
	id := plainData[appIDStart : int(appIDStart) + len(appId)]
	if len(id) <= 0 || string(id) != appId {
		
		return errors.New("appId mismatch! " + string(id) + "correct: " + appId)
		
	}
	if err := xml.Unmarshal([]byte(plainData[20 : (20 + length)]), r); err != nil {
		
		return err
		
	}
	
	return nil
	
}
func (r * WeiXinRequest) checkSignature(token string, req *http.Request) bool {
	
	ss := sort.StringSlice{
							token, 
							req.FormValue("timestamp"), 
							req.FormValue("nonce"),
							
							}
							
	sort.Strings(ss)
	s := strings.Join(ss, "")
	h := sha1.New()
    io.WriteString(h, s)
	
	return fmt.Sprintf("%x", h.Sum(nil)) == req.FormValue("signature")
	
}
func (r * WeiXinRequest) checkSignatureForEncrypted(token string, timestamp string, nonce string, originSignature string) bool {
	
	ss := sort.StringSlice{
							token, 
							timestamp, 
							nonce,
							
							}
							
	sort.Strings(ss)
	s := strings.Join(ss, "")
	h := sha1.New()
    io.WriteString(h, s)
	
	return fmt.Sprintf("%x", h.Sum(nil)) == originSignature
	
}
func (r *WeiXinRequest) TryUnpackWeiXinRequestForEncrypted(appId,token string, w http.ResponseWriter,plainData []byte, timestamp, nonce, originSignature string) bool {
	
	if !r.checkSignatureForEncrypted(token, timestamp, nonce, originSignature) {
		
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return false
		
	}
	
	if err := r.unpackRequestForEncrypted(plainData, appId); err != nil {
		
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return false
		
	}
	
	return true
	
}
func (r *WeiXinRequest) TryUnpackWeiXinRequest(token string, w http.ResponseWriter, req *http.Request) bool {
	
	if req.Method != "POST" {
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(req.FormValue("echostr")))
		return false
		
	}
	
	if !r.checkSignature(token, req) {
		
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return false
		
	}
	
	if err := r.unpackRequest(req); err != nil {
		
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return false
		
	}
	
	return true
	
}


