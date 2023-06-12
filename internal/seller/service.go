package seller

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("seller not found")
	ErrCidAlreadyExists =  errors.New("cid already registered")
	ErrSaveSeller = errors.New("error saving seller")
)

type Service interface{
	GetAll(c context.Context) ([]domain.Seller, error)
	Save(c context.Context, s domain.Seller) (domain.Seller, error)
}

type service struct{	
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func(s *service) GetAll(c context.Context) ([]domain.Seller, error){
	users, err := s.repository.GetAll(c)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) Save(c context.Context, seller domain.Seller)(domain.Seller, error){
	cidAlreadyExists := s.repository.Exists(c, seller.CID)
	if cidAlreadyExists {
		return domain.Seller{}, ErrCidAlreadyExists
	}
	userID, err := s.repository.Save(c, seller)
	if err != nil {
		return domain.Seller{}, ErrSaveSeller
	}
	seller.ID = userID
	return seller, nil
}	