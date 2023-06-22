package warehouse_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/warehouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWarehouse(t *testing.T) {
	t.Run("create valid warehouse", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod1",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		repositoryMock.On("Exists", mock.Anything, "cod1").Return(false)
		repositoryMock.On("Save", mock.Anything, expectedWarehouse).Return(1, nil)

		received, err := svc.Create(context.TODO(), expectedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, received)

	})

}

type RepositoryWarehouseMock struct {
	mock.Mock
}

func (r *RepositoryWarehouseMock) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Warehouse), args.Error(1)
}

func (r *RepositoryWarehouseMock) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Warehouse), args.Error(1)
}

func (r *RepositoryWarehouseMock) Exists(ctx context.Context, WarehouseCode string) bool {
	args := r.Called(ctx, WarehouseCode)
	return args.Get(0).(bool)
}

func (r *RepositoryWarehouseMock) Save(ctx context.Context, s domain.Warehouse) (int, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryWarehouseMock) Update(ctx context.Context, s domain.Warehouse) error {
	args := r.Called(ctx, s)
	return args.Error(1)
}

func (r *RepositoryWarehouseMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(1)
}
