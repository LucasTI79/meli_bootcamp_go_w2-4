package warehouse

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("warehouse not found")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error)
	// Update(w domain.Warehouse) (domain.Warehouse, error)
	// Delete(id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {

	wcode := s.repository.Exists(ctx, w.WarehouseCode)

	if wcode {
		return domain.Warehouse{}, errors.New("error saving warehouse")
	}

	id, err := s.repository.Save(ctx, w)
	if err != nil {
		return domain.Warehouse{}, errors.New("error saving warehouse")
	}

	w.ID = id

	return w, nil
}

func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	ware, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return ware, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	w, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Warehouse{}, ErrNotFound
	}

	return w, nil
}
