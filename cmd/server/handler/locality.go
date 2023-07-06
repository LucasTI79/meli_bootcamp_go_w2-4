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
