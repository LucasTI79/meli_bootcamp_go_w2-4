package section

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrNotFound             = errors.New("section not found")
	ErrInvalidSectionNumber = errors.New("section number alredy exists")
	ErrSavingSection        = errors.New("error saving section")
	ErrIdNotFound           = errors.New("section not found")
)

type Service interface {
	Save(ctx context.Context, section domain.CreateSection) (domain.Section, error)
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, section domain.Section, dto domain.CreateSection) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Save(ctx context.Context, createSection domain.CreateSection) (domain.Section, error) {
	existsSectionNumber := s.repository.Exists(ctx, *createSection.SectionNumber)
	if existsSectionNumber {
		return domain.Section{}, ErrInvalidSectionNumber
	}
	i, err := s.repository.Save(ctx, createSection)
	if err != nil {
		return domain.Section{}, ErrSavingSection
	}

	section := domain.Section{
		ID:                 i,
		SectionNumber:      *createSection.SectionNumber,
		CurrentTemperature: *createSection.CurrentTemperature,
		MinimumTemperature: *createSection.MinimumTemperature,
		CurrentCapacity:    *createSection.CurrentCapacity,
		MinimumCapacity:    *createSection.MinimumCapacity,
		MaximumCapacity:    *createSection.MaximumCapacity,
		WarehouseID:        *createSection.WarehouseID,
		ProductTypeID:      *createSection.ProductTypeID,
	}

	return section, nil

}

func (s *service) GetAll(ctx context.Context) ([]domain.Section, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) Get(ctx context.Context, id int) (domain.Section, error) {
	section, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Section{}, ErrIdNotFound
	}

	return section, nil

}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) Update(ctx context.Context, section domain.Section, dto domain.CreateSection) error {
	existsSectionNumber := s.repository.Exists(ctx, section.SectionNumber)
	if existsSectionNumber && section.SectionNumber != *dto.SectionNumber {
		return ErrInvalidSectionNumber
	}
	return s.repository.Update(ctx, section)
}
