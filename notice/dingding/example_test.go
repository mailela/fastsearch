package dingding

import (
	"fmt"
	"testing"
)

func Test_Notice(t *testing.T) {
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=YOUR_ACCESS_TOKEN"
	message := DingTalkMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: "Hello, this is a message from Go program!",
		},
	}

	err := SendDingTalkMessage(webhookURL, message)
	if err != nil {
		fmt.Println("Error sending message to DingTalk:", err)
	}
}
