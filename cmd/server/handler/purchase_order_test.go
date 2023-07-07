package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	purchaseorder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/purchase_order"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const PURCHASE_ORDER_URL = "/purchase-orders"

func TestPurchaseOrderCreate(t *testing.T) {
	t.Run("Returns 201 if purchase order is created successfully", func(t *testing.T) {
		svc := PurchaseOrderServiceMock{}
		h := handler.NewPurchaseOrder(&svc)
		server := getPurchaseOrderServer(h)

		dto := handler.PurchaseOrderRequest{
			OrderNumber:     testutil.ToPtr("12345"),
			OrderDate:       testutil.ToPtr("2022-12-03"),
			TrackingCode:    testutil.ToPtr("12345"),
			BuyerID:         testutil.ToPtr(1),
			ProductRecordID: testutil.ToPtr(2),
			OrderStatusID:   testutil.ToPtr(1),
		}
		date, _ := time.Parse("2006-01-02", *dto.OrderDate)
		expected := domain.PurchaseOrder{
			ID:              1,
			OrderNumber:     "12345",
			OrderDate:       date,
			TrackingCode:    "12345",
			BuyerID:         1,
			ProductRecordID: 2,
			OrderStatusID:   1,
		}
		svc.On("Create", mock.Anything, mock.Anything).Return(expected, nil)

		req, res := testutil.MakeRequest(http.MethodPost, PURCHASE_ORDER_URL, dto)
		server.ServeHTTP(res, req)

		var response testutil.SuccessResponse[domain.PurchaseOrder]
		json.Unmarshal(res.Body.Bytes(), &response)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, expected, response.Data)
	})
	t.Run("Returns 422 if purchase order date is invalid", func(t *testing.T) {
		svc := PurchaseOrderServiceMock{}
		h := handler.NewPurchaseOrder(&svc)
		server := getPurchaseOrderServer(h)

		dto := handler.PurchaseOrderRequest{
			OrderNumber:     testutil.ToPtr("12345"),
			OrderDate:       testutil.ToPtr("2022-120-300"),
			TrackingCode:    testutil.ToPtr("12345"),
			BuyerID:         testutil.ToPtr(1),
			ProductRecordID: testutil.ToPtr(2),
			OrderStatusID:   testutil.ToPtr(1),
		}
		req, res := testutil.MakeRequest(http.MethodPost, PURCHASE_ORDER_URL, dto)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})
	t.Run("Returns 409 if purchase order already exists", func(t *testing.T) {
		svc := PurchaseOrderServiceMock{}
		h := handler.NewPurchaseOrder(&svc)
		server := getPurchaseOrderServer(h)

		dto := handler.PurchaseOrderRequest{
			OrderNumber:     testutil.ToPtr("12345"),
			OrderDate:       testutil.ToPtr("2022-12-03"),
			TrackingCode:    testutil.ToPtr("12345"),
			BuyerID:         testutil.ToPtr(1),
			ProductRecordID: testutil.ToPtr(2),
			OrderStatusID:   testutil.ToPtr(1),
		}
		svc.On("Create", mock.Anything, mock.Anything).Return(domain.PurchaseOrder{}, purchaseorder.ErrAlreadyExists)

		req, res := testutil.MakeRequest(http.MethodPost, PURCHASE_ORDER_URL, dto)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
	t.Run("Returns 409 if purchase order already exists", func(t *testing.T) {
		svc := PurchaseOrderServiceMock{}
		h := handler.NewPurchaseOrder(&svc)
		server := getPurchaseOrderServer(h)

		dto := handler.PurchaseOrderRequest{
			OrderNumber:     testutil.ToPtr("12345"),
			OrderDate:       testutil.ToPtr("2022-12-03"),
			TrackingCode:    testutil.ToPtr("12345"),
			BuyerID:         testutil.ToPtr(1),
			ProductRecordID: testutil.ToPtr(2),
			OrderStatusID:   testutil.ToPtr(1),
		}
		svc.On("Create", mock.Anything, mock.Anything).Return(domain.PurchaseOrder{}, purchaseorder.ErrAlreadyExists)

		req, res := testutil.MakeRequest(http.MethodPost, PURCHASE_ORDER_URL, dto)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
	t.Run("Returns 500 if repository fails", func(t *testing.T) {
		svc := PurchaseOrderServiceMock{}
		h := handler.NewPurchaseOrder(&svc)
		server := getPurchaseOrderServer(h)

		dto := handler.PurchaseOrderRequest{
			OrderNumber:     testutil.ToPtr("12345"),
			OrderDate:       testutil.ToPtr("2022-12-03"),
			TrackingCode:    testutil.ToPtr("12345"),
			BuyerID:         testutil.ToPtr(1),
			ProductRecordID: testutil.ToPtr(2),
			OrderStatusID:   testutil.ToPtr(1),
		}
		svc.On("Create", mock.Anything, mock.Anything).Return(domain.PurchaseOrder{}, purchaseorder.ErrInternalServerError)

		req, res := testutil.MakeRequest(http.MethodPost, PURCHASE_ORDER_URL, dto)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})
}

func getPurchaseOrderServer(h *handler.PurchaseOrder) *gin.Engine {
	s := testutil.CreateServer()
	rg := s.Group(PURCHASE_ORDER_URL)
	{
		rg.POST("", middleware.Body[handler.PurchaseOrderRequest](), h.Create())
	}
	return s
}

type PurchaseOrderServiceMock struct {
	mock.Mock
}

func (r *PurchaseOrderServiceMock) Create(c context.Context, purchaseOrder purchaseorder.PurchaseOrderDTO) (domain.PurchaseOrder, error) {
	args := r.Called(c, purchaseOrder)
	return args.Get(0).(domain.PurchaseOrder), args.Error(1)
}
