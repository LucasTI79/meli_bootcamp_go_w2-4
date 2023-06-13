package employee

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("employee not found")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Employee, error)
	Create(ctx context.Context, w domain.Employee) (domain.Employee, error)
	Get(ctx context.Context, id int) (domain.Employee, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(ctx context.Context, w domain.Employee) (domain.Employee, error) {
	eid := s.repository.Exists(ctx, w.CardNumberID)

	if eid {
		return domain.Employee{}, errors.New("employee id already exists")
	}

	id, err := s.repository.Save(ctx, w)
	if err != nil {
		return domain.Employee{}, errors.New("error saving employee")
	}

	w.ID = id

	return w, nil
}

func (s *service) GetAll(ctx context.Context) ([]domain.Employee, error) {
	empl, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, ErrNotFound
	}
	return empl, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Employee, error) {
	e, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Employee{}, ErrNotFound
	}

	return e, nil
}
