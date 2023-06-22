package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var PRODUCTS_URL = "/products"

func TestProductCreate(t *testing.T) {
	t.Run("Returns 201 when creation succeds", func(t *testing.T) {
		mockSvc := ProductServiceMock{}
		h := handler.NewProduct(&mockSvc)
		server := getProductServer(h)

		body := handler.CreateRequest{
			Desc:       testutil.ToPtr("Sweet potato"),
			ExpR:       testutil.ToPtr(3),
			FreezeR:    testutil.ToPtr(1),
			Height:     testutil.ToPtr[float32](200),
			Length:     testutil.ToPtr[float32](40),
			NetW:       testutil.ToPtr[float32](10),
			Code:       testutil.ToPtr("SWP-1"),
			FreezeTemp: testutil.ToPtr[float32](20),
			Width:      testutil.ToPtr[float32](100),
			TypeID:     testutil.ToPtr(1),
			SellerID:   testutil.ToPtr(1),
		}
		created := domain.Product{
			ID:             0,
			Description:    "Sweet potato",
			ExpirationRate: 3,
			FreezingRate:   1,
			Height:         200,
			Length:         40,
			Netweight:      10,
			ProductCode:    "SWP-1",
			RecomFreezTemp: 20,
			Width:          100,
			ProductTypeID:  1,
			SellerID:       1,
		}

		mockSvc.On("Create", mock.Anything, mock.Anything).Return(created, nil)

		req, res := testutil.MakeRequest(http.MethodPost, "/products/", body)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})
	t.Run("Does not fail when fields are zero", func(t *testing.T) {
		mockSvc := ProductServiceMock{}
		h := handler.NewProduct(&mockSvc)
		server := getProductServer(h)

		body := handler.CreateRequest{
			Desc:       testutil.ToPtr(""),
			ExpR:       testutil.ToPtr(0),
			FreezeR:    testutil.ToPtr(0),
			Height:     testutil.ToPtr[float32](200),
			Length:     testutil.ToPtr[float32](40),
			NetW:       testutil.ToPtr[float32](10),
			Code:       testutil.ToPtr("SWP-1"),
			FreezeTemp: testutil.ToPtr[float32](0),
			Width:      testutil.ToPtr[float32](100),
			TypeID:     testutil.ToPtr(1),
			SellerID:   testutil.ToPtr(1),
		}

		mockSvc.On("Create", mock.Anything, mock.Anything).Return(domain.Product{}, nil)

		req, res := testutil.MakeRequest(http.MethodPost, "/products/", body)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})
	t.Run("Returns 422 when required fields are omitted", func(t *testing.T) {
		mockSvc := ProductServiceMock{}
		h := handler.NewProduct(&mockSvc)
		server := getProductServer(h)

		body := map[string]any{
			"description":  "",
			"product_code": "SWP-1",
			"seller_id":    1,
		}

		req, res := testutil.MakeRequest(http.MethodPost, "/products/", body)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})
	t.Run("Returns 409 when product code is not unique", func(t *testing.T) {
		mockSvc := ProductServiceMock{}
		h := handler.NewProduct(&mockSvc)
		server := getProductServer(h)

		body := handler.CreateRequest{
			Desc:       testutil.ToPtr(""),
			ExpR:       testutil.ToPtr(0),
			FreezeR:    testutil.ToPtr(0),
			Height:     testutil.ToPtr[float32](200),
			Length:     testutil.ToPtr[float32](40),
			NetW:       testutil.ToPtr[float32](10),
			Code:       testutil.ToPtr("SWP-1"),
			FreezeTemp: testutil.ToPtr[float32](0),
			Width:      testutil.ToPtr[float32](100),
			TypeID:     testutil.ToPtr(1),
			SellerID:   testutil.ToPtr(1),
		}

		mockSvc.On("Create", mock.Anything, mock.Anything).Return(domain.Product{}, product.NewErrInvalidProductCode(*body.Code))

		req, res := testutil.MakeRequest(http.MethodPost, "/products/", body)
		server.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
	})
}

func TestProductRead(t *testing.T) {
	t.Run("Returns all products on GetAll", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 204 when GetAll returns no products", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns existing product on Get by ID", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 400 when ID is not an int", func(t *testing.T) {
		t.Skip()
	})
}

func TestProductUpdate(t *testing.T) {
	t.Run("Returns 200 when update succeds", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Does not fail when updated value is zero", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 404 when ID is not found", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 409 when updated code exists", func(t *testing.T) {
		t.Skip()
	})
}

func TestProductDelete(t *testing.T) {

	t.Run("Returns 200 when delete succeds", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns 404 when ID is not found", func(t *testing.T) {
		t.Skip()
	})
}

func getProductServer(h *handler.Product) *gin.Engine {
	server := testutil.CreateServer()

	productRG := server.Group(PRODUCTS_URL)
	{
		productRG.POST("/", middleware.Body[handler.CreateRequest](), h.Create())
		productRG.GET("/", h.GetAll())
		productRG.GET("/:id", middleware.IntPathParam(), h.Get())
		productRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[handler.UpdateRequest](), h.Update())
		productRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
	}

	return server
}

type ProductServiceMock struct {
	mock.Mock
}

func (s *ProductServiceMock) Create(c context.Context, product product.CreateDTO) (domain.Product, error) {
	args := s.Called(c, product)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (s *ProductServiceMock) GetAll(c context.Context) ([]domain.Product, error) {
	args := s.Called(c)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (s *ProductServiceMock) Get(c context.Context, id int) (domain.Product, error) {
	args := s.Called(c, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (s *ProductServiceMock) Update(c context.Context, id int, updates product.UpdateDTO) (domain.Product, error) {
	args := s.Called(c, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (s *ProductServiceMock) Delete(c context.Context, id int) error {
	args := s.Called(c, id)
	return args.Error(0)
}
