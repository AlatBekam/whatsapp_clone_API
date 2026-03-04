package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type channel struct {
	ChannelId	string `json:"channel_id"`
	ChannelName  string `json:"channel_name"`
	ChannelType string `json:"channel_type"`
	Description string `json:"description"`
}

type createChannelInput struct {
	ChannelName  string `json:"channel_name" binding:"required"`
	ChannelType string `json:"channel_type" binding:"required"`
	Description string `json:"description"`
}

// fetch channel (all)
func getChannel(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, channels)

	response := gin.H{"response": channels}

	fmt.Println("===== RESPONSE =====")
	fmt.Println(response)
	fmt.Println("====================")
}

// Tambah Channel
func addChannel(c *gin.Context) {
	var req createChannelInput

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newChannel := channel {
		ChannelId: generateChannelID(),
		ChannelName: req.ChannelName,
		ChannelType: req.ChannelType,
		Description: req.Description,
	}

	channels = append(channels, newChannel)
	c.IndentedJSON(http.StatusCreated, newChannel)
	data, _ := json.MarshalIndent(channels, "", "  ")
	os.WriteFile("data/channels.json", data, 0644)
}
