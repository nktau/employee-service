package applayer

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nktau/employee-service/internal/common"
	"github.com/nktau/employee-service/internal/errors"
	"github.com/nktau/employee-service/internal/storagelayer"
	"go.uber.org/zap"
	"reflect"
	"strconv"
	"strings"
)

type App struct {
	store  storagelayer.Storage
	logger *zap.Logger
}

func New(store storagelayer.Storage, logger *zap.Logger) App {
	return App{store: store, logger: logger}
}

func (app App) CreateEmployee(body []byte) (id int, err error) {
	var employee common.EmployeeScheme // move to applayer
	err = json.Unmarshal(body, &employee)
	if err != nil {
		app.logger.Info("unmarshal body error",
			zap.String("client data", string(body)),
			zap.Error(err),
		)
		return 0, errors.ErrInvalidJsonData
	}
	validate := validator.New()
	err = validate.Struct(employee)
	if err != nil {
		return 0, errors.ErrWrongEmployeeData
	}
	id, err = app.store.Create(employee)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (app App) DeleteEmployee(id string) (err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.ErrInvalidEmployeeID
	}
	if idInt < 0 {
		return errors.ErrInvalidEmployeeID
	}
	err = app.store.Delete(idInt)
	return err
}

func (app App) GetEmployeesByCompanyId(id string) ([]byte, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.ErrInvalidCompanyID
	}
	employees, err := app.store.GetEmployeesByCompanyId(idInt)
	if err != nil {
		return nil, err
	}
	if len(employees) == 0 {
		return nil, errors.ErrNotExistsCompanyID
	}
	employeesJSON, err := json.Marshal(employees)
	if err != nil {
		return nil, err
	}
	return employeesJSON, nil
}

func (app App) GetEmployeesByCompanyAndDepartment(companyId, departmentName string) ([]byte, error) {
	companyIdInt, err := strconv.Atoi(companyId)
	if err != nil {
		return nil, errors.ErrInvalidCompanyID
	}
	employees, err := app.store.GetEmployeesByCompanyAndDepartment(companyIdInt, departmentName)
	if err != nil {
		return nil, err
	}
	if len(employees) == 0 {
		return nil, errors.ErrNotExistsCompanyORDepartment
	}
	employeesJSON, err := json.Marshal(employees)
	if err != nil {
		return nil, err
	}
	return employeesJSON, nil
}

func (app App) UpdateEmployee(id string, body []byte) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.ErrInvalidEmployeeID
	}
	if idInt < 0 {
		return errors.ErrInvalidEmployeeID
	}

	var employee common.EmployeeSchemeParticularUpdate
	err = json.Unmarshal(body, &employee)
	if err != nil {
		app.logger.Info("unmarshal body error",
			zap.String("client data", string(body)),
			zap.Error(err),
		)
		return errors.ErrInvalidJsonData
	}
	err = app.store.UpdateById(idInt, app.getEmployeeFieldsToUpdate(employee))

	return err
}

func (app App) getEmployeeFieldsToUpdate(employee common.EmployeeSchemeParticularUpdate) string {
	employeeFieldsToUpdate := ""
	val := reflect.ValueOf(employee)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if !valueField.IsZero() {

			nameField := val.Type().Field(i).Tag.Get("json")

			employeeFieldsToUpdate += fmt.Sprintf("%v='%v',", nameField, valueField)
		}
	}
	employeeFieldsToUpdate = strings.TrimSuffix(employeeFieldsToUpdate, ",")
	fmt.Println("employeeFieldsToUpdate", employeeFieldsToUpdate)
	return employeeFieldsToUpdate
}
