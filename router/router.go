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

	// User routes
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	// Entry routes
	r.POST("/entries", controllers.CreateEntry)
	r.GET("/entries/:id", controllers.GetEntry)
	r.GET("/users/:user_id/entries", controllers.GetEntriesByUser)
	r.PUT("/entries/:id", controllers.UpdateEntry)
	r.DELETE("/entries/:id", controllers.DeleteEntry)

	// Chat routes
	r.POST("/chat", controllers.ChatWithJournal)
	r.POST("/chat/sessions", controllers.CreateChatSession)
	r.GET("/chat/sessions/:id", controllers.GetChatSession)
	r.GET("/users/:user_id/chat/sessions", controllers.GetChatSessionsByUser)
	r.GET("/chat/sessions/:id/messages", controllers.GetChatHistory)
	r.DELETE("/chat/sessions/:id", controllers.DeleteChatSession)

	return r
}
