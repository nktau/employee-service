package storagelayer

import "go.uber.org/zap"

type employeeScheme struct {
	Id              int
	Name            string
	Surname         string
	Phone           string
	CompanyId       int
	PassportType    string
	PassportNumber  string
	DepartmentName  string
	DepartmentPhone string
}

type Storage interface {
	Add(employee employeeScheme) (int, error)
	Delete(id int) error
	GetById(id int) ([]employeeScheme, error)
	GetByCompanyIdAndDepartmentName(companyId int) ([]employeeScheme, error)
	UpdateById(id int, employee employeeScheme) error
}

type store struct {
	logger *zap.Logger
	DBDSN  string
}

func New(logger *zap.Logger, DBDSN string) Storage {
	store := store{logger: logger, DBDSN: DBDSN}
	return store
}

func (s store) Add(employee employeeScheme) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s store) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func (s store) GetById(id int) ([]employeeScheme, error) {
	//TODO implement me
	panic("implement me")
}

func (s store) GetByCompanyIdAndDepartmentName(companyId int) ([]employeeScheme, error) {
	//TODO implement me
	panic("implement me")
}

func (s store) UpdateById(id int, employee employeeScheme) error {
	//TODO implement me
	panic("implement me")
}
