package repository

import (
	"database/sql"
	"fmt"
	"time"
	"traingolang/internal/util"
)

type Post struct {
	ID          int64
	ImageID     sql.NullInt64
	Name        string
	Description string
	Topic       string
	Prompt      sql.NullString
	HotLevel    int8
	HotAt       sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostRepo interface {
	CreatePost(post *Post) (int64, error)
	SearchPosts(query string, hotLevel *int8, page, pageSize int) (*util.PaginatedResponse[Post], error)
	GetByID(id int64) (*Post, error)
	UpdatePost(post *Post) error
	DeletePost(id int64) error
}

type postRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) PostRepo {
	return &postRepo{db: db}
}

// ------------------ Create ------------------
func (r *postRepo) CreatePost(post *Post) (int64, error) {
	query := `
			INSERT INTO posts (image_id, name, description, topic, prompt, hot_level, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
	var id int64
	err := r.db.QueryRow(
		query,
		post.ImageID, post.Name, post.Description, post.Topic, post.Prompt, post.HotLevel, post.CreatedAt, post.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	post.ID = id
	return id, nil
}

// ------------------ Search với phân trang ------------------
func (r *postRepo) SearchPosts(name string, hotLevel *int8, page, pageSize int) (*util.PaginatedResponse[Post], error) {
	offset, limit := util.NewPagination(page, pageSize)

	where := "1=1"
	args := []interface{}{}

	// filter name
	if name != "" {
		where += fmt.Sprintf(" AND name ILIKE $%d", len(args)+1)
		args = append(args, "%"+name+"%")
	}

	// filter isHot
	if hotLevel != nil {
		where += fmt.Sprintf(" AND hot_level = $%d", len(args)+1)
		args = append(args, *hotLevel)
	}

	// LIMIT/OFFSET
	limitIdx := len(args) + 1
	offsetIdx := len(args) + 2
	query := fmt.Sprintf(`
			SELECT id, image_id, name, description, topic, prompt, hot_level, hot_at, created_at, updated_at
			FROM posts
			WHERE %s
			ORDER BY created_at DESC
			LIMIT $%d OFFSET $%d
		`, where, limitIdx, offsetIdx)
	// countQuery chỉ dùng args filter, bỏ limit/offset
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM posts WHERE %s", where)
	return util.Paginate[Post](
		r.db,
		query,
		countQuery,
		args,
		offset,
		limit,
		func(rows *sql.Rows) (*Post, error) {
			var p Post
			err := rows.Scan(
				&p.ID,
				&p.ImageID,
				&p.Name,
				&p.Description,
				&p.Topic,
				&p.Prompt,
				&p.HotLevel,
				&p.HotAt,
				&p.CreatedAt,
				&p.UpdatedAt,
			)
			return &p, err
		},
	)
}
func (r *postRepo) GetByID(id int64) (*Post, error) {
	query := `
			SELECT id, image_id, name, description, topic, prompt, hot_level, hot_at, created_at, updated_at
			FROM posts
			WHERE id = $1
		`

	var p Post
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.ImageID,
		&p.Name,
		&p.Description,
		&p.Topic,
		&p.Prompt,
		&p.HotLevel,
		&p.HotAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}
func (r *postRepo) DeletePost(id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
func (r *postRepo) UpdatePost(post *Post) error {
	query := `
		UPDATE posts
		SET
			image_id = $1,
			name = $2,
			description = $3,
			topic = $4,
			prompt = $5,
			hot_level = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(
		query,
		post.ImageID,
		post.Name,
		post.Description,
		post.Topic,
		post.Prompt,
		post.HotLevel,
		post.UpdatedAt,
		post.ID,
	)

	return err
}
