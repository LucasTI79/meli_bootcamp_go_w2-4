package routes

import (
	"database/sql"

	_ "github.com/extmatperez/meli_bootcamp_go_w2-4/docs"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/batches"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/buyer"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/carrier"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/domain"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"
	inboundOrder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/inbound_order"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/localities"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/product"
	purchaseorder "github.com/extmatperez/meli_bootcamp_go_w2-4/internal/purchase_order"
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
	r.buildBatchRoutes()
	r.buildInboundOrderRoutes()
	r.buildCarrierRoutes()
	r.buildLocalityRoutes()
	r.buildPurchaseOrderRoutes()
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

	r.rg.POST("/product-records/", middleware.Body[handler.CreateRequestRecord](), h.CreateRecord())
	productRG := r.rg.Group("/products")
	{
		productRG.POST("/", middleware.Body[handler.CreateRequest](), h.Create())
		productRG.GET("/", h.GetAll())
		productRG.GET("/:id", middleware.IntPathParam(), h.Get())
		productRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[handler.UpdateRequest](), h.Update())
		productRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
		productRG.GET("/report-records", h.GetRecords())
		productRG.GET("/report-records/:id", middleware.IntPathParam(), h.GetRecords())
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
		sec.GET("/report-products", h.GetAllReportProducts())
		sec.GET("/report-products/:id", middleware.IntPathParam(), h.GetReportProducts())
	}
}

func (r *router) buildWarehouseRoutes() {
	repo := warehouse.NewRepository(r.db)
	service := warehouse.NewService(repo)
	h := handler.NewWarehouse(service)

	rg := r.rg.Group("/warehouses")
	{
		rg.POST("", middleware.Body[domain.Warehouse](), h.Create())
		rg.GET("", h.GetAll())
		rg.GET("/:id", middleware.IntPathParam(), h.Get())
		rg.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Warehouse](), h.Update())
		rg.DELETE("/:id", middleware.IntPathParam(), h.Delete())
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
		employeeRG.GET("/report-inbound-orders/:id", middleware.IntPathParam(), h.GetInboundReport())
		employeeRG.GET("/report-inbound-orders/", h.GetInboundReport())
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
		buyerRG.GET("/report-purchase-orders/:id", middleware.IntPathParam(), h.PurchaseOrderReport())
		buyerRG.GET("/report-purchase-orders/", h.PurchaseOrderReport())
		buyerRG.DELETE("/:id", middleware.IntPathParam(), h.Delete())
		buyerRG.PATCH("/:id", middleware.IntPathParam(), middleware.Body[domain.Buyer](), h.Update())
	}
}

func (r *router) buildBatchRoutes() {
	repo := batches.NewRepository(r.db)
	service := batches.NewService(repo)
	h := handler.NewBatches(service)

	batchRG := r.rg.Group("/product-batches")
	{
		batchRG.POST("", middleware.Body[handler.CreateBatchesRequest](), h.Create())
	}
}

func (r *router) buildInboundOrderRoutes() {
	repo := inboundOrder.NewRepository(r.db)
	service := inboundOrder.NewService(repo)
	h := handler.NewInboundOrder(service)

	buyerRG := r.rg.Group("/inbound-orders")
	{
		buyerRG.POST("", middleware.Body[handler.InboundOrderRequest](), h.Create())
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

func (r *router) buildLocalityRoutes() {
	repo := localities.NewRepository(r.db)
	service := localities.NewService(repo)
	h := handler.NewLocality(service)

	rg := r.rg.Group("/localities")
	{
		rg.POST("", middleware.Body[localities.CreateDTO](), h.Create())
		rg.GET("/report-sellers", h.SellerReport())
		rg.GET("/report-sellers/:id", middleware.IntPathParam(), h.SellerReport())
		rg.GET("/report-carriers", h.CarrierReport())
		rg.GET("/report-carriers/:id", middleware.IntPathParam(), h.CarrierReport())
	}
}

func (r *router) buildPurchaseOrderRoutes() {
	repo := purchaseorder.NewRepository(r.db)
	service := purchaseorder.NewService(repo)
	h := handler.NewPurchaseOrder(service)

	purchaseOrderRG := r.rg.Group("/purchase-orders")
	{
		purchaseOrderRG.POST("", middleware.Body[handler.PurchaseOrderRequest](), h.Create())
	}
}
