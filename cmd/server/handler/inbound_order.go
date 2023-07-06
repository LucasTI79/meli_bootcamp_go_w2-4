package handler

import (
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	inboundOrder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/inbound_order"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
)

type InboundOrder struct {
	inboundOrderService inboundOrder.Service
}

func NewInboundOrder(i inboundOrder.Service) *InboundOrder {
	return &InboundOrder{
		inboundOrderService: i,
	}
}

func (i *InboundOrder) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		inboundOrder := middleware.GetBody[InboundOrderRequest](c)

		if inboundOrder.OrderDate == nil {
			web.Error(c, http.StatusBadRequest, "order must have a date")
			return
		}
		if inboundOrder.OrderNumber == nil {
			web.Error(c, http.StatusBadRequest, "order must have a number")
			return
		}
		if inboundOrder.EmployeeID == nil {
			web.Error(c, http.StatusBadRequest, "order must have a employee associated with")
			return
		}
		if inboundOrder.ProductBatchID == nil {
			web.Error(c, http.StatusBadRequest, "order must have a product batch associated with")
			return
		}
		if inboundOrder.WarehouseID == nil {
			web.Error(c, http.StatusBadRequest, "order must have a warehouse associated with")
			return
		}

		inboundValues := domain.InboundOrder{
			OrderDate:      *inboundOrder.OrderDate,
			OrderNumber:    *inboundOrder.OrderNumber,
			EmployeeID:     *inboundOrder.EmployeeID,
			ProductBatchID: *inboundOrder.ProductBatchID,
			WarehouseID:    *inboundOrder.WarehouseID,
		}

		res, err := i.inboundOrderService.Create(c.Request.Context(), inboundValues)
		if err != nil {
			web.Error(c, http.StatusConflict, "error creating inbound order")
			return
		}
		web.Success(c, http.StatusCreated, res)
	}
}

type InboundOrderRequest struct {
	OrderDate      *time.Time `json:"order_date" `
	OrderNumber    *string    `json:"order_number" `
	EmployeeID     *int       `json:"employee_id" `
	ProductBatchID *int       `json:"product_batch_id" `
	WarehouseID    *int       `json:"warehouse_id" `
}
