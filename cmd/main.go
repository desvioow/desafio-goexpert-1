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
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
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
	w.Write([]byte("Hello, world! Limit my / endpoint plz"))
}

// melhorar gerenciamento de conexao com redis
// melhorar injecao de dependencias
// melhorar logs de erros
// implementar algoritmo Token Bucket (Balde de Tokens) ou Moving Window
// terminar HOJE esse projeto
/*
Recomendação para Seu Caso
Escolha Token Bucket se:
Prioridade é simplicidade
Baixa latência é crítica
Possui recursos computacionais limitados
Aceita pequenas variações no rate limiting

Escolha Moving Window se:
Precisão é fundamental
Controle rigoroso de requisições
Prevenir abuse de API é crítico
Possui infraestrutura robusta

token bucket e mais simples e menos preciso
moving window e mais preciso e mais complexo
*/
