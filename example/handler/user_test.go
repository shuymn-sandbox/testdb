package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shuymn-sandbox/testdb/example/handler"
	repository "github.com/shuymn-sandbox/testdb/example/repository/mysql"
)

func TestUserHandlerCreateUser(t *testing.T) {
	userRepository := repository.NewUserRepository(db)
	h := handler.NewUserHandler(userRepository)

	b := strings.NewReader(`{"username":"username","email":"email"}`)
	req, err := http.NewRequest(http.MethodPost, "/users", b)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	h.CreateUser(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("got %v, want %v", status, http.StatusCreated)
	}

	got := rr.Body.String()
	want := `,"username":"username","email":"email","createdAt":"`
	if !strings.Contains(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
