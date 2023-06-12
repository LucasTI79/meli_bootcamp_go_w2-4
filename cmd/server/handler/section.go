package handler

import (
	"fmt"
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
)

type Section struct {
	sectionService section.Service
}

func NewSection(s section.Service) *Section {
	return &Section{
		sectionService: s,
	}
}

func (s *Section) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Section) Get() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Section) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var section domain.CreateSection
		if err := c.ShouldBindJSON(&section); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(section)
		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "section number invalid")
		}

		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "current temperature invalid")
		}

		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "minimum temperature invalid")

		}

		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "current capacity invalid")

		}

		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "minimum capacity invalid")

		}

		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "maximum capacity invalid")

		}

		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "warehouse id invalid")

		}

		if section.CurrentCapacity == 0 {
			web.Error(c, http.StatusBadRequest, "product type id invalid")

		}

	}
}

func (s *Section) Update() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func (s *Section) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
