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
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var BASE_URL_BUYER = "/api/v1/buyers"
var BASE_URL_BUYER_ID = "/api/v1/buyers/:id"

func TestBuyerCreate(t *testing.T) {
	t.Run("Returns 201 if successful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := testutil.CreateServer()
		server.POST(BASE_URL_BUYER, buyerHandler.Create())
		expected := domain.BuyerCreate{
			ID:           1,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		svcMock.On("Create", mock.Anything, expected).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPost, BASE_URL_BUYER, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.BuyerCreate]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("Returns 422 if receives missing field type", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := testutil.CreateServer()
		server.POST(BASE_URL_BUYER, buyerHandler.Create())

		body := map[string]any{
			"first_name": "nome",
			"last_name":  "sobrenome",
		}
		request, response := testutil.MakeRequest(http.MethodPost, BASE_URL_BUYER, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("Returns 409 if Card number already exists", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := testutil.CreateServer()
		server.POST(BASE_URL_BUYER, buyerHandler.Create())

		expected := domain.BuyerCreate{
			ID:           1,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		svcMock.On("Create", mock.Anything, expected).Return(domain.BuyerCreate{}, errors.New("buyer already exists"))

		request, response := testutil.MakeRequest(http.MethodPost, BASE_URL_BUYER, expected)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusConflict, response.Code)
	})
}

func TestBuyerGet(t *testing.T) {
	t.Run("Returns 200 if successful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := testutil.CreateServer()
		server.GET(BASE_URL_BUYER, buyerHandler.GetAll())
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

		request, response := testutil.MakeRequest(http.MethodGet, BASE_URL_BUYER, mock.Anything)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[[]domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Returns 404 if id is not existing", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := testutil.CreateServer()
		server.GET(BASE_URL_BUYER, buyerHandler.Get())
		id := 12
		urlId := fmt.Sprintf("%s/%d", BASE_URL_BUYER, id)
		svcMock.On("Get", mock.Anything, 1).Return(domain.Buyer{}, errors.New("buyer not found"))
		request, response := testutil.MakeRequest(http.MethodGet, urlId, mock.Anything)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("Returns 200 if id exists", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := testutil.CreateServer()
		server.GET(BASE_URL_BUYER_ID, buyerHandler.Get())
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		urlId := fmt.Sprintf("%s/%d", BASE_URL_BUYER, expected.ID)
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
		server := testutil.CreateServer()
		server.DELETE(BASE_URL_BUYER_ID, buyerHandler.Delete())
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		urlId := fmt.Sprintf("%s/%d", BASE_URL_BUYER, expected.ID)
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
		server := testutil.CreateServer()
		server.DELETE(BASE_URL_BUYER_ID, buyerHandler.Delete())
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		urlId := fmt.Sprintf("%s/%d", BASE_URL_BUYER, expected.ID)
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
		server := testutil.CreateServer()
		server.PATCH(BASE_URL_BUYER_ID, buyerHandler.Update())
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "crash",
			LastName:     "bandicoot",
		}
		urlId := fmt.Sprintf("%s/%d", BASE_URL_BUYER, expected.ID)
		svcMock.On("Update", mock.Anything, expected).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPatch, urlId, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expected, received.Data)
	})

	t.Run("Returns 404 if update is unsuccessful", func(t *testing.T) {
		svcMock := ServiceMockBuyer{}
		buyerHandler := handler.NewBuyer(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(BASE_URL_BUYER_ID, buyerHandler.Update())
		expected := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "crash",
			LastName:     "bandicoot",
		}
		urlId := fmt.Sprintf("%s/%d", BASE_URL_BUYER, expected.ID)
		svcMock.On("Update", mock.Anything, expected).Return(domain.Buyer{}, errors.New("buyer not updated"))

		request, response := testutil.MakeRequest(http.MethodPatch, urlId, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Buyer]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
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

func (svc *ServiceMockBuyer) Create(c context.Context, s domain.BuyerCreate) (domain.BuyerCreate, error) {
	args := svc.Called(c, s)
	return args.Get(0).(domain.BuyerCreate), args.Error(1)
}

func (svc *ServiceMockBuyer) Update(ctx context.Context, s domain.Buyer, id int) (domain.Buyer, error) {
	args := svc.Called(ctx, s)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (svc *ServiceMockBuyer) Delete(ctx context.Context, id int) error {
	args := svc.Called(ctx, id)
	return args.Error(0)
}
