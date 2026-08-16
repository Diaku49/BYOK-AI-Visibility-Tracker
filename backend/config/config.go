package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port      string `envconfig:"PORT" default:"8080"`
	JWTSecret string `envconfig:"JWT_SECRET" default:"mysecretnotavailable"`
	MasterKey string `envconfig:"MASTER_KEY"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("AITracker", &cfg)
	if err != nil {
		return nil, fmt.Errorf("process configuration: %w", err)
	}

	return &cfg, nil
}
