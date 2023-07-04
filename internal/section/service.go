package section

import (
	"context"
	"errors"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

type CreateSection struct {
	SectionNumber      int `binding:"required" json:"section_number"`
	CurrentTemperature int `binding:"required" json:"current_temperature"`
	MinimumTemperature int `binding:"required" json:"minimum_temperature"`
	CurrentCapacity    int `binding:"required" json:"current_capacity"`
	MinimumCapacity    int `binding:"required" json:"minimum_capacity"`
	MaximumCapacity    int `binding:"required" json:"maximum_capacity"`
	WarehouseID        int `binding:"required" json:"warehouse_id"`
	ProductTypeID      int `binding:"required" json:"product_type_id"`
}

type UpdateSection struct {
	SectionNumber      *int `json:"section_number"`
	CurrentTemperature *int `json:"current_temperature"`
	MinimumTemperature *int `json:"minimum_temperature"`
	CurrentCapacity    *int `json:"current_capacity"`
	MinimumCapacity    *int `json:"minimum_capacity"`
	MaximumCapacity    *int `json:"maximum_capacity"`
	WarehouseID        *int `json:"warehouse_id"`
	ProductTypeID      *int `json:"product_type_id"`
}

// Errors
var (
	ErrNotFound             = errors.New("section not found")
	ErrInvalidSectionNumber = errors.New("section number alredy exists")
	ErrSavingSection        = errors.New("error saving section")
	ErrGetSections          = errors.New("error getting sections")
)

type Service interface {
	Create(ctx context.Context, section CreateSection) (domain.Section, error)
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Update(ctx context.Context, dto UpdateSection, id int) (domain.Section, error)
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

func (s *service) Create(ctx context.Context, createSection CreateSection) (domain.Section, error) {
	existsSectionNumber := s.repository.Exists(ctx, createSection.SectionNumber)
	if existsSectionNumber {
		return domain.Section{}, ErrInvalidSectionNumber
	}
	section := mapCreateToDomain(&createSection)
	i, err := s.repository.Save(ctx, *section)
	if err != nil {
		return domain.Section{}, ErrSavingSection
	}

	sec := domain.Section{
		ID:                 i,
		SectionNumber:      section.SectionNumber,
		CurrentTemperature: section.CurrentTemperature,
		MinimumTemperature: section.MinimumTemperature,
		CurrentCapacity:    section.CurrentCapacity,
		MinimumCapacity:    section.MinimumCapacity,
		MaximumCapacity:    section.MaximumCapacity,
		WarehouseID:        section.WarehouseID,
		ProductTypeID:      section.ProductTypeID,
	}

	return sec, nil
}

func (s *service) GetAll(ctx context.Context) ([]domain.Section, error) {
	sec, err := s.repository.GetAll(ctx)
	if err != nil {
		return []domain.Section{}, ErrGetSections
	}
	return sec, nil
}

func (s *service) Get(ctx context.Context, id int) (domain.Section, error) {
	section, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.Section{}, ErrNotFound
	}

	return section, nil

}

func (s *service) Delete(ctx context.Context, id int) error {
	_, err := s.Get(ctx, id)
	if err != nil {
		return ErrNotFound
	}
	return s.repository.Delete(ctx, id)
}

func (s *service) Update(ctx context.Context, dto UpdateSection, id int) (domain.Section, error) {
	sec, err := s.Get(ctx, id)
	if err != nil {
		return domain.Section{}, ErrNotFound
	}
	if dto.SectionNumber != nil {
		existsSectionNumber := s.repository.Exists(ctx, *dto.SectionNumber)
		if existsSectionNumber && sec.SectionNumber != *dto.SectionNumber {
			return domain.Section{}, ErrInvalidSectionNumber
		}
	}
	applyValues(&sec, dto)
	sec, err = s.repository.Update(ctx, sec)
	return sec, err
}

func applyValues(sec *domain.Section, dto UpdateSection) {
	if dto.SectionNumber != nil {
		sec.SectionNumber = *dto.SectionNumber
	}

	if dto.CurrentTemperature != nil {
		sec.CurrentTemperature = *dto.CurrentTemperature
	}

	if dto.MinimumTemperature != nil {
		sec.MinimumTemperature = *dto.MinimumTemperature
	}

	if dto.CurrentCapacity != nil {
		sec.CurrentCapacity = *dto.CurrentCapacity
	}

	if dto.MinimumCapacity != nil {
		sec.MinimumCapacity = *dto.MinimumCapacity
	}

	if dto.MaximumCapacity != nil {
		sec.MaximumCapacity = *dto.MaximumCapacity
	}

	if dto.WarehouseID != nil {
		sec.WarehouseID = *dto.WarehouseID
	}

	if dto.ProductTypeID != nil {
		sec.ProductTypeID = *dto.ProductTypeID
	}
}

func mapCreateToDomain(section *CreateSection) *domain.Section {
	return &domain.Section{
		SectionNumber:      section.SectionNumber,
		CurrentTemperature: section.CurrentTemperature,
		MinimumTemperature: section.MinimumTemperature,
		CurrentCapacity:    section.CurrentCapacity,
		MinimumCapacity:    section.MinimumCapacity,
		MaximumCapacity:    section.MaximumCapacity,
		WarehouseID:        section.WarehouseID,
		ProductTypeID:      section.ProductTypeID,
	}
}
