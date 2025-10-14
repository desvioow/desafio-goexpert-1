package main

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	config.Load()

	r := chi.NewRouter()
	r.Use(middleware.RateLimiterMiddleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO
		// criar handler que vai chamar strategy de persistencia para adicionar k-v no redis de token e ip da request
		w.Write([]byte("Hello, world! Limit my endpoint plz"))
	})
	http.ListenAndServe(":8000", r)
}
