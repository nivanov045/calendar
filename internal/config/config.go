package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func BuildConfig() (Config, error) {
	var cfg Config
	cfg.buildFromFlags()
	err := cfg.buildFromEnv()
	return cfg, err
}

func (cfg *Config) buildFromFlags() {
	flag.StringVar(&cfg.Address, "a", "127.0.0.1:8080", "address")
	flag.Parse()
}
func (cfg *Config) buildFromEnv() error {
	err := env.Parse(cfg)
	if err != nil {
		log.Error().Err(err).Stack()
	}
	return err
}
