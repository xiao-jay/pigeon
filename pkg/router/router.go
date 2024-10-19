package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var validChannelType = map[string]string{
	"feishu":   "true",
	"fangtang": "true",
}

func InitRouter() {
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
		c.JSON(http.StatusBadRequest, gin.H{"status": "fial", "message": fmt.Sprintf("not support:%s", json.ChannelType)})
	}

	// Respond back to the client
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": json.Message})

}
