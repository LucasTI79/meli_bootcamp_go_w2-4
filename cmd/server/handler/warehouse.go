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
