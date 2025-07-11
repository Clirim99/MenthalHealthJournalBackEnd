package controllers

import (
	//"log"
	//"net/http"
	//"strings"

	"menthalhealthjournal/services"

	"github.com/gin-gonic/gin"
	//"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context){
    services.RegisterUser(c)
}
func LoginUser (c *gin.Context){
    services.LoginUser(c)
}