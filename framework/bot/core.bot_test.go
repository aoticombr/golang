package bot

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aoticombr/golang/lib"
)

var chamadasCtx atomic.Int32
var chamadasClassico atomic.Int32

type botCtx struct{}

func (b *botCtx) ExecutarCtx(ctx context.Context) {
	if chamadasCtx.Add(1) == 1 {
		panic("primeiro ciclo falha") // o runner deve recuperar e seguir
	}
}

type botClassico struct{}

func (b *botClassico) Executar() { chamadasClassico.Add(1) }

func TestCiclo(t *testing.T) {
	if (BotClass{Tempo: 5}).Ciclo() != 5*time.Minute {
		t.Fatal("Tempo em minutos")
	}
	if (BotClass{Tempo: 5, Intervalo: 10 * time.Second}).Ciclo() != 10*time.Second {
		t.Fatal("Intervalo tem precedência")
	}
	if (BotClass{}).Ciclo() != time.Minute {
		t.Fatal("zero deve virar 1 minuto (antes: panic no NewTicker)")
	}
}

func TestExecutarComCtxERecover(t *testing.T) {
	lib.NewLog().LogNivel = lib.N_NENHUM // não gravar arquivo de log durante o teste
	reg := GetRegistraBotInstance()
	antigo := reg.RegisteredClasses
	reg.RegisteredClasses = make(RegisterBotClasses)
	defer func() { reg.RegisteredClasses = antigo }()

	RegistrarBotIntervalo("ctx", botCtx{}, 20*time.Millisecond)
	RegistrarBotIntervalo("classico", botClassico{}, 20*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	c := &CoreBot{}
	if err := c.Executar(ctx); err != nil {
		t.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	cancel()

	if chamadasCtx.Load() < 2 {
		t.Fatalf("bot com ctx deveria seguir após panic; chamadas=%d", chamadasCtx.Load())
	}
	if chamadasClassico.Load() < 2 {
		t.Fatalf("bot clássico chamadas=%d", chamadasClassico.Load())
	}
}

type naoEhBot struct{}

func TestBotInvalido(t *testing.T) {
	reg := GetRegistraBotInstance()
	antigo := reg.RegisteredClasses
	reg.RegisteredClasses = make(RegisterBotClasses)
	defer func() { reg.RegisteredClasses = antigo }()
	RegistrarBot("x", naoEhBot{}, 1)
	if err := (&CoreBot{}).Executar(context.Background()); err == nil {
		t.Fatal("esperava erro para struct sem Executar")
	}
}
