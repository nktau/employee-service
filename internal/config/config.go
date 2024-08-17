package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenAddress string
	DatabaseURL   string
}

func New() Config {
	cfg := Config{}
	cfg.parseFlags()
	cfg.parseEnv()
	return cfg

}
func (cfg *Config) parseFlags() {
	flag.StringVar(&cfg.ListenAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseURL, "d", "postgres://postgres:mysecretpassword@localhost:5433/employee_service?sslmode=disable",
		"database dsn")
	flag.Parse()
}

func (cfg *Config) parseEnv() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ListenAddress = envRunAddr
	}
	if value, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseURL = value
	}
}
