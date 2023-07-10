package buyer

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Error definitions
var (
	ErrNotFound            = errors.New("buyer not found")
	ErrGeneric             = errors.New("")
	ErrInternalServerError = errors.New("internal server error")
	ErrAlreadyExists       = errors.New("buyer already exists")
	ErrSavingBuyer         = errors.New("error saving buyer")
)

type CountByBuyer struct {
	ID           int
	CardNumberID int
	FirstName    string
	LastName     string
	Count        int
}

// Service is the buyer service interface
type Service interface {
	Create(ctx context.Context, b domain.BuyerCreate) (domain.Buyer, error)
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	Get(ctx context.Context, id int) (domain.Buyer, error)
	Update(ctx context.Context, b domain.Buyer, id int) (domain.Buyer, error)
	Delete(ctx context.Context, id int) error
	CountPurchaseOrders(ctx context.Context, id int) ([]CountByBuyer, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		r,
	}
}

func (s *service) Create(ctx context.Context, b domain.BuyerCreate) (domain.Buyer, error) {
	ex := s.repository.Exists(ctx, b.CardNumberID)
	if ex {
		return domain.Buyer{}, ErrAlreadyExists
	}
	buyerDomain := *mapCreateToDomain(b)
	id, err := s.repository.Save(ctx, buyerDomain)
	if err != nil {
		return domain.Buyer{}, ErrSavingBuyer
	}
	buyerDomain.ID = id
	return buyerDomain, nil
}

func (s *service) Update(ctx context.Context, b domain.Buyer, id int) (domain.Buyer, error) {
	buyer, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Buyer{}, errors.New("error getting buyer")
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

func (s *service) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	b, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, ErrNotFound
	}

	return b, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Buyer, error) {
	b, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Buyer{}, ErrNotFound
	}

	return b, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	return nil
}

func (s *service) CountPurchaseOrders(ctx context.Context, id int) ([]CountByBuyer, error) {
	if id == 0 {
		e, err := s.repository.GetAllPurchaseOrders(ctx)
		if err != nil {
			return []CountByBuyer{}, err
		}
		return e, nil

	}
	e, err := s.repository.GetPurchaseOrderByID(ctx, id)
	if err != nil {
		return []CountByBuyer{}, err
	}

	return []CountByBuyer{e}, nil
}

func mapCreateToDomain(b domain.BuyerCreate) *domain.Buyer {
	return &domain.Buyer{
		CardNumberID: b.CardNumberID,
		FirstName:    b.FirstName,
		LastName:     b.LastName,
	}
}
