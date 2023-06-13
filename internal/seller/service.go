package seller

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound             = errors.New("seller not found")
	ErrCidAlreadyExists     = errors.New("cid already registered")
	ErrSaveSeller           = errors.New("error saving seller")
	ErrSellersNotRegistered = errors.New("there are no registered sellers")
)

type Service interface {
	GetAll(c context.Context) ([]domain.Seller, error)
	Get(ctx context.Context, id int) (domain.Seller, error)
	Save(c context.Context, s domain.Seller) (domain.Seller, error)
	Update(ctx context.Context, idn int, s domain.Seller) (domain.Seller, error)
	Delete(ctx context.Context, id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) GetAll(c context.Context) ([]domain.Seller, error) {
	sellers, err := s.repository.GetAll(c)
	if err != nil {
		return []domain.Seller{}, ErrSellersNotRegistered
	}
	return sellers, nil
}

func (s *service) Get(c context.Context, id int) (domain.Seller, error) {
	seller, err := s.repository.Get(c, id)
	if err != nil {
		return domain.Seller{}, ErrNotFound
	}
	return seller, nil
}

func (s *service) Save(c context.Context, seller domain.Seller) (domain.Seller, error) {
	cidAlreadyExists := s.repository.Exists(c, seller.CID)
	if cidAlreadyExists {
		return domain.Seller{}, ErrCidAlreadyExists
	}
	sellerID, err := s.repository.Save(c, seller)
	if err != nil {
		return domain.Seller{}, ErrSaveSeller
	}
	seller.ID = sellerID
	return seller, nil
}

func (s *service) Update(c context.Context, id int, newSeller domain.Seller) (domain.Seller, error) {
	seller, err := s.repository.Get(c, id)
	if err != nil {
		return domain.Seller{}, ErrNotFound
	}
	// caso o CID enviado
	if newSeller.CID != 0{
		// se o cid for diferente do anterior
		if newSeller.CID != seller.CID {
			cidAlreadyExists := s.repository.Exists(c, newSeller.CID)
			// valida se o cid está disponivel
			if cidAlreadyExists {
				return domain.Seller{}, ErrCidAlreadyExists
			}
		}
		seller.CID = newSeller.CID
	}
	if newSeller.Address != "" {
		seller.Address = newSeller.Address 
	}
	if newSeller.CompanyName != "" {
		seller.CompanyName = newSeller.CompanyName 
	}
	if newSeller.Telephone != "" {
		seller.Telephone = newSeller.Telephone 
	}

	errUpdate:= s.repository.Update(c, seller)
	if errUpdate != nil {
		return domain.Seller{}, errUpdate
	}
	return seller, nil
}

func (s *service) Delete(c context.Context, id int) error{
	err := s.repository.Delete(c, id)	
	if err != nil {
		return ErrNotFound
	}
	return nil
}