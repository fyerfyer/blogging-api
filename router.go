package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *Application) Routers() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/posts", app.Post.CreatePost).Methods("POST")
	r.HandleFunc("/posts/{id}", app.Post.UpdatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", app.Post.DeletePost).Methods("DELETE")
	r.HandleFunc("/posts/{id}", app.Post.GetPost).Methods("GET")
	r.HandleFunc("/posts", app.Post.GetAllPosts).Methods("GET")

	return r
}
