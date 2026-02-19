package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL            *url.URL
	alive          bool
	mu             sync.RWMutex
	ReverseProxy   *httputil.ReverseProxy
	HealthCheckURL string
	Endpoints      []string
}

func NewBackend(backendURL string, healthCheckUrl string, endpoints []string) (*Backend, error) {
	url, err := url.Parse(backendURL)

	if err != nil {
		return nil, err
	}

	return &Backend{
		URL:            url,
		alive:          true,
		mu:             sync.RWMutex{},
		ReverseProxy:   httputil.NewSingleHostReverseProxy(url),
		HealthCheckURL: healthCheckUrl,
		Endpoints:      endpoints,
	}, nil
}

func (backendServer *Backend) HealthCheck() (bool, string, error) {
	// TODO: add support for other layers (currently only L7)
	// TODO: add support for timeouts
	resp, err := http.Get(backendServer.URL.String() + backendServer.HealthCheckURL)
	if err != nil {
		return false, "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO: add support for other status codes and logging
		return false, resp.Status, nil
	}
	return true, resp.Status, nil
}

func (backendServer *Backend) SetAlive(alive bool) {
	backendServer.mu.Lock()
	backendServer.alive = alive
	defer backendServer.mu.Unlock()
}

func (backendServer *Backend) IsAlive() bool {
	backendServer.mu.RLock()
	defer backendServer.mu.RUnlock()

	return backendServer.alive
}
