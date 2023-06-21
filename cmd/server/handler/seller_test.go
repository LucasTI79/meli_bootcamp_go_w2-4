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
	t.Run("Returns 201 when successful", func(t *testing.T) {
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
