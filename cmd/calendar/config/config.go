package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
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
		log.Println("config::buildFromEnv::error: in env parsing:", err)
	}
	return err
}
