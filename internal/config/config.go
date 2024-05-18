package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenAddress string
	DatabaseDSN   string
}

func New() Config {
	cfg := Config{}
	cfg.parseFlags()
	cfg.parseEnv()
	return cfg

}
func (cfg *Config) parseFlags() {
	flag.StringVar(&cfg.ListenAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseDSN, "d", "",
		"database dsn")
	flag.Parse()
}

func (cfg *Config) parseEnv() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ListenAddress = envRunAddr
	}
	if value, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = value
	}
}
