package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
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

func (b *Buyer) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
			return
		}
		buyer, err := b.buyerService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		web.Success(c, http.StatusOK, buyer)
	}
}

func (b *Buyer) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
			return
		}
		err = b.buyerService.Delete(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		web.Success(c, http.StatusOK, err)
	}
}

func (b *Buyer) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		buyers, err := b.buyerService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Buyer not found")
			return
		}
		if len(buyers) == 0 {
			web.Success(c, http.StatusNoContent, buyers)
			return
		}
		web.Success(c, http.StatusOK, buyers)
	}
}

func (b *Buyer) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buyer domain.BuyerCreate
		if err := c.ShouldBindJSON(&buyer); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "buyer not created")
			return
		}
		buyerF, err := b.buyerService.Create(c.Request.Context(), buyer)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "buyer not created")
			return
		}
		web.Response(c, http.StatusCreated, buyerF)
	}
}

func (b *Buyer) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buyer domain.Buyer
		if err := c.ShouldBindJSON(&buyer); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "buyer not created")
			return
		}
		/*cardID := c.Param(("card_number_id"))
		fmt.Println(" ", cardID)
		if cardID != "" {
			fmt.Println("a ...any")
			web.Error(c, http.StatusUnprocessableEntity, "buyer not created")
			return
		}*/
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
			return
		}
		buyerUpdated, err := b.buyerService.Update(c, buyer, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "Buyer not updated")
			return
		}
		web.Success(c, http.StatusOK, buyerUpdated)
	}
}
