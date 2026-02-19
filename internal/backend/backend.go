package backend

import (
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
