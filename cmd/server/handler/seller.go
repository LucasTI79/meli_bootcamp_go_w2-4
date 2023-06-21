package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
)

type Seller struct {
	sellerService seller.Service
}

func NewSeller(s seller.Service) *Seller {
	return &Seller{
		sellerService: s,
	}
}

// GetAll retrieves all sellers.
//
//	@Summary		Get all sellers
//	@Description	Retrieves all sellers
//	@Tags			Sellers
//	@Produce		json
//	@Success		200	{array}	domain.Seller	"Successfully retrieved sellers"
//	@Success		204	"No Content"
//	@Failure		500	{object}	web.errorResponse	"Internal Server Error"
//	@Router			/api/v1/sellers [get]
func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		sellers, err := s.sellerService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		if len(sellers) == 0 {
			web.Success(c, http.StatusNoContent, sellers)
		}
		web.Success(c, http.StatusOK, sellers)
	}
}

// GetById retrieves a seller by ID.
//
//	@Summary		Get a seller by ID
//	@Description	Retrieves a seller based on the provided ID
//	@Produce		json
//	@Tags			Sellers
//	@Param			id	path		int					true	"Seller ID"
//	@Success		200	{object}	domain.Seller		"Successfully retrieved seller"
//	@Failure		400	{object}	web.errorResponse	"Bad Request"
//	@Failure		404	{object}	web.errorResponse	"Not Found"
//	@Router			/api/v1/sellers/{id} [get]
func (s *Seller) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		seller, errGetSeller := s.sellerService.Get(c, id)
		if errGetSeller != nil {
			web.Error(c, http.StatusNotFound, errGetSeller.Error())
			return
		}
		web.Success(c, http.StatusOK, seller)
	}
}

// Create creates a new seller.
//
//	@Summary		Create a new seller
//	@Description	Creates a new seller with the provided data
//	@Accept			json
//	@Produce		json
//	@Param			seller	body	domain.Seller	true	"Seller object"
//	@Tags			Sellers
//	@Success		201	{object}	domain.Seller		"Successfully created seller"
//	@Failure		404	{object}	web.errorResponse	"Not Found"
//	@Failure		422	{object}	web.errorResponse	"Unprocessable Entity"
//	@Failure		500	{object}	web.errorResponse	"Internal Server Error"
//	@Router			/api/v1/sellers [post]
func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.Seller
		if err := c.ShouldBind(&req); err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		if req.CID == 0 {
			web.Error(c, http.StatusUnprocessableEntity, "cid is required")
			return
		}
		if req.CompanyName == "" {
			web.Error(c, http.StatusUnprocessableEntity, "company name is required")
			return
		}
		if req.Address == "" {
			web.Error(c, http.StatusUnprocessableEntity, "address is required")
			return
		}
		if req.Telephone == "" {
			web.Error(c, http.StatusUnprocessableEntity, "phone is required")
			return
		}

		sellerSaved, err := s.sellerService.Save(c, req)
		if err != nil {
			if err == seller.ErrCidAlreadyExists {
				web.Error(c, http.StatusConflict, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
		}
		web.Success(c, http.StatusCreated, sellerSaved)
	}
}

// Update updates an existing seller.
//
//	@Summary		Update an existing seller
//	@Description	Updates an existing seller with the provided data
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int				true	"Seller ID"
//	@Param			seller	body	domain.Seller	true	"Seller object"
//	@Tags			Sellers
//	@Success		200	{object}	domain.Seller		"Successfully updated seller"
//	@Failure		400	{object}	web.errorResponse	"Bad Request"
//	@Failure		404	{object}	web.errorResponse	"Not Found"
//	@Failure		500	{object}	web.errorResponse	"Internal Server Error"
//	@Router			/api/v1/sellers/{id} [put]
func (s *Seller) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		var seller domain.Seller
		if err := c.ShouldBindJSON(&seller); err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		sellerUpdated, err := s.sellerService.Update(c, id, seller)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		web.Success(c, http.StatusOK, sellerUpdated)
	}
}

// Delete deletes a seller by ID.
//
//	@Summary		Delete a seller by ID
//	@Description	Deletes a seller based on the provided ID
//	@Param			id	path	int	true	"Seller ID"
//	@Tags			Sellers
//	@Success		204	"No Content"
//	@Failure		400	{object}	web.errorResponse	"Bad Request"
//	@Failure		404	{object}	web.errorResponse	"Not Found"
//	@Failure		500	{object}	web.errorResponse	"Internal Server Error"
//	@Router			/api/v1/sellers/{id} [delete]
func (s *Seller) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		errDelete := s.sellerService.Delete(c, id)
		if errDelete != nil {
			if errDelete == seller.ErrNotFound {
				web.Error(c, http.StatusNotFound, errDelete.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, errDelete.Error())
			}
		}
		web.Success(c, http.StatusNoContent, nil)
	}
}
