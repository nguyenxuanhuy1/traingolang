package repository

import (
	"database/sql"
	"time"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Đếm số token còn hạn của user
func (r *RefreshTokenRepository) CountValidByUser(userID int64) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM refresh_tokens
		WHERE user_id = $1
		  AND expires_at > now()
	`, userID).Scan(&count)

	return count, err
}

// Tạo refresh token
func (r *RefreshTokenRepository) Create(
	userID int64,
	token string,
	expiresAt time.Time,
) error {
	_, err := r.db.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expiresAt)

	return err
}

// Tìm user từ refresh token còn hạn
func (r *RefreshTokenRepository) Find(token string) (int64, error) {
	var userID int64
	err := r.db.QueryRow(`
		SELECT user_id
		FROM refresh_tokens
		WHERE token = $1
		  AND expires_at > now()
	`, token).Scan(&userID)

	return userID, err
}

// Logout 1 thiết bị
func (r *RefreshTokenRepository) Delete(token string) error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_tokens
		WHERE token = $1
	`)
	return err
}

// Khoá tài khoản → xoá toàn bộ token
func (r *RefreshTokenRepository) DeleteAllByUser(userID int64) error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`)
	return err
}

// Cleanup token hết hạn (cron job)
func (r *RefreshTokenRepository) CleanupExpired() error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_tokens
		WHERE expires_at < now()
	`)
	return err
}
