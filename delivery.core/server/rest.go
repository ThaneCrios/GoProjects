package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/web"
	"gitlab.com/faemproject/backend/delivery/delivery.core/handler"
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

	//auth
	g.GET("/courier/:phone_number", r.GetCourierTokenInfoByPhoneNumber)
	//g.GET("/user/auth/:login", r.AuthCourier)

	//orders
	g.GET("/orders/:uuid", r.GetOrderByUUID)
	g.GET("/orders/filter", r.OrdersFilter)
	g.POST("/orders", r.CreateOrder)
	g.PUT("/orders", r.UpdateOrder)
	g.DELETE("/order/delete/:uuid", r.DeleteOrder)
	//g.PUT("/orders", r.UpdateStatusOrder)
	g.PUT("/order/set_courier", r.SetCourierToOrder)
	g.PUT("/order/remove_courier", r.RemoveCourierFromOrder)

	// couriers
	g.GET("/couriers/:uuid", r.GetCourierByUUID)
	g.POST("/courier", r.CreateCourier)
	g.GET("/couriers/filter", r.CouriersFilter)
	g.PUT("/courier/:uuid", r.UpdateCourier)
	g.DELETE("/courier/delete/:uuid", r.DeleteCourier)
	g.PUT("/courier/set_status", r.UpdateStatusCourier)
	g.GET("/courier", r.GetCourierByChatID)
	g.PUT("/courier/update_coor", r.UpdateCourierCoordinates)
	g.GET("/courier/get_uuid/:chat_id", r.GetCourierUUIDByChatID)

	//tasks
	g.GET("/tasks/cour_tasks/status", r.GetCourierTasksByStatus)
	g.GET("/tasks/status", r.TasksListByStatus)
	g.PUT("/task/update", r.UpdateTask)
	g.GET("/tasks/courier_tasks/:uuid", r.CourierTasks)
	g.GET("/next_task/:uuid", r.NextTaskCourier)
	g.PUT("/task/finish/:uuid", r.FinishTask)
	g.GET("/tasks/order/:uuid", r.GetOrderTasks)
	g.PUT("/tasks/change_queue", r.ChangeQueue)

	//user
	g.POST("/user", r.CreateUser)
	g.GET("/user/:login", r.GetUserByLogin)
	g.PUT("/user/set_state", r.SetUserState)
	g.PUT("/user/delete/:uuid", r.MarkUserAsDeleted)

	//events
	g.GET("/events/courier/:uuid", r.GetCourierEvents)
	g.GET("/events/order/:uuid", r.GetOrderEvents)

	//payments
	g.GET("/payment/types", r.GetPaymentTypes)
	g.GET("/payment/states", r.GetPaymentStatuses)
	g.GET("/task/types", r.GetTaskTypes)
	g.GET("/task/states", r.GetTaskStates)
	g.GET("/courier/states", r.GetCourierStates)
	g.GET("/courier/types", r.GetCourierTypes)
	g.GET("/order/states", r.GetOrderStates)
	g.GET("/order/types", r.GetOrderTypes)
}
