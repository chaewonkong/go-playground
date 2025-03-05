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

type Config interface {
	// Keys() []string
	UnmarshalFromEnviron() error
}

type config struct {
	Port string `env:"PORT"`
	DSN  string `env:"DSN"`
}

func (c *config) UnmarshalFromEnviron() error {
	_, err := env.UnmarshalFromEnviron(c)
	return err
}

var _ Config = (*config)(nil)
