package handler

import (
	"net/http"
	"strconv"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/warehouse"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web"
	"github.com/gin-gonic/gin"
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
// @Summary Retrieve a warehouse
// @Description Get a warehouse by ID
// @Tags warehouses
// @Param id path int true "Warehouse ID"
// @Produce json
// @Success 200 {object} domain.Warehouse
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Warehouse not found"
// @Router /warehouses/{id} [get]
func (w *Warehouse) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusBadRequest, "invalid id")
			return
		}

		warehouse, err := w.warehouseService.Get(c, id)
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

// GetAll retrieves all warehouses.
// @Summary Retrieve all warehouses
// @Description Get all warehouses
// @Tags warehouses
// @Produce json
// @Success 200 {array} domain.Warehouse
// @Failure 400 {string} string "Warehouse not found"
// @Router /warehouses [get]
func (w *Warehouse) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		warehouses, err := w.warehouseService.GetAll(c)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "warehouse not found")
			return
		}
		web.Success(c, http.StatusOK, warehouses)
	}
}

// Create creates a new warehouse.
// @Summary Create a warehouse
// @Description Create a new warehouse
// @Tags warehouses
// @Accept json
// @Produce json
// @Param warehouse body domain.Warehouse true "Warehouse object"
// @Success 201 {object} domain.Warehouse
// @Failure 400 {string} string "Invalid request"
// @Failure 422 {string} string "Warehouse not created"
// @Router /warehouses [post]
func (w *Warehouse) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var warehouse domain.Warehouse
		if err := c.ShouldBindJSON(&warehouse); err != nil {
			web.Error(c, http.StatusBadRequest, "warehouse not created")
			return
		}
		if warehouse.WarehouseCode == "" {
			web.Error(c, http.StatusBadRequest, "warehouse need to be passed")
			return
		}

		warehouse, err := w.warehouseService.Create(c, warehouse)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, "warehouse not created")
			return
		}
		web.Success(c, http.StatusCreated, warehouse)
	}
}

// Update updates a warehouse.
// @Summary Update a warehouse
// @Description Update a warehouse by ID
// @Tags warehouses
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Param warehouse body domain.Warehouse true "Updated warehouse object"
// @Success 200 {object} domain.Warehouse
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Invalid ID"
// @Failure 405 {string} string "Warehouse not updated"
// @Router /warehouses/{id} [patch]
func (w *Warehouse) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		var warehouse domain.Warehouse
		if err := c.ShouldBindJSON(&warehouse); err != nil {
			web.Error(c, http.StatusBadRequest, "warehouse not updated")
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		warehouse.ID = id
		warehouse, err = w.warehouseService.Update(c, warehouse)
		if err != nil {
			web.Error(c, http.StatusMethodNotAllowed, "warehouse not updated")
			return
		}
		web.Success(c, http.StatusOK, warehouse)
	}
}

// Delete deletes a warehouse.
// @Summary Delete a warehouse
// @Description Delete a warehouse by ID
// @Tags warehouses
// @Param id path int true "Warehouse ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID"
// @Failure 405 {string} string "Warehouse not deleted"
// @Router /warehouses/{id} [delete]
func (w *Warehouse) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			web.Error(c, http.StatusNotFound, "invalid id")
			return
		}
		err = w.warehouseService.Delete(c, id)
		if err != nil {
			web.Error(c, http.StatusMethodNotAllowed, "warehouse not deleted")
			return
		}
		web.Success(c, http.StatusNoContent, nil)
	}
}
