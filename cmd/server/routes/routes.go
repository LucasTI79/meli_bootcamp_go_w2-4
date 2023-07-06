package routes

import (
	"database/sql"

	_ "github.com/extmatperez/meli_bootcamp_go_w2-4/docs"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/carrier"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
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
	r.buildCarrierRoutes()
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
		sellerGroup.GET("/:id", middleware.IntPathParam(), handler.Get())
		sellerGroup.POST("/", middleware.Body[domain.Seller](), handler.Create())
		sellerGroup.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Seller](), handler.Update())
		sellerGroup.DELETE("/:id", middleware.IntPathParam(), handler.Delete())
	}
}

func (r *router) buildProductRoutes() {
	repo := product.NewRepository(r.db)
	service := product.NewService(repo)
	h := handler.NewProduct(service)

	productRG := r.rg.Group("/products")
	{
		productRG.POST("/", middleware.Body[handler.CreateRequest](), h.Create())
		productRG.GET("/", h.GetAll())
		productRG.GET("/:id", middleware.IntPathParam(), h.Get())
		productRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[handler.UpdateRequest](), h.Update())
		productRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
	}
}

func (r *router) buildSectionRoutes() {
	repository := section.NewRepository(r.db)
	service := section.NewService(repository)
	h := handler.NewSection(service)
	sec := r.rg.Group("/sections")
	{
		sec.POST("", middleware.Body[section.CreateSection](), h.Create())
		sec.GET("", h.GetAll())
		sec.GET("/:id", middleware.IntPathParam(), h.Get())
		sec.DELETE("/:id", middleware.IntPathParam(), h.Delete())
		sec.PATCH("/:id", middleware.IntPathParam(), middleware.Body[section.UpdateSection](), h.Update())
	}
}

func (r *router) buildWarehouseRoutes() {
	repo := warehouse.NewRepository(r.db)
	service := warehouse.NewService(repo)
	h := handler.NewWarehouse(service)

	productRG := r.rg.Group("/warehouses")
	{
		productRG.POST("", middleware.Body[domain.Warehouse](), h.Create())
		productRG.GET("", h.GetAll())
		productRG.GET("/:id", middleware.IntPathParam(), h.Get())
		productRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Warehouse](), h.Update())
		productRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
	}
}

func (r *router) buildEmployeeRoutes() {
	repo := employee.NewRepository(r.db)
	svc := employee.NewService(repo)
	h := handler.NewEmployee(svc)

	employeeRG := r.rg.Group("/employees")
	{
		employeeRG.GET("", h.GetAll())
		employeeRG.POST("", middleware.Body[domain.Employee](), h.Create())
		employeeRG.GET("/:id", middleware.IntPathParam(), h.Get())
		employeeRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Employee](), h.Update())
		employeeRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
	}
}

func (r *router) buildBuyerRoutes() {
	repo := buyer.NewRepository(r.db)
	service := buyer.NewService(repo)
	h := handler.NewBuyer(service)

	buyerRG := r.rg.Group("/buyers")
	{
		buyerRG.GET("", h.GetAll())
		buyerRG.POST("", middleware.Body[domain.BuyerCreate](), h.Create())
		buyerRG.GET("/:id", middleware.IntPathParam(), h.Get())
		buyerRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
		buyerRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Buyer](), h.Update())
	}
}

func (r *router) buildCarrierRoutes() {
	repo := carrier.NewRepository(r.db)
	service := carrier.NewService(repo)
	h := handler.NewCarrier(service)

	productRG := r.rg.Group("/carrier")
	{
		productRG.POST("/", middleware.Body[handler.CarrierRequest](), h.Create())
	}
}
