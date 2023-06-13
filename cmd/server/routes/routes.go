package routes

import (
	"database/sql"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/gin-gonic/gin"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func NewRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{eng: eng, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()

	r.buildSellerRoutes()
	r.buildProductRoutes()
	r.buildSectionRoutes()
	r.buildWarehouseRoutes()
	r.buildEmployeeRoutes()
	r.buildBuyerRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildSellerRoutes() {
	// Example
	// repo := seller.NewRepository(r.db)
	// service := seller.NewService(repo)
	// handler := handler.NewSeller(service)
	// r.r.GET("/seller", handler.GetAll)
}

func (r *router) buildProductRoutes() {}

func (r *router) buildSectionRoutes() {}

func (r *router) buildWarehouseRoutes() {}

func (r *router) buildEmployeeRoutes() {}

func (r *router) buildBuyerRoutes() {
	repo := buyer.NewRepository(r.db)
	service := buyer.NewService(repo)
	handler := handler.NewBuyer(service)
	r.rg.GET("/buyers", handler.GetAll())
	r.rg.POST("/buyers", handler.Create())
	r.rg.GET("/buyers/:id", handler.Get())
	r.rg.DELETE("/buyers/:id", handler.Delete())
	r.rg.PATCH("/buyers/:id", handler.Update())
}
