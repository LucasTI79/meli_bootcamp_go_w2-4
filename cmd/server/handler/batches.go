package handler

import (
	"net/http"
	"time"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/batches"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
	"github.com/gin-gonic/gin"
)

type Batches struct {
	service batches.Service
}

type CreateBatchesRequest struct {
	BatchNumber        int    `binding:"required" json:"batch_number"`
	CurrentQuantity    int    `binding:"required" json:"current_quantity"`
	CurrentTemperature int    `binding:"required" json:"current_temperature"`
	DueDate            string `binding:"required" json:"due_date"`
	InitialQuantity    int    `binding:"required" json:"initial_quantity"`
	ManufacturingDate  string `binding:"required" json:"manufacturing_date"`
	ManufacturingHour  int    `binding:"required" json:"manufacturing_hour"`
	MinimumTemperature int    `binding:"required" json:"minimum_temperature"`
	ProductID          int    `binding:"required" json:"product_id"`
	SectionID          int    `binding:"required" json:"section_id"`
}

func ConvertDate(c CreateBatchesRequest) (batches.CreateBatches, error) {
	DueDate, err := time.Parse("2006-01-02", c.DueDate)
	if err != nil {
		return batches.CreateBatches{}, err
	}
	ManufacturingDate, err := time.Parse("2006-01-02", c.ManufacturingDate)
	if err != nil {
		return batches.CreateBatches{}, err
	}
	return batches.CreateBatches{
		BatchNumber:        c.BatchNumber,
		CurrentQuantity:    c.CurrentQuantity,
		CurrentTemperature: c.CurrentTemperature,
		DueDate:            DueDate,
		InitialQuantity:    c.CurrentQuantity,
		ManufacturingDate:  ManufacturingDate,
		ManufacturingHour:  c.ManufacturingHour,
		MinimumTemperature: c.CurrentTemperature,
		ProductID:          c.ProductID,
		SectionID:          c.SectionID,
	}, err
}

func NewBatches(s batches.Service) *Batches {
	return &Batches{
		service: s,
	}
}

// Create godoc
//
// @Summary	Create a new batch
// @Tags		Batches
// @Accept		json
// @Produce	json
// @Param		request	body	CreateBatchesRequest	true	"Batch data"
// @Success	201	{object}	web.response	"Created batch"
// @Failure	400	{object}	web.errorResponse	"Invalid request body"
// @Failure	409	{object}	web.errorResponse	"Batch number already exists"
// @Failure	500	{object}	web.errorResponse	"Failed to create batch"
// @Router	/api/v1/batches [post]
func (s *Batches) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		dto := middleware.GetBody[CreateBatchesRequest](c)
		convertdate, err := ConvertDate(dto)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, err.Error())
			return
		}
		batch, err := s.service.Create(c, convertdate)
		if err != nil {
			if err == batches.ErrInvalidBatchNumber {
				web.Error(c, http.StatusConflict, err.Error())
				return
			} else {
				web.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
		web.Success(c, http.StatusCreated, batch)
	}
}
