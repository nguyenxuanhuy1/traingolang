package handler

import (
	"net/http"
	"traingolang/internal/config"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}

	// Parse JSON từ body request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// Insert vào DB
	var id int
	err := config.DB.QueryRow(`
		INSERT INTO users (username, avatar)
		VALUES ($1, $2)
		RETURNING id
	`, req.Username, req.Avatar).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"username": req.Username,
	})
}
