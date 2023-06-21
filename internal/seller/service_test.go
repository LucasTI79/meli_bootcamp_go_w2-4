package seller_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSeller(t *testing.T) {
	t.Run("Create valid seller", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		seller := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		repositoryMock.On("Exists", mock.Anything, 123).Return(false)
		repositoryMock.On("Save", mock.Anything, seller).Return(1, nil)

		received, err := svc.Save(context.TODO(), seller)

		assert.NoError(t, err)
		assert.Equal(t, seller, received)
	})

	t.Run("Create seller with conflict", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		repositoryMock.On("Exists", mock.Anything, 123).Return(true)

		_, err := svc.Save(context.TODO(), expected)

		repositoryMock.AssertNumberOfCalls(t, "Save", 0)
		assert.ErrorIs(t, err, seller.ErrCidAlreadyExists)
	})
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Seller, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Seller), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Seller, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, cid int) bool {
	args := r.Called(ctx, cid)
	return args.Get(0).(bool)
}

func (r *RepositoryMock) Save(ctx context.Context, s domain.Seller) (int, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, s domain.Seller) error {
	args := r.Called(ctx, s)
	return args.Error(1)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(1)
}
