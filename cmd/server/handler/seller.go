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
// @Summary Get all sellers
// @Description Retrieves all sellers
// @Produce json
// @Success 200 {array} domain.Seller "Successfully retrieved sellers"
// @Success 204 "No Content"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /sellers [get]
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
// @Summary Get a seller by ID
// @Description Retrieves a seller based on the provided ID
// @Produce json
// @Param id path int true "Seller ID"
// @Success 200 {object} domain.Seller "Successfully retrieved seller"
// @Failure 400 {object} errorResponse "Bad Request"
// @Failure 404 {object} errorResponse "Not Found"
// @Router /sellers/{id} [get]
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
// @Summary Create a seller
// @Description Create a new seller with the given information
// @Accept json
// @Produce json
// @Param seller body domain.Seller true "Seller object"
// @Success 201 {object} domain.Seller "Successfully created seller"
// @Failure 404 {object} errorResponse "Not Found"
// @Failure 422 {object} errorResponse "Unprocessable Entity"
// @Failure 409 {object} errorResponse "Conflict"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /sellers [post]
func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.Seller
		if err := c.Bind(&req); err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		if req.CID == 0 {
			web.Error(c, http.StatusUnprocessableEntity, "cid é obrigatório")
			return
		}
		if req.CompanyName == "" {
			web.Error(c, http.StatusUnprocessableEntity, "nome da empresa é obrigatório")
			return
		}
		if req.Address == "" {
			web.Error(c, http.StatusUnprocessableEntity, "endereço é obrigatório")
			return
		}
		if req.Telephone == "" {
			web.Error(c, http.StatusUnprocessableEntity, "telefone é obrigatório")
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

// Update updates a seller by ID.
// @Summary Update a seller by ID
// @Description Updates a seller with the given ID and information
// @Accept json
// @Produce json
// @Param id path int true "Seller ID"
// @Param seller body domain.Seller true "Seller object"
// @Success 200 {object} domain.Seller "Successfully updated seller"
// @Failure 404 {object} errorResponse "Not Found"
// @Failure 500 {object} errorResponse "Internal Server Error"
// @Router /sellers/{id} [put]
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
// @Summary Delete a seller by ID
// @Description Deletes a seller with the given ID
// @Produce plain
// @Param id path int true "Seller ID"
// @Success 204 "No Content"
// @Failure 404 {object} errorResponse "Not Found"
// @Router /sellers/{id} [delete]
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
