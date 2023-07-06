package handler

import (
	"errors"
	"net/http"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/carrier"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

type Carrier struct {
	carrierService carrier.Service
}

type CarrierRequest struct {
	CID         *int    `binding:"required" json:"cid"`
	CompanyName *string `binding:"required" json:"company_name"`
	Address     *string `binding:"required" json:"address"`
	Telephone   *string `binding:"required" json:"telephone"`
	LocalityID  *int    `binding:"required" json:"locality_id"`
}

func NewCarrier(s carrier.Service) *Carrier {
	return &Carrier{
		carrierService: s,
	}
}

// Create carrier
//
//	@Summary	Create new carrier
//	@Tags		Carrier
//	@Accept		json
//	@Produce	json
//	@Param		product	body		CarrierRequest		true	"Carrier to be added"
//	@Success	201		{object}	web.response		"Returns created carrier"
//	@Failure	409		{object}	web.errorResponse	"`cid` is not unique or `locality_id` not found"
//	@Failure	422		{object}	web.errorResponse	"Missing fields or invalid field types"
//	@Failure	500		{object}	web.errorResponse	"Could not save carrier"
//	@Router		/api/v1/carrier [post]
func (i *Carrier) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := middleware.GetBody[CarrierRequest](c)
		dto := mapCarrierRequestToDTO(&req)
		p, err := i.carrierService.Create(c.Request.Context(), *dto)

		if err != nil {
			errStatus := checkError(err)
			web.Error(c, errStatus, err.Error())
			return
		}

		web.Success(c, http.StatusCreated, p)
	}
}

func checkError(err error) int {
	if errors.Is(err, carrier.ErrAlreadyExists) ||
		errors.Is(err, carrier.ErrLocalityIDNotFound) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func mapCarrierRequestToDTO(req *CarrierRequest) *carrier.CarrierDTO {
	return &carrier.CarrierDTO{
		CID:         *req.CID,
		CompanyName: *req.CompanyName,
		Address:     *req.Address,
		Telephone:   *req.Telephone,
		LocalityID:  *req.LocalityID,
	}
}
