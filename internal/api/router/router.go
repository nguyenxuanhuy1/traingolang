package router

import (
	"traingolang/internal/api/handler"
	"traingolang/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// API
	r.GET("/ping", handler.Ping)
	r.POST("/user/create", handler.CreateUser)

	// WebSocket: NHIỀU PHÒNG
	r.GET("/ws/:match_id", websocket.HandleWebSocket)

	return r
}
