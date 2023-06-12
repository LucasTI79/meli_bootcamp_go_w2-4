package handler

import (
	"net/http"

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
// @Tags Users
// @Description List all users
// @Accept  json
// @Produce  json
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Router /users [get]
func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.sellerService.GetAll(c.Request.Context())
		if err != nil {
			web.Error(c, http.StatusBadRequest, "não há vendedores cadastrados")
			
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func (s *Seller) Get() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// Create godoc
// @Summary Create
// @Tags Users
// @Description Create user
// @Accept  json
// @Produce  json
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Router /users [post]
func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req domain.Seller
		if err := c.Bind(&req); err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			 })
			return
		}
		if req.CID == 0{
			web.Error(c, http.StatusUnprocessableEntity, "cid é obrigatório")
			return
		}
		if req.CompanyName == ""{
			web.Error(c, http.StatusUnprocessableEntity, "nome da empresa é obrigatório")
			return
		}
		if req.Address == ""{
			web.Error(c, http.StatusUnprocessableEntity, "endereço é obrigatório")
			return
		}
		if req.Telephone == ""{
			web.Error(c, http.StatusUnprocessableEntity, "telefone é obrigatório")
			return
		}

		user, err := s.sellerService.Save(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error()})
			return
		}
		web.Success(c, http.StatusCreated, user)
	}
}

func (s *Seller) Update() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Seller) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
