package buyer_test

import (
	"context"
	"errors"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBuyer(t *testing.T) {
	t.Run("Create valid buyer", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		buyer := domain.Buyer{
			ID:           1,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		b := domain.BuyerCreate{
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}

		repositoryMock.On("Exists", mock.Anything, b.CardNumberID).Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(1, nil)

		received, err := svc.Create(context.TODO(), b)

		assert.NoError(t, err)
		assert.Equal(t, buyer, received)
	})

	t.Run("Create Buyer with conflict", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		buyerMock := domain.BuyerCreate{
			ID:           1,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}

		repositoryMock.On("Exists", mock.Anything, "123").Return(true)

		_, err := svc.Create(context.TODO(), buyerMock)

		repositoryMock.AssertNumberOfCalls(t, "Save", 0)
		assert.ErrorIs(t, err, buyer.ErrAlreadyExists)
	})

	t.Run("Create Buyer with error", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		buyerMock := domain.BuyerCreate{
			ID:           1,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}

		repositoryMock.On("Exists", mock.Anything, "123").Return(false)
		repositoryMock.On("Save", mock.Anything, mock.Anything).Return(0, buyer.ErrSavingBuyer)

		_, err := svc.Create(context.TODO(), buyerMock)

		assert.ErrorIs(t, err, buyer.ErrSavingBuyer)
	})
}

func TestGetbyIDBuyer(t *testing.T) {
	t.Run("Find by non existent Id", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		repositoryMock.On("Get", mock.Anything, 12).Return(domain.Buyer{}, buyer.ErrNotFound)

		returned, err := svc.Get(context.TODO(), 12)

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "buyer not found")
		assert.Equal(t, domain.Buyer{}, returned)
	})

	t.Run("Find by existent Id", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)
		buyerMock := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}

		repositoryMock.On("Get", mock.Anything, 12).Return(buyerMock, nil)

		returned, err := svc.Get(context.TODO(), 12)

		assert.NoError(t, err)
		assert.Equal(t, buyerMock, returned)
	})

}
func TestGetBuyer(t *testing.T) {
	t.Run("Find all buyer", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		buyerMock := []domain.Buyer{
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

		repositoryMock.On("GetAll", mock.Anything).Return(buyerMock, nil)

		received, err := svc.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.ElementsMatch(t, buyerMock, received)
	})

	t.Run("Find all buyer with error", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		repositoryMock.On("GetAll", mock.Anything).Return([]domain.Buyer{}, buyer.ErrNotFound)

		_, err := svc.GetAll(context.TODO())

		assert.ErrorIs(t, err, buyer.ErrNotFound)
	})
}

func TestUpdateBuyer(t *testing.T) {
	t.Run("Update existent buyer", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)
		buyerMock := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "nome",
			LastName:     "sobrenome",
		}
		buyerUpdate := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "lucas",
			LastName:     "ganda",
		}

		repositoryMock.On("Get", mock.Anything, 12).Return(buyerMock, nil)
		repositoryMock.On("Update", mock.Anything, buyerUpdate).Return(nil)

		returned, err := svc.Update(context.TODO(), buyerUpdate, 12)

		assert.NoError(t, err)
		assert.Equal(t, buyerUpdate, returned)
	})

	t.Run("Update non existent buyer", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)
		buyerUpdate := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "lucas",
			LastName:     "ganda",
		}

		repositoryMock.On("Get", mock.Anything, 12).Return(domain.Buyer{}, buyer.ErrNotFound)
		returned, err := svc.Update(context.TODO(), buyerUpdate, 12)

		repositoryMock.AssertNumberOfCalls(t, "Update", 0)
		assert.Error(t, err)
		assert.Equal(t, domain.Buyer{}, returned)
	})

	t.Run("Update buyer with error", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)
		buyerPreUpdate := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "usuario",
			LastName:     "de teste",
		}
		buyerUpdate := domain.Buyer{
			ID:           12,
			CardNumberID: "123",
			FirstName:    "lucas",
			LastName:     "ganda",
		}

		repositoryMock.On("Get", mock.Anything, 12).Return(buyerPreUpdate, nil)
		repositoryMock.On("Update", mock.Anything, buyerUpdate).Return(buyer.ErrNotFound)
		_, err := svc.Update(context.TODO(), buyerUpdate, 12)

		assert.ErrorIs(t, err, buyer.ErrNotFound)
	})
}

func TestDeleteBuyer(t *testing.T) {
	t.Run("Delete existent buyer", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		repositoryMock.On("Delete", mock.Anything, 12).Return(nil)

		err := svc.Delete(context.TODO(), 12)

		assert.NoError(t, err)
	})

	t.Run("Delete non existent buyer", func(t *testing.T) {
		repositoryMock := RepositoryMock{}
		svc := buyer.NewService(&repositoryMock)

		repositoryMock.On("Delete", mock.Anything, 12).Return(buyer.ErrNotFound)

		err := svc.Delete(context.TODO(), 12)

		assert.Error(t, err)
		assert.Equal(t, errors.New("buyer not found"), err)
	})
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Buyer), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Buyer, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, cardNumberID string) bool {
	args := r.Called(ctx, cardNumberID)
	return args.Get(0).(bool)
}

func (r *RepositoryMock) Save(ctx context.Context, s domain.Buyer) (int, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, s domain.Buyer) error {
	args := r.Called(ctx, s)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
