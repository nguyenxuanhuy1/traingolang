package handler

import (
	"database/sql"
	"net/http"

	"traingolang/internal/auth"
	"traingolang/internal/config"
	"traingolang/internal/repository"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	// 1️⃣ Validate image_type trước
	imageType := c.PostForm("image_type")
	if imageType == "" {
		c.JSON(400, gin.H{"error": "image_type required"})
		return
	}

	// 2️⃣ Lấy file
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}
	defer file.Close()

	// Upload Cloudinary
	result, err := service.UploadImageWithBlur(file, imageType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Lấy user_id từ token
	claimsAny, _ := c.Get(auth.ContextUserKey)
	claims := claimsAny.(*auth.Claims)

	imageRepo := repository.NewImageRepository(config.DB)

	img := repository.Image{
		ImageURL:    result.Image,
		BlurURL:     result.Blur,
		TinyBlurURL: result.TinyBlur,
		PublicID:    result.PublicID,
		ImageType:   imageType,
		OwnerID:     sql.NullInt64{Int64: claims.UserID, Valid: true},
	}

	if err := imageRepo.Create(&img); err != nil {
		c.JSON(500, gin.H{"error": "save image failed"})
		return
	}
}
