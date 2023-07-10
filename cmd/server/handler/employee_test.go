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
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var EMPLOYEE_URL = "/employees"

func TestCreateEmployee(t *testing.T) {
	t.Run("should return status 201 if sucessfull", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)
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
		server := getEmployeeServer(controller)
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
		server := getEmployeeServer(controller)
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
		server := getEmployeeServer(controller)
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
		server := getEmployeeServer(controller)
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
		server := getEmployeeServer(controller)
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
		server := getEmployeeServer(controller)

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
	t.Run("should return status 400 when not sucessfull", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)

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
		server := getEmployeeServer(controller)
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
		server := getEmployeeServer(controller)

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
		server := getEmployeeServer(controller)
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
	t.Run("should return status 422 when body is invalid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)

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
		server := getEmployeeServer(controller)

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
}

func TestDeleteEmployee(t *testing.T) {
	t.Run("should return status 204 when sucessfull", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)

		idToDelete := 1

		url := fmt.Sprintf("%s/%d", EMPLOYEE_URL, idToDelete)

		mockedService.On("Delete", mock.Anything, idToDelete).Return(nil)

		req, res := testutil.MakeRequest(http.MethodDelete, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.Employee]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusNoContent, res.Code)

	})
	t.Run("should return status 404 when employee does not exist", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)

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
func TestGetInboundReports(t *testing.T) {
	t.Run("should return status 200 when id is valid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)
		r := []domain.InboundReport{{
			ID:                 1,
			CardNumberID:       "130",
			FirstName:          "Mario",
			LastName:           "Kart",
			WarehouseID:        1,
			InboundOrdersCount: 1,
		}}
		url := fmt.Sprintf("%s/report-inbound-orders/%d", EMPLOYEE_URL, 1)

		mockedService.On("GetInboundReport", mock.Anything, 1).Return(r, nil)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[[]domain.InboundReport]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, r, received.Data)
	})
	t.Run("should return status 404 when id not found", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)

		url := fmt.Sprintf("%s/report-inbound-orders/%d", EMPLOYEE_URL, 1)

		mockedService.On("GetInboundReport", mock.Anything, 1).Return([]domain.InboundReport{}, employee.ErrNotFound)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
	t.Run("should return status 200 when id is valid", func(t *testing.T) {
		mockedService := EmployeeServiceMock{}
		controller := handler.NewEmployee(&mockedService)
		server := getEmployeeServer(controller)
		r := []domain.InboundReport{{
			ID:                 1,
			CardNumberID:       "130",
			FirstName:          "Mario",
			LastName:           "Kart",
			WarehouseID:        1,
			InboundOrdersCount: 1,
		},
			{
				ID:                 2,
				CardNumberID:       "122",
				FirstName:          "Me",
				LastName:           "Tony",
				WarehouseID:        22,
				InboundOrdersCount: 2,
			}}
		url := fmt.Sprintf("%s/report-inbound-orders/", EMPLOYEE_URL)

		mockedService.On("GetInboundReport", mock.Anything, 0).Return(r, nil)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[[]domain.InboundReport]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, r, received.Data)
	})
}

func getEmployeeServer(h *handler.Employee) *gin.Engine {
	s := testutil.CreateServer()

	employeeRG := s.Group(EMPLOYEE_URL)
	{
		employeeRG.GET("", h.GetAll())
		employeeRG.POST("", middleware.Body[domain.Employee](), h.Create())
		employeeRG.GET("/:id", middleware.IntPathParam(), h.Get())
		employeeRG.GET("/report-inbound-orders/:id", middleware.IntPathParam(), h.GetInboundReport())
		employeeRG.GET("/report-inbound-orders/", h.GetInboundReport())
		employeeRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Employee](), h.Update())
		employeeRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
	}

	return s
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
func (svc *EmployeeServiceMock) GetInboundReport(ctx context.Context, id int) ([]domain.InboundReport, error) {
	args := svc.Called(ctx, id)
	return args.Get(0).([]domain.InboundReport), args.Error(1)
}
