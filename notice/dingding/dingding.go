package dingding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DingTalkMessage 定义发送给钉钉的消息结构
type DingTalkMessage struct {
	MsgType string `json:"msgtype"` // 消息类型，这里使用text
	Text    struct {
		Content string `json:"content"` // 消息内容
	} `json:"text"`
	// 如果需要发送Markdown或其他类型消息，可以添加相应字段
}

// SendDingTalkMessage 发送消息到钉钉
func SendDingTalkMessage(webhookURL string, message DingTalkMessage) error {
	// 将消息对象序列化为JSON
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 打印响应体（可选）
	fmt.Println(string(responseBody))

	return nil
}
