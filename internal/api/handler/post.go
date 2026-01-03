package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"traingolang/internal/repository"

	"github.com/gin-gonic/gin"
)

type CreatePostRequest struct {
	Name        string `json:"name"`        // tên bài viết
	Description string `json:"description"` // mô tả
	Topic       string `json:"topic"`       // chủ đề
	Prompt      string `json:"prompt"`      // prompt nếu có
	HotLevel    int8   `json:"hotLevel"`
}
type SearchPostRequest struct {
	Name     string `json:"name"`
	hotLevel *int8  `json:"isHot"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

func CreatePost(postRepo repository.PostRepo, imageRepo repository.ImageRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy JSON data
		dataStr := c.PostForm("data")
		if dataStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "data is required"})
			return
		}

		var req CreatePostRequest
		if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data object"})
			return
		}

		// 2. Upload ảnh + lưu DB dùng helper
		img, err := UploadAndSaveImage(c, "post")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 3. Tạo post
		post := &repository.Post{
			ImageID:     sql.NullInt64{Int64: img.ID, Valid: true},
			Name:        req.Name,
			Description: req.Description,
			Topic:       req.Topic,
			Prompt:      sql.NullString{String: req.Prompt, Valid: req.Prompt != ""},
			HotLevel:    req.HotLevel,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		postID, err := postRepo.CreatePost(post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		post.ID = postID

		c.JSON(http.StatusCreated, gin.H{
			"message": "Tạo mới thành công",
		})
	}
}
func SearchPostsHandler(postRepo repository.PostRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SearchPostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		result, err := postRepo.SearchPosts(req.Name, req.hotLevel, req.Page, req.PageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
