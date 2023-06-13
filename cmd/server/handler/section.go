package handler

import (
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
			web.Success(c, http.StatusNoContent, sections)
			return
		}
		web.Success(c, http.StatusOK, sections)
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
		web.Success(c, http.StatusOK, sec)
		return
	}
}

func (s *Section) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto section.CreateSection
		if err := c.ShouldBindJSON(&dto); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
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
		web.Success(c, http.StatusCreated, sec)

	}
}

func (s *Section) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, "id should be a number")
		}

		var dto section.UpdateSection
		if err := c.ShouldBindJSON(&dto); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}

		sec, err := s.sectionService.Update(c.Request.Context(), dto, id)

		if err != nil {
			if err == section.ErrInvalidSectionNumber {
				web.Error(c, http.StatusConflict, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusOK, sec)
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

		web.Success(c, http.StatusNoContent, domain.Section{})

	}
}
