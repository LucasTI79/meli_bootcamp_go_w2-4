package purchaseorder_test

import (
	"context"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	purchaseOrder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/purchase_order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePurchaseOrder(t *testing.T) {
	t.Run("if fields are correct should create a purchase order", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := purchaseOrder.NewService(&mockedRepository)

		p := purchaseOrder.PurchaseOrderDTO{
			ID:              1,
			OrderNumber:     "125",
			OrderDate:       time.Now().AddDate(0, 0, 1),
			TrackingCode:    "124",
			BuyerID:         1,
			ProductRecordID: 1,
			OrderStatusID:   1,
		}

		expected := purchaseOrder.MapPurchaseOrderDTOToDomain(&p)
		mockedRepository.On("Exists", mock.Anything, p.OrderNumber).Return(false)
		mockedRepository.On("Create", mock.Anything, expected).Return(1, nil)

		expected.ID = 1
		purchaseOrder, err := s.Create(context.TODO(), p)
		assert.NoError(t, err)
		assert.Equal(t, expected, purchaseOrder)

	})
	t.Run("if order number already exist", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := purchaseOrder.NewService(&mockedRepository)

		p := purchaseOrder.PurchaseOrderDTO{
			ID:              1,
			OrderNumber:     "125",
			OrderDate:       time.Now().AddDate(0, 0, 1),
			TrackingCode:    "124",
			BuyerID:         1,
			ProductRecordID: 1,
			OrderStatusID:   1,
		}

		mockedRepository.On("Exists", mock.Anything, p.OrderNumber).Return(true)

		_, err := s.Create(context.TODO(), p)
		assert.ErrorIs(t, err, purchaseOrder.ErrAlreadyExists)
	})
	t.Run("if one of the foreign keys are not found", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := purchaseOrder.NewService(&mockedRepository)

		p := purchaseOrder.PurchaseOrderDTO{
			ID:              1,
			OrderNumber:     "125",
			OrderDate:       time.Now().AddDate(0, 0, 1),
			TrackingCode:    "124",
			BuyerID:         1,
			ProductRecordID: 1,
			OrderStatusID:   1,
		}

		expected := purchaseOrder.MapPurchaseOrderDTOToDomain(&p)
		mockedRepository.On("Exists", mock.Anything, p.OrderNumber).Return(false)
		mockedRepository.On("Create", mock.Anything, expected).Return(0, purchaseOrder.ErrFKNotFound)

		expected.ID = 1
		_, err := s.Create(context.TODO(), p)
		assert.ErrorIs(t, err, purchaseOrder.ErrFKNotFound)
	})
	t.Run("if product record is not found", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := purchaseOrder.NewService(&mockedRepository)

		p := purchaseOrder.PurchaseOrderDTO{
			ID:              1,
			OrderNumber:     "125",
			OrderDate:       time.Now().AddDate(0, 0, 1),
			TrackingCode:    "124",
			BuyerID:         1,
			ProductRecordID: 1,
			OrderStatusID:   1,
		}

		expected := purchaseOrder.MapPurchaseOrderDTOToDomain(&p)
		mockedRepository.On("Exists", mock.Anything, p.OrderNumber).Return(false)
		mockedRepository.On("Create", mock.Anything, expected).Return(0, purchaseOrder.ErrProductRecordIDNotFound)

		expected.ID = 1
		_, err := s.Create(context.TODO(), p)
		assert.ErrorIs(t, err, purchaseOrder.ErrProductRecordIDNotFound)
	})
	t.Run("if internal server error occurs", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := purchaseOrder.NewService(&mockedRepository)

		p := purchaseOrder.PurchaseOrderDTO{
			ID:              1,
			OrderNumber:     "125",
			OrderDate:       time.Now().AddDate(0, 0, 1),
			TrackingCode:    "124",
			BuyerID:         1,
			ProductRecordID: 1,
			OrderStatusID:   1,
		}

		expected := purchaseOrder.MapPurchaseOrderDTOToDomain(&p)
		mockedRepository.On("Exists", mock.Anything, p.OrderNumber).Return(false)
		mockedRepository.On("Create", mock.Anything, expected).Return(0, purchaseOrder.ErrInternalServerError)

		expected.ID = 1
		_, err := s.Create(context.TODO(), p)
		assert.ErrorIs(t, err, purchaseOrder.ErrInternalServerError)
	})
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Create(ctx context.Context, p domain.PurchaseOrder) (int, error) {
	args := r.Called(ctx, p)
	return args.Get(0).(int), args.Error(1)
}
func (r *RepositoryMock) Exists(ctx context.Context, orderNumber string) bool {
	args := r.Called(ctx, orderNumber)
	return args.Get(0).(bool)
}
