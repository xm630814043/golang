package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ReqBody struct {
	UUID  string `json:"uuid"`
	Event struct {
		AppID            string `json:"app_id"`
		ChatType         string `json:"chat_type"`
		IsMention        bool   `json:"is_mention"`
		LarkVersion      string `json:"lark_version"`
		MessageID        string `json:"message_id"`
		MsgType          string `json:"msg_type"`
		OpenChatID       string `json:"open_chat_id"`
		OpenID           string `json:"open_id"`
		OpenMessageID    string `json:"open_message_id"`
		ParentID         string `json:"parent_id"`
		RootID           string `json:"root_id"`
		TenantKey        string `json:"tenant_key"`
		Text             string `json:"text"`
		TextWithoutAtBot string `json:"text_without_at_bot"`
		Type             string `json:"type"`
		UserAgent        string `json:"user_agent"`
		UserOpenID       string `json:"user_open_id"`
	} `json:"event"`
	Token     string `json:"token"`
	Ts        string `json:"ts"`
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
}

type Challenge struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}
type Challenge1 struct {
	Challenge string `json:"challenge"`
}

type Tenant_access struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

type Send_meg struct {
	OpenID  string `json:"open_id"`
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type Access_token struct {
	App_id     string `json:"app_id"`
	App_secret string `json:"app_secret"`
}
type Send_callback struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		MessageID string `json:"message_id"`
	} `json:"data"`
}

// handle_request_url_verify 原样返回 challenge 字段内容
func handle_request_url_verify(post_obj_byte []byte) []byte {
	var post_obj Challenge
	var challenge Challenge1
	err = json.Unmarshal(post_obj_byte, &post_obj)
	challenge.Challenge = post_obj.Challenge
	rsp, _ := json.Marshal(challenge)
	return rsp
}

func get_tenant_access_token() string {
	var jsonData Access_token
	var bodyContent Tenant_access
	url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"

	jsonData.App_id = APP_ID
	jsonData.App_secret = APP_SECRET
	jsonStr, err := json.Marshal(jsonData)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &bodyContent)
	if bodyContent.Code != 0 {
		fmt.Printf("get tenant_access_token error, code =%v", bodyContent.Code)
	}
	return bodyContent.TenantAccessToken
}

func send_message(token string, open_id string, text string) {
	url := "http://open.feishu.cn/open-apis/message/v4/send/"
	var send_meg Send_meg
	var bodyContent Send_callback
	send_meg.OpenID = open_id
	send_meg.MsgType = "text"
	send_meg.Content.Text = text
	jsonStr, err := json.Marshal(send_meg)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &bodyContent)
	if bodyContent.Code != 0 {
		fmt.Printf("get tenant_access_token error, code =%v, msg =%v", bodyContent.Code, bodyContent.Msg)
	}
}

// handle_message 此处只处理 text 类型消息，其他类型消息忽略
func handle_message(Content ReqBody) {
	//var w http.ResponseWriter
	if Content.Event.Type != "text" {
		fmt.Printf("unknown msg_type =%v\n", Content.Event.Type)
		//return
	}
	// 调用发消息 API 之前，先要获取 API 调用凭证：tenant_access_token
	access_token := get_tenant_access_token()
	if access_token == "" {
		return
	}
	// 机器人 echo 收到的消息
	send_message(access_token, Content.Event.OpenID, Content.Event.Text)
	//response(w,[]byte(""))
	return
}
