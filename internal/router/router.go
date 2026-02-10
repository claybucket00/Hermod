package router

import (
	"hermod/internal/service"
	"net/http"
	"strings"
)

type Router struct {
	routes map[string]*service.Service
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*service.Service),
	}
}

func (r *Router) AddRoute(prefix string, svc *service.Service) {
	r.routes[prefix] = svc
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	for prefix, service := range r.routes {
		if strings.HasPrefix(path, prefix) {
			service.ServeHTTP(w, req)
			return
		}
	}

	http.NotFound(w, req)
}
