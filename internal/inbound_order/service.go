package inboundorder

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound            = errors.New("InboundOrder not found")
	ErrAlreadyExists       = errors.New("InboundOrder id already exists")
	ErrInternalServerError = errors.New("internal server error")
)

type Service interface {
	Create(ctx context.Context, e domain.InboundOrder) (domain.InboundOrder, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(ctx context.Context, i domain.InboundOrder) (domain.InboundOrder, error) {

	id, err := s.repository.Save(ctx, i)
	if err != nil {
		return domain.InboundOrder{}, ErrInternalServerError
	}

	i.ID = id

	return i, nil
}
