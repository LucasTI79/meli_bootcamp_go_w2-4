package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var BASE_URL = "/api/v1/sellers"

func TestCreate(t *testing.T) {
	t.Run("Returns 201 if successful", func(t *testing.T) {
		svcMock := ServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(BASE_URL, sellerHandler.Create())

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		svcMock.On("Save", mock.Anything, expected).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPost, BASE_URL, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, expected, received.Data)
		fmt.Println(received)
	})

	t.Run("Returns 400 if receives invalid field type", func(t *testing.T) {
		svcMock := ServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(BASE_URL, sellerHandler.Create())

		body := map[string]any{
			"cid":          "123", // passing CID as string instead of int
			"company_name": "TEST",
			"address":      "test street",
			"telephone":    "9999999",
		}
		request, response := testutil.MakeRequest(http.MethodPost, BASE_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Returns 422 if receives missing field type", func(t *testing.T) {
		svcMock := ServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(BASE_URL, sellerHandler.Create())

		body := map[string]any{
			"telephone": "9999999",
		}
		request, response := testutil.MakeRequest(http.MethodPost, BASE_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
}

type ServiceMock struct {
	mock.Mock
}

func (svc *ServiceMock) GetAll(c context.Context) ([]domain.Seller, error) {
	args := svc.Called(c)
	return args.Get(0).([]domain.Seller), args.Error(1)
}

func (svc *ServiceMock) Get(ctx context.Context, id int) (domain.Seller, error) {
	args := svc.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (svc *ServiceMock) Save(c context.Context, s domain.Seller) (domain.Seller, error) {
	args := svc.Called(c, s)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (svc *ServiceMock) Update(ctx context.Context, id int, s domain.Seller) (domain.Seller, error) {
	args := svc.Called(ctx, s)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (svc *ServiceMock) Delete(ctx context.Context, id int) error {
	args := svc.Called(ctx, id)
	return args.Error(0)
}
