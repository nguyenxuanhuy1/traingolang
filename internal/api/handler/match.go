package handler

import (
	"traingolang/internal/config"

	"github.com/gin-gonic/gin"
)

func CreateMatch(c *gin.Context) {
	var req struct {
		MaxPlayers int `json:"max_players"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	var id int
	config.DB.QueryRow(`
		INSERT INTO matches (status, max_players)
		VALUES ('waiting', $1)
		RETURNING id
	`, req.MaxPlayers).Scan(&id)

	c.JSON(200, gin.H{"match_id": id})
}

func JoinMatch(c *gin.Context) {
	var req struct {
		MatchID int `json:"match_id"`
		UserID  int `json:"user_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	_, err := config.DB.Exec(`
		INSERT INTO match_players (match_id, user_id)
		VALUES ($1, $2)
	`, req.MatchID, req.UserID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "joined"})
}
