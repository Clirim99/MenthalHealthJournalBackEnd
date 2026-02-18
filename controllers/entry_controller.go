package controllers

import (
	"menthalhealthjournal/services"

	"github.com/gin-gonic/gin"
)

func CreateEntry(c *gin.Context) {
	services.CreateEntry(c)
}

func GetEntry(c *gin.Context) {
	services.GetEntry(c)
}

func GetEntriesByUser(c *gin.Context) {
	services.GetEntriesByUser(c)
}

func UpdateEntry(c *gin.Context) {
	services.UpdateEntry(c)
}

func DeleteEntry(c *gin.Context) {
	services.DeleteEntry(c)
}
