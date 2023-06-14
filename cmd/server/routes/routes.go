package routes

import (
	"database/sql"

	// _ "github.com/extmatperez/meli_bootcamp_go_w2-4/docs"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/warehouse"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.buildDocumentationRoutes()

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

func (r *router) buildDocumentationRoutes() {
	r.rg.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *router) buildSellerRoutes() {
}

func (r *router) buildProductRoutes() {
	repo := product.NewRepository(r.db)
	service := product.NewService(repo)
	h := handler.NewProduct(service)

	productRG := r.rg.Group("/products")
	{
		productRG.POST("/", middleware.JSONMapper[handler.CreateRequest](), h.Create())
		productRG.GET("/", h.GetAll())
		productRG.GET("/:id", h.Get())
		productRG.PATCH("/:id", middleware.JSONMapper[handler.UpdateRequest](), h.Update())
		productRG.DELETE("/:id", h.Delete())
	}
}

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

func (r *router) buildWarehouseRoutes() {
	repo := warehouse.NewRepository(r.db)
	service := warehouse.NewService(repo)
	handler := handler.NewWarehouse(service)
	r.rg.POST("/warehouses", handler.Create())
	r.rg.GET("/warehouses", handler.GetAll())
	r.rg.GET("/warehouses/:id", handler.Get())
	r.rg.PATCH("/warehouses/:id", handler.Update())
	r.rg.DELETE("/warehouses/:id", handler.Delete())
}

func (r *router) buildEmployeeRoutes() {}

func (r *router) buildBuyerRoutes() {}
