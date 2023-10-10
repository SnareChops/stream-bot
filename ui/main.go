package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {

	router := chi.NewRouter()

	router.Get("/", file("ui/index.html"))
	router.Get("/index.js", file("ui/index.js"))
	router.Get("/index.css", file("ui/index.css"))

	return router
}

func file(path string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, path)
	})
}
