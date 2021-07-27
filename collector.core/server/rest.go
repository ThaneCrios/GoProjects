package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/web"
	"gitlab.com/faemproject/backend/delivery/collector.core/handler"
)

const (
	apiPrefix = "/api/v3"
)

type Rest struct {
	Router  *echo.Echo
	Handler *handler.Handler
}

// Route defines all the application rest endpoints
func (r *Rest) Route() {
	web.UseHealthCheck(r.Router)

	g := r.Router.Group(apiPrefix)

	//orders
	g.GET("/orders/free_orders", r.GetFreeOrders)
	g.PUT("/orders/grab_order", r.GrabOrder)
	g.GET("/orders/:uuid", r.GetOrderByUUID)
	g.POST("/order/create", r.CreateOrder)
	g.PUT("/order/mark_product", r.MarkProduct)
	g.PUT("/order/change_product", r.ChangeProduct)
	g.GET("/order/current_order_products", r.GetProductsFromOrderByFilter)
	g.PUT("/order/finish", r.FinishCollectOrder)

	//collectors
	g.GET("/collector/:uuid", r.GetCollectorByUUID)
	g.POST("/collector/create", r.CreateCollector)

	//products
	g.POST("/product/create", r.CreateProduct)
	g.GET("/product/get_by_barcode/:barCode", r.GetProductByBarCode)
	//g.GET("/products/compare", r.CompareTwoProducts)
	g.GET("/product/:uuid", r.GetProductByUUID)
	g.POST("/products/appoint_barcode", r.AppointBarCode)
}
