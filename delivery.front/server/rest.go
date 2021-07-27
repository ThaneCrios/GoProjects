package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/web"
	"gitlab.com/faemproject/backend/delivery/delivery.front/handler"
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
	g.GET("/orders/:uuid", r.Handler.GetOrderByUUID)
	g.GET("/orders/filter", r.Handler.OrdersFilter)
	g.POST("/orders", r.Handler.CreateOrder)
	g.PUT("/orders/:uuid", r.Handler.UpdateOrder)
	g.PUT("/order/delete/:uuid", r.Handler.MarkOrderAsDeleted)
	g.PUT("/orders", r.Handler.UpdateStatusOrder)
	g.PUT("/order/set_courier", r.Handler.SetCourierToOrder)
	g.GET("/order/states", r.Handler.OrderStates)
	g.PUT("/order/remove_courier", r.Handler.OrderRemoveCourier)
	g.PUT("/order/delete/:uuid", r.Handler.OrderDeleted)

	//courier
	g.GET("/couriers/:uuid", r.Handler.GetCourierByUUID)
	g.POST("/courier", r.Handler.CreateCourier)
	g.GET("/couriers/filter", r.Handler.GetFilteredCouriers)
	g.PUT("/courier/:uuid", r.Handler.UpdateCourier)
	g.DELETE("/courier/:uuid", r.Handler.MarkCourierAsDeleted)
	g.PUT("/courier/set_status", r.Handler.UpdateStatusCourier)
	g.PUT("/courier/delete/:uuid", r.Handler.DeleteCourier)
	g.GET("/courier/states", r.Handler.CourierStates)
	g.PUT("/courier/update_coor", r.Handler.UpdateCoordinate)
	g.GET("/courier_info/:phoneNumber", r.Handler.GetCourierTokenInfoByPhoneNumber)

	//tasks
	g.GET("/tasks/courier_tasks/:uuid", r.Handler.GetCourierTasks)
	g.PUT("/tasks/change_queue", r.Handler.ChangeQueue)
	g.GET("/tasks/order/:uuid", r.Handler.GetOrderTasks)
	g.PUT("/task/update", r.Handler.UpdateTask)
	g.PUT("/task/finish/:uuid", r.Handler.FinishTask)
	g.GET("/next_task/:uuid", r.Handler.NextTaskCourier)

}
