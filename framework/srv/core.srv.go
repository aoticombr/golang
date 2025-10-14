package srv

import (
	"time"

	"github.com/aoticombr/golang/config"
	lib "github.com/aoticombr/golang/framework/lib"
)

type CoreSrv struct {
	Srv   *config.Service
	Certs []*config.Cert
	Dbs   []*config.Database
}

func (srv *CoreSrv) LogDebug(v ...interface{}) {
	lib.NewLog().Debug(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogInfo(v ...interface{}) {
	lib.NewLog().Info(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogError(v ...interface{}) {
	lib.NewLog().Error(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogWarning(v ...interface{}) {
	lib.NewLog().Warning(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogCritical(v ...interface{}) {
	lib.NewLog().Critical(srv.Srv.Name, v...)
}
func (srv *CoreSrv) LogFatal(v ...interface{}) {
	lib.NewLog().Fatal(srv.Srv.Name, v...)
}

func NewCoreSrv(certs []*config.Cert, dbs []*config.Database, service *config.Service) *CoreSrv {
	srv := &CoreSrv{
		Certs: certs,
		Dbs:   dbs,
		Srv:   service,
	}
	return srv
}

func (c *CoreSrv) Executar() error {
	for name := range GetRegistraSrvInstance().RegisteredClasses {
		Controller, err := GetRegistraSrvInstance().FindSrvClassByKeyAndNewAsObject(name)
		if err != nil {
			return err
		}

		go func(name string, cl IControllerSrv) {

			for {
				lib.NewLog().Debug("[Srv]", "Core.Srv.Executar(", name, ")")
				go cl.Executar()
				time.Sleep(1 * time.Minute)
			}

		}(name, Controller)

	}
	return nil
}
