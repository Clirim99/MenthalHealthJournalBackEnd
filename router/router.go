package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"menthalhealthjournal/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS middleware (optional, but usually needed)
	r.Use(cors.Default())

	// Register POST /register route using Gin handler func directly
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)


	return r
}
