package response


import (
	
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"mime/multipart"
	"time"
	"log"
	"github.com/yujinliang/wechat/mp/accesstoken"
	
)

// response from weixinmp when we call weixinmap`api to send or get something.
type Response struct {
	// error fields
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// token fields
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	// media fields
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
	// ticket fields
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
	QrCodeUrl     string `json:"url"`
	//template message
	TemplateId   string `json:"template_id"`
	MsgId		 string `json:"msgid"`
	MassMsgId	 string	`json:"msg_id"`
	MsgStatus	 string `json:"msg_status"`
	//long url to short url
	ShortURL string `json:"short_url"`
	
}
//适应用正常消息， 与出错消息， 共用一个消息格式时， 才可调用此接口，记得在response.Response加入返回的正常消息字段！！！
func SendPostRequest(url string, retryMaxN int, c chan accesstoken.AccessToken, data []byte) ( rc *Response, err error) {
	
	for i := 0; i < retryMaxN; i++ {
		
		token := <- c
		if time.Since(token.Expires).Seconds() < 0 {
			
			rc, err = post(url + token.Token, "application/json; charset=utf-8", bytes.NewBuffer(data))
			if err != nil {
				
				log.Println("SendPostRequest: ", err)
				return
				
			} else {
				
				log.Println("SendPostRequest: ",rc.ErrCode, rc.ErrMsg)
				switch rc.ErrCode {
					
					case 0: {
						
						return rc, nil
						
					}
					case 42001: { 
						
						continue
						
					}
					default: {
						
						return nil, errors.New(fmt.Sprintf("Weixin send post request reply[%d]: %s", rc.ErrCode, rc.ErrMsg))
						
					}
				}
			}
			
		}
	}
	
	return nil, errors.New("Weixin post request too many times: " + url)
	
}
func post(url string, bodyType string, body *bytes.Buffer) (*Response, error) {
	
	resp , err := http.Post(url, bodyType, body)
	if err != nil {
		
		return nil, err
		
	}
	defer resp.Body.Close()
	
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		
		return nil, err
		
	}
	
	var rc Response
	if err := json.Unmarshal(data, &rc); err != nil {
		
		return nil, err
		
	}
	
	return &rc, nil
	
}
//注意：正常情况消息，与出错情况（即errCode消息格式，）两种消息不共用同一个消息格式，在此种情况下需要用WithByteResult结尾的接口，因为要分别解析返回结果。
//然而对于正常情况与出错情况共用一个消息格式，只是靠errCode == 0来区分是否出错的时候， 则以下WithByteResult结尾的接口不再适用了！因为此接口
//只要发现消息体中有errCode就认为出错了！！！只返回空和错误！！！ 
func SendPostRequestWithByteResult(url string, retryMaxN int, c chan accesstoken.AccessToken, data []byte) ([]byte, error) {
	
		for i := 0; i < retryMaxN; i++ {
		
		token := <- c
		if time.Since(token.Expires).Seconds() < 0 {
			
			rc , err := http.Post(url + token.Token, "application/json; charset=utf-8", bytes.NewBuffer(data))
			if err != nil {
				
				return nil, err
				
			}
			defer rc.Body.Close()
			
			reply, err := ioutil.ReadAll(rc.Body)
			if err != nil {
				
				return nil, err
				
			}
			
			if reply != nil {
		
				if bytes.Contains(reply, []byte("errcode")) {
				
				//handle error.
				var err_rc Response
				if err := json.Unmarshal(reply, &err_rc); err != nil {
				
					return nil, err
				
				} else {
				
						switch err_rc.ErrCode {
						
							case 42001: {
							
								continue
							
							}
							default: {	
						
								return nil, errors.New(fmt.Sprintf("error[%d]: %s", err_rc.ErrCode, err_rc.ErrMsg))
							
							}
					
						}
				
					}
				
				}
				
				return reply, nil
					
			}
			
		}
	
	}
		
	return nil, errors.New("WeiXin Post Request too many times: " + url)
	
}
func SendGetRequestWithByteResult(url string, tryMaxN int, c chan accesstoken.AccessToken) ([]byte, error) {
	
	for i := 0; i < tryMaxN; i++ {
		
		token := <- c
		if time.Since(token.Expires).Seconds() < 0 {
			
			rc, err := http.Get(url + token.Token)
			if err != nil {
				
				return nil, err
				
			}
			defer rc.Body.Close()
			
			reply, err := ioutil.ReadAll(rc.Body)
			if err != nil {
				
				return nil, err
				
			}
			
			if reply != nil {
		
				if bytes.Contains(reply, []byte("errcode")) {
				
				//handle error.
				var err_rc Response
				if err := json.Unmarshal(reply, &err_rc); err != nil {
				
					return nil, err
				
				} else {
				
						switch err_rc.ErrCode {
						
							case 42001: {
							
								continue
							
							}
							default: {	
						
								return nil, errors.New(fmt.Sprintf("error[%d]: %s", err_rc.ErrCode, err_rc.ErrMsg))
							
							}
					
						}
				
					}
				
				}
				
				return reply, nil
					
			}
			
		}
	
	}
		
	return nil, errors.New("WeiXin Get Request too many times: " + url)
	
}
func SendGetRequest(url string, retryMaxN int, c chan accesstoken.AccessToken) (rc *Response, err error) {
	
	for i := 0; i < retryMaxN; i++ {
		
		token := <- c
		if time.Since(token.Expires).Seconds() < 0 {
			
			rc, err = get(url + token.Token)
			if err != nil {
				
				return
				
			} else {
				
				switch rc.ErrCode {
					
					case 0: {
						
						return rc, nil
						
					}
					case 42001: { 
						
						continue
						
					}
					default: {
						
						return nil, errors.New(fmt.Sprintf("Weixin send get  request reply[%d]: %s", rc.ErrCode, rc.ErrMsg))
						
					}
				}
			}
		}
	}
	
	return nil, errors.New("WeiXin send get reuqest too many times: " + url)
	
}
func get(url string) (*Response, error) {
	
	resp ,err := http.Get(url)
	if err != nil {
		
		return nil, err
		
	}
	defer resp.Body.Close()
	
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		
		return nil, err
		
	}
	
	var rc Response
	if err := json.Unmarshal(data, &rc); err != nil {
		
		return nil, err
		
	}

	return &rc, nil
	
}
func UploadMedia(url string, retryMaxN int, filename string, description []byte, c chan accesstoken.AccessToken, reader io.Reader) (string, error) {
	
	for i := 0;  i < retryMaxN; i++ {
		
		token := <- c
		if time.Since(token.Expires).Seconds() < 0 {
			
			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			fileWriter, err := bodyWriter.CreateFormFile("media", filename)
			if err != nil {
				
				return "", err
				
			}

			if _, err = io.Copy(fileWriter, reader); err != nil {
				
				return "", err
				
			}
			
			//for description.
			if description != nil && len(description) > 0 {
				
				descWriter, descErr := bodyWriter.CreateFormField("description")
				if descErr != nil {
				
					return "", descErr
				
				}
				if _, descErr = io.Copy(descWriter,bytes.NewReader(description)); descErr != nil {
				
					return "", descErr
				
				}
			
			}
			
			if err = bodyWriter.Close(); err != nil {
				
				return "", err
				
			}
			
			contentType := bodyWriter.FormDataContentType()
			rc, err := http.Post(url+token.Token, contentType, bodyBuf)
			if err != nil {
				
				return "", err
				
			}
			defer rc.Body.Close()
			reply, err := ioutil.ReadAll(rc.Body)
			if err != nil {
				
				return "", err
				
			}
			
			var result Response
			err = json.Unmarshal(reply, &result)
			if err != nil {
				
				return "", err
				
			} else {
				
				switch result.ErrCode {
					
					case 0: {
						
						return result.MediaId, nil
						
					}
					case 42001: {
						
						continue
						
					}
					default: {
						
						return "", errors.New(fmt.Sprintf("WeiXin upload[%d]: %s", result.ErrCode, result.ErrMsg))
					}
				}
				
			}
				
		}
		
	}
	
	return "", errors.New(fmt.Sprintf("Weixin upload media too many times: %s", url))
}
func DownloadMedia(url string, retryMaxN int, c chan accesstoken.AccessToken,writer io.Writer) ([]byte, error) {
	
	for i := 0; i < retryMaxN; i++ {
		
		token := <- c
		if time.Since(token.Expires).Seconds() < 0 {
			
			rc, err := http.Get(url + token.Token)
			if err != nil {
				
				return nil, err
				
			}
			defer rc.Body.Close()
			
			var respBegin [15]byte // {"errcode": or {"errmsg":"

			n, err := io.ReadFull(rc.Body, respBegin[:])
			switch {
				
				case err == nil: {
				
					break
				
				}
				case err == io.ErrUnexpectedEOF: {
					
					_, err = writer.Write(respBegin[:n])
					return nil, err
			
				}
				case err == io.EOF: {
					
					return nil, err //??????
					
				}
				default: {
					
					return nil, err
					
				}
				
			}
			
			gatherRespBody := io.MultiReader(bytes.NewReader(respBegin[:]), rc.Body)
			
			//if rc.Header.Get("Content-Type") != "text/plain" {
			if !bytes.Contains(respBegin[:], []byte("errcode")) {
				
				//如是为图文， 或视频， 则直接返回[]byte做进一步的处理， 
				if bytes.Contains(respBegin[:], []byte("news_item")) || bytes.Contains(respBegin[:], []byte("down_url")){
					
					buf := new(bytes.Buffer)
					buf.ReadFrom(gatherRespBody)
					return buf.Bytes(), nil
					
				} 
				//如果为文件流， 则直接存入writer.
				_, err := io.Copy(writer, gatherRespBody)
				return nil, err
				
			} else {
				
				reply, err := ioutil.ReadAll(gatherRespBody)
				if err != nil {
					
					return nil, err
					
				} 
				
				var result Response
				if err := json.Unmarshal(reply, &result); err != nil {
					
					return nil, err
					
				} else {
					
					switch result.ErrCode {
						
						case 0: {
						
							return nil, nil
						
						}
						case 42001: {// access_token timeout and retry
					
							continue
						
						}
						default: {
						
							return nil, errors.New(fmt.Sprintf("WeiXin download[%d]: %s", result.ErrCode, result.ErrMsg))
						
						}
					
					}
					
				}
				
			}
			
		}
		
	}
	
	return nil, errors.New("Weixin DownloadMedia too many times")
	
}
//下载永久素材
func DownloadMediaByPost(url string, retryMaxN int, c chan accesstoken.AccessToken,writer io.Writer, data []byte) ([]byte, error) {
	
	for i := 0; i < retryMaxN; i++ {
		
		token := <- c
		if time.Since(token.Expires).Seconds() < 0 {
			
			rc , err := http.Post(url, "application/json; charset=utf-8", bytes.NewReader(data))
			if err != nil {
				
				return nil, err
				
			}
			defer rc.Body.Close()
			
			var respBegin [15]byte // {"errcode": or {"errmsg":"

			n, err := io.ReadFull(rc.Body, respBegin[:])
			switch {
				
				case err == nil: {
				
					break
				
				}
				case err == io.ErrUnexpectedEOF: {
					
					_, err = writer.Write(respBegin[:n])
					return nil, err
			
				}
				case err == io.EOF: {
					
					return nil, err //??????
					
				}
				default: {
					
					return nil, err
					
				}
				
			}
			
			gatherRespBody := io.MultiReader(bytes.NewReader(respBegin[:]), rc.Body)
			
			//if rc.Header.Get("Content-Type") != "text/plain" {
			if !bytes.Contains(respBegin[:], []byte("errcode")) {
				
				//如是为图文， 或视频， 则直接返回[]byte做进一步的处理， 
				if bytes.Contains(respBegin[:], []byte("news_item")) || bytes.Contains(respBegin[:], []byte("down_url")){
					
					buf := new(bytes.Buffer)
					buf.ReadFrom(gatherRespBody)
					return buf.Bytes(), nil
					
				} 
				//如果为文件流， 则直接存入writer.
				_, err := io.Copy(writer, gatherRespBody)
				return nil, err
				
			} else {
				
				reply, err := ioutil.ReadAll(gatherRespBody)
				if err != nil {
					
					return nil, err
					
				} 
				
				var result Response
				if err := json.Unmarshal(reply, &result); err != nil {
					
					return nil, err
					
				} else {
					
					switch result.ErrCode {
						
						case 0: {
						
							return nil, nil
						
						}
						case 42001: {// access_token timeout and retry
					
							continue
						
						}
						default: {
						
							return nil, errors.New(fmt.Sprintf("WeiXin download[%d]: %s", result.ErrCode, result.ErrMsg))
						
						}
					
					}
					
				}
				
			}
			
		}
		
	}
	
	return nil, errors.New("Weixin DownloadMedia too many times")
	
}
