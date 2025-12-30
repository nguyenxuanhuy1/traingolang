package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidLogin = errors.New("invalid credentials")
)

type User struct {
	ID       int64
	Username string
	Password string
	Role     string
	Locked   bool
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Tạo user mới
func (r *UserRepository) Create(username, password string) (int64, error) {
	var id int64

	err := r.db.QueryRow(`
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id
	`, username, password).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, ErrUserExists
		}
		return 0, err
	}

	return id, nil
}

// Tìm user theo username
func (r *UserRepository) FindByUsername(username string) (*User, error) {
	var u User

	err := r.db.QueryRow(`
		SELECT id, username, password, role, locked
		FROM users
		WHERE username = $1
	`, username).Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.Role,
		&u.Locked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidLogin
		}
		return nil, err
	}

	return &u, nil
}
