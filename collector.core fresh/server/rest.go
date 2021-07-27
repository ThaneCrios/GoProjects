package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/web"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/handler"
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
	g.GET("/orders/free_orders/:uuid", r.GetFreeOrders)
	g.GET("/order/current_order_products", r.GetProductsFromOrderByFilter)
	g.GET("/orders/:uuid", r.GetOrderByUUID)
	g.GET("/orders/duplicate/:uuid", r.GetDuplicateOrderByUUID)
	g.GET("/orders/my_orders/:uuid", r.GetMyOrders)
	g.POST("/order/create", r.CreateOrder)
	g.PUT("/orders/grab_order", r.GrabOrder)
	g.PUT("/order/mark_product", r.MarkProduct)
	g.PUT("/order/change_product", r.ChangeProduct)
	g.PUT("/order/finish", r.FinishCollectOrder)
	g.PUT("/order/cancel", r.CancelOrder)
	g.PUT("/order/remove_product", r.RemoveProductFromOrder)
	g.PUT("/order/add_product", r.AddProductToOrder)

	//collectors
	g.GET("/collector/:uuid", r.GetCollectorByUUID)
	g.POST("/collector/create", r.CreateCollector)

	//products
	g.GET("/products/:barcode", r.GetProductByBarCode)
	g.GET("/product/:uuid", r.GetProductByUUID)
	g.GET("/products/uuids_with_barcodes", r.GetProductUUIDsWithBarCodes)
	g.POST("/product/create", r.CreateProduct)
	g.PUT("/products/appoint_barcode", r.AppointBarCode)

}
