package seller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var ErrRepository = errors.New("error in the repository layer")

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
	t.Run("return domain error when repository fails ", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		repositoryMock.On("Exists", mock.Anything, 123).Return(false)
		repositoryMock.On("Save", mock.Anything, expected).Return(0, ErrRepository)

		_, err := svc.Save(context.TODO(), expected)
		assert.ErrorIs(t, err, seller.ErrRepository)
	})
}
func TestDelete(t *testing.T) {
	t.Run("returns error not found when seller does not exist ", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)
		idToDelete := 1

		repositoryMock.On("Get", mock.Anything, idToDelete).Return(domain.Seller{}, seller.ErrNotFound)
		err := svc.Delete(context.TODO(), idToDelete)

		assert.ErrorIs(t, err, seller.ErrNotFound)
	})
	t.Run("returns no error when sucessfull", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)
		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		repositoryMock.On("Get", mock.Anything, expected.ID).Return(expected, nil)
		repositoryMock.On("Delete", mock.Anything, expected.ID).Return(nil)

		err := svc.Delete(context.TODO(), expected.ID)

		assert.NoError(t, err)

	})
	t.Run("returns domain error when error occurs on repository", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)
		expected := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		repositoryMock.On("Get", mock.Anything, expected.ID).Return(expected, nil)
		repositoryMock.On("Delete", mock.Anything, expected.ID).Return(ErrRepository)

		err := svc.Delete(context.TODO(), expected.ID)

		assert.ErrorIs(t, err, seller.ErrRepository)

	})
}
func TestUpdateSeller(t *testing.T) {
	t.Run("Update valid seller", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		seller := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		sellerUpdate := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "Meli",
			Address:     "Osasco",
			Telephone:   "1134489093",
		}
		repositoryMock.On("Get", mock.Anything, 1).Return(seller, nil)
		repositoryMock.On("Update", mock.Anything, sellerUpdate).Return(nil)

		received, err := svc.Update(context.TODO(), 1, sellerUpdate)

		assert.NoError(t, err)
		assert.Equal(t, sellerUpdate, received)
	})

	t.Run("Update non existent seller", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		sellerMock := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}
		repositoryMock.On("Get", mock.Anything, 1).Return(domain.Seller{}, seller.ErrNotFound)

		_, err := svc.Update(context.TODO(), 1, sellerMock)

		assert.ErrorIs(t, err, seller.ErrNotFound)
	})
}

func TestGetSeller(t *testing.T) {
	t.Run("get valids sellers", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		sellerMock := []domain.Seller{
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

		repositoryMock.On("GetAll", mock.Anything).Return(sellerMock, nil)

		received, err := svc.GetAll(context.TODO())

		assert.ElementsMatch(t, sellerMock, received)
		assert.NoError(t, err)
	})
	t.Run("get invalids sellers", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		repositoryMock.On("GetAll", mock.Anything).Return([]domain.Seller{}, seller.ErrFindSellers)
		_, err := svc.GetAll(context.TODO())

		assert.ErrorIs(t, err, seller.ErrFindSellers)
	})
	t.Run("get valid seller", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		sellerMock := domain.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "TEST",
			Address:     "test street",
			Telephone:   "9999999",
		}

		repositoryMock.On("Get", mock.Anything, 1).Return(sellerMock, nil)

		received, err := svc.Get(context.TODO(), 1)

		assert.Equal(t, sellerMock, received)
		assert.NoError(t, err)
	})

	t.Run("get invalid seller", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := seller.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, 1).Return(domain.Seller{}, seller.ErrNotFound)
		_, err := svc.Get(context.TODO(), 1)

		assert.ErrorIs(t, err, seller.ErrNotFound)
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
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
