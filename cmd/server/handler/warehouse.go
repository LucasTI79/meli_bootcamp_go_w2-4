package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/warehouse"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
)

var (
	ErrAlrearyExist        = "warehouse can be alreary exist"
	ErrWarehouseCodeUnique = "warehouse code must be unique"
	ErrUnprocessableEntity = "action could not be processed correctly due to invalid data provided"
	ErrServerInternalError = "something went wrong with the request"
	ErrWarehouseEmpty      = "warehousecode need to be passed, it can't be empty"
	ErrWarehouseNotFound   = "Warehouse not found"
	ErrInvalidID           = "Invalid ID"
	ErrWarehouseNotDeleted = "Warehouse not deleted"
)

type Warehouse struct {
	warehouseService warehouse.Service
}

func NewWarehouse(w warehouse.Service) *Warehouse {
	return &Warehouse{
		warehouseService: w,
	}
}

// Get retrieves a warehouse by ID.
//
//	@Summary		Retrieve a warehouse
//	@Description	Get a warehouse by ID
//	@Tags			Warehouses
//	@Param			id	path	int	true	"Warehouse ID"
//	@Produce		json
//	@Success		200	{object}	domain.Warehouse
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Warehouse not found"
//	@Router			/api/v1/warehouses/{id} [get]
func (w *Warehouse) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, ErrInvalidID)
			return
		}

		warehouse, err := w.warehouseService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, ErrWarehouseNotFound)
			return
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

// GetAll retrieves all warehouses.
//
//	@Summary		Retrieve all warehouses
//	@Description	Get all warehouses
//	@Tags			Warehouses
//	@Produce		json
//	@Success		200	{array}	domain.Warehouse
//	@Success		204	"warehouses is empty"
//	@Failure		500	{string}	string	"something went wrong with the request"
//	@Router			/api/v1/warehouses [get]
func (w *Warehouse) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouses, err := w.warehouseService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, ErrServerInternalError)
			return
		}
		if len(warehouses) == 0 {
			web.Success(c, http.StatusNoContent, warehouses)
			return
		}
		web.Success(c, http.StatusOK, warehouses)
	}
}

// Create creates a new warehouse.
//
//	@Summary		Create a warehouse
//	@Description	Create a new warehouse
//	@Tags			Warehouses
//	@Accept			json
//	@Produce		json
//	@Param			warehouse	body		domain.Warehouse	true	"Warehouse object"
//	@Success		201			{object}	domain.Warehouse
//	@Failure		422			{string}	string	"warehousecode need to be passed, it can't be empty"
//	@Failure		500			{string}	string	err.Error()
//	@Failure		409			{string}	string	"warehouse can be alreary exist"
//	@Router			/api/v1/warehouses [post]
func (w *Warehouse) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var warehouse domain.Warehouse
		if err := c.ShouldBindJSON(&warehouse); err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			fmt.Println(err.Error())
			return
		}
		if warehouse.WarehouseCode == "" {
			web.Error(c, http.StatusUnprocessableEntity, ErrWarehouseEmpty)
			return
		}

		warehouse, err := w.warehouseService.Create(c, warehouse)
		if err != nil {
			web.Error(c, http.StatusConflict, ErrAlrearyExist)
			return
		}
		web.Success(c, http.StatusCreated, warehouse)
	}
}

// Update updates a warehouse.
//
//	@Summary		Update a warehouse
//	@Description	Update a warehouse by ID
//	@Tags			Warehouses
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int					true	"Warehouse ID"
//	@Param			warehouse	body		domain.Warehouse	true	"Updated warehouse object"
//	@Success		200			{object}	domain.Warehouse
//	@Failure		422			{string}	string	"action could not be processed correctly due to invalid data provided"
//	@Failure		404			{string}	string	"Invalid ID"
//	@Failure		409			{string}	string	"Warehouse code must be unique"
//	@Router			/api/v1/warehouses/{id} [patch]
func (w *Warehouse) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		var warehouse domain.Warehouse
		if err := c.ShouldBindJSON(&warehouse); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrUnprocessableEntity)
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusNotFound, ErrInvalidID)
			return
		}
		warehouse.ID = id
		warehouse, err = w.warehouseService.Update(c, warehouse)
		if err != nil {
			web.Error(c, http.StatusConflict, ErrWarehouseCodeUnique)
			return
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

// Delete deletes a warehouse.
//
//	@Summary		Delete a warehouse
//	@Description	Delete a warehouse by ID
//	@Tags			Warehouses
//	@Param			id	path	int	true	"Warehouse ID"
//	@Success		204	"No Content"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		405	{string}	string	"Warehouse not deleted"
//	@Router			/api/v1/warehouses/{id} [delete]
func (w *Warehouse) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusNotFound, ErrInvalidID)
			return
		}
		err = w.warehouseService.Delete(c, id)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, ErrWarehouseNotDeleted)
			return
		}
		web.Success(c, http.StatusNoContent, nil)
	}
}
