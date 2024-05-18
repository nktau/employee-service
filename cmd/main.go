package main

import (
	"github.com/employee-service/internal/applayer"
	"github.com/employee-service/internal/httplayer"
	log "github.com/employee-service/internal/logger"
	"github.com/employee-service/internal/storagelayer"
	"log"
)

func main() {
	logger := log.InitLogger()
	storage := storagelayer.New(logger)
	app := applayer.New(storage)
	api := httplayer.New(app)
	err := api.Start()
	if err != nil {
		log.Fatal(err)
	}
}
