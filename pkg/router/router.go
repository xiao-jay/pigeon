package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pigeon/config"
	channels "pigeon/pkg/channel"
)

var validChannelType = map[string]func(msgs []config.Msg, extra any) error{
	"feishu":   channels.Feishu{}.SendMessage,
	"fangtang": channels.FangTang{}.SendMessage,
}

var config3 config.Config

func InitRouter(config2 *config.Config) {
	config3 = *config2
	r := gin.New()
	r.POST("/send_messages", sendMessages)
	err := r.Run(":7030")
	if err != nil {
		return
	}
}

// sendMessages handles the POST request for /send_messages
func sendMessages(c *gin.Context) {
	var json struct {
		Title       string `json:"title"`
		Message     string `json:"message" binding:"required"` // Expecting a message field in the JSON body
		ChannelType string `json:"channelType"`
	}

	// Bind the JSON body to the struct
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Here you can process the message as needed
	log.Printf("Received message: %s\n\n", json.Message)

	if _, found := validChannelType[json.ChannelType]; !found {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": fmt.Sprintf("not support:%s", json.ChannelType)})
	}
	msg := config.Msg{
		Title:       json.Title,
		Description: json.Message,
		Channel:     0,
	}
	msg_list := append([]config.Msg{}, msg)
	if json.ChannelType == "feishu" {
		if err := validChannelType[json.ChannelType](msg_list, config3.FeishuWebHooks); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": fmt.Sprintf("err is :%s", err)})
		}
	}

	// Respond back to the client
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": json.Message})

}
