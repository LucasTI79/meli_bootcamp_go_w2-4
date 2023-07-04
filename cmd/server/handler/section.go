package handler

import (
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
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

// GetAll godoc
//
//	@Summary	Get all sections
//	@Tags		Sections
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	web.response		"Returns all sections"
//	@Success	204	{object}	web.response		"No sections to retrieve"
//	@Failure	500	{object}	web.errorResponse	"Could not fetch sections"
//	@Router		/api/v1/sections [get]
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

// Get godoc
//
//	@Summary	Get section by ID
//	@Tags		Sections
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Section ID"
//	@Success	200	{object}	web.response		"Returns section"
//	@Failure	400	{object}	web.errorResponse	"Invalid ID type"
//	@Failure	404	{object}	web.errorResponse	"Could not find section"
//	@Router		/api/v1/sections/{id} [get]
func (s *Section) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		sec, err := s.sectionService.Get(c.Request.Context(), id)

		if err != nil {
			web.Error(c, http.StatusNotFound, err.Error())
			return
		}
		web.Success(c, http.StatusOK, sec)
	}
}

// Create godoc
//
//	@Summary	Create new section
//	@Tags		Sections
//	@Accept		json
//	@Produce	json
//	@Param		product	body		section.CreateSection	true	"section to be added"
//	@Success	201		{object}	web.response			"Returns created section"
//	@Failure	409		{object}	web.errorResponse		"`section_number` is not unique"
//	@Failure	422		{object}	web.errorResponse		"Missing fields or invalid field types"
//	@Failure	500		{object}	web.errorResponse		"Could not save section"
//	@Router		/api/v1/sections [post]
func (s *Section) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		dto := middleware.GetBody[section.CreateSection](c)

		sec, err := s.sectionService.Create(c, dto)
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

// Update godoc
//
//	@Summary	Updates existing section
//	@Tags		Sections
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"Section ID"
//	@Param		section	body		section.UpdateSection	true	"Fields to update"
//	@Success	200		{object}	web.response			"Returns updated section"
//	@Failure	400		{object}	web.errorResponse		"Invalid ID type"
//	@Failure	404		{object}	web.errorResponse		"Could not find section"
//	@Failure	409		{object}	web.errorResponse		"`section_number` is not unique"
//	@Failure	422		{object}	web.errorResponse		"Invalid field types"
//	@Failure	500		{object}	web.errorResponse		"Could not save section"
//	@Router		/api/v1/sections/{id} [patch]
func (s *Section) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		dto := middleware.GetBody[section.UpdateSection](c)

		sec, err := s.sectionService.Update(c.Request.Context(), dto, id)

		if err != nil {
			if err == section.ErrInvalidSectionNumber {
				web.Error(c, http.StatusConflict, err.Error())
			} else if err == section.ErrNotFound {
				web.Error(c, http.StatusNotFound, err.Error())
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
			}
			return
		}

		web.Success(c, http.StatusOK, sec)
	}
}

// Delete godoc
//
//	@Summary	Delete section by ID
//	@Tags		Sections
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Section ID"
//	@Success	200	{object}	web.response		"Section deleted successfully"
//	@Failure	400	{object}	web.errorResponse	"Invalid ID type"
//	@Failure	404	{object}	web.errorResponse	"Could not find section"
//	@Failure	500	{object}	web.errorResponse	"Could not delete section"
//	@Router		/api/v1/sections/{id} [delete]
func (s *Section) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")

		err := s.sectionService.Delete(c, id)

		if err != nil {
			web.Error(c, http.StatusNotFound, "id %d not found", id)
			return
		}

		web.Success(c, http.StatusNoContent, domain.Section{})

	}
}
