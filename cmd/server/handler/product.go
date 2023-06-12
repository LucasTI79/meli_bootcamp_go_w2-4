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

type CreateProductRequest struct {
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

func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func (p *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (p *Product) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := middleware.ParsedRequest[CreateProductRequest](c)
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

func (p *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
