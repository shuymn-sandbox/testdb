package repository

import (
	"context"
	"errors"

	"github.com/shuymn-sandbox/testdb/example/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, name, email string) (*model.User, error)
	Find(ctx context.Context, id int) (*model.User, error)
}
