package srv

import (
	"context"
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

func (c *CoreSrv) Executar(ctx context.Context) error {
	for name := range GetRegistraSrvInstance().RegisteredClasses {
		controller, err := GetRegistraSrvInstance().FindSrvClassByKeyAndNewAsObject(name)
		if err != nil {
			return err
		}

		go func(name string, cl IControllerSrv) {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					lib.NewLog().Debug("[Srv]", "Core.Srv.Executar(", name, ") encerrando por shutdown")
					return
				default:
				}
				lib.NewLog().Debug("[Srv]", "Core.Srv.Executar(", name, ")")
				cl.Executar()
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}(name, controller)
	}
	return nil
}
