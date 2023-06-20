package handler

import (
	"errors"
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/middleware"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/optional"
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

// UpdateRequest contains pointers so that the Handler is able to
// distinguish between omitted (nil) and given (not-nil) fields.
// This does not affect the way the user passes the Request body.
type UpdateRequest struct {
	Desc       *string  `json:"description"`
	ExpR       *int     `json:"expiration_rate"`
	FreezeR    *int     `json:"freezing_rate"`
	Height     *float32 `json:"height"`
	Length     *float32 `json:"length"`
	NetW       *float32 `json:"netweight"`
	Code       *string  `json:"product_code"`
	FreezeTemp *float32 `json:"recommended_freezing_temperature"`
	Width      *float32 `json:"width"`
	TypeID     *int     `json:"product_type_id"`
	SellerID   *int     `json:"seller_id"`
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
//	@Success	200	{object}	web.response		"Returns all products"
//	@Success	204	{object}	web.response		"No products to retrieve"
//	@Failure	500	{object}	web.errorResponse	"Could not fetch products"
//	@Router		/api/v1/products [get]
func (p *Product) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		ps, err := p.productService.GetAll(c.Request.Context())

		if err != nil {
			errStatus := mapErrorToStatus(err)
			web.Error(c, errStatus, err.Error())
			return
		}

		if len(ps) == 0 {
			web.Success(c, http.StatusNoContent, ps)
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
//	@Param		id	path		int					true	"Product ID"
//	@Success	200	{object}	web.response		"Returns product"
//	@Failure	400	{object}	web.errorResponse	"Invalid ID type"
//	@Failure	404	{object}	web.errorResponse	"Could not find product"
//	@Router		/api/v1/products/{id} [get]
func (p *Product) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := web.GetIntParam(c, "id")
		if err != nil {
			web.Error(c, http.StatusBadRequest, "id path parameter should be an int")
			return
		}

		p, err := p.productService.Get(c.Request.Context(), id)

		if err != nil {
			errStatus := mapErrorToStatus(err)
			web.Error(c, errStatus, err.Error())
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
//	@Success	201		{object}	web.response		"Returns created product"
//	@Failure	409		{object}	web.errorResponse	"`product_code` is not unique"
//	@Failure	422		{object}	web.errorResponse	"Missing fields or invalid field types"
//	@Failure	500		{object}	web.errorResponse	"Could not save product"
//	@Router		/api/v1/products [post]
func (p *Product) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := middleware.ParsedRequest[CreateRequest](c)
		dto := mapCreateRequestToDTO(&req)
		p, err := p.productService.Create(c.Request.Context(), *dto)

		if err != nil {
			errStatus := mapErrorToStatus(err)
			web.Error(c, errStatus, err.Error())
			return
		}

		web.Success(c, http.StatusCreated, p)
	}
}

// Update godoc
//
//	@Summary	Updates existing product
//	@Tags		Products
//	@Accept		json
//	@Param		id	path	int	true	"Product ID"
//	@Produce	json
//	@Param		product	body		UpdateRequest		true	"Fields to update"
//	@Success	200		{object}	web.response		"Returns updated product"
//	@Failure	400		{object}	web.errorResponse	"Invalid ID type"
//	@Failure	404		{object}	web.errorResponse	"Could not find product"
//	@Failure	409		{object}	web.errorResponse	"`product_code` is not unique"
//	@Failure	422		{object}	web.errorResponse	"Invalid field types"
//	@Failure	500		{object}	web.errorResponse	"Could not save product"
//	@Router		/api/v1/products/{id} [patch]
func (p *Product) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := web.GetIntParam(c, "id")
		if err != nil {
			web.Error(c, http.StatusBadRequest, "id path parameter should be an int")
			return
		}

		req := middleware.ParsedRequest[UpdateRequest](c)
		dto := mapUpdateRequestToDTO(&req)
		p, err := p.productService.Update(c.Request.Context(), id, *dto)

		if err != nil {
			errStatus := mapErrorToStatus(err)
			web.Error(c, errStatus, err.Error())
			return
		}

		web.Success(c, http.StatusOK, p)
	}
}

// Delete godoc
//
//	@Summary	Delete product by ID
//	@Tags		Products
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Product ID"
//	@Success	200	{object}	web.response		"Product deleted successfully"
//	@Failure	400	{object}	web.errorResponse	"Invalid ID type"
//	@Failure	404	{object}	web.errorResponse	"Could not find product"
//	@Failure	500	{object}	web.errorResponse	"Could not delete product"
//	@Router		/api/v1/products/{id} [delete]
func (p *Product) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := web.GetIntParam(c, "id")
		if err != nil {
			web.Error(c, http.StatusBadRequest, "id path parameter should be an int")
			return
		}

		err = p.productService.Delete(c.Request.Context(), id)

		if err != nil {
			errStatus := mapErrorToStatus(err)
			web.Error(c, errStatus, err.Error())
			return
		}
		web.Success(c, http.StatusOK, nil)
	}
}

func mapErrorToStatus(err error) int {
	var invalidProductCode *product.ErrInvalidProductCode
	var notFound *product.ErrNotFound

	if errors.As(err, &invalidProductCode) {
		return http.StatusConflict
	}
	if errors.As(err, &notFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func mapUpdateRequestToDTO(req *UpdateRequest) *product.UpdateDTO {
	dto := product.UpdateDTO{}
	dto.Desc = *optional.FromPtr(req.Desc)
	dto.ExpR = *optional.FromPtr(req.ExpR)
	dto.FreezeR = *optional.FromPtr(req.FreezeR)
	dto.Height = *optional.FromPtr(req.Height)
	dto.Length = *optional.FromPtr(req.Length)
	dto.NetW = *optional.FromPtr(req.NetW)
	dto.Code = *optional.FromPtr(req.Code)
	dto.FreezeTemp = *optional.FromPtr(req.FreezeTemp)
	dto.Width = *optional.FromPtr(req.Width)
	dto.TypeID = *optional.FromPtr(req.TypeID)
	dto.SellerID = *optional.FromPtr(req.SellerID)
	return &dto
}

func mapCreateRequestToDTO(req *CreateRequest) *product.CreateDTO {
	return &product.CreateDTO{
		Desc:       req.Desc,
		ExpR:       req.ExpR,
		FreezeR:    req.FreezeR,
		Height:     req.Height,
		Length:     req.Length,
		NetW:       req.NetW,
		Code:       req.Code,
		FreezeTemp: req.FreezeTemp,
		Width:      req.Width,
		TypeID:     req.TypeID,
		SellerID:   req.SellerID,
	}
}
