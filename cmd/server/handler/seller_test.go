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
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var SELLER_URL = "/api/v1/sellers"
var SELLER_URL_ID_PATH = "/api/v1/sellers/:id"

func TestCreateSeller(t *testing.T) {
	t.Run("Returns 201 if successful", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		svcMock.On("Save", mock.Anything, expected).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, expected, received.Data)
		fmt.Println(received)
	})

	t.Run("Returns 400 if receives invalid field type", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		body := map[string]any{
			"cid":          "123", // passing CID as string instead of int
			"company_name": "TEST",
			"address":      "test street",
			"telephone":    "9999999",
		}
		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("Returns 422 if receives missing field type", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		body := map[string]any{
			"cid": 123,
		}
		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
	t.Run("Returns 422 if receives missing field type", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		body := map[string]any{
			"cid":          123,
			"company_name": "company",
		}
		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
	t.Run("Returns 422 if receives missing field type", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		body := map[string]any{
			"cid":          123,
			"company_name": "company",
			"address":      "address",
		}
		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
	t.Run("Returns 422 if receives missing field type", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		body := map[string]any{
			"company_name": "company",
			"address":      "address",
			"telephone":    "9999999",
		}
		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, body)
		server.ServeHTTP(response, request)

		fmt.Println(response.Code)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("Returns 409 if CID already exists", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		svcMock.On("Save", mock.Anything, expected).Return(domain.Seller{}, seller.ErrCidAlreadyExists)

		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, expected)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})
	t.Run("Returns 500 if api is not valid", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.POST(SELLER_URL, sellerHandler.Create())

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		svcMock.On("Save", mock.Anything, expected).Return(domain.Seller{}, errors.New(""))

		request, response := testutil.MakeRequest(http.MethodPost, SELLER_URL, expected)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
func TestDeleteSeller(t *testing.T) {
	t.Run("returns 404 when id does not exist", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.DELETE(SELLER_URL_ID_PATH, sellerHandler.Delete())

		idToDelete := 1
		url := fmt.Sprintf("%s/%d", SELLER_URL, idToDelete)

		svcMock.On("Delete", mock.Anything, idToDelete).Return(seller.ErrNotFound)

		request, response := testutil.MakeRequest(http.MethodDelete, url, nil)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)

	})
	t.Run("returns 400 if receives invalid field type", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.DELETE(SELLER_URL, sellerHandler.Delete())

		request, response := testutil.MakeRequest(http.MethodDelete, SELLER_URL, nil)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)

	})
	t.Run("returns if api is not valid", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.DELETE(SELLER_URL_ID_PATH, sellerHandler.Delete())

		idToDelete := 1
		url := fmt.Sprintf("%s/%d", SELLER_URL, idToDelete)

		svcMock.On("Delete", mock.Anything, idToDelete).Return(errors.New(""))

		request, response := testutil.MakeRequest(http.MethodDelete, url, nil)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)

	})
	t.Run("returns 204 when sucessfull", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.DELETE(SELLER_URL_ID_PATH, sellerHandler.Delete())

		idToDelete := 1
		url := fmt.Sprintf("%s/%d", SELLER_URL, idToDelete)

		svcMock.On("Delete", mock.Anything, idToDelete).Return(nil)

		request, response := testutil.MakeRequest(http.MethodDelete, url, nil)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNoContent, response.Code)

	})
}
func TestUpdateSeller(t *testing.T) {
	t.Run("Returns 200 if update is successful", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(SELLER_URL_ID_PATH, sellerHandler.Update())

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		url := fmt.Sprintf("%s/%d", SELLER_URL, expected.ID)
		svcMock.On("Update", mock.Anything, expected.ID, expected).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodPatch, url, expected)
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expected, received.Data)
	})
	t.Run("Returns 409  if cid already exist", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(SELLER_URL_ID_PATH, sellerHandler.Update())

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		url := fmt.Sprintf("%s/%d", SELLER_URL, expected.ID)

		svcMock.On("Update", mock.Anything, expected.ID, expected).Return(domain.Seller{}, seller.ErrCidAlreadyExists)

		request, response := testutil.MakeRequest(http.MethodPatch, url, expected)
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusConflict, response.Code)
	})
	t.Run("Returns 400 if receives invalid field type", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(SELLER_URL_ID_PATH, sellerHandler.Update())

		url := fmt.Sprintf("%s/%s", SELLER_URL, "a")
		request, response := testutil.MakeRequest(http.MethodPatch, url, "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("Returns 422 if receives invalid body", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(SELLER_URL_ID_PATH, sellerHandler.Update())

		url := fmt.Sprintf("%s/%d", SELLER_URL, 1)
		request, response := testutil.MakeRequest(http.MethodPatch, url, "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})
	t.Run("Returns 404 when id does not exist", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(SELLER_URL_ID_PATH, sellerHandler.Update())

		url := fmt.Sprintf("%s/%d", SELLER_URL, 1)
		svcMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(domain.Seller{}, seller.ErrNotFound)

		request, response := testutil.MakeRequest(http.MethodPatch, url, domain.Seller{})
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
	t.Run("Returns 500 api is not valid", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(SELLER_URL_ID_PATH, sellerHandler.Update())

		url := fmt.Sprintf("%s/%d", SELLER_URL, 1)
		svcMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(domain.Seller{}, errors.New(""))

		request, response := testutil.MakeRequest(http.MethodPatch, url, domain.Seller{})
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}
func TestReadSeller(t *testing.T) {
	t.Run("returns 200 if getAll is successful", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.GET(SELLER_URL, sellerHandler.GetAll())

		expected := []domain.Seller{
			{
				ID:          1,
				CID:         123,
				CompanyName: "TEST",
				Address:     "test street",
				Telephone:   "9999999",
			},
			{
				ID:          1,
				CID:         1234,
				CompanyName: "TESTE",
				Address:     "test street",
				Telephone:   "8888888",
			},
		}

		svcMock.On("GetAll", mock.Anything).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodGet, SELLER_URL, "")
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[[]domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.ElementsMatch(t, expected, received.Data)
	})
	t.Run("returns error when the api is not valid", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.GET(SELLER_URL, sellerHandler.GetAll())

		svcMock.On("GetAll", mock.Anything).Return([]domain.Seller{}, errors.New(""))

		request, response := testutil.MakeRequest(http.MethodGet, SELLER_URL, "")
		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
	t.Run("returns 204 when there's no sellers", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.GET(SELLER_URL, sellerHandler.GetAll())

		svcMock.On("GetAll", mock.Anything).Return([]domain.Seller{}, nil)

		request, response := testutil.MakeRequest(http.MethodGet, SELLER_URL, "")
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[[]domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNoContent, response.Code)
		assert.Len(t, received.Data, 0)
	})
	t.Run("returns 200 if get is successful", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.GET(SELLER_URL_ID_PATH, sellerHandler.GetById())

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		url := fmt.Sprintf("%s/%d", SELLER_URL, expected.ID)
		svcMock.On("Get", mock.Anything, expected.ID).Return(expected, nil)

		request, response := testutil.MakeRequest(http.MethodGet, url, "")
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expected, received.Data)
	})

	t.Run("returns error when the id is not valid", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.GET(SELLER_URL_ID_PATH, sellerHandler.GetById())

		url := fmt.Sprintf("%s/%s", SELLER_URL, "aa")

		request, response := testutil.MakeRequest(http.MethodGet, url, "")
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("returns error when the api is not valid", func(t *testing.T) {
		svcMock := SellerServiceMock{}
		sellerHandler := handler.NewSeller(&svcMock)
		server := testutil.CreateServer()
		server.GET(SELLER_URL_ID_PATH, sellerHandler.GetById())

		url := fmt.Sprintf("%s/%d", SELLER_URL, 1)

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		svcMock.On("Get", mock.Anything, expected.ID).Return(domain.Seller{}, seller.ErrNotFound)
		request, response := testutil.MakeRequest(http.MethodGet, url, "")
		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Seller]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

type SellerServiceMock struct {
	mock.Mock
}

func (svc *SellerServiceMock) GetAll(c context.Context) ([]domain.Seller, error) {
	args := svc.Called(c)
	return args.Get(0).([]domain.Seller), args.Error(1)
}

func (svc *SellerServiceMock) Get(ctx context.Context, id int) (domain.Seller, error) {
	args := svc.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (svc *SellerServiceMock) Save(c context.Context, s domain.Seller) (domain.Seller, error) {
	args := svc.Called(c, s)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (svc *SellerServiceMock) Update(ctx context.Context, id int, s domain.Seller) (domain.Seller, error) {
	args := svc.Called(ctx, id, s)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (svc *SellerServiceMock) Delete(ctx context.Context, id int) error {
	args := svc.Called(ctx, id)
	return args.Error(0)
}
