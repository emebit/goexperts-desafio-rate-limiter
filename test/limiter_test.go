package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"Desafios/3/goexperts-desafio-rate-limiter/internal/limiter"
	"Desafios/3/goexperts-desafio-rate-limiter/internal/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	store := limiter.NewRedisStore(rdb)

	rl := limiter.NewRateLimiter(store, 5, 10, 10, 10)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	limitedHandler := middleware.RateLimiter(rl, handler)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
		req.RemoteAddr = "172.31.240.1" //"192.168.1.1"
		resp := httptest.NewRecorder()

		limitedHandler.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	}

	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	req.RemoteAddr = "172.31.240.1" //"192.168.1.1"
	resp := httptest.NewRecorder()

	limitedHandler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusTooManyRequests, resp.Code)

	time.Sleep(11 * time.Second)

	req = httptest.NewRequest("GET", "http://localhost:8080/", nil)
	req.RemoteAddr = "172.31.240.1" //"192.168.1.1"
	resp = httptest.NewRecorder()

	limitedHandler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
