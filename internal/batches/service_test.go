package batches_test

import (
	"context"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/batches"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	t.Run("should return error when batch number already exists", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := batches.NewService(&repositoryMock)

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(true)

		batch, err := svc.Create(context.Background(), batches.CreateBatches{})
		assert.Equal(t, domain.Batches{}, batch)
		assert.Equal(t, batches.ErrInvalidBatchNumber, err)
	})

	t.Run("create a batches is a successfully", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := batches.NewService(&repositoryMock)

		fakeStruct := batches.CreateBatches{
			BatchNumber:        113,
			CurrentQuantity:    200,
			CurrentTemperature: 20,
			DueDate:            time.Date(2020, time.April, 4, 10, 0, 0, 0, time.UTC),
			InitialQuantity:    10,
			ManufacturingDate:  time.Date(2020, time.April, 4, 10, 0, 0, 0, time.UTC),
			ManufacturingHour:  10,
			MinimumTemperature: 5,
			ProductID:          2,
			SectionID:          1,
		}

		expected := domain.Batches{
			BatchNumber:        113,
			CurrentQuantity:    200,
			CurrentTemperature: 20,
			DueDate:            time.Date(2020, time.April, 4, 10, 0, 0, 0, time.UTC),
			InitialQuantity:    10,
			ManufacturingDate:  time.Date(2020, time.April, 4, 10, 0, 0, 0, time.UTC),
			ManufacturingHour:  10,
			MinimumTemperature: 5,
			ProductID:          2,
			SectionID:          1,
		}

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(0, nil)
		repositoryMock.On("Create", mock.Anything, mock.Anything).Return(fakeStruct, nil)

		batch, err := svc.Create(context.Background(), fakeStruct)
		assert.Equal(t, expected, batch)
		assert.Equal(t, nil, err)
	})

	t.Run("error when save the creste", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := batches.NewService(&repositoryMock)

		fakeStruct := batches.CreateBatches{}

		repositoryMock.On("Exists", mock.Anything, mock.Anything).Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(0, batches.ErrSavingBatch)

		batch, err := svc.Create(context.Background(), fakeStruct)
		assert.Equal(t, domain.Batches{}, batch)
		assert.Equal(t, batches.ErrSavingBatch, err)
	})
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Exists(ctx context.Context, batchNumber int) bool {
	args := r.Called(ctx, batchNumber)
	return args.Bool(0)
}

func (r *RepositoryMock) Create(ctx context.Context, batch domain.Batches) (domain.Batches, error) {
	args := r.Called(ctx, batch)
	return args.Get(0).(domain.Batches), args.Error(1)
}

func (r *RepositoryMock) Save(ctx context.Context, batch domain.Batches) (int, error) {
	args := r.Called(ctx, batch)
	return args.Int(0), args.Error(1)
}
