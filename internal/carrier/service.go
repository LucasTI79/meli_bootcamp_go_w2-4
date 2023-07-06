package carrier

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

var (
	ErrAlreadyExists       = errors.New("cid already exists")
	ErrInternalServerError = errors.New("internal server error")
	ErrLocalityIDNotFound  = errors.New("locality_id not found")
)

type CarrierDTO struct {
	CID         int
	CompanyName string
	Address     string
	Telephone   string
	LocalityID  int
}

type Service interface {
	Create(c context.Context, carrier CarrierDTO) (domain.Carrier, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(c context.Context, carrier CarrierDTO) (domain.Carrier, error) {
	if s.repo.Exists(c, carrier.CID) {
		return domain.Carrier{}, ErrAlreadyExists
	}

	i := mapCarrierDTOToDomain(&carrier)
	id, err := s.repo.Create(c, *i)
	if err != nil {
		var errMsg error
		switch {
		case errors.Is(err, ErrAlreadyExists):
			errMsg = ErrAlreadyExists
		case errors.Is(err, ErrLocalityIDNotFound):
			errMsg = ErrLocalityIDNotFound
		default:
			errMsg = ErrInternalServerError
		}

		return domain.Carrier{}, errMsg
	}

	i.ID = id
	return *i, nil
}

func mapCarrierDTOToDomain(carrier *CarrierDTO) *domain.Carrier {
	return &domain.Carrier{
		CID:         carrier.CID,
		CompanyName: carrier.CompanyName,
		Address:     carrier.Address,
		Telephone:   carrier.Telephone,
		LocalityID:  carrier.LocalityID,
	}
}
