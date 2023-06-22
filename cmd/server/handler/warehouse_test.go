package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var WAREHOUSE_URL = "/api/v1/warehouses"

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
