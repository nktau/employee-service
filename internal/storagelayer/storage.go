package storagelayer

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nktau/employee-service/internal/common"
	"github.com/nktau/employee-service/internal/errors"
	"go.uber.org/zap"
	"time"
)

type EmployeeDBScheme struct {
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
	Create(employee common.EmployeeScheme) (int, error)
	Delete(id int) error
	GetEmployeesByCompanyId(id int) ([]EmployeeDBScheme, error)
	GetEmployeesByCompanyAndDepartment(companyId int, departmentName string) ([]EmployeeDBScheme, error)
	UpdateById(int, string) error
}

type store struct {
	logger       *zap.Logger
	dbConnection *sql.DB
}

func New(logger *zap.Logger, databaseURL string) Storage {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		logger.Fatal("can't create db connection", zap.Error(err))
	}
	store := store{logger: logger, dbConnection: db}

	return store
}

func (s store) Create(employee common.EmployeeScheme) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	insertEmployeeQuery := `INSERT INTO employees (name, surname, phone, companyId, 
                       passportType, passportNumber, departmentName, departmentPhone) 
					  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`
	row := s.dbConnection.QueryRowContext(ctx, insertEmployeeQuery,
		employee.Name, employee.Surname, employee.Phone, employee.CompanyId,
		employee.PassportType, employee.PassportNumber, employee.DepartmentName,
		employee.DepartmentPhone,
	)
	var employeeId int
	err := row.Scan(&employeeId)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return 0, err
	}
	return employeeId, nil
}

func (s store) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res, err := s.dbConnection.ExecContext(ctx, "DELETE FROM employees WHERE id=$1", id)
	if err != nil {
		s.logger.Error("can't delete user", zap.Error(err))
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return err
	}
	if count == 0 {
		return errors.ErrNotExistsEmployeeID
	}
	return nil
}

func (s store) GetEmployeesByCompanyId(id int) ([]EmployeeDBScheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	GetEmployeesByCompanyIDQuery := `SELECT * FROM employees WHERE companyId=$1;`
	rows, err := s.dbConnection.QueryContext(ctx, GetEmployeesByCompanyIDQuery, id)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var employees []EmployeeDBScheme
	for rows.Next() {
		var employee EmployeeDBScheme
		err = rows.Scan(&employee.Id, &employee.Name, &employee.Surname,
			&employee.Phone, &employee.CompanyId, &employee.PassportType,
			&employee.PassportNumber, &employee.DepartmentName, &employee.DepartmentPhone)
		if err != nil {
			s.logger.Error("", zap.Error(err))
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (s store) GetEmployeesByCompanyAndDepartment(companyId int, departmentName string) ([]EmployeeDBScheme, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	GetEmployeesByCompanyIDQuery := `SELECT * FROM employees WHERE companyId=$1 and departmentName=$2;`
	rows, err := s.dbConnection.QueryContext(ctx, GetEmployeesByCompanyIDQuery, companyId, departmentName)
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var employees []EmployeeDBScheme
	for rows.Next() {
		var employee EmployeeDBScheme
		err = rows.Scan(&employee.Id, &employee.Name, &employee.Surname,
			&employee.Phone, &employee.CompanyId, &employee.PassportType,
			&employee.PassportNumber, &employee.DepartmentName, &employee.DepartmentPhone)
		if err != nil {
			s.logger.Error("", zap.Error(err))
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (s store) UpdateById(id int, employee string) error {
	// данные пользователя формируются внутри кавычек, поэтому sql инъекции тут нет.
	//
	//запрос curl localhost:8080/employee/8 -X PATCH -d '{"Name":"NewName123;drop table employees"}'
	//сформируется в UPDATE employees SET name='NewName123;drop table employees' WHERE id=8
	sqlQuery := fmt.Sprintf("UPDATE employees SET %s WHERE id=%d", employee, id)
	s.logger.Debug(sqlQuery)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res, err := s.dbConnection.ExecContext(ctx, sqlQuery)
	if err != nil {
		s.logger.Error("can't update user", zap.Error(err))
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		s.logger.Error("", zap.Error(err))
		return err
	}
	fmt.Println("count: ", count)
	if count == 0 {
		return errors.ErrNotExistsEmployeeID
	}
	return nil
}
