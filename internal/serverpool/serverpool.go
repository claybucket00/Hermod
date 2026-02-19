package serverpool

import (
	"hermod/internal/proxy"
	"net/http"
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	backends []*proxy.Backend
	current  uint64
	mu       sync.RWMutex
}

func NewServerPool(backends []*proxy.Backend) *ServerPool {
	return &ServerPool{
		backends: backends,
		current:  0,
		mu:       sync.RWMutex{},
	}
}

func (s *ServerPool) NextIndex() int {
	atomic.AddUint64(&s.current, 1)
	s.current = s.current % uint64(len(s.backends))
	return int(s.current)
}

func lb(s *ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peer := s.GetNextPeer()
		if peer != nil {
			peer.ReverseProxy.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
	}
}

// Round Robin
func (s *ServerPool) GetNextPeer() *proxy.Backend {
	next := s.NextIndex()
	l := next + len(s.backends)
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.AddUint64(&s.current, 1)
				s.current = s.current % uint64(len(s.backends))
			}
			return s.backends[idx]
		}
	}
	return nil
}
