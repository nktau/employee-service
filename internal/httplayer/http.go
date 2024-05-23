package httplayer

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/nktau/employee-service/internal/applayer"
	employeeerrors "github.com/nktau/employee-service/internal/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var MessageInternalServerError = "internal server error"
var MessageInvalidJSON = "invalid request json data"
var MessageEmployeeDeletedSuccessfully = "employee deleted successfully"
var MessageEmployeeUpdatedSuccessfully = "employee updated successfully"

type Api struct {
	app    applayer.App
	mux    *chi.Mux
	logger *zap.Logger
}

func New(app applayer.App, logger *zap.Logger) Api {
	api := Api{app: app, mux: chi.NewRouter(), logger: logger}
	api.mux.Use(api.withLogging)
	api.mux.Middlewares()
	api.mux.Post("/employee", api.createEmployee)
	api.mux.Delete("/employee/{id}", api.deleteEmployee)
	api.mux.Patch("/employee/{id}", api.updateEmployee)
	api.mux.Get("/company/{id}", api.getEmployeesByCompanyID)
	api.mux.Get("/company/{id}/department/{name}", api.getEmployeesByCompanyAndDepartment)

	return api
}

func (api Api) Start() error {
	err := http.ListenAndServe(`:8080`, api.mux)
	return err
}

func (api Api) createEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Content-Type:", w.Header().Get("Content-Type"))
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.logger.Error("read body error", zap.Error(err))
		http.Error(w, MessageInternalServerError, http.StatusInternalServerError)
		return
	}

	id, err := api.app.CreateEmployee(body)
	if err != nil {
		switch {
		case errors.Is(err, employeeerrors.ErrWrongEmployeeData):
			http.Error(w, employeeerrors.ErrWrongEmployeeData.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, employeeerrors.ErrInvalidJsonData):
			http.Error(w, employeeerrors.ErrInvalidJsonData.Error(), http.StatusBadRequest)
			return
		default:
			http.Error(w, MessageInternalServerError, http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(fmt.Sprint(id)))
	if err != nil {
		api.logger.Error("write response body error", zap.Error(err))
	}

}

func (api Api) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	err := api.app.DeleteEmployee(chi.URLParam(r, "id"))
	if err != nil {
		switch {
		case errors.Is(err, employeeerrors.ErrInvalidEmployeeID):
			http.Error(w, employeeerrors.ErrInvalidEmployeeID.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, employeeerrors.ErrNotExistsEmployeeID):
			http.Error(w, employeeerrors.ErrNotExistsEmployeeID.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, MessageInternalServerError, http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte(MessageEmployeeDeletedSuccessfully))
	if err != nil {
		api.logger.Error("write response body error", zap.Error(err))
	}
}

func (api Api) updateEmployee(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.logger.Error("read body error", zap.Error(err))
		http.Error(w, MessageInternalServerError, http.StatusInternalServerError)
		return
	}
	err = api.app.UpdateEmployee(chi.URLParam(r, "id"), body)
	if err != nil {
		switch {
		case errors.Is(err, employeeerrors.ErrInvalidEmployeeID):
			http.Error(w, employeeerrors.ErrInvalidEmployeeID.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, employeeerrors.ErrNotExistsEmployeeID):
			http.Error(w, employeeerrors.ErrNotExistsEmployeeID.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, MessageInternalServerError, http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(MessageEmployeeUpdatedSuccessfully))
	if err != nil {
		api.logger.Error("write response body error", zap.Error(err))
	}

}

func (api Api) getEmployeesByCompanyID(w http.ResponseWriter, r *http.Request) {
	employees, err := api.app.GetEmployeesByCompanyId(chi.URLParam(r, "id"))
	if err != nil {
		switch {
		case errors.Is(err, employeeerrors.ErrInvalidCompanyID):
			http.Error(w, employeeerrors.ErrInvalidCompanyID.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, employeeerrors.ErrNotExistsCompanyID):
			http.Error(w, employeeerrors.ErrNotExistsCompanyID.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, MessageInternalServerError, http.StatusInternalServerError)
			return
		}
	}
	_, err = w.Write(employees)
	if err != nil {
		api.logger.Error("write response body error", zap.Error(err))
	}
}

func (api Api) getEmployeesByCompanyAndDepartment(w http.ResponseWriter, r *http.Request) {
	employees, err := api.app.GetEmployeesByCompanyAndDepartment(
		chi.URLParam(r, "id"),
		chi.URLParam(r, "name"),
	)
	if err != nil {
		switch {
		case errors.Is(err, employeeerrors.ErrInvalidCompanyID):
			http.Error(w, employeeerrors.ErrInvalidCompanyID.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, employeeerrors.ErrNotExistsCompanyORDepartment):
			http.Error(w, employeeerrors.ErrNotExistsCompanyORDepartment.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, MessageInternalServerError, http.StatusInternalServerError)
			return
		}
	}
	_, err = w.Write(employees)
	if err != nil {
		api.logger.Error("write response body error", zap.Error(err))
	}
}
