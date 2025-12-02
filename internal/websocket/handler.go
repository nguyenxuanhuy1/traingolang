package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var hub = NewHub()

func init() {
	go hub.Run()
}

func HandleWebSocket(c *gin.Context) {
	matchID, _ := strconv.Atoi(c.Param("match_id"))
	userID, _ := strconv.Atoi(c.Query("user_id"))

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := NewClient(hub, conn, matchID, userID)
	hub.Register <- client

	go client.ReadPump()
	go client.WritePump()

	joinMsg := []byte(`{"event":"player_joined","user_id":` + strconv.Itoa(userID) + `}`)
	hub.Broadcast <- Message{MatchID: matchID, Data: joinMsg}
}
