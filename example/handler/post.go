package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/shuymn-sandbox/testdb/example/repository"
)

type postHandler struct {
	userRepository repository.UserRepository
	postRepository repository.PostRepository
}

func NewPostHandler(
	userRepository repository.UserRepository,
	postRepository repository.PostRepository,
) *postHandler {
	return &postHandler{
		userRepository: userRepository,
		postRepository: postRepository,
	}
}

func (h *postHandler) CreatePost(w http.ResponseWriter, r *http.Request, userIDStr string) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.userRepository.Find(ctx, int(userID))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := h.postRepository.Create(ctx, user, input.Title, input.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *postHandler) ListPosts(w http.ResponseWriter, r *http.Request, userIDStr string) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.userRepository.Find(ctx, int(userID))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := h.postRepository.ListByUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
