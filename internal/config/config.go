package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type BackendConfig struct {
	Address   string   `yaml:"address"`
	Endpoints []string `yaml:"endpoints"`
}

type Config struct {
	Backends []BackendConfig `yaml:"backends"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
