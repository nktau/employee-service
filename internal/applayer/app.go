package applayer

import "github.com/employee-service/internal/storagelayer"

type App struct {
	store storagelayer.Storage
}

func New(store storagelayer.Storage) App {
	return App{store: store}
}
