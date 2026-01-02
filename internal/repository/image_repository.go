package repository

import "database/sql"

type Image struct {
	ID          int64
	ImageURL    string
	BlurURL     string
	TinyBlurURL string
	PublicID    string
	ImageType   string
	OwnerID     sql.NullInt64
}

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Create(img *Image) error {
	query := `
		INSERT INTO images
		(image_url, blur_url, tiny_blur_url, public_id, image_type, owner_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	return r.db.QueryRow(
		query,
		img.ImageURL,
		img.BlurURL,
		img.TinyBlurURL,
		img.PublicID,
		img.ImageType,
		img.OwnerID,
	).Scan(&img.ID)
}
