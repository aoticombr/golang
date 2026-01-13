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

// NewRateLimiter cria um novo rate limiter personalizado
// Parâmetros:
//   - requestsPerSecond: número de requests permitidas por segundo (ex: 10.0)
//   - burstSize: número máximo de requests simultâneas permitidas (ex: 20)
//   - cleanupInterval: intervalo para limpeza de limiters inativos (ex: 1 minuto)
//
// Retorna uma instância configurada de RateLimiter
// Exemplo: NewRateLimiter(5.0, 10, time.Minute) = 5 req/seg, burst de 10, cleanup a cada minuto
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

// getLimiter retorna o rate limiter específico para um IP
// Se o IP não existir no mapa, cria um novo limiter para ele
// Thread-safe: usa mutex para proteger acesso concorrente
//
// Parâmetros:
//   - ip: endereço IP do cliente
//
// Retorna o rate.Limiter associado ao IP
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

// cleanupRoutine executa limpeza periódica de rate limiters inativos
// Goroutine que roda em background removendo limiters de IPs que não fazem requests há muito tempo
// Critério: remove limiters que não foram usados nos últimos 5 minutos
//
// Objetivo: economizar memória evitando acúmulo infinito de limiters
// Executa automaticamente no intervalo definido em rl.cleanup
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

// getClientIP extrai o endereço IP real do cliente da requisição HTTP
// Considera headers de proxy e load balancers para obter o IP original
//
// Ordem de prioridade:
// 1. X-Forwarded-For (pega o primeiro IP da lista, ignorando proxies intermediários)
// 2. X-Real-IP (header alternativo usado por alguns proxies)
// 3. RemoteAddr (IP direto da conexão, fallback)
//
// Parâmetros:
//   - r: ponteiro para http.Request
//
// Retorna o IP do cliente como string
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

// RateLimitMiddleware middleware que aplica rate limiting por IP e rejeita requests excessivos
// Comportamento: se o limite for excedido, retorna erro 429 (Too Many Requests) imediatamente
//
// Como funciona:
// 1. Extrai o IP do cliente
// 2. Obtém/cria o rate limiter para esse IP
// 3. Verifica se o request é permitido
// 4. Se sim: passa para o próximo handler
// 5. Se não: retorna erro 429
//
// Parâmetros:
//   - next: próximo handler na cadeia de middlewares
//
// Retorna http.Handler que pode ser usado na cadeia de middlewares
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

// RateLimitWithWait middleware que aplica rate limiting mas aguarda em vez de rejeitar
// Comportamento: se o limite for excedido, aguarda até que seja possível processar o request
//
// Como funciona:
// 1. Extrai o IP do cliente
// 2. Obtém/cria o rate limiter para esse IP
// 3. Aguarda até que um "token" esteja disponível (respeitando timeout do contexto)
// 4. Se timeout: retorna erro 408 (Request Timeout)
// 5. Se conseguir: passa para o próximo handler
//
// Parâmetros:
//   - next: próximo handler na cadeia de middlewares
//
// Retorna http.Handler que pode ser usado na cadeia de middlewares
// Cuidado: pode causar requests lentos se muitos clientes estiverem aguardando
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

// Instâncias globais de rate limiter pré-configuradas
var (
	// DefaultRateLimiter: configuração balanceada para uso geral
	// 10 requests por segundo, burst de 20, limpeza a cada minuto
	// Adequado para APIs REST comuns
	DefaultRateLimiter = NewRateLimiter(10, 20, time.Minute*1)

	// StrictRateLimiter: configuração mais restritiva para endpoints sensíveis
	// 2 requests por segundo, burst de 5, limpeza a cada minuto
	// Adequado para endpoints de login, cadastro, operações críticas
	StrictRateLimiter = NewRateLimiter(2, 5, time.Minute*1)
)

// RateLimit retorna um middleware de rate limiting com configuração padrão
// Usa DefaultRateLimiter (10 req/seg, burst 20)
// Rejeita requests excessivos com status 429
//
// Exemplo de uso:
//
//	router.Use(middlewares.RateLimit())
//
// Retorna função que pode ser usada como middleware
func RateLimit() func(http.Handler) http.Handler {
	return DefaultRateLimiter.RateLimitMiddleware
}

// StrictRateLimit retorna um middleware de rate limiting mais restritivo
// Usa StrictRateLimiter (2 req/seg, burst 5)
// Adequado para endpoints sensíveis (login, registro, etc.)
//
// Exemplo de uso:
//
//	authRouter.Use(middlewares.StrictRateLimit())
//
// Retorna função que pode ser usada como middleware
func StrictRateLimit() func(http.Handler) http.Handler {
	return StrictRateLimiter.RateLimitMiddleware
}

// RateLimitWithWait retorna um middleware que aguarda em vez de rejeitar
// Usa DefaultRateLimiter mas com comportamento de espera
// Pode causar latência alta se muitos clientes estiverem aguardando
//
// Exemplo de uso:
//
//	uploadRouter.Use(middlewares.RateLimitWithWait())
//
// Retorna função que pode ser usada como middleware
func RateLimitWithWait() func(http.Handler) http.Handler {
	return DefaultRateLimiter.RateLimitWithWait
}
