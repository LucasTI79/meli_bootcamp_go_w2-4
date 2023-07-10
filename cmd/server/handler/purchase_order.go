package handler

import (
	"errors"
	"net/http"
	"time"

	purchaseOrder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/purchase_order"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

type PurchaseOrder struct {
	purchaseOrderService purchaseOrder.Service
}

type PurchaseOrderRequest struct {
	OrderNumber     *string `binding:"required" json:"order_number"`
	OrderDate       *string `binding:"required" json:"order_date"`
	TrackingCode    *string `binding:"required" json:"tracking_code"`
	BuyerID         *int    `binding:"required" json:"buyer_id"`
	ProductRecordID *int    `binding:"required" json:"product_record_id"`
	OrderStatusID   *int    `binding:"required" json:"order_status_id"`
}

func NewPurchaseOrder(s purchaseOrder.Service) *PurchaseOrder {
	return &PurchaseOrder{
		purchaseOrderService: s,
	}
}

// Create purchase order
//
//	@Summary	Create new purchase order
//	@Tags		Purchase order
//	@Accept		json
//	@Produce	json
//	@Param		purchaseOrder	body		PurchaseOrderRequest		true	"purchase order to be added"
//	@Success	201		{object}	web.response		"Returns created purchase order"
//	@Failure	409		{object}	web.errorResponse	"`id` is not unique or `locality_id` not found"
//	@Failure	422		{object}	web.errorResponse	"Missing fields or invalid field types"
//	@Failure	500		{object}	web.errorResponse	"Could not save purchase order"
//	@Router		/api/v1/purchase-orders [post]
func (i *PurchaseOrder) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := middleware.GetBody[PurchaseOrderRequest](c)
		dto, err := mapPurchaseOrderRequestToDTO(&req)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		i, err := i.purchaseOrderService.Create(c.Request.Context(), *dto)

		if err != nil {
			errStatus := checkErrorStatusPurchaseOrder(err)
			web.Error(c, errStatus, err.Error())
			return
		}

		web.Success(c, http.StatusCreated, i)
	}
}

func checkErrorStatusPurchaseOrder(err error) int {
	if errors.Is(err, purchaseOrder.ErrAlreadyExists) ||
		errors.Is(err, purchaseOrder.ErrFKNotFound) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func mapPurchaseOrderRequestToDTO(req *PurchaseOrderRequest) (*purchaseOrder.PurchaseOrderDTO, error) {
	orderDate, err := time.Parse("2006-01-02", *req.OrderDate)
	if err != nil {
		return &purchaseOrder.PurchaseOrderDTO{}, err
	}

	return &purchaseOrder.PurchaseOrderDTO{
		OrderNumber:     *req.OrderNumber,
		OrderDate:       orderDate,
		TrackingCode:    *req.TrackingCode,
		BuyerID:         *req.BuyerID,
		ProductRecordID: *req.ProductRecordID,
		OrderStatusID:   *req.OrderStatusID,
	}, err
}
