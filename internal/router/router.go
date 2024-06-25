package router

import (
	"net/http"
	"strings"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (ro *Router) AddRoute(method, path string, handler http.HandlerFunc) {
	ro.routes = append(ro.routes, Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range ro.routes {
		if r.Method == route.Method && matchPath(r.URL.Path, route.Path) {
			route.Handler(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func matchPath(requestPath, routePath string) bool {
	requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")
	routeParts := strings.Split(strings.Trim(routePath, "/"), "/")
	if len(requestParts) != len(routeParts) {
		return false
	}
	for i, part := range routeParts {
		if part != requestParts[i] && part != "{}" {
			return false
		}
	}
	return true
}
