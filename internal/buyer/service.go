package buyer

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("buyer not found")
	ErrGeneric  = errors.New("")
)

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

//arrumar as validaçoes
//arrumar os errors fazendo um pkg
//fazer uma struct pro update com card nao required
//swagger

func (s *service) Create(ctx context.Context, b domain.BuyerCreate) (domain.BuyerCreate, error) {
	ex := s.repository.Exists(ctx, b.CardNumberID)
	if ex {
		return domain.BuyerCreate{}, errors.New("Error saving buyer")
	}
	id, err := s.repository.Save(ctx, b)
	if err != nil {
		return domain.BuyerCreate{}, errors.New("Error saving buyer")
	}
	b.ID = id
	return b, nil
}

func (s *service) Update(ctx context.Context, b domain.Buyer, id int) (domain.Buyer, error) {
	buyer, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Buyer{}, errors.New("Error geting buyer")
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

func NewService(r Repository) Service {
	return &service{
		r,
	}
}
