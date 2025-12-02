package websocket

import (
	"fmt"
)

// Lưu vị trí tạm trong RAM
var PlayerPositions = make(map[int]map[int]map[string]float64)

// PlayerPositions[matchID][userID]["x"] = 100

func handleMoveEvent(c *Client, msg map[string]interface{}) {
	userID := int(msg["user_id"].(float64))
	x := msg["x"].(float64)
	y := msg["y"].(float64)

	if PlayerPositions[c.MatchID] == nil {
		PlayerPositions[c.MatchID] = make(map[int]map[string]float64)
	}

	PlayerPositions[c.MatchID][userID] = map[string]float64{
		"x": x,
		"y": y,
	}

	// Broadcast vị trí mới cho tất cả người trong phòng
	moveMsg := fmt.Sprintf(
		`{"event":"move_update","user_id":%d,"x":%f,"y":%f}`,
		userID, x, y,
	)

	c.Hub.Broadcast <- Message{
		MatchID: c.MatchID,
		Data:    []byte(moveMsg),
	}
}
