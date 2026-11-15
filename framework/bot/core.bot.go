package bot

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/aoticombr/golang/config"
	"github.com/aoticombr/golang/lib"
)

type CoreBot struct {
	Bot   *config.Bot
	Certs []*config.Cert
	Dbs   []*config.Database
}

func (bot *CoreBot) LogDebug(v ...interface{}) {
	lib.NewLog().Debug(bot.Bot.Name, v...)
}
func (bot *CoreBot) LogInfo(v ...interface{}) {
	lib.NewLog().Info(bot.Bot.Name, v...)
}
func (bot *CoreBot) LogError(v ...interface{}) {
	lib.NewLog().Error(bot.Bot.Name, v...)
}
func (bot *CoreBot) LogWarning(v ...interface{}) {
	lib.NewLog().Warning(bot.Bot.Name, v...)
}
func (bot *CoreBot) LogCritical(v ...interface{}) {
	lib.NewLog().Critical(bot.Bot.Name, v...)
}
func (bot *CoreBot) LogFatal(v ...interface{}) {
	lib.NewLog().Fatal(bot.Bot.Name, v...)
}

func NewCoreBot(certs []*config.Cert, dbs []*config.Database, bot *config.Bot) *CoreBot {
	srv := &CoreBot{
		Certs: certs,
		Dbs:   dbs,
		Bot:   bot,
	}
	return srv
}

func (c *CoreBot) Executar(ctx context.Context) error {
	for name, bot := range GetRegistraBotInstance().RegisteredClasses {
		instancia, err := GetRegistraBotInstance().NovaInstancia(name)
		if err != nil {
			return err
		}
		_, temCtx := instancia.(IControllerBotCtx)
		_, temClassico := instancia.(IControllerBot)
		if !temCtx && !temClassico {
			return fmt.Errorf("bot %s não implementa Executar() nem ExecutarCtx(ctx) no ponteiro", name)
		}

		go func(name string, ciclo time.Duration, cl any) {
			ticker := time.NewTicker(ciclo)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					lib.NewLog().Debug("[Bot]", "Core.Bot.Executar(", name, ") encerrando por shutdown")
					return
				default:
				}
				lib.NewLog().Debug("[Bot]", "Core.Bot.Executar(", name, ")", "Ciclo:", ciclo.String())
				// Ciclos são sequenciais (o próximo só começa depois que este termina),
				// então não há sobreposição; um panic é registrado e o bot continua.
				executarCiclo(ctx, name, cl)
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}(name, bot.Ciclo(), instancia)
	}

	return nil
}

func executarCiclo(ctx context.Context, name string, cl any) {
	defer func() {
		if r := recover(); r != nil {
			lib.NewLog().Error("[Bot]", "Core.Bot.Executar(", name, ") panic:", fmt.Sprint(r), string(debug.Stack()))
		}
	}()
	if c, ok := cl.(IControllerBotCtx); ok {
		c.ExecutarCtx(ctx)
		return
	}
	cl.(IControllerBot).Executar()
}
