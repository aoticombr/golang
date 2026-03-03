package middlewares

/*
   RATE LIMITING COM REDIS - EXEMPLO

   Para usar Redis rate limiting, siga os passos:

   1. Adicione a dependência:
      go get github.com/go-redis/redis/v8

   2. Descomente o código abaixo

   3. Configure sua conexão Redis

   Vantagens do Redis:
   - Rate limiting distribuído (múltiplos servidores)
   - Persiste entre reinicializações
   - Melhor performance para alta escala


import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRateLimiter struct {
	client   *redis.Client
	requests int
	window   time.Duration
	prefix   string
}

func NewRedisRateLimiter(redisClient *redis.Client, requests int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:   redisClient,
		requests: requests,
		window:   window,
		prefix:   "rate_limit:",
	}
}

func (r *RedisRateLimiter) IsAllowed(ctx context.Context, ip string) (bool, error) {
	key := r.prefix + ip

	script := `
		local key = KEYS[1]
		local window = tonumber(ARGV[1])
		local limit = tonumber(ARGV[2])
		local current_time = tonumber(ARGV[3])

		redis.call('ZREMRANGEBYSCORE', key, 0, current_time - window)

		local current = redis.call('ZCARD', key)

		if current < limit then
			redis.call('ZADD', key, current_time, current_time)
			redis.call('EXPIRE', key, window)
			return 1
		else
			return 0
		end
	`

	result, err := r.client.Eval(ctx, script, []string{key},
		r.window.Milliseconds(),
		r.requests,
		time.Now().UnixMilli()).Result()

	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

func (r *RedisRateLimiter) RedisRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip := getClientIP(req)
		ctx := req.Context()

		allowed, err := r.IsAllowed(ctx, ip)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(r.requests))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, req)
	})
}

*/
