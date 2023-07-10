package handler

import (
	"errors"
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

type Buyer struct {
	buyerService buyer.Service
}

func NewBuyer(b buyer.Service) *Buyer {
	return &Buyer{
		buyerService: b,
	}
}

// Buyer GoDoc
//
//	@Summary		Get a buyer by ID
//	@Description	Get a buyer by ID
//	@Tags			Buyers
//	@Param			id	path		int	true	"Buyer ID"
//	@Success		200	{object}	domain.Buyer
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Buyer not found"
//	@Router			/api/v1/buyers/{id} [get]
func (b *Buyer) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")

		buyer, err := b.buyerService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "buyer not found")
			return
		}
		web.Success(c, http.StatusOK, buyer)
	}
}

// Buyer GoDoc
//
//	@Summary		Delete a buyer by ID
//	@Description	Delete a buyer by ID
//	@Tags			Buyers
//	@Param			id	path		int		true	"Buyer ID"
//	@Success		200	{string}	string	"Buyer deleted"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Buyer not found"
//	@Router			/api/v1/buyers/{id} [delete]
func (b *Buyer) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")

		err := b.buyerService.Delete(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "buyer not found")
			return
		}
		web.Success(c, http.StatusOK, "buyer deleted")
	}
}

// Buyer GoDoc
//
//	@Summary		Get all buyers
//	@Description	Get all buyers
//	@Tags			Buyers
//	@Success		200	{array}		domain.Buyer
//	@Failure		500	{string}	string	"Buyer not found"
//	@Failure		204	{string}	string	"No buyers found"
//	@Router			/api/v1/buyers [get]
func (b *Buyer) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		buyers, err := b.buyerService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "buyer not found")
			return
		}
		if len(buyers) == 0 {
			web.Success(c, http.StatusNoContent, buyers)
			return
		}
		web.Success(c, http.StatusOK, buyers)
	}
}

// Buyer GoDoc
//
//	@Summary		Create a new buyer
//	@Description	Create a new buyer
//	@Tags			Buyers
//	@Accept			json
//	@Param			buyer	body		domain.BuyerCreate	true	"Buyer object"
//	@Success		201		{object}	domain.Buyer
//	@Failure		422		{string}	string	"Buyer not created"
//	@Router			/api/v1/buyers [post]
func (b *Buyer) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		buyer := middleware.GetBody[domain.BuyerCreate](c)

		buyerF, err := b.buyerService.Create(c.Request.Context(), buyer)
		if err != nil {
			web.Error(c, http.StatusConflict, "buyer not created")
			return
		}
		web.Response(c, http.StatusCreated, buyerF)
	}
}

// Buyer GoDoc
//
//	@Summary		Update a buyer by ID
//	@Description	Update a buyer by ID
//	@Tags			Buyers
//	@Accept			json
//	@Param			id		path		int				true	"Buyer ID"
//	@Param			buyer	body		domain.Buyer	true	"Buyer object"
//	@Success		200		{object}	domain.Buyer
//	@Failure		400		{string}	string	"Invalid ID"
//	@Failure		404		{string}	string	"Buyer not updated"
//	@Router			/api/v1/buyers/{id} [patch]
func (b *Buyer) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		buyer := middleware.GetBody[domain.Buyer](c)

		buyerUpdated, err := b.buyerService.Update(c, buyer, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "buyer not updated")
			return
		}
		web.Success(c, http.StatusOK, buyerUpdated)
	}
}

// ReportAllPurchaseOrders godoc
//
//	@Summary	Return purchaseOrder count for each buyer
//	@Description	Return purchaseOrder count for each buyer
//	@Tags		Buyers
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	web.response		"Returns purchaseOrder count for buyers"
//	@Success	204	{object}	web.response		"No content was found"
//	@Failure	500	{object}	web.errorResponse	"Could not generate report"
//	@Router		/api/v1/buyers/report-purchase-orders [get]
func (h *Buyer) PurchaseOrderReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		if _, exists := c.Params.Get("id"); exists {
			id = c.GetInt("id")
		}

		report, err := h.buyerService.CountPurchaseOrders(c, id)

		if err != nil {
			status := checkErrorStatusReportPurchaseOrder(err)
			web.Error(c, status, err.Error())
			return
		}

		web.Success(c, http.StatusOK, report)
	}
}

// ReportPurchaseOrdersByID godoc
//
//	@Summary	Return purchaseOrder count for given buyer
//	@Description	Return purchaseOrder count for given buyer
//	@Tags		Buyers
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Buyer ID"
//	@Success	200	{object}	web.response		"Returns purchaseOrder count for buyer"
//	@Failure	404	{object}	web.response		"ID was not found"
//	@Failure	500	{object}	web.errorResponse	"Could not generate report"
//	@Router		/api/v1/buyers/report-purchase-orders/{id} [get]
func _() {} // Implementation is in the PurchaseOrderReport function

func checkErrorStatusReportPurchaseOrder(err error) int {
	if errors.Is(err, buyer.ErrAlreadyExists) {
		return http.StatusConflict
	}
	if errors.Is(err, buyer.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
