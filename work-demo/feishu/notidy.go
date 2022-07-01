// 飞书机器人消息推送 参考文档：https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN
package feishu

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

type FSResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type TextItem struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type NotifyData struct {
	MsgType string `json:"msg_type"`
	Content `json:"content"`
}

type Content struct {
	Post PostData `json:"post"`
}

type PostData struct {
	ZhCn zhCn `json:"zh_cn"`
}

type zhCn struct {
	Title   string       `json:"title"`
	Content [][]TextItem `json:"content"`
}

const WebhookUrl = "https://open.feishu.cn/open-apis/bot/v2/hook/de809d97-575a-42fc-a498-05b89e501f4a" // 飞书报警地址

// @Title FSErrorNotify
// @Description  飞书报警
// @Author zhaopengtao 2022-06-24 09:38:41
// @Param webhookUrl
// @Param title
// @Param errParam
// @Param message
// @Return error

func FSErrorNotify(webhookUrl, title string, errParam error, message interface{}) error {
	// 参数
	msgStr, _ := json.Marshal(message)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	var lineStr = strconv.Itoa(line)
	var errMsg = "正常输出"
	if errParam != nil {
		errMsg = errParam.Error()
	}
	var msgItem = []TextItem{
		{
			Tag:  "text",
			Text: "环境：" + os.Getenv("RUN_TIME") + "\n",
		},
		{
			Tag:  "text",
			Text: "位置：" + file + ":" + lineStr + "\n",
		},
		{
			Tag:  "text",
			Text: "函数：" + f.Name() + "\n",
		},
		{
			Tag:  "text",
			Text: "错误：" + errMsg + "\n",
		},
		{
			Tag:  "text",
			Text: "调用链：" + string(debug.Stack()) + "\n",
		},
		{
			Tag:  "text",
			Text: "参数：" + string(msgStr) + "\n",
		},
	}

	var msg = NotifyData{
		MsgType: "post",
		Content: Content{PostData{ZhCn: zhCn{
			Title: title,
			Content: [][]TextItem{
				msgItem,
			},
		}}},
	}

	params, err := json.Marshal(msg)

	req, err := http.NewRequest("POST", webhookUrl, strings.NewReader(string(params)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	result, err := new(http.Client).Do(req)
	if err != nil {
		return err
	}

	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}

	var response = FSResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	if response.Code != 0 {
		return errors.New(fmt.Sprintf("GeoNotifyRobot StatusCode != 0; StatusCode ==%v ", response.Code))
	}

	//return errors.New("111")
	return errors.New("111")
}
