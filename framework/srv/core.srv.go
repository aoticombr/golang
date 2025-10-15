package srv

import (
	"time"

	"github.com/aoticombr/golang/config"
	lib "github.com/aoticombr/golang/lib"
)

type CoreSrv struct {
	Srv   *config.Service
	Certs []*config.Cert
	Dbs   []*config.Database
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
