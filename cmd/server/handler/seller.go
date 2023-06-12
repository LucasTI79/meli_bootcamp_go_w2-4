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

// GetAll godoc
// @Summary GetAll
// @Tags Sellers
// @Description List all sellers
// @Accept  json
// @Produce  json
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Router /sellers [get]
func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		sellers, err := s.sellerService.GetAll(c.Request.Context())
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, sellers)
	}
}

// GetById godoc
// @Summary GetById
// @Tags Sellers
// @Description List sellers by id
// @Accept  json
// @Produce  json
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Router /sellers [get]
func (s *Seller) GetById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}

		seller, err := s.sellerService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		web.Success(c, http.StatusOK, seller)
	}
}

// Create godoc
// @Summary Create
// @Tags Sellers
// @Description Create sellet
// @Accept  json
// @Produce  json
// @Param token header string true "token"
// @Success 200 {object} web.Response
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

		sellerSaved, err := s.sellerService.Save(c.Request.Context(), req)
		if err != nil {
			if err == seller.ErrCidAlreadyExists {
				web.Error(c, http.StatusConflict, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		web.Success(c, http.StatusCreated, sellerSaved)
	}
}

func (s *Seller) Update() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Seller) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
