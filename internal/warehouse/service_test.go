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

	t.Run("create warehouse with conflict (warehousecode alread exixt)", func(t *testing.T) {
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

		repositoryMock.On("Exists", mock.Anything, "cod1").Return(true)

		_, err := svc.Create(context.TODO(), expectedWarehouse)

		repositoryMock.AssertNumberOfCalls(t, "Save", 0)

		assert.ErrorIs(t, err, warehouse.ErrInvalidWarehouseCode)
	})

	t.Run("create error on save warehouse", func(t *testing.T) {
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
		repositoryMock.On("Save", mock.Anything, expectedWarehouse).Return(0, warehouse.ErrorSavingWarehouse)

		_, err := svc.Create(context.TODO(), expectedWarehouse)

		assert.ErrorIs(t, err, warehouse.ErrorSavingWarehouse)

	})
}

func TestGetAllWarehouse(t *testing.T) {
	t.Run("test get all warehouses", func(t *testing.T) {
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

		repositoryMock.On("GetAll", mock.Anything).Return([]domain.Warehouse{expectedWarehouse}, nil)

		received, err := svc.GetAll(context.TODO())

		assert.True(t, len(received) == 1)
		assert.NoError(t, err)
	})

	t.Run("test get all if return some error", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		repositoryMock.On("GetAll", mock.Anything).Return([]domain.Warehouse{}, warehouse.ErrorProcessedData)

		_, err := svc.GetAll(context.TODO())

		assert.ErrorIs(t, err, warehouse.ErrorProcessedData)
	})
}

func TestGetWarehouse(t *testing.T) {
	t.Run("test get warehouse by id", func(t *testing.T) {
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

		repositoryMock.On("Get", mock.Anything, 1).Return(expectedWarehouse, nil)

		received, err := svc.Get(context.TODO(), 1)

		assert.Equal(t, expectedWarehouse, received)
		assert.NoError(t, err)
	})

	t.Run("test det warehouse by id not found", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		expectedWarehouse := domain.Warehouse{}

		repositoryMock.On("Get", mock.Anything, 1).Return(expectedWarehouse, warehouse.ErrNotFound)

		_, err := svc.Get(context.TODO(), 1)

		assert.ErrorIs(t, err, warehouse.ErrNotFound)
	})

}

func TestUpdateWarehouse(t *testing.T) {
	t.Run("test update warehouse", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		expectedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod1",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
			LocalityID:         2,
		}

		repositoryMock.On("Get", mock.Anything, expectedWarehouse.ID).Return(expectedWarehouse, nil)
		repositoryMock.On("Update", mock.Anything, expectedWarehouse).Return(nil)
		repositoryMock.On("Exists", mock.Anything, expectedWarehouse.WarehouseCode).Return(false)

		received, err := svc.Update(context.TODO(), expectedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, expectedWarehouse, received)
	})

	t.Run("test update warehouse with duplicate code", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		existingWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod1",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		updatedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod2",
			Address:            "Rua da Hora",
			Telephone:          "2222222",
			MinimumCapacity:    1,
			MinimumTemperature: 20,
		}

		repositoryMock.On("Get", mock.Anything, existingWarehouse.ID).Return(existingWarehouse, nil)
		repositoryMock.On("Exists", mock.Anything, updatedWarehouse.WarehouseCode).Return(true)
		repositoryMock.On("Update", mock.Anything, updatedWarehouse).Return(warehouse.ErrInvalidWarehouseCode)

		_, err := svc.Update(context.TODO(), updatedWarehouse)

		assert.ErrorIs(t, err, warehouse.ErrInvalidWarehouseCode)
	})

	t.Run("test update warehouse with diferentes code", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		existingWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod1",
			Address:            "Rua da Hora",
			Telephone:          "11111111",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		updatedWarehouse := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod2",
			Address:            "Rua da Hora",
			Telephone:          "2222222",
			MinimumCapacity:    1,
			MinimumTemperature: 20,
		}

		repositoryMock.On("Get", mock.Anything, existingWarehouse.ID).Return(existingWarehouse, nil)
		repositoryMock.On("Exists", mock.Anything, updatedWarehouse.WarehouseCode).Return(false)
		repositoryMock.On("Update", mock.Anything, updatedWarehouse).Return(nil)

		received, err := svc.Update(context.TODO(), updatedWarehouse)

		assert.NoError(t, err)
		assert.Equal(t, updatedWarehouse, received)
	})

	t.Run("test update warehouse ErrorProcessedData", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		expectedWarehouseReceived := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod1",
			Address:            "Rua da Hora",
			Telephone:          "2222222",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		repositoryMock.On("Get", mock.Anything, expectedWarehouseReceived.ID).Return(expectedWarehouseReceived, nil)
		repositoryMock.On("Exists", mock.Anything, expectedWarehouseReceived.WarehouseCode).Return(false)
		repositoryMock.On("Update", mock.Anything, expectedWarehouseReceived).Return(warehouse.ErrorProcessedData)

		received, err := svc.Update(context.TODO(), expectedWarehouseReceived)

		assert.ErrorIs(t, err, warehouse.ErrorProcessedData)
		assert.Equal(t, domain.Warehouse{}, received)
	})

	t.Run("test update warehouse not found", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		expectedWarehouseReceived := domain.Warehouse{
			ID:                 1,
			WarehouseCode:      "cod1",
			Address:            "Rua da Hora",
			Telephone:          "2222222",
			MinimumCapacity:    10,
			MinimumTemperature: 2,
		}

		repositoryMock.On("Get", mock.Anything, expectedWarehouseReceived.ID).Return(domain.Warehouse{}, warehouse.ErrNotFound)

		received, err := svc.Update(context.TODO(), expectedWarehouseReceived)

		assert.ErrorIs(t, err, warehouse.ErrNotFound)
		assert.Equal(t, domain.Warehouse{}, received)

	})
}

func TestDeleteWarehouse(t *testing.T) {
	t.Run("test delete warehouse", func(t *testing.T) {
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

		repositoryMock.On("Delete", context.TODO(), expectedWarehouse.ID).Return(nil)

		err := svc.Delete(context.TODO(), expectedWarehouse.ID)

		assert.NoError(t, err)
	})

	t.Run("test delete warehouse, id not found", func(t *testing.T) {
		repositoryMock := RepositoryWarehouseMock{}
		svc := warehouse.NewService(&repositoryMock)

		expectedWarehouse := domain.Warehouse{}

		repositoryMock.On("Delete", context.TODO(), expectedWarehouse.ID).Return(warehouse.ErrNotFound)

		err := svc.Delete(context.TODO(), expectedWarehouse.ID)

		assert.ErrorIs(t, err, warehouse.ErrNotFound)
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
	return args.Error(0)
}

func (r *RepositoryWarehouseMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
