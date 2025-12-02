package websocket

import (
	"fmt"
)

func handleShootEvent(c *Client, msg map[string]interface{}) {
	userID := int(msg["user_id"].(float64))
	x := msg["x"].(float64)
	y := msg["y"].(float64)
	dirX := msg["dirX"].(float64)
	dirY := msg["dirY"].(float64)

	arrow := fmt.Sprintf(
		`{"event":"arrow","user_id":%d,"x":%f,"y":%f,"dirX":%f,"dirY":%f}`,
		userID, x, y, dirX, dirY,
	)

	c.Hub.Broadcast <- Message{
		MatchID: c.MatchID,
		Data:    []byte(arrow),
	}
}
