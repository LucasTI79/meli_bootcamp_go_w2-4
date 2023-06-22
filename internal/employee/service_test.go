package employee_test

import (
	"context"
	"testing"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEmployee(t *testing.T) {
	t.Run("if fields are correct should create a employee ", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)

		e := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Melo",
			WarehouseID:  1,
		}

		mockedRepository.On("Exists", mock.Anything, "126").Return(false)
		mockedRepository.On("Save", mock.Anything, e).Return(1, nil)

		createdEmployee, err := s.Create(context.TODO(), e)
		assert.NoError(t, err)
		assert.Equal(t, e, createdEmployee)

	})
	t.Run("if card_number_id already exists should return a conflict error", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)

		e := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Melo",
			WarehouseID:  1,
		}

		mockedRepository.On("Exists", mock.Anything, "126").Return(true)
		mockedRepository.On("Save", mock.Anything, e).Return(nil, employee.ErrAlreadyExists)

		_, err := s.Create(context.TODO(), e)
		assert.ErrorIs(t, err, employee.ErrAlreadyExists)

	})
}

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) GetAll(ctx context.Context) ([]domain.Employee, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Employee), args.Error(1)
}

func (r *RepositoryMock) Get(ctx context.Context, id int) (domain.Employee, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(domain.Employee), args.Error(1)
}

func (r *RepositoryMock) Exists(ctx context.Context, cid string) bool {
	args := r.Called(ctx, cid)
	return args.Get(0).(bool)
}

func (r *RepositoryMock) Save(ctx context.Context, s domain.Employee) (int, error) {
	args := r.Called(ctx, s)
	return args.Get(0).(int), args.Error(1)
}

func (r *RepositoryMock) Update(ctx context.Context, s domain.Employee) error {
	args := r.Called(ctx, s)
	return args.Error(1)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(1)
}
