package warehouse

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound             = errors.New("warehouse not found")
	ErrInvalidWarehouseCode = errors.New("warehouse code has to be unique")
	ErrorSavingWarehouse    = errors.New("error saving warehouse")
	ErrorProcessedData      = errors.New("action could not be processed correctly due to invalid data provided")
)

// Service is the interface for warehouse operations.
type Service interface {
	Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error)
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
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
//
//	@summary	Creates a new warehouse.
//	@param		warehouse	body	domain.Warehouse	true	"Warehouse object"
//	@return		201 {object} domain.Warehouse
//	@return		400 {object} BadRequestError "Invalid request"
//	@tags		Warehouse
func (s *service) Create(ctx context.Context, w domain.Warehouse) (domain.Warehouse, error) {
	wcode := s.repository.Exists(ctx, w.WarehouseCode)
	if wcode {
		return domain.Warehouse{}, ErrInvalidWarehouseCode
	}

	id, err := s.repository.Save(ctx, w)
	if err != nil {
		return domain.Warehouse{}, ErrorSavingWarehouse
	}

	w.ID = id

	return w, nil
}

// GetAll retrieves all warehouses.
//
//	@summary	Retrieves all warehouses.
//	@return		200 {array} Warehouse
//	@tags		Warehouse
func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	ware, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, ErrorProcessedData
	}

	return ware, nil
}

// Get retrieves a warehouse by its ID.
//
//	@summary	Retrieves a warehouse by ID.
//	@param		id	path	int	true	"Warehouse ID"
//	@return		200 {object} domain.Warehouse
//	@return		404 {object} NotFoundError "Warehouse not found"
//	@tags		Warehouse
func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	w, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Warehouse{}, ErrNotFound
	}

	return w, nil
}

// Update updates an existing warehouse.
//
//	@summary	Updates an existing warehouse.
//	@param		id			path	int					true	"Warehouse ID"
//	@param		warehouse	body	domain.Warehouse	true	"Updated warehouse object"
//	@return		200 {object} domain.Warehouse
//	@return		400 {object} BadRequestError "Invalid request"
//	@return		404 {object} NotFoundError "Warehouse not found"
//	@tags		Warehouse
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
		if w.WarehouseCode != currentWarehouse.WarehouseCode {
			if s.repository.Exists(ctx, w.WarehouseCode) {
				return domain.Warehouse{}, ErrInvalidWarehouseCode
			}
			currentWarehouse.WarehouseCode = w.WarehouseCode
		}
	}

	if w.MinimumCapacity != 0 {
		currentWarehouse.MinimumCapacity = w.MinimumCapacity
	}
	if w.MinimumTemperature != 0 {
		currentWarehouse.MinimumTemperature = w.MinimumTemperature
	}
	if w.LocalityID != 0 {
		currentWarehouse.LocalityID = w.LocalityID
	}

	err = s.repository.Update(ctx, currentWarehouse)
	if err != nil {
		return domain.Warehouse{}, ErrorProcessedData
	}

	return currentWarehouse, nil
}

// Delete deletes a warehouse by its ID.
//
//	@summary	Deletes a warehouse by ID.
//	@param		id	path	int	true	"Warehouse ID"
//	@return		204 "No Content"
//	@return		404 {object} NotFoundError "Warehouse not found"
//	@tags		Warehouse
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	return nil
}
