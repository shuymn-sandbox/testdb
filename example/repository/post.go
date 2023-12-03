package repository

import (
	"context"

	"github.com/shuymn-sandbox/testdb/example/model"
)

type PostRepository interface {
	Create(ctx context.Context, user *model.User, title, content string) (*model.Post, error)
	ListByUser(ctx context.Context, user *model.User) ([]*model.Post, error)
}
