package employee

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound            = errors.New("employee not found")
	ErrAlreadyExists       = errors.New("employee id already exists")
	ErrInternalServerError = errors.New("internal server error")
)

// Service define a interface para o serviço de funcionários.
type Service interface {
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Create(ctx context.Context, e domain.Employee) (domain.Employee, error)
	Get(ctx context.Context, id int) (domain.Employee, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, e domain.Employee) (domain.Employee, error)
	GetInboundReport(ctx context.Context, id int) ([]domain.InboundReport, error)
}

type service struct {
	repository Repository
}

// NewService cria uma nova instância do serviço de funcionários.
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// Create cria um novo funcionário.
func (s *service) Create(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	eid := s.repository.Exists(ctx, e.CardNumberID)

	if eid {
		return domain.Employee{}, ErrAlreadyExists
	}

	id, err := s.repository.Save(ctx, e)
	if err != nil {
		return domain.Employee{}, ErrInternalServerError
	}

	e.ID = id

	return e, nil
}

// GetAll obtém todas as informações dos funcionários.
func (s *service) GetAll(ctx context.Context) ([]domain.Employee, error) {
	empl, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, ErrNotFound
	}
	return empl, nil
}

// Get obtém as informações de um funcionário pelo ID.
func (s *service) Get(ctx context.Context, id int) (domain.Employee, error) {
	e, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Employee{}, ErrNotFound
	}

	return e, nil
}

// Delete remove um funcionário.
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	return nil
}

// Update atualiza as informações de um funcionário.
func (s *service) Update(ctx context.Context, e domain.Employee) (domain.Employee, error) {
	currentEmployee, err := s.repository.Get(ctx, e.ID)

	if err != nil {
		return domain.Employee{}, ErrNotFound
	}

	if e.FirstName != "" {
		currentEmployee.FirstName = e.FirstName
	}

	if e.LastName != "" {
		currentEmployee.LastName = e.LastName
	}

	if e.CardNumberID != "" {

		ecode := s.repository.Exists(ctx, e.CardNumberID)
		if ecode {
			return domain.Employee{}, ErrAlreadyExists
		}
		currentEmployee.CardNumberID = e.CardNumberID
	}

	if e.WarehouseID != 0 {
		currentEmployee.WarehouseID = e.WarehouseID
	}

	err = s.repository.Update(ctx, currentEmployee)
	if err != nil {
		return domain.Employee{}, ErrNotFound
	}
	return currentEmployee, nil
}

func (s *service) GetInboundReport(ctx context.Context, id int) ([]domain.InboundReport, error) {
	if id == 0 {
		e, err := s.repository.GetAllInboundReports(ctx)
		if err != nil {
			return []domain.InboundReport{}, err
		}
		return e, nil

	}
	e, err := s.repository.GetInboundReport(ctx, id)
	if err != nil {
		return []domain.InboundReport{}, err
	}

	return []domain.InboundReport{e}, nil
}
