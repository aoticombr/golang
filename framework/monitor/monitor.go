// Package monitor expõe um servidor HTTP de observabilidade com login,
// gráficos em tempo real (via SSE), histórico em ring buffer, motor de
// alertas configurável e endpoints pprof — tudo plugável via config.json.
package monitor

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/aoticombr/golang/config"
	"github.com/aoticombr/golang/lib"
)

type Monitor struct {
	cfg       *config.Config // referência ao config global (para Save)
	mcfg      *config.Monitor
	appName   string
	server    *http.Server
	collector *Collector
	alerts    *Alerts
	auth      *Auth

	startedAt time.Time
	hostname  string

	// fan-out de novas amostras para SSE
	subsMu sync.RWMutex
	subs   map[chan Snapshot]struct{}

	stopCh chan struct{}
}

// New monta o monitor mas não inicia o servidor; Start faz isso.
func New(cfg *config.Config, appName string) (*Monitor, error) {
	if cfg == nil || cfg.Monitor == nil {
		return nil, errors.New("config.Monitor não definido")
	}
	mc := cfg.Monitor
	if !mc.Ativo {
		return nil, errors.New("monitor desativado")
	}

	host, _ := os.Hostname()

	m := &Monitor{
		cfg:       cfg,
		mcfg:      mc,
		appName:   appName,
		startedAt: time.Now(),
		hostname:  host,
		subs:      make(map[chan Snapshot]struct{}),
		stopCh:    make(chan struct{}),
	}

	// Bootstrap de credenciais e segredo de sessão na primeira execução.
	dirty := false
	if mc.Pass == "" {
		hash, err := HashPassword("admin")
		if err != nil {
			return nil, err
		}
		mc.Pass = hash
		if mc.User == "" {
			mc.User = "admin"
		}
		dirty = true
		lib.NewLog().Info("[Monitor]", "Senha padrão 'admin' gerada — troque após o primeiro login.")
	}
	if mc.SessionKey == "" {
		mc.SessionKey = GenerateSecret()
		dirty = true
	}
	if dirty {
		if err := cfg.Save(); err != nil {
			lib.NewLog().Error("[Monitor]", "Falha ao salvar config.json:", err.Error())
		}
	}

	m.collector = NewCollector(mc.RetentionMin, mc.IntervalSec)
	m.alerts = NewAlerts(mc.Alerts)
	m.auth = NewAuth(mc.SessionKey, mc.User, mc.Pass)

	return m, nil
}

// Start inicia o ticker de coleta e o servidor HTTP. Bloqueia até erro.
func (m *Monitor) Start(ctx context.Context) error {
	go m.collectLoop(ctx)

	mux := m.routes()

	addr := net.JoinHostPort(m.mcfg.Host, m.mcfg.GetPortStr())
	if m.mcfg.Host == "" {
		addr = ":" + m.mcfg.GetPortStr()
	}

	m.server = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
		// Sem WriteTimeout: SSE precisa manter conexões longas abertas.
	}

	lib.NewLog().Info("[Monitor]", "Dashboard em http://"+addr+"/")

	err := m.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Shutdown faz o stop gracioso: fecha SSE subscribers e o servidor.
func (m *Monitor) Shutdown(ctx context.Context) error {
	close(m.stopCh)

	m.subsMu.Lock()
	for ch := range m.subs {
		close(ch)
		delete(m.subs, ch)
	}
	m.subsMu.Unlock()

	if m.server != nil {
		return m.server.Shutdown(ctx)
	}
	return nil
}

func (m *Monitor) collectLoop(ctx context.Context) {
	interval := time.Duration(m.mcfg.IntervalSec) * time.Second
	if interval <= 0 {
		interval = 2 * time.Second
	}
	t := time.NewTicker(interval)
	defer t.Stop()

	// primeira coleta imediata para preencher Last()
	m.tick()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopCh:
			return
		case <-t.C:
			m.tick()
		}
	}
}

func (m *Monitor) tick() {
	s := m.collector.Sample()
	hist := m.collector.History()
	fired := m.alerts.Evaluate(s, hist)

	snap := m.snapshot(s, fired)
	m.broadcast(snap)
}

func (m *Monitor) snapshot(s Sample, fired []Alert) Snapshot {
	return Snapshot{
		Sample:     s,
		Uptime:     int64(time.Since(m.startedAt).Seconds()),
		GoVersion:  runtime.Version(),
		NumCPU:     runtime.NumCPU(),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
		PID:        os.Getpid(),
		Hostname:   m.hostname,
		AppName:    m.appName,
		Alerts:     fired,
	}
}

func (m *Monitor) subscribe() chan Snapshot {
	ch := make(chan Snapshot, 8)
	m.subsMu.Lock()
	m.subs[ch] = struct{}{}
	m.subsMu.Unlock()
	return ch
}

func (m *Monitor) unsubscribe(ch chan Snapshot) {
	m.subsMu.Lock()
	if _, ok := m.subs[ch]; ok {
		delete(m.subs, ch)
		close(ch)
	}
	m.subsMu.Unlock()
}

// broadcast envia snapshot para todos os subscribers; um subscriber lento
// não bloqueia os demais (drop-on-full).
func (m *Monitor) broadcast(s Snapshot) {
	m.subsMu.RLock()
	defer m.subsMu.RUnlock()
	for ch := range m.subs {
		select {
		case ch <- s:
		default:
			// canal cheio — descarta para não atrasar coletor nem demais clientes
		}
	}
}
