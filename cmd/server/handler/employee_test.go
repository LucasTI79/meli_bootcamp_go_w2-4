package handler_test

import (
	"context"
	"encoding/json"
	"errors"
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
var EMPLOYEE_URL_ID_PATH = "/api/v1/employees/:id"

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
	t.Run("should return status 400 when missing first name", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.POST(EMPLOYEE_URL, controller.Create())
		e := domain.Employee{
			ID:           1,
			CardNumberID: "24",
			LastName:     "Kell",
			WarehouseID:  1,
		}
		req, res := testutil.MakeRequest(http.MethodPost, EMPLOYEE_URL, e)
		server.ServeHTTP(res, req)

		var received testutil.ErrorResponse
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
	t.Run("should return status 400 when missing card number id", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.POST(EMPLOYEE_URL, controller.Create())
		e := domain.Employee{
			ID:          1,
			FirstName:   "24",
			LastName:    "Kell",
			WarehouseID: 1,
		}
		req, res := testutil.MakeRequest(http.MethodPost, EMPLOYEE_URL, e)
		server.ServeHTTP(res, req)

		var received testutil.ErrorResponse
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
	t.Run("should return status 400 when missing last name", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.POST(EMPLOYEE_URL, controller.Create())
		e := domain.Employee{
			ID:           1,
			CardNumberID: "24",
			FirstName:    "Joel",
			WarehouseID:  1,
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
	t.Run("should return status 400 when  not sucessfull", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.GET(EMPLOYEE_URL, controller.GetAll())

		mockedService.On("GetAll", mock.Anything).Return([]domain.Employee{}, errors.New("employees not found"))

		req, res := testutil.MakeRequest(http.MethodGet, EMPLOYEE_URL, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}
func TestGetByIdEmployee(t *testing.T) {
	t.Run("should return status 200 when id is valid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.GET(EMPLOYEE_URL_ID_PATH, controller.Get())
		e := domain.Employee{
			ID:           1,
			CardNumberID: "125",
			FirstName:    "Mario",
			LastName:     "Kart",
			WarehouseID:  1,
		}
		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, e.ID)

		mockedService.On("Get", mock.Anything, 1).Return(e, nil)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Employee]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, e, received.Data)
	})
	t.Run("should return status 404 when employee does not exist", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.GET(EMPLOYEE_URL_ID_PATH, controller.Get())

		nonExistentId := 120

		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, nonExistentId)

		mockedService.On("Get", mock.Anything, nonExistentId).Return(domain.Employee{}, employee.ErrNotFound)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("should return status 200 and updated object when id is valid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.PATCH(EMPLOYEE_URL_ID_PATH, controller.Update())
		e := domain.Employee{
			ID:           1,
			CardNumberID: "125",
			FirstName:    "Mario",
			LastName:     "Kart",
			WarehouseID:  1,
		}
		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, e.ID)

		mockedService.On("Update", mock.Anything, e).Return(e, nil)

		req, res := testutil.MakeRequest(http.MethodPatch, url, e)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Employee]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, e, received.Data)
	})
	t.Run("should return status 400 when id is invalid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.PATCH(EMPLOYEE_URL_ID_PATH, controller.Update())

		invalidId := "oi"
		url := fmt.Sprintf("%s/%s", EMPLOYEE_URL, invalidId)

		req, res := testutil.MakeRequest(http.MethodPatch, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
	t.Run("should return status 422 when body is invalid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.PATCH(EMPLOYEE_URL_ID_PATH, controller.Update())

		e := map[string]any{
			"id":             1,
			"card_number_id": 125,
			"first_name":     "Mario",
			"last_name":      "Kart",
			"warehouse_id":   1,
		}
		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, e["id"])

		req, res := testutil.MakeRequest(http.MethodPatch, url, e)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})
	t.Run("should return status 404 when employee does not exist", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.PATCH(EMPLOYEE_URL_ID_PATH, controller.Update())

		e := domain.Employee{
			ID:           1,
			CardNumberID: "125",
			FirstName:    "Mario",
			LastName:     "Kart",
			WarehouseID:  1,
		}

		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, e.ID)

		mockedService.On("Update", mock.Anything, e).Return(domain.Employee{}, employee.ErrNotFound)

		req, res := testutil.MakeRequest(http.MethodPatch, url, e)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
	t.Run("should return status 400 when id is invalid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.GET(EMPLOYEE_URL_ID_PATH, controller.Get())

		invalidId := "oi"

		url := fmt.Sprintf("%s/%s", EMPLOYEE_URL, invalidId)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}

func TestDeleteEmployee(t *testing.T) {
	t.Run("should return status 204 when sucessfull", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.DELETE(EMPLOYEE_URL_ID_PATH, controller.Delete())

		idToDelete := 1

		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, idToDelete)

		mockedService.On("Delete", mock.Anything, idToDelete).Return(nil)

		req, res := testutil.MakeRequest(http.MethodDelete, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Employee]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusNoContent, res.Code)

	})
	t.Run("should return status 400 when id is invalid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.DELETE(EMPLOYEE_URL_ID_PATH, controller.Delete())

		idToDelete := "oi"

		url := fmt.Sprintf("%s/%s", EMPLOYEE_URL, idToDelete)

		req, res := testutil.MakeRequest(http.MethodDelete, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

	})
	t.Run("should return status 404 when employee does not exist", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := testutil.CreateServer()
		server.DELETE(EMPLOYEE_URL_ID_PATH, controller.Delete())

		idToDelete := 1

		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, idToDelete)

		mockedService.On("Delete", mock.Anything, idToDelete).Return(employee.ErrNotFound)

		req, res := testutil.MakeRequest(http.MethodDelete, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.ErrorResponse
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusNotFound, res.Code)
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
