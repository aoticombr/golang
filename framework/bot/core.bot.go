package bot

import (
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

func (c *CoreBot) Executar() error {
	for name, bot := range GetRegistraBotInstance().RegisteredClasses {
		Controller, err := GetRegistraBotInstance().FindBotClassByKeyAndNewAsObject(name)
		if err != nil {
			return err
		}

		go func(name string, tempo int, cl IControllerBot) {

			for {
				lib.NewLog().Debug("[Bot]", "Core.Bot.Executar(", name, ")", "Ciclo:", tempo, "minutos")
				go cl.Executar()
				time.Sleep(time.Duration(tempo) * time.Minute)
			}

		}(name, bot.Tempo, Controller)

	}

	return nil
}
