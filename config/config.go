package config

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/aoticombr/golang/lib"
)

var onceConfig sync.Once
var InstanceConfig *Config

type Config struct {
	Dbs      []*Database `json:"dbs"`
	Services []*Service  `json:"services"`
	Bots     []*Bot      `json:"bots"`
	Apis     []*Api      `json:"apis"`
	Jwt      []*Jwt      `json:"jwt"`
	Parans   Params      `json:"parans"`
	Certs    Certs       `json:"certs"`
	Log      Log         `json:"log"`
	Path     string      `json:"path"`
	JsonFile string
}

func NewConfig() *Config {
	onceConfig.Do(
		func() {

			var jsonFile *os.File
			var err error

			if lib.FileExists("config.json") {
				jsonFile, err = os.Open("config.json")
			} else {
				InstanceConfig := &Config{
					Dbs: []*Database{
						{
							Name:     "db1",
							Host:     "localhost",
							Port:     1521,
							User:     "usuario",
							Pass:     "senha",
							Schema:   "schema",
							Sid:      "service_name",
							PoolSize: 20,
							MaxConn:  100,
							Ativo:    true,
							Trace: Trace{
								Ativo: false,
								Path:  "trace",
							},
						},
					},
					Bots: []*Bot{
						{
							Name:  lib.GetApplicationName(),
							Dbs:   []string{"db1"},
							Ativo: false,
						},
					},
					Services: []*Service{
						{
							Name:  lib.GetApplicationName(),
							Dbs:   []string{"db1"},
							Ativo: false,
						},
					},

					Apis: []*Api{
						{
							Name: lib.GetApplicationName(),
							Host: "localhost",
							Port: 8080,
							Path: "api",
							Swagger: Swagger{
								AuthJwt: false,
								Ativo:   true,
							},
							Cors: Cors{
								MaxAge:           300,
								AllowCredentials: false,
								AllowHeaders:     []string{"*"},
								ExposedHeaders:   []string{"*"},
								AllowMethods:     []string{"*"},
								AllowOrigins:     []string{"*"},
								Ativo:            true,
							},
							Https: Https{
								Ativo: false,
							},
							Gateway: Gateway{
								Ativo: false,
							},
							Dbs:   []string{"db1"},
							Ativo: true,
						},
					},

					Log: Log{
						Nivel:  5,
						Screen: true,
					},
					Jwt: []*Jwt{
						{
							Name:           "name",
							ExpirationTime: 60,
							SecretKey:      "minha_senha",
						},
					},
					Parans: []*Param{},
					Certs:  []*Cert{},
					Path:   "",
				}
				jsonFile, _ := json.MarshalIndent(InstanceConfig, "", "  ")
				err = os.WriteFile("config.json", jsonFile, 0644)
			}

			if err != nil {
				panic("Erro ao carregar arquivo de configuração config.json. " + err.Error())
			}

			defer jsonFile.Close()

			byteValue, _ := io.ReadAll(jsonFile)

			InstanceConfig = &Config{}

			json.Unmarshal(byteValue, InstanceConfig)
		})

	return InstanceConfig
}

func (c *Config) Save() error {
	jsonFile, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := json.MarshalIndent(c, "", "  ")
	jsonFile.Write(byteValue)

	return nil
}
func (c *Config) AddDb(db Database) {
	c.Dbs = append(c.Dbs, &db)
}
func (c *Config) AddService(service Service) {
	c.Services = append(c.Services, &service)
}
func (c *Config) AddApi(api Api) {
	c.Apis = append(c.Apis, &api)
}
func (c *Config) RemoveDb(db Database) {
	for i, v := range c.Dbs {
		if strings.ToUpper(v.Name) == strings.ToUpper(db.Name) {
			c.Dbs = append(c.Dbs[:i], c.Dbs[i+1:]...)
			break
		}
	}
}
func (c *Config) RemoveService(service Service) {
	for i, v := range c.Services {
		if strings.ToUpper(v.Name) == strings.ToUpper(service.Name) {
			c.Services = append(c.Services[:i], c.Services[i+1:]...)
			break
		}
	}
}
func (c *Config) RemoveApi(api Api) {
	for i, v := range c.Apis {
		if strings.ToUpper(v.Name) == strings.ToUpper(api.Name) {
			c.Apis = append(c.Apis[:i], c.Apis[i+1:]...)
			break
		}
	}
}
func (c *Config) GetDB(name string) (bool, *Database) {
	for _, v := range c.Dbs {
		if strings.EqualFold(v.Name, name) {
			if v.Ativo {
				return true, v
			}
		}
	}
	return false, &Database{}
}

func (c *Config) GetCert(name string) (bool, *Cert) {
	for _, v := range c.Certs {
		if strings.ToUpper(v.Name) == strings.ToUpper(name) {
			if v.Ativo {
				return true, v
			}
		}
	}
	return false, &Cert{}
}

func (c *Config) GetService(name string) (bool, *Service) {
	for _, v := range c.Services {
		if strings.ToUpper(v.Name) == strings.ToUpper(name) {
			if v.Ativo {
				return true, v
			}
		}
	}
	return false, nil
}
func (c *Config) GetBot(name string) (bool, *Bot) {
	for _, v := range c.Bots {
		if strings.ToUpper(v.Name) == strings.ToUpper(name) {
			if v.Ativo {
				return true, v
			}
		}
	}
	return false, nil
}
func (c *Config) GetApi(name string) (bool, *Api) {
	for _, v := range c.Apis {
		if strings.ToUpper(v.Name) == strings.ToUpper(name) {
			if v.Ativo {
				return true, v
			}
		}
	}
	return false, nil
}
func (c *Config) GetJwt(name string) (bool, *Jwt) {
	for _, v := range c.Jwt {
		if strings.ToUpper(v.Name) == strings.ToUpper(name) {

			return true, v

		}
	}
	return false, nil
}
