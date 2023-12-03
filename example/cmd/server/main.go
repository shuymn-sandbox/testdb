package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/shuymn-sandbox/testdb/example/handler"
	repository "github.com/shuymn-sandbox/testdb/example/repository/mysql"
)

func main() {
	config := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db := sqlx.MustOpen("mysql", config.FormatDSN())

	userRepository := repository.NewUserRepository(db)
	postRepository := repository.NewPostRepository(db)

	userHandler := handler.NewUserHandler(userRepository)
	postHandler := handler.NewPostHandler(userRepository, postRepository)

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		// user
		r.Post("/users", userHandler.CreateUser)
		r.Route("/users/{userID:[0-9]+}", func(r chi.Router) {
			// post
			r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
				postHandler.ListPosts(w, r, chi.URLParam(r, "userID"))
			})
			r.Post("/posts", func(w http.ResponseWriter, r *http.Request) {
				postHandler.CreatePost(w, r, chi.URLParam(r, "userID"))
			})
		})
	})

	http.ListenAndServe(":8080", r)
}
