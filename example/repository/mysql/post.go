package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shuymn-sandbox/testdb/example/model"
	"github.com/shuymn-sandbox/testdb/example/repository"
)

type postRepository struct {
	db *sqlx.DB
}

var _ repository.PostRepository = (*postRepository)(nil)

func NewPostRepository(db *sqlx.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, user *model.User, title, content string) (*model.Post, error) {
	const q = "INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)"

	now := time.Now()
	res, err := r.db.ExecContext(ctx, q, user.ID, title, content, now)
	if err != nil {
		return nil, fmt.Errorf("failed to insert post: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return &model.Post{
		ID:        int(id),
		UserID:    user.ID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
	}, nil
}

func (r *postRepository) ListByUser(ctx context.Context, user *model.User) ([]*model.Post, error) {
	var posts []*model.Post
	err := r.db.SelectContext(ctx, &posts, "SELECT * FROM posts WHERE user_id = ?", user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	return posts, nil
}
