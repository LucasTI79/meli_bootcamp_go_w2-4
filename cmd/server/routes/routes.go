package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/extmatperez/meli_bootcamp_go_w2-4/cmd/server/handler"
	"github.com/extmatperez/meli_bootcamp_go_w2-4/internal/employee"

	_ "github.com/extmatperez/meli_bootcamp_go_w2-4/docs"
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
	// Seller routes
	// Example:
	// repo := seller.NewRepository(r.db)
	// service := seller.NewService(repo)
	// handler := handler.NewSeller(service)
	// r.rg.GET("/seller", handler.GetAll)
}

func (r *router) buildProductRoutes() {
	// Product routes
}

func (r *router) buildSectionRoutes() {
	// Section routes
}

func (r *router) buildWarehouseRoutes() {
	// Warehouse routes
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
	// Buyer routes
}
