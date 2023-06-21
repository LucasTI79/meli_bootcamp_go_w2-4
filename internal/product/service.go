package product

import (
	"context"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/optional"
)

type CreateDTO struct {
	Desc       string
	ExpR       int
	FreezeR    int
	Height     float32
	Length     float32
	NetW       float32
	Code       string
	FreezeTemp float32
	Width      float32
	TypeID     int
	SellerID   int
}

type UpdateDTO struct {
	Desc       optional.Opt[string]
	ExpR       optional.Opt[int]
	FreezeR    optional.Opt[int]
	Height     optional.Opt[float32]
	Length     optional.Opt[float32]
	NetW       optional.Opt[float32]
	Code       optional.Opt[string]
	FreezeTemp optional.Opt[float32]
	Width      optional.Opt[float32]
	TypeID     optional.Opt[int]
	SellerID   optional.Opt[int]
}

type Service interface {
	Create(c context.Context, product CreateDTO) (domain.Product, error)
	GetAll(c context.Context) ([]domain.Product, error)
	Get(c context.Context, id int) (domain.Product, error)
	Update(c context.Context, id int, updates UpdateDTO) (domain.Product, error)
	Delete(c context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(c context.Context, product CreateDTO) (domain.Product, error) {
	ps, err := s.repo.GetAll(c)
	if err != nil {
		return domain.Product{}, NewErrGeneric("error fetching products")
	}

	if !isUniqueProductCode(product.Code, ps) {
		return domain.Product{}, NewErrInvalidProductCode(product.Code)
	}

	p := mapCreateToDomain(&product)
	id, err := s.repo.Save(c, *p)
	if err != nil {
		return domain.Product{}, NewErrGeneric("error saving product")
	}

	p.ID = id
	return *p, nil
}

func (s *service) GetAll(c context.Context) ([]domain.Product, error) {
	ps, err := s.repo.GetAll(c)
	if err != nil {
		return nil, NewErrGeneric("could not fetch products")
	}
	return ps, nil
}

func (s *service) Get(c context.Context, id int) (domain.Product, error) {
	p, err := s.repo.Get(c, id)
	if err != nil {
		// TODO: Properly handle DB communication error differently
		return domain.Product{}, NewErrNotFound(id)
	}
	return p, nil
}

func (s *service) Update(c context.Context, id int, updates UpdateDTO) (domain.Product, error) {
	ps, err := s.repo.GetAll(c)
	if err != nil {
		return domain.Product{}, NewErrGeneric("could not fetch products")
	}

	p, err := s.repo.Get(c, id)
	if err != nil {
		return domain.Product{}, NewErrNotFound(id)
	}

	if code, hasVal := updates.Code.Value(); hasVal && p.ProductCode != code && !isUniqueProductCode(code, ps) {
		return domain.Product{}, NewErrInvalidProductCode(updates.Code.Val)
	}

	updated := applyUpdates(p, updates)
	if err := s.repo.Update(c, updated); err != nil {
		return domain.Product{}, NewErrGeneric("could not save changes")
	}

	return updated, nil
}

func (s *service) Delete(c context.Context, id int) error {
	err := s.repo.Delete(c, id)
	if err != nil {
		switch err.(type) {
		case *ErrNotFound:
			return NewErrNotFound(id)
		default:
			return NewErrGeneric("could not delete product")
		}
	}
	return nil
}

func isUniqueProductCode(code string, ps []domain.Product) bool {
	for _, p := range ps {
		if p.ProductCode == code {
			return false
		}
	}
	return true
}

func mapCreateToDomain(product *CreateDTO) *domain.Product {
	return &domain.Product{
		Description:    product.Desc,
		ExpirationRate: product.ExpR,
		FreezingRate:   product.FreezeR,
		Height:         product.Height,
		Length:         product.Length,
		Netweight:      product.NetW,
		ProductCode:    product.Code,
		RecomFreezTemp: product.FreezeTemp,
		Width:          product.Width,
		ProductTypeID:  product.TypeID,
		SellerID:       product.SellerID,
	}
}

func applyUpdates(p domain.Product, updates UpdateDTO) domain.Product {
	p.Description = updates.Desc.Or(p.Description)
	p.ExpirationRate = updates.ExpR.Or(p.ExpirationRate)
	p.FreezingRate = updates.FreezeR.Or(p.FreezingRate)
	p.Height = updates.Height.Or(p.Height)
	p.Length = updates.Length.Or(p.Length)
	p.Netweight = updates.NetW.Or(p.Netweight)
	p.ProductCode = updates.Code.Or(p.ProductCode)
	p.RecomFreezTemp = updates.FreezeTemp.Or(p.RecomFreezTemp)
	p.Width = updates.Width.Or(p.Width)
	p.ProductTypeID = updates.TypeID.Or(p.ProductTypeID)
	p.SellerID = updates.SellerID.Or(p.SellerID)
	return p
}
