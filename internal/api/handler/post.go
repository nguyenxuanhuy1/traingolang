package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Post struct {
	ID          int64
	ImageID     sql.NullInt64
	Name        string
	Description string
	Topic       string
	Prompt      sql.NullString
	IsHot       bool
	HotAt       sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ----- Struct request -----
type CreatePostRequest struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
	Topic       string `form:"topic"`
	Prompt      string `form:"prompt"`
	IsHot       bool   `form:"isHot"`
}

func CreatePostHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Parse các field thường
		var req CreatePostRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. Lấy file ảnh
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
			return
		}

		// 3. Lưu file vào folder local (hoặc thay bằng Cloudinary)
		dst := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 4. Lưu vào bảng images
		result, err := db.Exec("INSERT INTO images (url, created_at) VALUES (?, ?)", dst, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imageID, _ := result.LastInsertId()

		// 5. Tạo post với ImageID
		post := Post{
			ImageID:     sql.NullInt64{Int64: imageID, Valid: true},
			Name:        req.Name,
			Description: req.Description,
			Topic:       req.Topic,
			Prompt:      sql.NullString{String: req.Prompt, Valid: req.Prompt != ""},
			IsHot:       req.IsHot,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// 6. Insert post vào DB
		resultPost, err := db.Exec(
			`INSERT INTO posts (image_id, name, description, topic, prompt, is_hot, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			post.ImageID, post.Name, post.Description, post.Topic, post.Prompt, post.IsHot, post.CreatedAt, post.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		postID, _ := resultPost.LastInsertId()
		post.ID = postID

		c.JSON(http.StatusOK, gin.H{"data": post})
	}
}
