package main

import (
	"fmt"

	"github.com/Netflix/go-env"
)

func main() {
	cfg := new(config)
	if err := cfg.UnmarshalFromEnviron(); err != nil {
		panic(err)
	}

	fmt.Printf("Port: %s, DSN: %s\n", cfg.Port, cfg.DSN)
}

type config struct {
	ClusterConfig
	Port string `env:"PORT"`
	DSN  string `env:"DSN"`
}

type ClusterConfig struct {
	ClusterID string `env:"CLUSTER_ID"`
	SentryDSN string `env:"SENTRY_DSN"`
}

func (c *config) UnmarshalFromEnviron() error {
	_, err := env.UnmarshalFromEnviron(c)
	return err
}
