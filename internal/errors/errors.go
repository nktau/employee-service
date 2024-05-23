package errors

import "errors"

var ErrNotExistsEmployeeID = errors.New("no such employee in database")
var ErrWrongEmployeeData = errors.New("wrong employee data")
var ErrInvalidEmployeeID = errors.New("employee id must be a non-negative integer")
var ErrInvalidCompanyID = errors.New("employee id must be an integer")
var ErrNotExistsCompanyID = errors.New("no such company in database")
var ErrNotExistsCompanyORDepartment = errors.New("no such company or department in database")
var ErrInvalidJsonData = errors.New("invalid request json data")
