package seller

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("seller not found")
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
		// web.Error(c, http.StatusConflict, "cid inserido já existe")
		return domain.Seller{}, errors.New("cid já cadastrado")
	}
	userID, err := s.repository.Save(c, seller)
	if err != nil {
		return domain.Seller{}, err
	}
	seller.ID = userID
	return seller, nil
}	