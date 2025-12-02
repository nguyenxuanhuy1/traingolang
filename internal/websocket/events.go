package websocket

import (
	"fmt"
	"traingolang/internal/config"
)

// Xử lý event HIT
func handleHitEvent(c *Client, msg map[string]interface{}) {
	shooter := int(msg["shooter"].(float64))
	target := int(msg["target"].(float64))
	damage := int(msg["damage"].(float64))

	// 1. Trừ máu trong DB
	config.DB.Exec(`
		UPDATE match_players 
		SET hp = hp - $1, damage = damage + $1
		WHERE match_id=$2 AND user_id=$3
	`, damage, c.MatchID, target)

	// 2. Lấy HP mới
	var hp int
	config.DB.QueryRow(`
		SELECT hp FROM match_players 
		WHERE match_id=$1 AND user_id=$2
	`, c.MatchID, target).Scan(&hp)

	// 3. Broadcast HP mới
	hpMsg := fmt.Sprintf(`{"event":"hp_update","user_id":%d,"hp":%d}`, target, hp)
	c.Hub.Broadcast <- Message{MatchID: c.MatchID, Data: []byte(hpMsg)}

	// 4. Nếu chết
	if hp <= 0 {
		// Update DB: chết
		config.DB.Exec(`
			UPDATE match_players 
			SET alive=false, deaths=deaths+1 
			WHERE match_id=$1 AND user_id=$2
		`, c.MatchID, target)

		// Thêm kills cho shooter
		config.DB.Exec(`
			UPDATE match_players 
			SET kills = kills + 1
			WHERE match_id=$1 AND user_id=$2
		`, c.MatchID, shooter)

		// Broadcast chết
		deadMsg := fmt.Sprintf(`{"event":"dead","user_id":%d}`, target)
		c.Hub.Broadcast <- Message{MatchID: c.MatchID, Data: []byte(deadMsg)}

		// Check thắng
		checkWinner(c)
	}
}

// Kiểm tra còn 1 người alive → thắng
func checkWinner(c *Client) {
	rows, _ := config.DB.Query(`
		SELECT user_id FROM match_players 
		WHERE match_id=$1 AND alive=true
	`, c.MatchID)

	alivePlayers := []int{}
	for rows.Next() {
		var uid int
		rows.Scan(&uid)
		alivePlayers = append(alivePlayers, uid)
	}

	// Nếu còn 1 alive → người đó thắng
	if len(alivePlayers) == 1 {
		winner := alivePlayers[0]

		// Update DB
		config.DB.Exec(`
			UPDATE match_players 
			SET is_winner=true, final_rank=1
			WHERE match_id=$1 AND user_id=$2
		`, c.MatchID, winner)

		config.DB.Exec(`
			UPDATE matches 
			SET status='finished', ended_at=NOW()
			WHERE id=$1
		`, c.MatchID)

		// Broadcast winner
		winMsg := fmt.Sprintf(`{"event":"winner","user_id":%d}`, winner)
		c.Hub.Broadcast <- Message{MatchID: c.MatchID, Data: []byte(winMsg)}
	}
}
