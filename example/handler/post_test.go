package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/shuymn-sandbox/testdb/example/handler"
	repository "github.com/shuymn-sandbox/testdb/example/repository/mysql"
)

func TestPostHandlerCreatePost(t *testing.T) {
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

	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)
	h := handler.NewPostHandler(userRepository, postRepository)

	b := strings.NewReader(`{"title":"title","content":"content"}`)
	req, err := http.NewRequest(http.MethodPost, "/users/1/posts", b)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	h.CreatePost(rr, req, "1")
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("got %v, want %v", status, http.StatusCreated)
	}

	got := rr.Body.String()
	want := `,"userId":1,"title":"title","content":"content","createdAt":"`
	if !strings.Contains(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPostHandlerListPosts(t *testing.T) {
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

	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)
	h := handler.NewPostHandler(userRepository, postRepository)

	req, err := http.NewRequest(http.MethodGet, "/users/1/posts", http.NoBody)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	h.ListPosts(rr, req, "1")
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("got %v, want %v", status, http.StatusOK)
	}

	got := rr.Body.String()
	want := `[{"id":1,"userId":1,"title":"Hello, world!","content":"This is a test post.","createdAt":"2016-01-01T00:00:00Z"},{"id":2,"userId":1,"title":"Hello, world! 2","content":"This is a test post 2.","createdAt":"2016-01-02T00:00:00Z"}]` + "\n"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
