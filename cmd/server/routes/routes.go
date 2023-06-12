package routes

import (
	"database/sql"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
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

func (r *router) buildSectionRoutes() {
	repository := section.NewRepository(r.db)
	service := section.NewService(repository)
	handler := handler.NewSection(service)
	sec := r.rg.Group("sections")
	{
		sec.POST("/", handler.Create())
		sec.GET("/", handler.GetAll())
		sec.GET("/:id", handler.Get())
		sec.DELETE("/:id", handler.Delete())
		sec.PATCH(":id", handler.Update())
	}

}

func (r *router) buildWarehouseRoutes() {}

func (r *router) buildEmployeeRoutes() {}

func (r *router) buildBuyerRoutes() {}
