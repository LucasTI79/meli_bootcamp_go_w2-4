package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var EMPLOYEE_URL = "/api/v1/employees"

func TestCreateEmployee(t *testing.T) {
	t.Run("should return status 201 if sucessfull", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.POST(EMPLOYEE_URL, controller.Create())
		e := domain.Employee{
			ID:           1,
			CardNumberID: "125",
			FirstName:    "Mario",
			LastName:     "Kart",
			WarehouseID:  1,
		}
		mockedService.On("Create", mock.Anything, e).Return(e, nil)
		req, res := testutil.MakeRequest(http.MethodPost, EMPLOYEE_URL, e)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Employee]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, e, received.Data)
	})
	t.Run("should return status 400 when receives missing fields", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.POST(EMPLOYEE_URL, controller.Create())
		e := domain.Employee{
			LastName: "Kart",
		}
		req, res := testutil.MakeRequest(http.MethodPost, EMPLOYEE_URL, e)
		server.ServeHTTP(res, req)

		var received testutil.ErrorResponse
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
	t.Run("should return status 422 when receives invalid field type", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.POST(EMPLOYEE_URL, controller.Create())
		e := map[string]any{
			"id":             1,
			"card_number_id": 125,
			"first_name":     "Mario",
			"last_name":      "Kart",
			"warehouse_id":   1,
		}
		req, res := testutil.MakeRequest(http.MethodPost, EMPLOYEE_URL, e)
		server.ServeHTTP(res, req)

		var received testutil.ErrorResponse
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})
	t.Run("should return status 409 when card number already exist", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.POST(EMPLOYEE_URL, controller.Create())
		e := domain.Employee{
			ID:           1,
			CardNumberID: "125",
			FirstName:    "Mario",
			LastName:     "Kart",
			WarehouseID:  1,
		}
		mockedService.On("Create", mock.Anything, e).Return(domain.Employee{}, employee.ErrAlreadyExists)
		req, res := testutil.MakeRequest(http.MethodPost, EMPLOYEE_URL, e)
		server.ServeHTTP(res, req)

		var received testutil.ErrorResponse
		json.Unmarshal(res.Body.Bytes(), &received)
		fmt.Println(received)
		assert.Equal(t, http.StatusConflict, res.Code)
	})
}
func TestGetAllEmployees(t *testing.T) {
	t.Run("should return status 200 when sucessfull", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.GET(EMPLOYEE_URL, controller.GetAll())

		es := []domain.Employee{{
			ID:           1,
			CardNumberID: "125",
			FirstName:    "Mario",
			LastName:     "Kart",
			WarehouseID:  1,
		},
			{
				ID:           2,
				CardNumberID: "126",
				FirstName:    "Peter",
				LastName:     "Parker",
				WarehouseID:  2,
			}}
		mockedService.On("GetAll", mock.Anything).Return(es, nil)

		req, res := testutil.MakeRequest(http.MethodGet, EMPLOYEE_URL, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[[]domain.Employee]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, es, received.Data)
	})
}
func TestGetByIdEmployee(t *testing.T) {
	t.Run("should return status 200 when id is valid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		e := domain.Employee{
			ID:           1,
			CardNumberID: "125",
			FirstName:    "Mario",
			LastName:     "Kart",
			WarehouseID:  1,
		}
		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, e.ID)
		server.GET(url, controller.Get())

		mockedService.On("Get", mock.Anything).Return(e, nil)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Employee]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, e, received.Data)
	})
}

type EmployeeServiceMock struct {
	mock.Mock
}

func (svc *EmployeeServiceMock) GetAll(c context.Context) ([]domain.Employee, error) {
	args := svc.Called(c)
	return args.Get(0).([]domain.Employee), args.Error(1)
}

func (svc *EmployeeServiceMock) Get(ctx context.Context, id int) (domain.Employee, error) {
	args := svc.Called(ctx, id)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (svc *EmployeeServiceMock) Create(c context.Context, s domain.Employee) (domain.Employee, error) {
	args := svc.Called(c, s)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (svc *EmployeeServiceMock) Update(ctx context.Context, s domain.Employee) (domain.Employee, error) {
	args := svc.Called(ctx, s)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (svc *EmployeeServiceMock) Delete(ctx context.Context, id int) error {
	args := svc.Called(ctx, id)
	return args.Error(0)
}
