package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	MatchID int // phòng
	UserID  int // người chơi
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan []byte
}

// Constructor tạo client
func NewClient(hub *Hub, conn *websocket.Conn, matchID, userID int) *Client {
	return &Client{
		MatchID: matchID,
		UserID:  userID,
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan []byte, 256),
	}
}

// Đọc message từ client → gửi vào Hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, rawMsg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse JSON
		var msg map[string]interface{}
		json.Unmarshal(rawMsg, &msg)

		if msg["event"] == "move" {
			handleMoveEvent(c, msg)
			continue
		}
		if msg["event"] == "shoot" {
			handleShootEvent(c, msg)
			continue
		}
		// Xử lý event HIT
		if msg["event"] == "hit" {
			handleHitEvent(c, msg)
			continue
		}

		// Nếu không phải hit → broadcast bình thường
		c.Hub.Broadcast <- Message{
			MatchID: c.MatchID,
			Data:    rawMsg,
		}
	}

}

// Gửi message từ Hub → client
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
