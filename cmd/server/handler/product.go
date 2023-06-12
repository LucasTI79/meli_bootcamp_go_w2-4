package handler

import (
	"errors"
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/middleware"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
)

type Product struct {
	productService product.Service
}

type CreateRequest struct {
	Desc       string  `binding:"required" json:"description"`
	ExpR       int     `binding:"required" json:"expiration_rate"`
	FreezeR    int     `binding:"required" json:"freezing_rate"`
	Height     float32 `binding:"required" json:"height"`
	Length     float32 `binding:"required" json:"length"`
	NetW       float32 `binding:"required" json:"netweight"`
	Code       string  `binding:"required" json:"product_code"`
	FreezeTemp float32 `binding:"required" json:"recommended_freezing_temperature"`
	Width      float32 `binding:"required" json:"width"`
	TypeID     int     `binding:"required" json:"product_type_id"`
	SellerID   int     `json:"seller_id"`
}

func NewProduct(s product.Service) *Product {
	return &Product{
		productService: s,
	}
}

// GetAll godoc
//
//	@Summary	Get all products
//	@Tags		Products
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	responses.Response	"Returns all products"
//	@Failure	500		{object}	responses.Response	"Could not fetch products"
//	@Router		/api/v1/products [get]
func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ps, err := p.productService.GetAll(c.Request.Context())
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		web.Success(c, http.StatusOK, ps)
	}
}

// Get godoc
//
//	@Summary	Get product by ID
//	@Tags		Products
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	responses.Response	"Returns product"
//	@Failure	400		{object}	responses.Response	"Invalid ID type"
//	@Failure	404		{object}	responses.Response	"Could not find product"
//	@Router		/api/v1/products/:id [get]
func (p *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := web.GetIntParam(c, "id")
		if err != nil {
			web.Error(c, http.StatusBadRequest, "id path parameter should be an int")
			return
		}
		p, err := p.productService.Get(c.Request.Context(), id)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		web.Success(c, http.StatusOK, p)
	}
}

// Create godoc
//
//	@Summary	Create new product
//	@Tags		Products
//	@Accept		json
//	@Produce	json
//	@Param		product	body		CreateRequest		true	"Product to be added"
//	@Success	200		{object}	responses.Response	"Returns created product"
//	@Failure	409		{object}	responses.Response	"`product_code` is not unique"
//	@Failure	422		{object}	responses.Response	"Missing fields or invalid field types"
//	@Failure	500		{object}	responses.Response	"Could not save product"
//	@Router		/api/v1/products [post]
func (p *Product) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := middleware.ParsedRequest[CreateRequest](c)
		p, err := p.productService.Create(
			c.Request.Context(),
			req.Desc,
			req.ExpR,
			req.FreezeR,
			req.Height,
			req.Length,
			req.NetW,
			req.Code,
			req.FreezeTemp,
			req.Width,
			req.TypeID,
			req.SellerID,
		)

		if err != nil {
			if errors.Is(err, product.ErrInvalidProductCode{}) {
				web.Error(c, http.StatusConflict, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusCreated, p)
	}
}

func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// Delete godoc
//
//	@Summary	Delete product by ID
//	@Tags		Products
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	responses.Response	"Product deleted successfully"
//	@Failure	400		{object}	responses.Response	"Invalid ID type"
//	@Failure	404		{object}	responses.Response	"Could not find product"
//	@Failure	500		{object}	responses.Response	"Could not delete product"
//	@Router		/api/v1/products/:id [delete]
func (p *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := web.GetIntParam(c, "id")
		if err != nil {
			web.Error(c, http.StatusBadRequest, "id path parameter should be an int")
			return
		}
		err = p.productService.Delete(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, product.ErrNotFound{}) {
				web.Error(c, http.StatusNotFound, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		web.Success(c, http.StatusOK, nil)
	}
}
