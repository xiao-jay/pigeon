package channels

import (
	"fmt"
	"pigeon/config"
	"testing"
)

func TestSendMessageFeiShu(t *testing.T) {
	Config, err := config.GetConf("../../config/config.yaml")
	if err != nil {
		panic(err)
	}

	msgs := []config.Msg{
		{
			Title:       "测试消息",
			Description: "这是一个来自 Go 的测试消息。",
			Channel:     1, // 根据你的需求设置渠道
		},
	}

	if err := SendMessageFeiShu(msgs, Config.FeishuWebHooks); err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	} else {
		fmt.Println("Message sent successfully!")
	}
}
