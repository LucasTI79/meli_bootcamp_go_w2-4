package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var INBOUND_URL = "/inbound-orders"

func TestCreateInboundOrder(t *testing.T) {
	t.Run("should return status 201 when successfull", func(t *testing.T) {
		mockedService := InboundOrdersServiceMock{}
		controller := handler.NewInboundOrder(&mockedService)
		server := getInboundServer(controller)
		i := domain.InboundOrder{
			OrderDate:      time.Now().AddDate(0, 0, 1),
			OrderNumber:    "125",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		mockedService.On("Create", mock.Anything, i).Return(i, nil)
		req, res := testutil.MakeRequest(http.MethodPost, INBOUND_URL, i)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[domain.InboundOrder]
		json.Unmarshal(res.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, i, received.Data)

	})
	t.Run("should return status 409 when a error occurs when creating the order", func(t *testing.T) {
		mockedService := InboundOrdersServiceMock{}
		controller := handler.NewInboundOrder(&mockedService)
		server := getInboundServer(controller)
		i := domain.InboundOrder{
			OrderDate:      time.Now().AddDate(0, 0, 1),
			OrderNumber:    "125",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		mockedService.On("Create", mock.Anything, i).Return(domain.InboundOrder{}, errors.New("generic error"))
		req, res := testutil.MakeRequest(http.MethodPost, INBOUND_URL, i)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)

	})
	t.Run("should return status 400 when missing date", func(t *testing.T) {
		mockedService := InboundOrdersServiceMock{}
		controller := handler.NewInboundOrder(&mockedService)
		server := getInboundServer(controller)

		orderNumber := "125"
		employeeID := 1
		productBatchID := 1
		warehouseID := 1

		i := handler.InboundOrderRequest{
			OrderNumber:    &orderNumber,
			EmployeeID:     &employeeID,
			ProductBatchID: &productBatchID,
			WarehouseID:    &warehouseID,
		}

		mockedService.On("Create", mock.Anything, i).Return(domain.InboundOrder{}, errors.New("generic error"))
		req, res := testutil.MakeRequest(http.MethodPost, INBOUND_URL, i)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

	})
	t.Run("should return status 400 when missing order number", func(t *testing.T) {
		mockedService := InboundOrdersServiceMock{}
		controller := handler.NewInboundOrder(&mockedService)
		server := getInboundServer(controller)

		date := time.Now().AddDate(0, 0, 1)
		employeeID := 1
		productBatchID := 1
		warehouseID := 1

		i := handler.InboundOrderRequest{
			OrderDate:      &date,
			EmployeeID:     &employeeID,
			ProductBatchID: &productBatchID,
			WarehouseID:    &warehouseID,
		}

		mockedService.On("Create", mock.Anything, i).Return(domain.InboundOrder{}, errors.New("generic error"))
		req, res := testutil.MakeRequest(http.MethodPost, INBOUND_URL, i)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

	})
	t.Run("should return status 400 when missing employeeId", func(t *testing.T) {
		mockedService := InboundOrdersServiceMock{}
		controller := handler.NewInboundOrder(&mockedService)
		server := getInboundServer(controller)

		orderNumber := "125"
		date := time.Now().AddDate(0, 0, 1)
		productBatchID := 1
		warehouseID := 1

		i := handler.InboundOrderRequest{
			OrderDate:      &date,
			OrderNumber:    &orderNumber,
			ProductBatchID: &productBatchID,
			WarehouseID:    &warehouseID,
		}

		mockedService.On("Create", mock.Anything, i).Return(domain.InboundOrder{}, errors.New("generic error"))
		req, res := testutil.MakeRequest(http.MethodPost, INBOUND_URL, i)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

	})
	t.Run("should return status 400 when missing warehouse id", func(t *testing.T) {
		mockedService := InboundOrdersServiceMock{}
		controller := handler.NewInboundOrder(&mockedService)
		server := getInboundServer(controller)

		orderNumber := "125"
		date := time.Now().AddDate(0, 0, 1)
		productBatchID := 1
		employeeID := 1

		i := handler.InboundOrderRequest{
			OrderDate:      &date,
			OrderNumber:    &orderNumber,
			ProductBatchID: &productBatchID,
			EmployeeID:     &employeeID,
		}

		mockedService.On("Create", mock.Anything, i).Return(domain.InboundOrder{}, errors.New("generic error"))
		req, res := testutil.MakeRequest(http.MethodPost, INBOUND_URL, i)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

	})
	t.Run("should return status 400 when missing produtch batch id", func(t *testing.T) {
		mockedService := InboundOrdersServiceMock{}
		controller := handler.NewInboundOrder(&mockedService)
		server := getInboundServer(controller)

		orderNumber := "125"
		date := time.Now().AddDate(0, 0, 1)
		warehouseID := 1
		employeeID := 1

		i := handler.InboundOrderRequest{
			OrderDate:   &date,
			OrderNumber: &orderNumber,
			WarehouseID: &warehouseID,
			EmployeeID:  &employeeID,
		}

		mockedService.On("Create", mock.Anything, i).Return(domain.InboundOrder{}, errors.New("generic error"))
		req, res := testutil.MakeRequest(http.MethodPost, INBOUND_URL, i)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)

	})
}

func getInboundServer(h *handler.InboundOrder) *gin.Engine {
	s := testutil.CreateServer()

	inboundOrdersRG := s.Group(INBOUND_URL)
	{
		inboundOrdersRG.POST("", middleware.Body[handler.InboundOrderRequest](), h.Create())
	}

	return s
}

type InboundOrdersServiceMock struct {
	mock.Mock
}

func (svc *InboundOrdersServiceMock) Create(c context.Context, i domain.InboundOrder) (domain.InboundOrder, error) {
	args := svc.Called(c, i)
	return args.Get(0).(domain.InboundOrder), args.Error(1)
}
