package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port      string `envconfig:"PORT" default:"8080"`
	JWTSecret string `envconfig:"JWT_SECRET" default:"mysecretnotavailable"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("AITracker", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
