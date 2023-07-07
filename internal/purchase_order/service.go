package purchaseorder

import (
	"context"
	"errors"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
)

var (
	ErrAlreadyExists           = errors.New("order_number already exists")
	ErrInternalServerError     = errors.New("internal server error")
	ErrFKNotFound              = errors.New("buyer_id or order_status_id not found")
	ErrProductRecordIDNotFound = errors.New("product_record_id not found")
)

type PurchaseOrderDTO struct {
	ID              int
	OrderNumber     string
	OrderDate       time.Time
	TrackingCode    string
	BuyerID         int
	ProductRecordID int
	OrderStatusID   int
}

type Service interface {
	Create(c context.Context, purchaseOrder PurchaseOrderDTO) (domain.PurchaseOrder, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(c context.Context, purchaseOrder PurchaseOrderDTO) (domain.PurchaseOrder, error) {
	if s.repo.Exists(c, purchaseOrder.OrderNumber) {
		return domain.PurchaseOrder{}, ErrAlreadyExists
	}

	i := mapPurchaseOrderDTOToDomain(&purchaseOrder)
	id, err := s.repo.Create(c, i)
	if err != nil {
		if errors.Is(err, ErrFKNotFound) {
			return domain.PurchaseOrder{}, ErrFKNotFound
		}
		if errors.Is(err, ErrProductRecordIDNotFound) {
			return domain.PurchaseOrder{}, ErrProductRecordIDNotFound
		}
		return domain.PurchaseOrder{}, ErrInternalServerError
	}

	i.ID = id
	return i, nil
}

func mapPurchaseOrderDTOToDomain(purchaseOrder *PurchaseOrderDTO) domain.PurchaseOrder {
	return domain.PurchaseOrder{
		OrderNumber:     purchaseOrder.OrderNumber,
		OrderDate:       purchaseOrder.OrderDate,
		TrackingCode:    purchaseOrder.TrackingCode,
		BuyerID:         purchaseOrder.BuyerID,
		ProductRecordID: purchaseOrder.ProductRecordID,
		OrderStatusID:   purchaseOrder.OrderStatusID,
	}
}
