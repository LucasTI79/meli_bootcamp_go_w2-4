package routes

import (
	"database/sql"

	_ "github.com/extmatperez/meli_bootcamp_go_w2-4/docs"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/section"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/seller"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/warehouse"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/pkg/web/middleware"
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

func (r *router) buildDocumentationRoutes() {
	r.rg.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildSellerRoutes() {
	repo := seller.NewRepository(r.db)
	service := seller.NewService(repo)
	handler := handler.NewSeller(service)

	sellerGroup := r.rg.Group("/sellers")
	{
		sellerGroup.GET("/", handler.GetAll())
		sellerGroup.GET("/:id", handler.GetById())
		sellerGroup.POST("/", handler.Create())
		sellerGroup.PATCH("/:id", handler.Update())
		sellerGroup.DELETE("/:id", handler.Delete())
	}
}

func (r *router) buildProductRoutes() {
	repo := product.NewRepository(r.db)
	service := product.NewService(repo)
	h := handler.NewProduct(service)

	productRG := r.rg.Group("/products")
	{
		productRG.POST("/", middleware.JSONMapper[handler.CreateRequest](), h.Create())
		productRG.GET("/", h.GetAll())
		productRG.GET("/:id", middleware.IntPathParam(), h.Get())
		productRG.PATCH("/:id", middleware.IntPathParam(), middleware.JSONMapper[handler.UpdateRequest](), h.Update())
		productRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
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

func (r *router) buildEmployeeRoutes() {
	// Employee routes
	// Example:
	repo := employee.NewRepository(r.db)
	service := employee.NewService(repo)
	handler := handler.NewEmployee(service)

	r.rg.GET("/employees", handler.GetAll())
	r.rg.POST("/employees", handler.Create())
	r.rg.GET("/employees/:id", handler.Get())
	r.rg.PATCH("/employees/:id", handler.Update())
	r.rg.DELETE("/employees/:id", handler.Delete())
}

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
