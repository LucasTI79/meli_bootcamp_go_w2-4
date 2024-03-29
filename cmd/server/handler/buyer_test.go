package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/gin-gonic/gin"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var BUYER_URL = "/buyers"

func TestBuyerCreate(t *testing.T) {
	t.Run("Returns 201 if successful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := domain.Buyer{
			ID:           1,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		svcMock.On("Create", mock.Anything, mock.Anything).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPost, BUYER_URL, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.BuyerCreate]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("Returns 422 if receives missing field type", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)

		body := map[string]any{
			"first_name": "nome",
			"last_name":  "sobrenome",
		}
		request, response := testutil.MakeRequest(http.MethodPost, BUYER_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("Returns 409 if Card number already exists", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)

		b := domain.Buyer{
			ID:           1,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		svcMock.On("Create", mock.Anything, mock.Anything).Return(domain.Buyer{}, errors.New("buyer already exists"))

		request, response := testutil.MakeRequest(http.MethodPost, BUYER_URL, b)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusConflict, response.Code)
	})
}

func TestBuyerGet(t *testing.T) {
	t.Run("Returns 200 if successful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := []domain.Buyer{
			{
				ID:           1,
				CardNumberID: "123",
				FirstName:    "nome",
				LastName:     "sobrenome"},
			{
				ID:           1,
				CardNumberID: "123",
				FirstName:    "nome",
				LastName:     "sobrenome",
			},
		}
		svcMock.On("GetAll", mock.Anything).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodGet, BUYER_URL, mock.Anything)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[[]domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("Returns 500 if buyers not found", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)

		svcMock.On("GetAll", mock.Anything).Return([]domain.Buyer{}, buyer.ErrNotFound)

		request, response := testutil.MakeRequest(http.MethodGet, BUYER_URL, mock.Anything)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
	t.Run("Returns 204 if buyers length is zero", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)

		svcMock.On("GetAll", mock.Anything).Return([]domain.Buyer{}, nil)

		request, response := testutil.MakeRequest(http.MethodGet, BUYER_URL, mock.Anything)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNoContent, response.Code)
	})

	t.Run("Returns 404 if id is not existing", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		id := 12
		urlId := fmt.Sprintf("%s/%d", BUYER_URL, id)
		svcMock.On("Get", mock.Anything, id).Return(domain.Buyer{}, buyer.ErrNotFound)
		request, response := testutil.MakeRequest(http.MethodGet, urlId, mock.Anything)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Returns 200 if id exists", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		urlId := fmt.Sprintf("%s/%d", BUYER_URL, expected.ID)
		svcMock.On("Get", mock.Anything, expected.ID).Return(expected, nil)
		request, response := testutil.MakeRequest(http.MethodGet, urlId, mock.Anything)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expected, received.Data)
	})
}

func TestBuyerDelete(t *testing.T) {
	t.Run("Returns 204 if successful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		urlId := fmt.Sprintf("%s/%d", BUYER_URL, expected.ID)
		svcMock.On("Delete", mock.Anything, expected.ID).Return(nil)
		request, response := testutil.MakeRequest(http.MethodDelete, urlId, mock.Anything)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("Returns 404 if not existent", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		urlId := fmt.Sprintf("%s/%d", BUYER_URL, expected.ID)
		svcMock.On("Delete", mock.Anything, expected.ID).Return(errors.New("buyer not found"))
		request, response := testutil.MakeRequest(http.MethodDelete, urlId, mock.Anything)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestBuyerUpdate(t *testing.T) {
	t.Run("Returns 200 if update is successful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "crash",
			LastName:     "bandicoot",
		}
		urlId := fmt.Sprintf("%s/%d", BUYER_URL, expected.ID)
		svcMock.On("Update", mock.Anything, expected).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPatch, urlId, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expected, received.Data)
	})
	t.Run("Returns 422 if update json is invalid", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := domain.Buyer{
			ID: 12,
		}
		urlId := fmt.Sprintf("%s/%d", BUYER_URL, expected.ID)
		svcMock.On("Update", mock.Anything, expected).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPatch, urlId, expected)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
	t.Run("Returns 404 if update is unsuccessful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "crash",
			LastName:     "bandicoot",
		}
		urlId := fmt.Sprintf("%s/%d", BUYER_URL, expected.ID)
		svcMock.On("Update", mock.Anything, expected).Return(domain.Buyer{}, errors.New("buyer not updated"))

		request, response := testutil.MakeRequest(http.MethodPatch, urlId, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetPurchaseOrderReports(t *testing.T) {
	t.Run("should return status 200 when id is valid", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)
		expected := []buyer.CountByBuyer{{
			ID:           10,
			CardNumberID: 10,
			FirstName:    "meli",
			LastName:     "osasco",
			Count:        10,
		}}
		url := fmt.Sprintf("%s/report-purchase-orders/%d", BUYER_URL, 1)

		svcMock.On("CountPurchaseOrders", mock.Anything, mock.Anything).Return(expected, nil)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		var received testutil.SuccessResponse[[]buyer.CountByBuyer]
		json.Unmarshal(res.Body.Bytes(), &received)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.ElementsMatch(t, expected, received.Data)
	})
	t.Run("should return status 404 when id not found", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)

		url := fmt.Sprintf("%s/report-purchase-orders/%d", BUYER_URL, 1)

		svcMock.On("CountPurchaseOrders", mock.Anything, mock.Anything).Return([]buyer.CountByBuyer{}, buyer.ErrNotFound)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})
	t.Run("should return status 409 when id alredy exists", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)

		url := fmt.Sprintf("%s/report-purchase-orders/%d", BUYER_URL, 1)

		svcMock.On("CountPurchaseOrders", mock.Anything, mock.Anything).Return([]buyer.CountByBuyer{}, buyer.ErrAlreadyExists)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
	t.Run("should return status 500 when bd is off", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := getBuyerServer(buyerHandler)

		url := fmt.Sprintf("%s/report-purchase-orders/%d", BUYER_URL, 1)

		svcMock.On("CountPurchaseOrders", mock.Anything, mock.Anything).Return([]buyer.CountByBuyer{}, buyer.ErrInternalServerError)

		req, res := testutil.MakeRequest(http.MethodGet, url, nil)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func getBuyerServer(h *handler.Buyer) *gin.Engine {
	s := testutil.CreateServer()

	buyerRG := s.Group(BUYER_URL)
	{
		buyerRG.GET("", h.GetAll())
		buyerRG.POST("", middleware.Body[domain.BuyerCreate](), h.Create())
		buyerRG.GET("/:id", middleware.IntPathParam(), h.Get())
		buyerRG.GET("/report-purchase-orders/:id", middleware.IntPathParam(), h.PurchaseOrderReport())
		buyerRG.GET("/report-purchase-orders/", h.PurchaseOrderReport())
		buyerRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
		buyerRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Buyer](), h.Update())
	}

	return s
}

type ServiceMockBuyer struct {
	mock.Mock
}

func (svc *ServiceMockBuyer) GetAll(c context.Context) ([]domain.Buyer, error) {
	args := svc.Called(c)
	return args.Get(0).([]domain.Buyer), args.Error(1)
}

func (svc *ServiceMockBuyer) Get(ctx context.Context, id int) (domain.Buyer, error) {
	args := svc.Called(ctx, id)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (svc *ServiceMockBuyer) Create(c context.Context, s domain.BuyerCreate) (domain.Buyer, error) {
	args := svc.Called(c, s)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (svc *ServiceMockBuyer) Update(ctx context.Context, s domain.Buyer, id int) (domain.Buyer, error) {
	args := svc.Called(ctx, s)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (svc *ServiceMockBuyer) Delete(ctx context.Context, id int) error {
	args := svc.Called(ctx, id)
	return args.Error(0)
}

func (svc *ServiceMockBuyer) CountPurchaseOrders(ctx context.Context, id int) ([]buyer.CountByBuyer, error) {
	args := svc.Called(ctx, id)
	return args.Get(0).([]buyer.CountByBuyer), args.Error(1)
}
