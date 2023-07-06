package handler

import (
	"errors"
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/localities"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/optional"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

type Locality struct {
	locService localities.Service
}

func NewLocality(svc localities.Service) *Locality {
	return &Locality{
		locService: svc,
	}
}

// Create godoc
//
//	@Summary	Create new locality
//	@Tags		Localities
//	@Accept		json
//	@Produce	json
//	@Param		locality	body		localities.CreateDTO	true	"Locality to be added"
//	@Success	201			{object}	web.response			"Returns created locality"
//	@Failure	409			{object}	web.errorResponse		"Locality is not unique"
//	@Failure	422			{object}	web.errorResponse		"Missing fields or invalid field types"
//	@Failure	500			{object}	web.errorResponse		"Could not save locality"
//	@Router		/api/v1/localities [post]
func (h *Locality) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		dto := middleware.GetBody[localities.CreateDTO](c)

		loc, err := h.locService.Create(c.Request.Context(), dto)
		if err != nil {
			status := mapLocalityErrToStatus(err)
			web.Error(c, status, err.Error())
			return
		}

		web.Success(c, http.StatusCreated, loc)
	}
}

// ReportAll godoc
//
//	@Summary	Return seller count for each locality
//	@Tags		Localities
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	web.response		"Returns seller count for localities"
//	@Success	204	{object}	web.response		"No content was found"
//	@Failure	500	{object}	web.errorResponse	"Could not save locality"
//	@Router		/api/v1/localities/report-sellers [get]
func (h *Locality) SellerReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id optional.Opt[int]
		if _, hasId := c.Params.Get("id"); hasId {
			id = *optional.FromVal(c.GetInt("id"))
		}

		report, err := h.locService.CountSellers(c, id)

		if err != nil {
			status := mapLocalityErrToStatus(err)
			web.Error(c, status, err.Error())
			return
		}

		if len(report) == 0 {
			web.Success(c, http.StatusNoContent, report)
			return
		}
		web.Success(c, http.StatusOK, report)
	}
}

// ReportByID godoc
//
//	@Summary	Return seller count for given locality
//	@Tags		Localities
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Locality ID"
//	@Success	200	{object}	web.response		"Returns seller count for locality"
//	@Failure	404	{object}	web.response		"ID was not found"
//	@Failure	500	{object}	web.errorResponse	"Could not save locality"
//	@Router		/api/v1/localities/report-sellers/{id} [get]x
func _() {} // Implementation is in the SellerReport function

func mapLocalityErrToStatus(err error) int {
	var invalidLocality *localities.ErrInvalidLocality
	var notFound *localities.ErrNotFound

	if errors.As(err, &invalidLocality) {
		return http.StatusConflict
	}
	if errors.As(err, &notFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
