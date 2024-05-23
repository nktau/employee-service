package httplayer

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/nktau/employee-service/internal/applayer"
	"github.com/nktau/employee-service/internal/common"
	"github.com/nktau/employee-service/internal/errors"
	"github.com/nktau/employee-service/internal/logger"
	"github.com/nktau/employee-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestAddEmployee(t *testing.T) {
	employee := common.EmployeeScheme{
		Name:            "employee1Name",
		Surname:         "employee1Surname",
		Phone:           "employee1Phone",
		CompanyId:       1,
		PassportType:    "employee1PassportType",
		PassportNumber:  "employee1PassportNumber",
		DepartmentName:  "employee1DepartmentName",
		DepartmentPhone: "employee1DepartmentPhone",
	}
	invalidEmployee := `{"Name": "invalid employee", "Phone": "12345678912"}`
	invalidJson := `"Name": "invalid employee", "Phone": "12345678912"`

	type want struct {
		statusCode   int
		responseBody string
	}
	tests := []struct {
		name     string
		path     string
		employee interface{}
		want     want
	}{
		{
			name: "positive test #1",
			want: want{
				statusCode:   201,
				responseBody: "1",
			},
			path:     "/employee",
			employee: employee,
		},
		{
			name: "invalid employee #2",
			want: want{
				statusCode:   400,
				responseBody: errors.ErrWrongEmployeeData.Error(),
			},
			path:     "/employee",
			employee: invalidEmployee,
		},
		{
			name: "invalid json #3",
			want: want{
				statusCode:   400,
				responseBody: MessageInvalidJSON,
			},
			path:     "/employee",
			employee: invalidJson,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storageMock := mocks.NewMockStorage(ctrl)
	employeeID := 1
	storageMock.EXPECT().Create(employee).Return(employeeID, nil)
	logger := logger.InitLogger()
	app := applayer.New(storageMock, logger)
	api := New(app, logger)
	ts := httptest.NewServer(api.mux)

	client := resty.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonEmployee, _ := json.Marshal(tt.employee)
			fmt.Println(string(jsonEmployee))

			result, err := client.R().SetBody(tt.employee).Post(ts.URL + tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			fmt.Println(result.String())
			assert.Equal(t, tt.want.responseBody, result.String())
		})
	}
}

func TestDeleteEmployee(t *testing.T) {
	type want struct {
		statusCode   int
		responseBody string
	}
	validEmployeeID := 1
	notExistsEmployeeID := 0
	invalidEmployeeID := -1
	tests := []struct {
		name       string
		path       string
		employeeId int
		want       want
	}{
		{
			name: "positive test #1",
			want: want{
				statusCode:   202,
				responseBody: MessageEmployeeDeletedSuccessfully,
			},
			path: fmt.Sprintf("/employee/%d", validEmployeeID),
		},
		{
			name: "not exists employee #2",
			want: want{
				statusCode:   404,
				responseBody: errors.ErrNotExistsEmployeeID.Error(),
			},
			path: fmt.Sprintf("/employee/%d", notExistsEmployeeID),
		},
		{
			name: "invalid employee id #3",
			want: want{
				statusCode:   400,
				responseBody: errors.ErrInvalidEmployeeID.Error(),
			},
			path: fmt.Sprintf("/employee/%d", invalidEmployeeID),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storageMock := mocks.NewMockStorage(ctrl)
	storageMock.EXPECT().Delete(validEmployeeID).Return(nil)
	storageMock.EXPECT().Delete(notExistsEmployeeID).Return(errors.ErrNotExistsEmployeeID)
	logger := logger.InitLogger()
	app := applayer.New(storageMock, logger)
	api := New(app, logger)
	ts := httptest.NewServer(api.mux)

	client := resty.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.employeeId)
			result, err := client.R().Delete(ts.URL + tt.path)
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode())
			fmt.Println(result.String())
			assert.Equal(t, tt.want.responseBody, result.String())
		})
	}
}

// Надо бы написать тесты на остальные хендлеры,
//но не очень хочется тратить время на это в тестовом задании
