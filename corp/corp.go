package corp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Err 微信返回错误
type Err struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// AccessToken 微信企业号请求Token
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Err
	ExpiresInTime time.Time
}

// Client 微信企业号应用配置信息
type Client struct {
	CorpID      string
	AgentID     int
	AgentSecret string
	Token       AccessToken
}

// Result 发送消息返回结果
type Result struct {
	Err
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"infvalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

// Content 文本消息内容
type Content struct {
	Content string `json:"content"`
}

// Message 消息主体参数
type Message struct {
	ToUser  string  `json:"touser"`
	MsgType string  `json:"msgtype"`
	Text    Content `json:"text"`
}

// Send 发送信息
func Send(msg Message) error {
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + msg.ToUser

	resultByte, err := jsonPost(url, msg)
	if err != nil {
		return fmt.Errorf("invoke send api fail: %v", err)
	}

	result := Result{}
	err = json.Unmarshal(resultByte, &result)
	if err != nil {
		return fmt.Errorf("parse send api response fail: %v", err)
	}

	if result.ErrCode != 0 {
		err = fmt.Errorf("invoke send api return ErrCode = %d", result.ErrCode)
	}

	if result.InvalidUser != "" || result.InvalidParty != "" || result.InvalidTag != "" {
		err = fmt.Errorf("invoke send api partial fail, invalid user: %s, invalid party: %s, invalid tag: %s", result.InvalidUser, result.InvalidParty, result.InvalidTag)
	}

	return err
}

func jsonPost(url string, data interface{}) ([]byte, error) {
	jsonBody, err := encodeJSON(data)
	if err != nil {
		return nil, err
	}

	r, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	if r.Body == nil {
		return nil, fmt.Errorf("response body of %s is nil", url)
	}

	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}

func encodeJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
