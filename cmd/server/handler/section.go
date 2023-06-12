package handler

import (
	"errors"
	"net/http"
	"strconv"

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
	return func(c *gin.Context) {
		sections, err := s.sectionService.GetAll(c.Request.Context())
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		if len(sections) == 0 {
			web.Response(c, http.StatusNoContent, sections)
			return
		}
		web.Response(c, http.StatusOK, sections)
	}
}

func (s *Section) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
		}
		sec, err := s.sectionService.Get(c.Request.Context(), id)

		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		web.Response(c, http.StatusOK, sec)
		return
	}
}

func (s *Section) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto domain.CreateSection
		if err := c.ShouldBindJSON(&dto); err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "invalid section number")
			return
		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "current temperature invalid")
			return
		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "minimum temperature invalid")
			return
		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "current capacity invalid")
			return
		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "minimum capacity invalid")
			return

		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "maximum capacity invalid")
			return
		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "warehouse id invalid")
			return
		}

		if dto.CurrentCapacity == nil {
			web.Error(c, http.StatusBadRequest, "product type id invalid")
			return
		}

		sec, err := s.sectionService.Save(c, dto)
		if err != nil {
			if err == section.ErrInvalidSectionNumber {
				web.Error(c, http.StatusConflict, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
		web.Response(c, http.StatusCreated, sec)

	}
}

func (s *Section) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, errors.New("id should be a number").Error())
		}
		sec, err := s.sectionService.Get(c.Request.Context(), id)

		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}

		var dto domain.CreateSection
		if err := c.ShouldBindJSON(&dto); err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		if dto.SectionNumber != nil {
			sec.SectionNumber = *dto.SectionNumber
		}

		if dto.CurrentTemperature != nil {
			sec.CurrentTemperature = *dto.CurrentTemperature
		}

		if dto.MinimumTemperature != nil {
			sec.MinimumTemperature = *dto.MinimumTemperature
		}

		if dto.CurrentCapacity != nil {
			sec.CurrentCapacity = *dto.CurrentCapacity
		}

		if dto.MinimumCapacity != nil {
			sec.MinimumCapacity = *dto.MinimumCapacity
		}

		if dto.MaximumCapacity != nil {
			sec.MaximumCapacity = *dto.MaximumCapacity
		}

		if dto.WarehouseID != nil {
			sec.WarehouseID = *dto.WarehouseID
		}

		if dto.ProductTypeID != nil {
			sec.ProductTypeID = *dto.ProductTypeID
		}

		err = s.sectionService.Update(c.Request.Context(), sec, dto)

		if err != nil {
			if err == section.ErrInvalidSectionNumber {
				web.Error(c, http.StatusConflict, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Response(c, http.StatusOK, sec)
	}
}

func (s *Section) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		err = s.sectionService.Delete(c, id)

		if err != nil {
			web.Error(c, http.StatusNotFound, "id %d not found", id)
			return
		}

		web.Response(c, http.StatusNoContent, domain.Section{})

	}
}
