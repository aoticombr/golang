package api

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// handlersRaiz são handlers montados fora do prefixo (ex.: SPA em "/").
var handlersRaiz []handlerRaiz

type handlerRaiz struct {
	pattern string
	handler http.Handler
}

// RegistrarHandlerRaiz monta um handler fora do prefixo, direto no router
// (ex.: RegistrarHandlerRaiz("/*", spa) para servir uma UI web em "/").
// Só é usado por StartPrefixo (AddApiPrefixo).
func RegistrarHandlerRaiz(pattern string, h http.Handler) {
	handlersRaiz = append(handlersRaiz, handlerRaiz{pattern: pattern, handler: h})
}

// StartPrefixo sobe a API usando um prefixo fixo no lugar de "/{db}".
// Diferenças em relação a Start:
//   - todas as rotas ficam em "/<prefixo><pattern>" (prefixo vazio = sem prefixo);
//   - escuta em Host:Port (Host vazio = todas as interfaces);
//   - aplica os timeouts do config somente quando informados (> 0), para não
//     cortar WebSocket/downloads longos quando o JSON não define valores;
//   - inclui middleware.Recoverer (panic em handler vira 500).
func (ca *CoreApi) StartPrefixo(prefixo string) error {
	prefixo = "/" + strings.Trim(prefixo, "/")
	if prefixo == "/" {
		prefixo = ""
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	if ca.Api.Cors.Ativo {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   ca.Api.Cors.AllowOrigins,
			AllowedMethods:   ca.Api.Cors.AllowMethods,
			AllowedHeaders:   ca.Api.Cors.AllowHeaders,
			ExposedHeaders:   ca.Api.Cors.ExposedHeaders,
			AllowCredentials: ca.Api.Cors.AllowCredentials,
		}))
	}

	for _, RouteGroup := range GetRegistraInstance().RouteGroups {
		ca.LogInfo("#########################################")
		ca.LogInfo(RouteGroup.Name)
		ca.LogInfo("#########################################")
		rg := RouteGroup
		r.Group(func(r chi.Router) {
			for _, m := range rg.Middleware {
				r.Use(m)
			}
			for _, Route := range rg.Routes {
				Pattern := prefixo + Route.Pattern
				switch Route.Method {
				case GET:
					r.Get(Pattern, Route.HandlerFunc)
				case POST:
					r.Post(Pattern, Route.HandlerFunc)
				case PUT:
					r.Put(Pattern, Route.HandlerFunc)
				case DELETE:
					r.Delete(Pattern, Route.HandlerFunc)
				case PATCH:
					r.Patch(Pattern, Route.HandlerFunc)
				case WEBSOCKET:
					r.HandleFunc(Pattern, Route.HandlerFunc)
				}
				ca.LogInfo(Route.Method, Pattern)
			}
		})
	}
	for _, h := range handlersRaiz {
		r.Handle(h.pattern, h.handler)
	}

	addr := net.JoinHostPort(ca.Api.Host, ca.Api.GetPortStr())
	t := ca.Api.Timeouts
	srv := &http.Server{Addr: addr, Handler: r, ReadHeaderTimeout: 30 * time.Second}
	if t.ReadHeaderTimeout > 0 {
		srv.ReadHeaderTimeout = time.Duration(t.ReadHeaderTimeout) * time.Second
	}
	if t.ReadTimeout > 0 {
		srv.ReadTimeout = time.Duration(t.ReadTimeout) * time.Second
	}
	if t.WriteTimeout > 0 {
		srv.WriteTimeout = time.Duration(t.WriteTimeout) * time.Second
	}
	if t.IdleTimeout > 0 {
		srv.IdleTimeout = time.Duration(t.IdleTimeout) * time.Second
	}
	ca.server = srv

	if ca.Api.Https.Ativo {
		ca.LogInfo("Iniciando API HTTPS em", addr, "prefixo", prefixo)
		return srv.ListenAndServeTLS(ca.Api.Https.Cert, ca.Api.Https.Key)
	}
	ca.LogInfo("Iniciando API HTTP em", addr, "prefixo", prefixo)
	return srv.ListenAndServe()
}
