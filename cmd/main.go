package main

import (
	"github.com/nktau/employee-service/internal/applayer"
	"github.com/nktau/employee-service/internal/config"
	"github.com/nktau/employee-service/internal/httplayer"
	log "github.com/nktau/employee-service/internal/logger"
	"github.com/nktau/employee-service/internal/storagelayer"
	"github.com/nktau/employee-service/migrations"
)

func main() {
	logger := log.InitLogger()
	cfg := config.New()
	migrations.RunPostgresMigration(logger, "file://migrations/postgres", cfg.DatabaseURL)
	storage := storagelayer.New(logger, cfg.DatabaseURL)
	app := applayer.New(storage, logger)
	api := httplayer.New(app, logger)
	err := api.Start()
	if err != nil {
		logger.Fatal("failed to start api ")
	}
}
