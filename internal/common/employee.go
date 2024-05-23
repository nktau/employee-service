package common

type EmployeeScheme struct {
	Name            string `json:"name" validate:"required"`
	Surname         string `json:"surname" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	CompanyId       int    `json:"companyId" validate:"required"`
	PassportType    string `json:"passportType" validate:"required"`
	PassportNumber  string `json:"passportNumber" validate:"required"`
	DepartmentName  string `json:"departmentName" validate:"required"`
	DepartmentPhone string `json:"departmentPhone" validate:"required"`
}

type EmployeeSchemeParticularUpdate struct {
	Name            string `json:"name" validate:"omitempty"`
	Surname         string `json:"surname" validate:"omitempty"`
	Phone           string `json:"phone" validate:"omitempty"`
	CompanyId       int    `json:"companyId" validate:"omitempty"`
	PassportType    string `json:"passportType" validate:"omitempty"`
	PassportNumber  string `json:"passportNumber" validate:"omitempty"`
	DepartmentName  string `json:"departmentName" validate:"omitempty"`
	DepartmentPhone string `json:"departmentPhone" validate:"omitempty"`
}
