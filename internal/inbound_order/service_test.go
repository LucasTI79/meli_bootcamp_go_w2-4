package inboundorder_test

import (
	"context"
	"testing"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	inboundorder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/inbound_order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEmployee(t *testing.T) {
	t.Run("if fields are correct should create a inbound order", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := inboundorder.NewService(&mockedRepository)

		i := domain.InboundOrder{
			ID:             1,
			OrderDate:      time.Now().AddDate(0, 0, 1),
			OrderNumber:    "140",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		mockedRepository.On("Save", mock.Anything, i).Return(1, nil)

		report, err := s.Create(context.TODO(), i)
		assert.NoError(t, err)
		assert.Equal(t, i, report)

	})
	t.Run("should return a error when fails to create a inbound order", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := inboundorder.NewService(&mockedRepository)

		i := domain.InboundOrder{
			ID:             1,
			OrderDate:      time.Now().AddDate(0, 0, 1),
			OrderNumber:    "140",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}

		mockedRepository.On("Save", mock.Anything, i).Return(0, inboundorder.ErrInternalServerError)

		_, err := s.Create(context.TODO(), i)
		assert.ErrorIs(t, err, inboundorder.ErrInternalServerError)

	})

}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Save(ctx context.Context, i domain.InboundOrder) (int, error) {
	args := r.Called(ctx, i)
	return args.Get(0).(int), args.Error(1)
}
