package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/pprof"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/aoticombr/golang/config"
	"github.com/aoticombr/golang/lib"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

func (m *Monitor) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(m.recoverMiddleware)

	// Estáticos públicos (CSS/JS de login).
	r.Get("/static/*", m.staticHandler)

	// Login público.
	r.Get("/login", m.loginPage)
	r.Post("/login", m.loginPost)
	r.Post("/logout", m.logoutPost)

	// Tudo o resto protegido.
	r.Group(func(pr chi.Router) {
		pr.Use(m.auth.Middleware)

		pr.Get("/", m.dashboardPage)

		pr.Get("/api/snapshot", m.snapshotHandler)
		pr.Get("/api/history", m.historyHandler)
		pr.Get("/api/stream", m.streamHandler)
		pr.Get("/api/alerts", m.alertsHandler)
		pr.Get("/api/rules", m.rulesGet)
		pr.Put("/api/rules", m.rulesPut)
		pr.Post("/api/password", m.passwordPost)
		pr.Post("/api/gc", m.gcHandler)
		pr.Post("/api/freeosmemory", m.freeOSHandler)
		pr.Get("/api/goroutines", m.goroutineDumpHandler)

		pr.Get("/debug/pprof/", http.HandlerFunc(pprof.Index))
		pr.Get("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		pr.Get("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		pr.Get("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		pr.Get("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		pr.Get("/debug/pprof/{name}", m.pprofByName)
	})

	return r
}

func (m *Monitor) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				lib.NewLog().Error("[Monitor]", "panic:", rec, "\n", string(debug.Stack()))
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Monitor) staticHandler(w http.ResponseWriter, r *http.Request) {
	sub, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		http.Error(w, "assets indisponíveis", http.StatusInternalServerError)
		return
	}
	prefix := strings.TrimPrefix(r.URL.Path, "/static/")
	f, err := sub.Open(prefix)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if stat.IsDir() {
		http.NotFound(w, r)
		return
	}
	rs, ok := f.(readSeeker)
	if !ok {
		http.Error(w, "asset não seekable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType(prefix))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeContent(w, r, prefix, stat.ModTime(), rs)
}

type readSeeker interface {
	io.Reader
	io.Seeker
}

func contentType(name string) string {
	switch {
	case strings.HasSuffix(name, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(name, ".js"):
		return "application/javascript; charset=utf-8"
	case strings.HasSuffix(name, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(name, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(name, ".png"):
		return "image/png"
	case strings.HasSuffix(name, ".ico"):
		return "image/x-icon"
	}
	return "application/octet-stream"
}

func (m *Monitor) loginPage(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.auth.Authenticated(r); ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	m.serveAsset(w, r, "login.html", "text/html; charset=utf-8")
}

func (m *Monitor) loginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	if err := m.auth.Login(w, r, user, pass); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		// Mostra erro inline; mantém o usuário digitado.
		fmt.Fprintf(w, loginErrorPage, escapeHTML(err.Error()), escapeHTML(user))
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (m *Monitor) logoutPost(w http.ResponseWriter, r *http.Request) {
	m.auth.Logout(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (m *Monitor) dashboardPage(w http.ResponseWriter, r *http.Request) {
	m.serveAsset(w, r, "index.html", "text/html; charset=utf-8")
}

func (m *Monitor) serveAsset(w http.ResponseWriter, r *http.Request, name, ct string) {
	data, err := assetsFS.ReadFile("assets/" + name)
	if err != nil {
		http.Error(w, "asset não encontrado: "+name, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ct)
	w.Write(data)
}

func (m *Monitor) snapshotHandler(w http.ResponseWriter, r *http.Request) {
	s := m.collector.Last()
	snap := m.snapshot(s, nil)
	writeJSON(w, snap)
}

func (m *Monitor) historyHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, m.collector.History())
}

func (m *Monitor) alertsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, m.alerts.History())
}

func (m *Monitor) rulesGet(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, m.mcfg.Alerts)
}

// rulesPut substitui o conjunto de regras e persiste no config.json.
func (m *Monitor) rulesPut(w http.ResponseWriter, r *http.Request) {
	var rules []config.MonitorRule
	if err := json.NewDecoder(r.Body).Decode(&rules); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.mcfg.Alerts = rules
	m.alerts.SetRules(rules)
	if err := m.cfg.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"ok": true, "count": len(rules)})
}

func (m *Monitor) passwordPost(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Current string `json:"current"`
		New     string `json:"new"`
		User    string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(body.New) < 4 {
		http.Error(w, "nova senha muito curta (mínimo 4 caracteres)", http.StatusBadRequest)
		return
	}
	if err := m.auth.verifyPassword(m.mcfg.User, body.Current); err != nil {
		http.Error(w, "senha atual incorreta", http.StatusUnauthorized)
		return
	}
	hash, err := HashPassword(body.New)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m.mcfg.Pass = hash
	if body.User != "" {
		m.mcfg.User = body.User
	}
	m.auth.Update(m.mcfg.User, m.mcfg.Pass)
	if err := m.cfg.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"ok": true})
}

func (m *Monitor) gcHandler(w http.ResponseWriter, r *http.Request) {
	runtime.GC()
	writeJSON(w, map[string]any{"ok": true})
}

func (m *Monitor) freeOSHandler(w http.ResponseWriter, r *http.Request) {
	debug.FreeOSMemory()
	writeJSON(w, map[string]any{"ok": true})
}

func (m *Monitor) goroutineDumpHandler(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(buf[:n])
}

func (m *Monitor) pprofByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	pprof.Handler(name).ServeHTTP(w, r)
}

// streamHandler implementa SSE: envia 'history' inicial e depois cada
// snapshot novo. Mantém-se aberto até o cliente desconectar.
func (m *Monitor) streamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming não suportado", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // proxy-friendly

	enc := json.NewEncoder(w)

	// snapshot inicial para preencher o dashboard imediatamente
	initial := m.snapshot(m.collector.Last(), nil)
	fmt.Fprint(w, "event: snapshot\ndata: ")
	enc.Encode(initial)
	fmt.Fprint(w, "\n")
	flusher.Flush()

	ch := m.subscribe()
	defer m.unsubscribe(ch)

	keepalive := time.NewTicker(20 * time.Second)
	defer keepalive.Stop()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-keepalive.C:
			fmt.Fprint(w, ": keepalive\n\n")
			flusher.Flush()
		case snap, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprint(w, "event: snapshot\ndata: ")
			if err := enc.Encode(snap); err != nil {
				return
			}
			fmt.Fprint(w, "\n")
			flusher.Flush()
		}
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

func escapeHTML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	)
	return r.Replace(s)
}

// verifyPassword expõe a verificação de bcrypt sem efeitos colaterais —
// usado pela troca de senha.
func (a *Auth) verifyPassword(user, pass string) error {
	a.mu.Lock()
	expectedUser := a.user
	expectedHash := a.passHash
	a.mu.Unlock()
	if !subtleEqual(user, expectedUser) {
		return errors.New("usuário inválido")
	}
	return bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(pass))
}
