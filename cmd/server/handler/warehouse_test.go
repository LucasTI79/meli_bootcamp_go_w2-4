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

var WAREHOUSE_URL = "/api/v1/warehouses"
var WAREHOUSE_URL_ID = "/api/v1/warehouses/:id"

func TestWarehouseCreate(t *testing.T) {
	t.Run("test create, if successful, return 201", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.POST(WAREHOUSE_URL, warehouseHandler.Create())

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}
		request, response := testutil.MakeRequest(http.MethodPost, WAREHOUSE_URL, expectedWarehouse)

		svcMock.On("Create", mock.Anything, expectedWarehouse).Return(expectedWarehouse, nil)

		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Warehouse]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusCreated)
		assert.Equal(t, expectedWarehouse, received.Data)

	})

	t.Run("test create, if warehouse is empty, return 422", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.POST(WAREHOUSE_URL, warehouseHandler.Create())

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		request, response := testutil.MakeRequest(http.MethodPost, WAREHOUSE_URL, expectedWarehouse)

		server.ServeHTTP(response, request)

		var received testutil.ErrorResponse
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusUnprocessableEntity)
		assert.Equal(t, received.Message, handler.ErrWarehouseEmpty)
	})

	// t.Run("test create, if warehouse is exist, return 409", func(t *testing.T) {
	// 	svcMock := ServiceWarehouseMock{}
	// 	warehouseHandler := handler.NewWarehouse(&svcMock)
	// 	server := testutil.CreateServer()
	// 	server.POST(WAREHOUSE_URL, warehouseHandler.Create())

	// 	expectedWarehouse := domain.Warehouse{
	// 		ID:                 1,
	// 		WarehouseCode:      "cod",
	// 		Address:            "Rua da Hora",
	// 		Telephone:          "11111111",
	// 		MinimumCapacity:    10,
	// 		MinimumTemperature: 2,
	// 	}

	// 	request, response := testutil.MakeRequest(http.MethodPost, WAREHOUSE_URL, expectedWarehouse)

	// 	svcMock.On("Create", mock.Anything, expectedWarehouse).Return(domain.Warehouse{}, handler.ErrAlrearyExist)

	// 	server.ServeHTTP(response, request)

	// 	var received testutil.ErrorResponse
	// 	json.Unmarshal(response.Body.Bytes(), &received)

	// 	assert.Equal(t, response.Code, http.StatusConflict)
	// 	assert.Equal(t, received.Message, "warehouse already exists")

	// })

}

func TestWarehouseGetAll(t *testing.T) {
	t.Run("test if getall return a list that warehouse", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.GET(WAREHOUSE_URL, warehouseHandler.GetAll())

		expectedWarehouse := []domain.Warehouse{
			{
				ID:                 1,
				WarehouseCode:      "cod",
				Address:            "Rua da Hora",
				Telephone:          "11111111",
				MinimumCapacity:    10,
				MinimumTemperature: 2,
			},
			{
				ID:                 2,
				WarehouseCode:      "cod2",
				Address:            "Rua da Hora",
				Telephone:          "11111111",
				MinimumCapacity:    10,
				MinimumTemperature: 2,
			},
		}

		request, response := testutil.MakeRequest(http.MethodGet, WAREHOUSE_URL, "")

		svcMock.On("GetAll", mock.Anything).Return(expectedWarehouse, nil)

		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Warehouse]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusOK)

	})

	t.Run("test if getall not return a error if a length was zero", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.GET(WAREHOUSE_URL, warehouseHandler.GetAll())

		expectedWarehouse := []domain.Warehouse{}

		request, response := testutil.MakeRequest(http.MethodGet, WAREHOUSE_URL, "")

		svcMock.On("GetAll", mock.Anything).Return(expectedWarehouse, nil)

		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Warehouse]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusNoContent)

	})

	t.Run("test if getall return a error 500", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.GET(WAREHOUSE_URL, warehouseHandler.GetAll())

		request, response := testutil.MakeRequest(http.MethodGet, WAREHOUSE_URL, "")

		svcMock.On("GetAll", mock.Anything).Return(nil, nil)

		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Warehouse]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusInternalServerError)

	})
}

func TestWarehouseGet(t *testing.T) {
	t.Run("test get, when the id is valid", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.GET(WAREHOUSE_URL_ID, warehouseHandler.Get())

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		url := fmt.Sprintf("%s/%d", WAREHOUSE_URL, expectedWarehouse.ID)

		svcMock.On("Get", mock.Anything, 1).Return(expectedWarehouse, nil)
		request, response := testutil.MakeRequest(http.MethodGet, url, nil)

		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Warehouse]
		json.Unmarshal(response.Body.Bytes(), &received)
		fmt.Println(request.Body)
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, expectedWarehouse, received.Data)

	})
	t.Run("test get, when the id is invalid - 400", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.GET(WAREHOUSE_URL, warehouseHandler.Get())

		request, response := testutil.MakeRequest(http.MethodGet, WAREHOUSE_URL, "")

		svcMock.On("Get", mock.Anything, 1).Return(nil, nil)

		server.ServeHTTP(response, request)

		var received testutil.ErrorResponse
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusBadRequest)
		assert.Equal(t, received.Message, handler.ErrInvalidID)

	})

	t.Run("test get, when the warehouse not exist return 404", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.GET(WAREHOUSE_URL_ID, warehouseHandler.Get())

		expectedWarehouse2 := domain.Warehouse{}

		svcMock.On("Get", mock.Anything, 2).Return(expectedWarehouse2, handler.ErrWarehouseNotFound)
		request, response := testutil.MakeRequest(http.MethodGet, WAREHOUSE_URL, "")

		server.ServeHTTP(response, request)
		var received testutil.ErrorResponse
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusNotFound)

	})
}

func TestWarehouseUpdate(t *testing.T) {
	t.Run("test update, when the id is valid", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(WAREHOUSE_URL_ID, warehouseHandler.Update())

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		url := fmt.Sprintf("%s/%d", WAREHOUSE_URL, expectedWarehouse.ID)

		svcMock.On("Update", mock.Anything, expectedWarehouse).Return(expectedWarehouse, nil)
		request, response := testutil.MakeRequest(http.MethodPatch, url, expectedWarehouse)

		server.ServeHTTP(response, request)

		var received testutil.SuccessResponse[domain.Warehouse]
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, expectedWarehouse, received.Data)

	})

	t.Run("test update, when the id is not found- 404", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(WAREHOUSE_URL_ID, warehouseHandler.Update())

		expectedWarehouse2 := domain.Warehouse{}

		svcMock.On("Get", mock.Anything, 2).Return(expectedWarehouse2, handler.ErrWarehouseNotFound)
		request, response := testutil.MakeRequest(http.MethodPatch, WAREHOUSE_URL, "")

		server.ServeHTTP(response, request)
		var received testutil.ErrorResponse
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	// t.Run("test update, when the id must to be unique - 409", func(t *testing.T) {
	// 	svcMock := ServiceWarehouseMock{}
	// 	warehouseHandler := handler.NewWarehouse(&svcMock)
	// 	server := testutil.CreateServer()
	// 	server.PATCH(WAREHOUSE_URL_ID, warehouseHandler.Update())

	// 	expectedWarehouse := domain.Warehouse{
	// 		ID:                 1,
	// 		WarehouseCode:      "cod",
	// 		Address:            "Rua da Hora",
	// 		Telephone:          "11111111",
	// 		MinimumCapacity:    10,
	// 		MinimumTemperature: 2,
	// 	}

	// 	url := fmt.Sprintf("%s/%d", WAREHOUSE_URL, expectedWarehouse.ID)

	// 	svcMock.On("Update", mock.Anything, expectedWarehouse).Return(expectedWarehouse, handler.ErrWarehouseCodeUnique)
	// 	svcMock.On("Exists", mock.Anything, expectedWarehouse.WarehouseCode).Return(true)
	// 	request, response := testutil.MakeRequest(http.MethodPatch, url, expectedWarehouse)

	// 	server.ServeHTTP(response, request)

	// 	var received testutil.ErrorResponse
	// 	json.Unmarshal(response.Body.Bytes(), &received)

	// 	assert.Equal(t, response.Code, http.StatusConflict)
	// 	assert.Equal(t, received.Message, handler.ErrWarehouseCodeUnique)

	// })

	t.Run("test update, when the return an error in the server 422", func(t *testing.T) {
		svcMock := ServiceWarehouseMock{}
		warehouseHandler := handler.NewWarehouse(&svcMock)
		server := testutil.CreateServer()
		server.PATCH(WAREHOUSE_URL_ID, warehouseHandler.Update())

		expectedWarehouse2 := domain.Warehouse{}

		svcMock.On("Get", mock.Anything, 2).Return(expectedWarehouse2, handler.ErrUnprocessableEntity)
		request, response := testutil.MakeRequest(http.MethodPatch, WAREHOUSE_URL_ID, "")

		server.ServeHTTP(response, request)
		var received testutil.ErrorResponse
		json.Unmarshal(response.Body.Bytes(), &received)

		assert.Equal(t, response.Code, http.StatusUnprocessableEntity)
	})
}

type ServiceWarehouseMock struct {
	mock.Mock
}

func (r *ServiceWarehouseMock) Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	args := r.Called(ctx, w)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (r *ServiceWarehouseMock) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (r *ServiceWarehouseMock) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (r *ServiceWarehouseMock) Exists(ctx context.Context, WarehouseCode string) bool {
	args := r.Called(ctx, WarehouseCode)
	return args.Get(0).(bool)
}

func (r *ServiceWarehouseMock) Save(ctx context.Context, s domain.Warehouse) (int, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}

func (r *ServiceWarehouseMock) Update(ctx context.Context, s domain.Warehouse) (domain.Warehouse, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (r *ServiceWarehouseMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
