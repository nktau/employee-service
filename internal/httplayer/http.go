package httplayer

import (
	"fmt"
	"github.com/employee-service/internal/applayer"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Api struct {
	app applayer.App
	mux *chi.Mux
}

func New(app applayer.App) Api {
	api := Api{app: app, mux: chi.NewRouter()}
	api.mux.Get("/{id}", api.getUrlByShorter)
	api.mux.Post("/", api.createShorterByUrl)
	return api
}

func (api Api) Start() error {
	err := http.ListenAndServe(`:8080`, api.mux)
	return err
}

func (api Api) createShorterByUrl(w http.ResponseWriter, r *http.Request) {
	bodyReader, err := r.GetBody()
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ok, I got your url:", body)
	w.Write([]byte("ok"))

}

func (api Api) getUrlByShorter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "123" {
		w.Write([]byte("google.com\n"))
	}

}
