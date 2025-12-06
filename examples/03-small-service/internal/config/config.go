// Package config
package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerAddr string
	Debug      bool
	MaxUsers   int
}

func Load() Config {
	cfg := Config{
		ServerAddr: ":8080",
		Debug:      false,
		MaxUsers:   100,
	}

	if addr := os.Getenv("SERVER_ADDR"); addr != "" {
		cfg.ServerAddr = addr
	}

	if os.Getenv("DEBUG") == "true" {
		cfg.Debug = true
	}

	if maxUsers := os.Getenv("MAX_USERS"); maxUsers != "" {
		if n, err := strconv.Atoi(maxUsers); err == nil {
			cfg.MaxUsers = n
		}
	}

	return cfg
}
