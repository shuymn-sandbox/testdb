package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shuymn-sandbox/testdb/example/model"
	"github.com/shuymn-sandbox/testdb/example/repository"
)

type userRepository struct {
	db *sqlx.DB
}

var _ repository.UserRepository = (*userRepository)(nil)

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, username, email string) (*model.User, error) {
	const q = "INSERT INTO users (username, email, created_at) VALUES (?, ?, ?)"

	now := time.Now()
	res, err := r.db.ExecContext(ctx, q, username, email, now)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return &model.User{
		ID:        int(id),
		Username:  username,
		Email:     email,
		CreatedAt: now,
	}, nil
}

func (r *userRepository) Find(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
