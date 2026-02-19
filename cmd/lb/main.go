package lb

import (
	"hermod/internal/config"
	"hermod/internal/router"
	"hermod/internal/serverpool"
	"log"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := router.NewRouter()
	for _, backend := range cfg.Backends {
		pool := serverpool.NewServerPool(backend.Address)
	}
}
