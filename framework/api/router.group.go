package api

import (
	"net/http"
)

type Route struct {
	Pattern     string
	Method      TipoMetodoEnum
	HandlerFunc http.HandlerFunc
}

type RouteGroup struct {
	Name       string
	Routes     []Route
	Middleware []func(http.Handler) http.Handler
}

func (rg *RouteGroup) FindRoute(pattern string, method TipoMetodoEnum) *Route {
	for _, r := range rg.Routes {
		if (r.Pattern == pattern) && (r.Method == method) {
			return &r
		}
	}
	return nil
}
func (rg *RouteGroup) RegisterRoute(pattern string, method TipoMetodoEnum, handler http.HandlerFunc) {
	R := rg.FindRoute(pattern, method)
	if R == nil {
		R = &Route{
			Method:      method,
			Pattern:     pattern,
			HandlerFunc: handler,
		}
		rg.Routes = append(rg.Routes, *R)
	}
}
func (rg *RouteGroup) Use(middleware func(http.Handler) http.Handler) {
	rg.Middleware = append(rg.Middleware, middleware)
}

type RouteGroups []*RouteGroup

type RegistraController struct {
	RouteGroups RouteGroups
}

var registraGlobal *RegistraController

func GetRegistraInstance() *RegistraController {
	if registraGlobal == nil {
		registraGlobal = &RegistraController{
			RouteGroups: RouteGroups{},
		}
	}
	return registraGlobal
}

func (rc *RegistraController) RegisterRouteGroup(rg *RouteGroup) {
	rc.RouteGroups = append(rc.RouteGroups, rg)
}
func (rc *RegistraController) FindRouteGroup(name string) *RouteGroup {
	for _, rg := range rc.RouteGroups {
		if rg.Name == name {
			return rg
		}
	}
	return nil
}

func RegisterRouterGroup(name string) *RouteGroup {
	Rg := GetRegistraInstance().FindRouteGroup(name)
	if Rg == nil {
		Rg = &RouteGroup{
			Name:   name,
			Routes: []Route{},
		}
		GetRegistraInstance().RegisterRouteGroup(Rg)
	}

	return Rg
}
