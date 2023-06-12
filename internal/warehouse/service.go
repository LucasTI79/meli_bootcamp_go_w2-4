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

// Service is the interface for warehouse operations.
type Service interface {
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error)
	Update(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

// NewService creates a new warehouse service with the given repository.
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// Create creates a new warehouse.
// @summary Creates a new warehouse.
// @param warehouse body Warehouse true "Warehouse object"
// @return 201 {object} Warehouse
// @return 400 {object} BadRequestError "Invalid request"
// @tags Warehouse
func (s *service) Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	wcode := s.repository.Exists(ctx, w.WarehouseCode)
	if wcode {
		return domain.Warehouse{}, errors.New("warehouse to be single")
	}

	id, err := s.repository.Save(ctx, w)
	if err != nil {
		return domain.Warehouse{}, errors.New("error saving warehouse")
	}

	w.ID = id

	return w, nil
}

// GetAll retrieves all warehouses.
// @summary Retrieves all warehouses.
// @return 200 {array} Warehouse
// @tags Warehouse
func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	ware, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, ErrNotFound
	}

	return ware, nil
}

// Get retrieves a warehouse by its ID.
// @summary Retrieves a warehouse by ID.
// @param id path int true "Warehouse ID"
// @return 200 {object} Warehouse
// @return 404 {object} NotFoundError "Warehouse not found"
// @tags Warehouse
func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	w, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Warehouse{}, ErrNotFound
	}

	return w, nil
}

// Update updates an existing warehouse.
// @summary Updates an existing warehouse.
// @param id path int true "Warehouse ID"
// @param warehouse body Warehouse true "Updated warehouse object"
// @return 200 {object} Warehouse
// @return 400 {object} BadRequestError "Invalid request"
// @return 404 {object} NotFoundError "Warehouse not found"
// @tags Warehouse
func (s *service) Update(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	currentWarehouse, err := s.repository.Get(ctx, w.ID)
	if err != nil {
		return domain.Warehouse{}, ErrNotFound
	}

	if w.Address != "" {
		currentWarehouse.Address = w.Address
	}
	if w.Telephone != "" {
		currentWarehouse.Telephone = w.Telephone
	}
	if w.WarehouseCode != "" {

		wcode := s.repository.Exists(ctx, w.WarehouseCode)
		if wcode {
			return domain.Warehouse{}, errors.New("warehouse must be unique")
		}
		currentWarehouse.WarehouseCode = w.WarehouseCode
	}
	if w.MinimumCapacity != 0 {
		currentWarehouse.MinimumCapacity = w.MinimumCapacity
	}
	if w.MinimumTemperature != 0 {
		currentWarehouse.MinimumTemperature = w.MinimumTemperature
	}

	err = s.repository.Update(ctx, currentWarehouse)
	if err != nil {
		return domain.Warehouse{}, ErrNotFound
	}

	return currentWarehouse, nil
}

// Delete deletes a warehouse by its ID.
// @summary Deletes a warehouse by ID.
// @param id path int true "Warehouse ID"
// @return 204 "No Content"
// @return 404 {object} NotFoundError "Warehouse not found"
// @tags Warehouse
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	return nil
}
