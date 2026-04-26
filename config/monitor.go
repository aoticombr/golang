package config

import "strconv"

type Monitor struct {
	Ativo        bool          `json:"ativo"`
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	User         string        `json:"user"`
	Pass         string        `json:"pass"`
	SessionKey   string        `json:"sessionKey"`
	IntervalSec  int           `json:"intervalSec"`
	RetentionMin int           `json:"retentionMin"`
	Alerts       []MonitorRule `json:"alerts"`
}

type MonitorRule struct {
	Name      string  `json:"name"`
	Metric    string  `json:"metric"`
	Op        string  `json:"op"`
	Threshold float64 `json:"threshold"`
	WindowSec int     `json:"windowSec"`
	Severity  string  `json:"severity"`
	Ativo     bool    `json:"ativo"`
}

func (m *Monitor) GetPortStr() string {
	return strconv.Itoa(m.Port)
}
