package monitor

import (
	"runtime"
	"sync"
	"time"
)

// Sample é uma fotografia do estado do runtime em um instante.
type Sample struct {
	T            int64   `json:"t"`            // unix millis
	Goroutines   int     `json:"goroutines"`
	Threads      int     `json:"threads"`
	CGoCalls     int64   `json:"cgoCalls"`
	HeapAlloc    uint64  `json:"heapAlloc"`
	HeapInuse    uint64  `json:"heapInuse"`
	HeapIdle     uint64  `json:"heapIdle"`
	HeapSys      uint64  `json:"heapSys"`
	HeapReleased uint64  `json:"heapReleased"`
	HeapObjects  uint64  `json:"heapObjects"`
	StackInuse   uint64  `json:"stackInuse"`
	StackSys     uint64  `json:"stackSys"`
	Sys          uint64  `json:"sys"`
	TotalAlloc   uint64  `json:"totalAlloc"`
	NumGC        uint32  `json:"numGC"`
	GCPauseMs    float64 `json:"gcPauseMs"`    // última pausa
	GCCPUPercent float64 `json:"gcCpuPercent"` // % de CPU em GC desde o início
	AllocRate    float64 `json:"allocRate"`    // bytes/s alocados desde o sample anterior
	NextGC       uint64  `json:"nextGC"`
	LastGCMs     int64   `json:"lastGCMs"` // ms desde a última GC
}

type Snapshot struct {
	Sample
	Uptime     int64    `json:"uptime"`
	GoVersion  string   `json:"goVersion"`
	NumCPU     int      `json:"numCPU"`
	GOMAXPROCS int      `json:"gomaxprocs"`
	PID        int      `json:"pid"`
	Hostname   string   `json:"hostname"`
	AppName    string   `json:"appName"`
	Alerts     []Alert  `json:"alerts,omitempty"`
}

// Collector coleta amostras do runtime numa cadência fixa e mantém um
// ring buffer com o histórico recente. É thread-safe para leituras
// concorrentes (SSE, /history) enquanto o ticker grava novos pontos.
type Collector struct {
	mu        sync.RWMutex
	buf       []Sample
	cap       int
	head      int // próxima posição a escrever
	size      int // quantos válidos no buffer
	last      Sample
	prevTotal uint64
	prevT     time.Time
}

func NewCollector(retentionMin, intervalSec int) *Collector {
	if intervalSec <= 0 {
		intervalSec = 2
	}
	if retentionMin <= 0 {
		retentionMin = 60
	}
	cap := (retentionMin * 60) / intervalSec
	if cap < 60 {
		cap = 60
	}
	return &Collector{
		buf: make([]Sample, cap),
		cap: cap,
	}
}

// Sample lê o runtime e adiciona uma amostra ao buffer.
func (c *Collector) Sample() Sample {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	now := time.Now()
	threads, _ := runtime.ThreadCreateProfile(nil)

	var lastPause float64
	if ms.NumGC > 0 {
		idx := (ms.NumGC + 255) % 256
		lastPause = float64(ms.PauseNs[idx]) / 1e6
	}

	var lastGCMs int64
	if ms.LastGC > 0 {
		lastGCMs = time.Since(time.Unix(0, int64(ms.LastGC))).Milliseconds()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var allocRate float64
	if !c.prevT.IsZero() {
		dt := now.Sub(c.prevT).Seconds()
		if dt > 0 {
			allocRate = float64(ms.TotalAlloc-c.prevTotal) / dt
		}
	}
	c.prevT = now
	c.prevTotal = ms.TotalAlloc

	s := Sample{
		T:            now.UnixMilli(),
		Goroutines:   runtime.NumGoroutine(),
		Threads:      threads,
		CGoCalls:     runtime.NumCgoCall(),
		HeapAlloc:    ms.HeapAlloc,
		HeapInuse:    ms.HeapInuse,
		HeapIdle:     ms.HeapIdle,
		HeapSys:      ms.HeapSys,
		HeapReleased: ms.HeapReleased,
		HeapObjects:  ms.HeapObjects,
		StackInuse:   ms.StackInuse,
		StackSys:     ms.StackSys,
		Sys:          ms.Sys,
		TotalAlloc:   ms.TotalAlloc,
		NumGC:        ms.NumGC,
		GCPauseMs:    lastPause,
		GCCPUPercent: ms.GCCPUFraction * 100,
		AllocRate:    allocRate,
		NextGC:       ms.NextGC,
		LastGCMs:     lastGCMs,
	}

	c.buf[c.head] = s
	c.head = (c.head + 1) % c.cap
	if c.size < c.cap {
		c.size++
	}
	c.last = s

	return s
}

// History retorna as amostras na ordem cronológica (mais antiga → mais nova).
func (c *Collector) History() []Sample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Sample, c.size)
	if c.size < c.cap {
		copy(out, c.buf[:c.size])
		return out
	}
	// buffer cheio: começa em head e dá a volta
	copy(out, c.buf[c.head:])
	copy(out[c.cap-c.head:], c.buf[:c.head])
	return out
}

// Last retorna a amostra mais recente, ou zero-value se ainda não houve coleta.
func (c *Collector) Last() Sample {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.last
}

// MetricValue extrai um valor numérico nomeado de uma Sample. Usado pelas
// regras de alerta para resolver "metric": "goroutines" → s.Goroutines.
func MetricValue(s Sample, name string) (float64, bool) {
	switch name {
	case "goroutines":
		return float64(s.Goroutines), true
	case "threads":
		return float64(s.Threads), true
	case "heapAlloc":
		return float64(s.HeapAlloc), true
	case "heapInuse":
		return float64(s.HeapInuse), true
	case "heapIdle":
		return float64(s.HeapIdle), true
	case "heapSys":
		return float64(s.HeapSys), true
	case "heapObjects":
		return float64(s.HeapObjects), true
	case "stackInuse":
		return float64(s.StackInuse), true
	case "sys":
		return float64(s.Sys), true
	case "numGC":
		return float64(s.NumGC), true
	case "gcPauseMs":
		return s.GCPauseMs, true
	case "gcCpuPercent":
		return s.GCCPUPercent, true
	case "allocRate":
		return s.AllocRate, true
	}
	return 0, false
}
