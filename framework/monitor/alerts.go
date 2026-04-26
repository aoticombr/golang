package monitor

import (
	"strconv"
	"sync"
	"time"

	"github.com/aoticombr/golang/config"
)

type Alert struct {
	T        int64   `json:"t"`
	Rule     string  `json:"rule"`
	Metric   string  `json:"metric"`
	Severity string  `json:"severity"` // info, warn, critical
	Value    float64 `json:"value"`
	Message  string  `json:"message"`
}

const alertHistoryCap = 200

// alertGate evita que a mesma regra dispare repetidamente em sequência. Só
// permite reentrada quando a condição volta a falsear no meio do caminho.
type alertGate struct {
	firing  bool
	lastFmt int64 // unix ms do último disparo
}

type Alerts struct {
	mu       sync.RWMutex
	rules    []config.MonitorRule
	gates    map[string]*alertGate
	history  []Alert
	hHead    int
	hSize    int
	listener func(Alert)
}

func NewAlerts(rules []config.MonitorRule) *Alerts {
	return &Alerts{
		rules:   rules,
		gates:   make(map[string]*alertGate),
		history: make([]Alert, alertHistoryCap),
	}
}

// SetRules permite atualizar as regras em runtime (ex: após edição do config.json).
func (a *Alerts) SetRules(rules []config.MonitorRule) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rules = rules
	a.gates = make(map[string]*alertGate)
}

// SetListener registra um callback chamado a cada novo disparo.
func (a *Alerts) SetListener(fn func(Alert)) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.listener = fn
}

// Evaluate roda todas as regras contra a amostra mais recente. Para regras
// do tipo "growth" usa o histórico para comparar contra um ponto de
// referência dentro da janela.
func (a *Alerts) Evaluate(latest Sample, history []Sample) []Alert {
	a.mu.Lock()
	rules := a.rules
	a.mu.Unlock()

	var fired []Alert
	for _, r := range rules {
		if !r.Ativo {
			continue
		}
		val, ok := MetricValue(latest, r.Metric)
		if !ok {
			continue
		}

		breached, msg := evalRule(r, val, latest, history)

		a.mu.Lock()
		g, ok := a.gates[r.Name]
		if !ok {
			g = &alertGate{}
			a.gates[r.Name] = g
		}
		if breached {
			if !g.firing {
				g.firing = true
				al := Alert{
					T:        latest.T,
					Rule:     r.Name,
					Metric:   r.Metric,
					Severity: defaultSeverity(r.Severity),
					Value:    val,
					Message:  msg,
				}
				g.lastFmt = al.T
				a.appendHistoryLocked(al)
				fired = append(fired, al)
				if a.listener != nil {
					l := a.listener
					a.mu.Unlock()
					l(al)
					a.mu.Lock()
				}
			}
		} else {
			g.firing = false
		}
		a.mu.Unlock()
	}

	return fired
}

func (a *Alerts) appendHistoryLocked(al Alert) {
	a.history[a.hHead] = al
	a.hHead = (a.hHead + 1) % alertHistoryCap
	if a.hSize < alertHistoryCap {
		a.hSize++
	}
}

// History retorna alertas em ordem cronológica decrescente (mais recente primeiro).
func (a *Alerts) History() []Alert {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Alert, a.hSize)
	for i := 0; i < a.hSize; i++ {
		idx := (a.hHead - 1 - i + alertHistoryCap) % alertHistoryCap
		out[i] = a.history[idx]
	}
	return out
}

func defaultSeverity(s string) string {
	switch s {
	case "info", "warn", "critical":
		return s
	}
	return "warn"
}

func evalRule(r config.MonitorRule, val float64, latest Sample, history []Sample) (bool, string) {
	switch r.Op {
	case ">":
		if val > r.Threshold {
			return true, fmtMsg(r, val)
		}
	case ">=":
		if val >= r.Threshold {
			return true, fmtMsg(r, val)
		}
	case "<":
		if val < r.Threshold {
			return true, fmtMsg(r, val)
		}
	case "<=":
		if val <= r.Threshold {
			return true, fmtMsg(r, val)
		}
	case "growth":
		if r.WindowSec <= 0 {
			return false, ""
		}
		ref := findRefInWindow(history, latest.T, r.WindowSec, r.Metric)
		if ref == 0 {
			return false, ""
		}
		pct := (val - ref) / ref * 100
		if pct >= r.Threshold {
			return true, fmtGrowthMsg(r, val, ref, pct)
		}
	}
	return false, ""
}

func findRefInWindow(history []Sample, nowMs int64, windowSec int, metric string) float64 {
	cutoff := nowMs - int64(windowSec)*1000
	for _, s := range history {
		if s.T >= cutoff {
			if v, ok := MetricValue(s, metric); ok {
				return v
			}
			return 0
		}
	}
	return 0
}

func fmtMsg(r config.MonitorRule, v float64) string {
	return r.Name + ": " + r.Metric + " " + r.Op + " " + fmtNum(r.Threshold) + " (atual: " + fmtNum(v) + ")"
}

func fmtGrowthMsg(r config.MonitorRule, cur, ref, pct float64) string {
	return r.Name + ": " + r.Metric + " cresceu " + fmtNum(pct) + "% em " + fmtSec(r.WindowSec) + " (de " + fmtNum(ref) + " para " + fmtNum(cur) + ")"
}

func fmtNum(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func fmtSec(s int) string {
	d := time.Duration(s) * time.Second
	return d.String()
}
