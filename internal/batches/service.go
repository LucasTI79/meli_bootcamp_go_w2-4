package batches

import (
	"context"
	"errors"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

// Errors
var (
	ErrInvalidBatchNumber = errors.New("batch number alredy exists")
	ErrSavingBatch        = errors.New("error saving batch")
)

type CreateBatches struct {
	BatchNumber        int       `binding:"required" json:"batch_number"`
	CurrentQuantity    int       `binding:"required" json:"current_quantity"`
	CurrentTemperature int       `binding:"required" json:"current_temperature"`
	DueDate            time.Time `binding:"required" json:"due_date"`
	InitialQuantity    int       `binding:"required" json:"initial_quantity"`
	ManufacturingDate  time.Time `binding:"required" json:"manufacturing_date"`
	ManufacturingHour  int       `binding:"required" json:"manufacturing_hour"`
	MinimumTemperature int       `binding:"required" json:"minimum_temperature"`
	ProductID          int       `binding:"required" json:"product_id"`
	SectionID          int       `binding:"required" json:"section_id"`
}

type Service interface {
	Create(ctx context.Context, batches CreateBatches) (domain.Batches, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) Create(ctx context.Context, b CreateBatches) (domain.Batches, error) {
	existsBatchNumber := s.repository.Exists(ctx, b.BatchNumber)
	if existsBatchNumber {
		return domain.Batches{}, ErrInvalidBatchNumber
	}

	batch := domain.Batches{
		BatchNumber:        b.BatchNumber,
		CurrentQuantity:    b.CurrentQuantity,
		CurrentTemperature: b.CurrentTemperature,
		DueDate:            b.DueDate,
		InitialQuantity:    b.InitialQuantity,
		ManufacturingDate:  b.ManufacturingDate,
		ManufacturingHour:  b.ManufacturingHour,
		MinimumTemperature: b.MinimumTemperature,
		ProductID:          b.ProductID,
		SectionID:          b.SectionID,
	}

	i, err := s.repository.Save(ctx, batch)
	if err != nil {
		return domain.Batches{}, ErrSavingBatch
	}
	batch.ID = i
	return batch, nil
}
