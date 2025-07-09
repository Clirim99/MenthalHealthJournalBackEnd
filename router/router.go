package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"menthalhealthjournal/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Optional: enable CORS
	r.Use(cors.Default())

	r.POST("/register", gin.WrapF(controllers.RegisterUser))

	return r
}
