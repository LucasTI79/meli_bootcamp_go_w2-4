package product_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	t.Run("Creates valid product", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		dto := product.CreateDTO{
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
		expected := *product.MapCreateToDomain(&dto)
		expected.ID = 1

		mockRepo.On("Exists", mock.Anything, mock.Anything).Return(false)
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(expected.ID, nil)

		p, err := svc.Create(context.TODO(), dto)

		assert.NoError(t, err)
		assert.Equal(t, expected, p)
	})
	t.Run("Doesn't create product if product code exists", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		dto := product.CreateDTO{
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
		var expectedErr *product.ErrInvalidProductCode

		mockRepo.On("Exists", mock.Anything, mock.Anything).Return(true)

		_, err := svc.Create(context.TODO(), dto)

		assert.ErrorAs(t, err, &expectedErr)
	})
}

func TestRead(t *testing.T) {
	t.Run("Gets all products", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		expected := getTestProducts()

		mockRepo.On("GetAll", mock.Anything).Return(expected, nil)
		ps, err := svc.GetAll(context.TODO())

		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, ps)
	})
	t.Run("Gets correct product by ID", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		expected := getTestProducts()[0]

		mockRepo.On("Get", mock.Anything, expected.ID).Return(expected, nil)
		p, err := svc.Get(context.TODO(), expected.ID)

		assert.NoError(t, err)
		assert.Equal(t, expected, p)
	})
	t.Run("Returns not found for nonexistent ID", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		p := getTestProducts()[0]
		var expectedErr *product.ErrNotFound

		mockRepo.On("Get", mock.Anything, p.ID).Return(domain.Product{}, product.NewErrNotFound(p.ID))
		_, err := svc.Get(context.TODO(), p.ID)

		assert.ErrorAs(t, err, &expectedErr)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Updates given fields for existing product", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		toUpdate := domain.Product{
			ID:             1,
			Description:    "Sweet potato",
			ExpirationRate: 3,
			FreezingRate:   1,
			Height:         200,
			Length:         40,
			Netweight:      10,
			ProductCode:    "SWP-1",
			RecomFreezTemp: 20,
			Width:          100,
			ProductTypeID:  1,
			SellerID:       1,
		}
		updates := product.UpdateDTO{
			Desc:     *optional.FromVal("Garlic"),
			Height:   *optional.FromVal[float32](42),
			SellerID: *optional.FromVal(10),
		}
		expected := domain.Product{
			ID:             1,
			Description:    "Garlic",
			ExpirationRate: 3,
			FreezingRate:   1,
			Height:         42,
			Length:         40,
			Netweight:      10,
			ProductCode:    "SWP-1",
			RecomFreezTemp: 20,
			Width:          100,
			ProductTypeID:  1,
			SellerID:       10,
		}

		mockRepo.On("Get", mock.Anything, toUpdate.ID).Return(toUpdate, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		received, err := svc.Update(context.TODO(), toUpdate.ID, updates)

		assert.NoError(t, err)
		assert.Equal(t, expected, received)
	})
	t.Run("Update fails if product code is not unique", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		toUpdate := getTestProducts()[1]
		updates := product.UpdateDTO{Code: *optional.FromVal("SWP-1")}
		var expectedErr *product.ErrInvalidProductCode

		mockRepo.On("Get", mock.Anything, toUpdate.ID).Return(toUpdate, nil)
		mockRepo.On("Exists", mock.Anything, updates.Code.Val).Return(true)

		_, err := svc.Update(context.TODO(), toUpdate.ID, updates)

		assert.ErrorAs(t, err, &expectedErr)
	})
	t.Run("Update succeds if product code doesn't change", func(t *testing.T) {
		mockRepo := RepositoryMock{}
		svc := product.NewService(&mockRepo)

		toUpdate := getTestProducts()[1]
		updates := product.UpdateDTO{Code: *optional.FromVal(toUpdate.ProductCode)}

		mockRepo.On("Get", mock.Anything, toUpdate.ID).Return(toUpdate, nil)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		received, err := svc.Update(context.TODO(), toUpdate.ID, updates)

		assert.NoError(t, err)
		assert.Equal(t, toUpdate, received)
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

func getTestProducts() []domain.Product {
	return []domain.Product{
		{
			ID:             1,
			Description:    "abc",
			ExpirationRate: 1,
			FreezingRate:   2,
			Height:         3,
			Length:         4,
			Netweight:      5,
			ProductCode:    "PRODUCT-1",
			RecomFreezTemp: 6,
			Width:          7,
			ProductTypeID:  8,
			SellerID:       9,
		},
		{
			ID:             2,
			Description:    "cde",
			ExpirationRate: 1,
			FreezingRate:   2,
			Height:         3,
			Length:         4,
			Netweight:      5,
			ProductCode:    "PRODUCT-2",
			RecomFreezTemp: 6,
			Width:          7,
			ProductTypeID:  8,
			SellerID:       9,
		},
	}
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
