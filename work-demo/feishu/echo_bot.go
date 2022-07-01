package feishu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	err                    error
	Port                   = ":8086"
	APP_ID                 = "cli_a123714348f9100c"
	APP_SECRET             = "trU95Sgfxgn9VwpioQvHghhuZPkEn22G"
	APP_VERIFICATION_TOKEN = "WWj6ni3eEs3XrtpYbiNZybjkBOF0EHcA"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var bodyContent ReqBody
	// 解析请求 body
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	} else {
		err = json.Unmarshal(bodyData, &bodyContent)
		// 校验 verification token 是否匹配，token 不匹配说明该回调并非来自开发平台
		if bodyContent.Token != APP_VERIFICATION_TOKEN {
			fmt.Printf("verification token not match, token =%v", bodyContent.Token)
			return
		}
		// 根据 type 处理不同类型事件
		if bodyContent.Type == "url_verification" {
			rsp := handle_request_url_verify(bodyData)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = w.Write(rsp)
			// 获取事件内容和类型，并进行相应处理，此处只关注给机器人推送的消息事件
		} else if bodyContent.Type == "event_callback" {
			if bodyContent.Event.Type == "message" {
				handle_message(bodyContent)
			}
		}
		fmt.Printf(string(bodyData))
	}
}
