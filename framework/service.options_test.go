package core

import (
	"testing"

	"github.com/kardianos/service"
)

func TestVerboServico(t *testing.T) {
	casos := []struct {
		args  []string
		verbo string
		ok    bool
	}{
		{[]string{"app"}, "", false},
		{[]string{"app", "install"}, "install", true},
		{[]string{"app", "UNINSTALL"}, "uninstall", true},
		{[]string{"app", "status"}, "status", true},
		{[]string{"app", "-migrate-only"}, "", false},
		{[]string{"app", "qualquer"}, "", false},
		{[]string{"app", "install", "-test.run", "X"}, "", false},
		{[]string{"app.test", "-test.v", "-test.run", "TestX"}, "", false},
	}
	for _, c := range casos {
		v, ok := VerboServico(c.args)
		if v != c.verbo || ok != c.ok {
			t.Errorf("VerboServico(%v) = %q,%v; quer %q,%v", c.args, v, ok, c.verbo, c.ok)
		}
	}
}

func TestAplicarServiceOptions(t *testing.T) {
	defer func() { serviceOptionsPendentes = nil }()
	app := &App{Name: "x"}

	cfg := &service.Config{Name: "x", DisplayName: "padrao"}
	app.aplicarServiceOptions(cfg)
	if cfg.DisplayName != "padrao" {
		t.Fatal("sem options não deve alterar o padrão")
	}

	AddServiceConfig(ServiceOptions{DisplayName: "Meu App", Description: "desc", Restart: "always", InicioManual: true})
	app.aplicarServiceOptions(cfg)
	if cfg.DisplayName != "Meu App" || cfg.Description != "desc" {
		t.Fatalf("options não aplicadas: %+v", cfg)
	}
	if cfg.Option["Restart"] != "always" || cfg.Option["StartType"] != "manual" {
		t.Fatalf("Option incorreto: %+v", cfg.Option)
	}
}
