package router

import (
	"traingolang/internal/api/handler"
	"traingolang/internal/auth"

	// "traingolang/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// PUBLIC ROUTES
	r.POST("/api/user/register", handler.Register)
	r.POST("/api/user/login", handler.Login)

	// PROTECTED ROUTES
	api := r.Group("/api")
	api.Use(auth.Middleware())
	{
		api.POST("/match/create", handler.CreateMatch)
		api.POST("/match/join", handler.JoinMatch)
		api.GET("/profile", handler.Profile)
		api.POST("/upload", handler.UploadHandler)
	}

	return r
}
