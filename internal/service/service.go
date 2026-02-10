package service

import (
	"hermod/internal/proxy"
	"net/http"
)

type Service struct {
	Name string
	Pool *proxy.ServerPool
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := s.Pool.GetNextPeer()
	if backend == nil {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}
	backend.ReverseProxy.ServeHTTP(w, r)
}
