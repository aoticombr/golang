package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter gerencia rate limiters por IP
type RateLimiter struct {
	// limiters armazena um mapa onde a chave é o IP do cliente e o valor é seu rate limiter individual
	// Cada IP tem seu próprio contador de requests independente dos outros IPs
	limiters map[string]*rate.Limiter

	// mu é um mutex de leitura/escrita para proteger acesso concorrente ao mapa limiters
	// Permite múltiplas leituras simultâneas mas apenas uma escrita por vez
	mu sync.RWMutex

	// rate define quantas requests por segundo são permitidas (ex: 10.0 = 10 requests/segundo)
	// Tipo rate.Limit do pacote golang.org/x/time/rate
	rate rate.Limit

	// burst define o número máximo de requests que podem ser feitas de uma só vez
	// Funciona como um "balde" que se reabastece na velocidade definida por 'rate'
	burst int

	// cleanup define o intervalo de tempo para executar limpeza automática
	// Remove rate limiters de IPs que não fazem requests há muito tempo para economizar memória
	cleanup time.Duration
}

// NewRateLimiter cria um novo rate limiter
func NewRateLimiter(requestsPerSecond float64, burstSize int, cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burstSize,
		cleanup:  cleanupInterval,
	}

	// Goroutine para limpeza periódica de limiters antigos
	go rl.cleanupRoutine()

	return rl
}

// getLimiter retorna o rate limiter para um IP específico
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// cleanupRoutine remove rate limiters antigos
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(rl.cleanup)
	for {
		<-ticker.C
		rl.mu.Lock()
		for ip, limiter := range rl.limiters {
			// Remove limiters que não foram usados nos últimos 5 minutos
			if limiter.TokensAt(time.Now().Add(-5*time.Minute)) >= float64(rl.burst) {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// getClientIP extrai o IP real do cliente
func getClientIP(r *http.Request) string {
	// Verifica headers de proxy
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Pega o primeiro IP da lista
		if idx := len(ip); idx > 0 {
			if commaIdx := 0; commaIdx < idx {
				for i, char := range ip {
					if char == ',' {
						commaIdx = i
						break
					}
				}
				if commaIdx > 0 {
					ip = ip[:commaIdx]
				}
			}
		}
		return ip
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fallback para RemoteAddr
	ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	return ip
}

// RateLimitMiddleware middleware para rate limiting por IP
func (rl *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitWithWait middleware que aguarda em vez de rejeitar
func (rl *RateLimiter) RateLimitWithWait(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := rl.getLimiter(ip)

		// Aguarda até que seja permitido (com timeout)
		ctx := r.Context()
		err := limiter.WaitN(ctx, 1)
		if err != nil {
			http.Error(w, "Request timeout due to rate limiting", http.StatusRequestTimeout)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Instância global de rate limiter
var (
	// Rate limiter padrão: 10 requests por segundo, burst de 20
	DefaultRateLimiter = NewRateLimiter(10, 20, time.Minute*1)

	// Rate limiter mais restritivo: 2 requests por segundo, burst de 5
	StrictRateLimiter = NewRateLimiter(2, 5, time.Minute*1)
)

// Middlewares prontos para uso
func RateLimit() func(http.Handler) http.Handler {
	return DefaultRateLimiter.RateLimitMiddleware
}

func StrictRateLimit() func(http.Handler) http.Handler {
	return StrictRateLimiter.RateLimitMiddleware
}

func RateLimitWithWait() func(http.Handler) http.Handler {
	return DefaultRateLimiter.RateLimitWithWait
}
