package mysql_test

import (
	"context"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shuymn-sandbox/testdb/example/model"
	"github.com/shuymn-sandbox/testdb/example/repository/mysql"
)

func TestPostRepositoryCreate(t *testing.T) {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("mysql"),
		testfixtures.Paths("testdata/users.yml"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		t.Fatalf("failed to create fixtures: %v", err)
	}
	if err = fixtures.Load(); err != nil {
		t.Fatalf("failed to load fixtures: %v", err)
	}

	ctx := context.Background()
	r := mysql.NewPostRepository(db)
	user := &model.User{ID: 1}
	got, err := r.Create(ctx, user, "title", "content")
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	if got.UserID != user.ID {
		t.Errorf("got %v, want %v", got.UserID, user.ID)
	}
	if got.Title != "title" {
		t.Errorf("got %v, want %v", got.Title, "title")
	}
	if got.Content != "content" {
		t.Errorf("got %v, want %v", got.Content, "content")
	}
}

func TestPostRepositoryListByUser(t *testing.T) {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db.DB),
		testfixtures.Dialect("mysql"),
		testfixtures.Paths("testdata/users.yml", "testdata/posts.yml"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		t.Fatalf("failed to create fixtures: %v", err)
	}
	if err = fixtures.Load(); err != nil {
		t.Fatalf("failed to load fixtures: %v", err)
	}

	ctx := context.Background()
	r := mysql.NewPostRepository(db)
	user := &model.User{ID: 1}
	got, err := r.ListByUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to list posts: %v", err)
	}
	// 適当
	if len(got) != 2 {
		t.Errorf("got %v, want %v", len(got), 2)
	}
}
