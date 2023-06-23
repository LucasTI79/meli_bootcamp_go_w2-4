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

		_, err := s.Create(context.TODO(), e)
		assert.ErrorIs(t, err, employee.ErrAlreadyExists)

	})
	t.Run("if there is an error creating the employee should return internal server error", func(t *testing.T) {
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
		mockedRepository.On("Save", mock.Anything, e).Return(0, employee.ErrInternalServerError)

		_, err := s.Create(context.TODO(), e)
		assert.ErrorIs(t, err, employee.ErrInternalServerError)

	})
}
func TestGetAllEmployees(t *testing.T) {
	t.Run("returns a list of employees", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)
		es := []domain.Employee{{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Melo",
			WarehouseID:  1,
		}, {
			ID:           2,
			CardNumberID: "128",
			FirstName:    "Mario",
			LastName:     "Melo",
			WarehouseID:  1,
		},
		}
		mockedRepository.On("GetAll", mock.Anything).Return(es, nil)
		employees, err := s.GetAll(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, es, employees)

	})
	t.Run("returns not found when there is a error", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)

		mockedRepository.On("GetAll", mock.Anything).Return([]domain.Employee{}, employee.ErrNotFound)
		_, err := s.GetAll(context.TODO())
		assert.ErrorIs(t, err, employee.ErrNotFound)

	})
}
func TestGetByIdEmployee(t *testing.T) {
	t.Run("return correct employee for valid id", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)

		e := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Melo",
			WarehouseID:  1,
		}

		mockedRepository.On("Get", mock.Anything, e.ID).Return(e, nil)
		employee, err := s.Get(context.TODO(), 1)
		assert.NoError(t, err)
		assert.Equal(t, e, employee)
	})
	t.Run("return error for invalid id", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)

		mockedRepository.On("Get", mock.Anything, 0).Return(domain.Employee{}, employee.ErrNotFound)
		_, err := s.Get(context.TODO(), 0)
		assert.ErrorIs(t, err, employee.ErrNotFound)
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("updates an employee correctly", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)
		currentE := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Melo",
			WarehouseID:  1,
		}
		newE := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Aragao",
			WarehouseID:  1,
		}

		mockedRepository.On("Get", mock.Anything, currentE.ID).Return(currentE, nil)
		mockedRepository.On("Exists", mock.Anything, newE.CardNumberID).Return(false)
		mockedRepository.On("Update", mock.Anything, newE).Return(nil)

		updatedEmployee, err := s.Update(context.TODO(), newE)
		assert.NoError(t, err)
		assert.Equal(t, newE, updatedEmployee)
	})

	t.Run("returns error when employee does not exist", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)
		e := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Aragao",
			WarehouseID:  1,
		}

		mockedRepository.On("Get", mock.Anything, e.ID).Return(domain.Employee{}, employee.ErrNotFound)

		_, err := s.Update(context.TODO(), e)
		assert.ErrorIs(t, err, employee.ErrNotFound)
	})

	t.Run("returns error when card number ID already exists", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)
		currentE := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Melo",
			WarehouseID:  1,
		}
		newE := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Aragao",
			WarehouseID:  1,
		}

		mockedRepository.On("Get", mock.Anything, currentE.ID).Return(currentE, nil)
		mockedRepository.On("Exists", mock.Anything, newE.CardNumberID).Return(true)

		_, err := s.Update(context.TODO(), newE)
		assert.ErrorIs(t, err, employee.ErrAlreadyExists)
	})
	t.Run("return a error when the updates fails ", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)
		currentE := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Melo",
			WarehouseID:  1,
		}
		newE := domain.Employee{
			ID:           1,
			CardNumberID: "126",
			FirstName:    "Lucas",
			LastName:     "Aragao",
			WarehouseID:  1,
		}

		mockedRepository.On("Get", mock.Anything, currentE.ID).Return(currentE, nil)
		mockedRepository.On("Exists", mock.Anything, newE.CardNumberID).Return(false)
		mockedRepository.On("Update", mock.Anything, newE).Return(employee.ErrNotFound)

		_, err := s.Update(context.TODO(), newE)
		assert.ErrorIs(t, err, employee.ErrNotFound)
	})
}
func TestDeleteEmployee(t *testing.T) {
	t.Run("sucessfuly deletes a employee", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)

		mockedRepository.On("Delete", mock.Anything, 1).Return(nil)
		err := s.Delete(context.TODO(), 1)
		assert.NoError(t, err)
	})
	t.Run("return error for invalid id", func(t *testing.T) {
		mockedRepository := RepositoryMock{}
		s := employee.NewService(&mockedRepository)

		mockedRepository.On("Delete", mock.Anything, 0).Return(employee.ErrNotFound)
		err := s.Delete(context.TODO(), 0)
		assert.ErrorIs(t, err, employee.ErrNotFound)
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
	return args.Error(0)
}

func (r *RepositoryMock) Delete(ctx context.Context, id int) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}
