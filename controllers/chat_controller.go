package controllers

import (
	"menthalhealthjournal/services"

	"github.com/gin-gonic/gin"
)

func CreateChatSession(c *gin.Context) {
	services.CreateChatSession(c)
}

func GetChatSession(c *gin.Context) {
	services.GetChatSession(c)
}

func GetChatSessionsByUser(c *gin.Context) {
	services.GetChatSessionsByUser(c)
}

func ChatWithJournal(c *gin.Context) {
	services.ChatWithJournal(c)
}

func GetChatHistory(c *gin.Context) {
	services.GetChatHistory(c)
}

func DeleteChatSession(c *gin.Context) {
	services.DeleteChatSession(c)
}
