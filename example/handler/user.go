package handler

import (
	"encoding/json"
	"net/http"

	"github.com/shuymn-sandbox/testdb/example/repository"
)

type userHandler struct {
	userRepository repository.UserRepository
}

func NewUserHandler(userRepository repository.UserRepository) *userHandler {
	return &userHandler{userRepository: userRepository}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userRepository.Create(r.Context(), input.Username, input.Email)
	if err != nil {
		// NOTE: エラーの種類によってステータスコードを変えるべき
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
