package api

import (
	"encoding/json"
	"net/http"

	"github.com/aoticombr/golang/config"
	"github.com/aoticombr/golang/dbconndataset"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type CoreApi struct {
	Api        *config.Api
	Certs      config.Certs
	Dbs        []*config.Database
	Parans     config.Params
	Connection *dbconndataset.ConnDataSet
	Processo   *Processo
}

func NewCoreApi(
	certs []*config.Cert,
	dbs []*config.Database,
	app *config.Api,
	parans config.Params) *CoreApi {
	api := &CoreApi{
		Certs:    certs,
		Dbs:      dbs,
		Api:      app,
		Parans:   parans,
		Processo: &Processo{},
	}

	return api
}

func (ca *CoreApi) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{

		AllowedOrigins:   ca.Api.Cors.AllowOrigins,     // Allow only this origin - http://localhost:3000
		AllowedMethods:   ca.Api.Cors.AllowMethods,     //[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   ca.Api.Cors.AllowHeaders,     //[]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   ca.Api.Cors.ExposedHeaders,   //[]string{"Link"},
		AllowCredentials: ca.Api.Cors.AllowCredentials, //true,
		MaxAge:           ca.Api.Cors.MaxAge,           //300, // Maximum value not to check for CORS request again (in seconds)
	}))

	RouteGroups := GetRegistraInstance().RouteGroups
	for _, RouteGroup := range RouteGroups {
		ca.LogInfo("#########################################")
		ca.LogInfo(RouteGroup.Name)
		ca.LogInfo("#########################################")
		r.Group(func(r chi.Router) {
			for _, Middleware := range RouteGroup.Middleware {
				r.Use(Middleware)
			}
			for _, Route := range RouteGroup.Routes {
				Pattern := "/{db}" + Route.Pattern
				PatterComplet := ""
				if ca.Api.Https.Ativo {
					PatterComplet += "https://"
				} else {
					PatterComplet += "http://"
				}
				PatterComplet += ca.Api.Host
				PatterComplet += ":" + ca.Api.GetPortStr()
				PatterComplet += "/{db}" + Route.Pattern
				switch Route.Method {
				case GET:
					r.Get(Pattern, Route.HandlerFunc)
					ca.LogInfo("GET", PatterComplet)
				case POST:
					r.Post(Pattern, Route.HandlerFunc)
					ca.LogInfo("POST", PatterComplet)
				case PUT:
					r.Put(Pattern, Route.HandlerFunc)
					ca.LogInfo("PUT", PatterComplet)
				case DELETE:
					r.Delete(Pattern, Route.HandlerFunc)
					ca.LogInfo("DELETE", PatterComplet)
				case PATCH:
					r.Patch(Pattern, Route.HandlerFunc)
					ca.LogInfo("PATCH", PatterComplet)
				case WEBSOCKET:
					r.HandleFunc(Pattern, Route.HandlerFunc)
					ca.LogInfo("WEBSOCKET", PatterComplet)
				}
			}
		})
	}

	// Inicie o servidor
	http.ListenAndServe(":"+ca.Api.GetPortStr(), r)

	return nil
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, code int, body []byte, err error) {
	if (code >= 200) && (code <= 299) {
		w.WriteHeader(code)
		if body != nil {
			w.Write(body)
		}
	} else {
		w.WriteHeader(code)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		} else {
			w.Write(body)
		}

	}
}
