package main

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/middleware"
	"desafio-goexpert-1/internal/strategy/impl"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	config.Load()

	redisStrategy := &impl.RedisPersistenceStrategy{
		Host:     config.RedisConfig.Host,
		Port:     config.RedisConfig.Port,
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.DB,
	}
	err := redisStrategy.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisStrategy.Disconnect()

	r := chi.NewRouter()
	r.Use(middleware.RateLimiterMiddleware(redisStrategy))
	r.Get("/", httpHandler)
	http.ListenAndServe(":8080", r)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>I am not limited yet, try overloading me with requests!</h1>"))
}
