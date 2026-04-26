package monitor

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	cookieName     = "aoti_mon"
	sessionTTL     = 12 * time.Hour
	loginThrottle  = 1500 * time.Millisecond // delay anti-bruteforce em falha
	maxFailures    = 8
	failureWindow  = 5 * time.Minute
	lockoutPeriod  = 15 * time.Minute
)

// Auth gerencia autenticação por cookie HMAC e bloqueio por tentativas.
// O segredo HMAC vem do config (gerado se vazio na primeira execução).
type Auth struct {
	mu        sync.Mutex
	secret    []byte
	user      string
	passHash  string
	failures  map[string]*failureCounter
}

type failureCounter struct {
	count    int
	first    time.Time
	lockedAt time.Time
}

func NewAuth(secret, user, passHash string) *Auth {
	return &Auth{
		secret:   []byte(secret),
		user:     user,
		passHash: passHash,
		failures: make(map[string]*failureCounter),
	}
}

// Update troca usuário/senha em runtime após alteração via dashboard.
func (a *Auth) Update(user, passHash string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.user = user
	a.passHash = passHash
}

// User retorna o usuário corrente (para exibir no header).
func (a *Auth) User() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.user
}

// HashPassword gera o bcrypt para gravar no config.json.
func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", errors.New("senha vazia")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// isBcryptHash detecta se a string já está no formato bcrypt
// (prefixos $2a$, $2b$, $2y$ — variantes geradas por libs distintas).
func isBcryptHash(s string) bool {
	if len(s) < 60 {
		return false
	}
	return strings.HasPrefix(s, "$2a$") ||
		strings.HasPrefix(s, "$2b$") ||
		strings.HasPrefix(s, "$2y$")
}

// GenerateSecret produz uma chave hex de 32 bytes para assinar sessões.
func GenerateSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return time.Now().Format(time.RFC3339Nano) // fallback fraco mas determinístico
	}
	return hex.EncodeToString(b)
}

// Login valida user/pass, aplica throttling e bloqueio. Devolve o cookie
// pronto para Set-Cookie em caso de sucesso.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request, user, pass string) error {
	ip := clientIP(r)

	a.mu.Lock()
	if a.locked(ip) {
		a.mu.Unlock()
		return errors.New("muitas tentativas — tente novamente em alguns minutos")
	}
	expectedUser := a.user
	expectedHash := a.passHash
	a.mu.Unlock()

	if expectedHash == "" {
		return errors.New("monitor sem senha configurada")
	}

	userOK := subtleEqual(user, expectedUser)
	passOK := bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(pass)) == nil

	if !userOK || !passOK {
		time.Sleep(loginThrottle)
		a.recordFailure(ip)
		return errors.New("usuário ou senha inválidos")
	}

	a.clearFailures(ip)
	a.setSessionCookie(w, r, user)
	return nil
}

func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   r.TLS != nil,
	})
}

// Authenticated verifica o cookie de sessão e retorna o usuário se válido.
func (a *Auth) Authenticated(r *http.Request) (string, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil || c.Value == "" {
		return "", false
	}
	parts := strings.SplitN(c.Value, ".", 3)
	if len(parts) != 3 {
		return "", false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "", false
	}
	a.mu.Lock()
	mac := hmac.New(sha256.New, a.secret)
	mac.Write(payload)
	mac.Write([]byte("."))
	mac.Write([]byte(parts[1]))
	expected := mac.Sum(nil)
	expectedUser := a.user
	a.mu.Unlock()

	if !hmac.Equal(sig, expected) {
		return "", false
	}
	expiry, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", false
	}
	if time.Now().Unix() > expiry {
		return "", false
	}
	user := string(payload)
	if user != expectedUser {
		// usuário foi renomeado — invalida a sessão antiga
		return "", false
	}
	return user, true
}

// Middleware redireciona para /login as requisições de páginas, e devolve
// 401 JSON para chamadas a /api/*. Páginas autenticadas ficam disponíveis
// para o handler envolvido.
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := a.Authenticated(r); ok {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/debug/pprof") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func (a *Auth) setSessionCookie(w http.ResponseWriter, r *http.Request, user string) {
	expiry := strconv.FormatInt(time.Now().Add(sessionTTL).Unix(), 10)
	payload := base64.RawURLEncoding.EncodeToString([]byte(user))

	a.mu.Lock()
	mac := hmac.New(sha256.New, a.secret)
	mac.Write([]byte(user))
	mac.Write([]byte("."))
	mac.Write([]byte(expiry))
	sig := mac.Sum(nil)
	a.mu.Unlock()

	value := payload + "." + expiry + "." + base64.RawURLEncoding.EncodeToString(sig)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(sessionTTL),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   r.TLS != nil,
	})
}

func (a *Auth) locked(ip string) bool {
	f, ok := a.failures[ip]
	if !ok {
		return false
	}
	if !f.lockedAt.IsZero() && time.Since(f.lockedAt) < lockoutPeriod {
		return true
	}
	if !f.lockedAt.IsZero() && time.Since(f.lockedAt) >= lockoutPeriod {
		delete(a.failures, ip)
	}
	return false
}

func (a *Auth) recordFailure(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	f, ok := a.failures[ip]
	if !ok || time.Since(f.first) > failureWindow {
		a.failures[ip] = &failureCounter{count: 1, first: time.Now()}
		return
	}
	f.count++
	if f.count >= maxFailures {
		f.lockedAt = time.Now()
	}
}

func (a *Auth) clearFailures(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.failures, ip)
}

func clientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		if i := strings.IndexByte(xf, ','); i > 0 {
			return strings.TrimSpace(xf[:i])
		}
		return strings.TrimSpace(xf)
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i > 0 {
		host = host[:i]
	}
	return host
}

// subtleEqual compara em tempo constante para o lado do usuário (a senha já
// passa pelo bcrypt que é constante por design).
func subtleEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
