package buyer

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Error definitions
var (
	ErrNotFound = errors.New("Buyer not found")
	ErrGeneric  = errors.New("")
)

// Service is the buyer service interface
type Service interface {
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	Get(ctx context.Context, id int) (domain.Buyer, error)
	Create(ctx context.Context, b domain.BuyerCreate) (domain.BuyerCreate, error)
	Update(ctx context.Context, b domain.Buyer, id int) (domain.Buyer, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

// Create creates a new buyer
// @Summary Create a new buyer
// @Description Creates a new buyer with the provided data
// @Tags buyers
// @Accept  json
// @Produce  json
// @Param buyer body domain.BuyerCreate true "Buyer data"
// @Success 200 {object} domain.BuyerCreate
// @Failure 400 {string} string "Error saving buyer"
// @Router /buyers [post]
func (s *service) Create(ctx context.Context, b domain.BuyerCreate) (domain.BuyerCreate, error) {
	ex := s.repository.Exists(ctx, b.CardNumberID)
	if ex {
		return domain.BuyerCreate{}, errors.New("Buyer already exists")
	}
	id, err := s.repository.Save(ctx, b)
	if err != nil {
		return domain.BuyerCreate{}, errors.New("Error saving buyer")
	}
	b.ID = id
	return b, nil
}

// Update updates an existing buyer
// @Summary Update an existing buyer
// @Description Updates an existing buyer with the provided data
// @Tags buyers
// @Accept  json
// @Produce  json
// @Param id path int true "Buyer ID"
// @Param buyer body domain.Buyer true "Buyer data"
// @Success 200 {object} domain.Buyer
// @Failure 400 {string} string "Error getting buyer"
// @Failure 404 {string} string "Buyer not found"
// @Router /buyers/{id} [put]
func (s *service) Update(ctx context.Context, b domain.Buyer, id int) (domain.Buyer, error) {
	buyer, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Buyer{}, errors.New("Error getting buyer")
	}
	if b.FirstName != "" {
		buyer.FirstName = b.FirstName
	}
	if b.LastName != "" {
		buyer.LastName = b.LastName
	}
	err = s.repository.Update(ctx, buyer)
	if err != nil {
		return domain.Buyer{}, ErrNotFound
	}

	return buyer, nil
}

// GetAll returns all buyers
// @Summary Get all buyers
// @Description Returns all buyers
// @Tags buyers
// @Produce  json
// @Success 200 {array} domain.Buyer
// @Failure 404 {string} string "Buyer not found"
// @Router /buyers [get]
func (s *service) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	b, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, ErrNotFound
	}

	return b, nil
}

// Get returns a buyer by ID
// @Summary Get a buyer by ID
// @Description Returns a buyer with the provided ID
// @Tags buyers
// @Param id path int true "Buyer ID"
// @Produce  json
// @Success 200 {object} domain.Buyer
// @Failure 404 {string} string "Buyer not found"
// @Router /buyers/{id} [get]
func (s *service) Get(ctx context.Context, id int) (domain.Buyer, error) {
	b, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Buyer{}, ErrNotFound
	}

	return b, nil
}

// Delete deletes a buyer by ID
// @Summary Delete a buyer by ID
// @Description Deletes a buyer with the provided ID
// @Tags buyers
// @Param id path int true "Buyer ID"
// @Produce  json
// @Success 200 {string} string "OK"
// @Failure 404 {string} string "Buyer not found"
// @Router /buyers/{id} [delete]
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	return nil
}

// NewService creates a new buyer service instance
func NewService(r Repository) Service {
	return &service{
		r,
	}
}
