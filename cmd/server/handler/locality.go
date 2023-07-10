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

type SellerReportEntry struct {
	LocalityID   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	SellerCount  int    `json:"sellers_count"`
}

type CarrierReportEntry struct {
	LocalityID   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	CarrierCount int    `json:"carriers_count"`
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

// ReportAllSellers godoc
//
//	@Summary	Return seller count for each locality
//	@Tags		Localities
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	web.response		"Returns seller count for localities"
//	@Success	204	{object}	web.response		"No content was found"
//	@Failure	500	{object}	web.errorResponse	"Could not generate report"
//	@Router		/api/v1/localities/report-sellers [get]
func (h *Locality) SellerReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id optional.Opt[int]
		if _, hasId := c.Params.Get("id"); hasId {
			id = *optional.FromVal(c.GetInt("id"))
		}

		report, err := h.locService.CountSellers(c, id)
		data := MapSellerReportToDTO(report)

		if err != nil {
			status := mapLocalityErrToStatus(err)
			web.Error(c, status, err.Error())
			return
		}

		if len(report) == 0 {
			web.Success(c, http.StatusNoContent, data)
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}

// ReportSellerByID godoc
//
//	@Summary	Return seller count for given locality
//	@Tags		Localities
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Locality ID"
//	@Success	200	{object}	web.response		"Returns seller count for locality"
//	@Failure	404	{object}	web.response		"ID was not found"
//	@Failure	500	{object}	web.errorResponse	"Could generate report"
//	@Router		/api/v1/localities/report-sellers/{id} [get]
func _() {} // Implementation is in the SellerReport function

// ReportAllCarriers godoc
//
//	@Summary	Return carrier count for each locality
//	@Tags		Localities
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	web.response		"Returns carrier count for localities"
//	@Success	204	{object}	web.response		"No content was found"
//	@Failure	500	{object}	web.errorResponse	"Could not generate report"
//	@Router		/api/v1/localities/report-carriers [get]
func (h *Locality) CarrierReport() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id optional.Opt[int]
		if _, hasId := c.Params.Get("id"); hasId {
			id = *optional.FromVal(c.GetInt("id"))
		}

		report, err := h.locService.CountCarriers(c, id)
		data := MapCarrierReportToDTO(report)

		if err != nil {
			status := mapLocalityErrToStatus(err)
			web.Error(c, status, err.Error())
			return
		}

		if len(report) == 0 {
			web.Success(c, http.StatusNoContent, data)
			return
		}
		web.Success(c, http.StatusOK, data)
	}
}

// ReportCarriersByID godoc
//
//	@Summary	Return carrier count for given locality
//	@Tags		Localities
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int					true	"Locality ID"
//	@Success	200	{object}	web.response		"Returns carrier count for locality"
//	@Failure	404	{object}	web.response		"ID was not found"
//	@Failure	500	{object}	web.errorResponse	"Could not generate report"
//	@Router		/api/v1/localities/report-carriers/{id} [get]
func _() {} // Implementation is in the CarrierReport function

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

func MapCarrierReportToDTO(report []localities.CountByLocality) []CarrierReportEntry {
	dtos := make([]CarrierReportEntry, 0)

	for _, entry := range report {
		dtos = append(dtos, mapCountToCarrierReport(entry))
	}

	return dtos
}

func MapSellerReportToDTO(report []localities.CountByLocality) []SellerReportEntry {
	dtos := make([]SellerReportEntry, 0)

	for _, entry := range report {
		dtos = append(dtos, mapCountToSellerReport(entry))
	}

	return dtos
}

func mapCountToCarrierReport(count localities.CountByLocality) CarrierReportEntry {
	return CarrierReportEntry{
		LocalityID:   count.ID,
		LocalityName: count.Name,
		CarrierCount: count.Count,
	}
}

func mapCountToSellerReport(count localities.CountByLocality) SellerReportEntry {
	return SellerReportEntry{
		LocalityID:   count.ID,
		LocalityName: count.Name,
		SellerCount:  count.Count,
	}
}
