package middleware

import (
	"fmt"
	"log"
	"net/http"
)

//TODO
// injetar RateLimiter no middleware

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Hello from RateLimiterMiddleware")

		next.ServeHTTP(w, r)
	})
}

func requestIsLimited(w http.ResponseWriter, r *http.Request) (bool, error) {
	log.Fatal("requestIsLimited not implemented")

	//TODO
	// chamar RateLimiter.limit() e retornar true ou false para tomar decisão de limitar ou não

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`))

	return true, nil
}
