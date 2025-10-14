package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLimit(t *testing.T) {
	//TODO
	// mockar dependencias
	// rateLimiter, request e response
	rateLimiter := RateLimiter{strategy: nil}
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := rateLimiter.CheckLimit(nil, req)
	if err != nil {
		return
	}
}
