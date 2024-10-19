package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pigeon/config"
)

// 定义目标 JSON 结构体
type Content struct {
	Post Post `json:"post"`
}

type Post struct {
	ZhCn ZhCn `json:"zh_cn"`
}

type ZhCn struct {
	Title   string          `json:"title"`
	Content [][]ContentItem `json:"content"`
}

type ContentItem struct {
	Tag    string `json:"tag"`
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

type FeishuMessage struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

// SendMessageFeiShu 发送消息到飞书
func SendMessageFeiShu(msgs []config.Msg, webHook_urls []string) error {

	if len(webHook_urls) == 0 {
		return fmt.Errorf("webhook urls is nil")
	}

	for _, webhookURL := range webHook_urls {

		// 遍历消息并发送
		for _, msg := range msgs {
			feishumsg := FeishuMessage{
				MsgType: "post",
				Content: Content{
					Post: Post{
						ZhCn: ZhCn{
							Title: msg.Title,
							Content: [][]ContentItem{
								{
									{
										Tag:  "text",
										Text: msg.Description,
									},
								},
							},
						},
					},
				},
			}
			// 将消息结构体编码为 JSON
			data, err := json.Marshal(feishumsg)
			if err != nil {
				return fmt.Errorf("error marshalling message: %v", err)
			}

			// 发送 POST 请求
			resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
			if err != nil {
				return fmt.Errorf("error sending message: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
			}
		}
	}
	return nil
}
