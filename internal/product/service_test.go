package product_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	t.Run("Creates valid product", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		productDTO := product.CreateDTO{
			Desc:       "Sweet potato",
			ExpR:       3,
			FreezeR:    1,
			Height:     200,
			Length:     40,
			NetW:       10,
			Code:       "SWP-1",
			FreezeTemp: 20,
			Width:      100,
			TypeID:     1,
			SellerID:   1,
		}
		expected := *product.MapCreateToDomain(&productDTO)
		expected.ID = 1

		mockRepo.On("Exists", mock.Anything, mock.Anything).Return(false)
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(expected.ID, nil)

		p, err := svc.Create(context.TODO(), productDTO)

		assert.NoError(t, err)
		assert.Equal(t, expected, p)
	})
	t.Run("Doesn't create product if product code exists", func(t *testing.T) {
		t.Skip()
	})
}

func TestRead(t *testing.T) {
	t.Run("Gets all products", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Gets correct product by ID", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns not found for nonexistent ID", func(t *testing.T) {
		t.Skip()
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Updates given fields for existing product", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns not found for nonexistent ID", func(t *testing.T) {
		t.Skip()
	})
}

func TestDelete(t *testing.T) {
	t.Run("Deletes existing product", func(t *testing.T) {
		t.Skip()
	})
	t.Run("Returns not found for nonexistent ID", func(t *testing.T) {
		t.Skip()
	})
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Product, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Product, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, productCode string) bool {
	args := r.Called(ctx, productCode)
	return args.Get(0).(bool)
}

func (r *RepositoryMock) Save(ctx context.Context, p domain.Product) (int, error) {
	args := r.Called(ctx, p)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, p domain.Product) error {
	args := r.Called(ctx, p)
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
