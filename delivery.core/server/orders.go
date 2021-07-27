package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

func (r *Rest) GetOrderByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	order, err := r.Handler.GetOrderByUUID(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	orderFront := proto.OrderFromCoreToFront(order)
	return c.JSON(http.StatusOK, orderFront)
}

func (r *Rest) OrdersFilter(c echo.Context) error {
	filter := new(proto.OrdersFilter)

	err := c.Bind(filter)
	var event models.Event
	order, err := r.Handler.OrdersFilter(c.Request().Context(), filter)
	if err != nil {
		res := logs.OutputRestError("", err)
		event = r.Handler.EventCreating("order filtering", err.Error(), "", "")
		err = r.Handler.CreateEvent(c.Request().Context(), &event)
		return c.JSON(http.StatusBadRequest, res)
	}

	event = r.Handler.EventCreating("order filtering", "done", "", "")
	err = r.Handler.CreateEvent(c.Request().Context(), &event)
	ordersFront := convertOrderArrayElements(order)
	return c.JSON(http.StatusOK, ordersFront)
}

/*func (r *Rest) CreateOrder(c echo.Context) error {
	var eventTask, eventOrder models.Event
	order := new(models.Order)
	err := c.Bind(order)
	if err != nil {
		res := logs.OutputRestError("bind error", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err, order = r.Handler.CreateOrder(c.Request().Context(), order)
	if err != nil {
		res := logs.OutputRestError("failed to create order", err)
		eventOrder = r.Handler.EventCreating("order creating", err.Error(), order.CourierUUID, order.UUID)
		err = r.Handler.CreateEvent(c.Request().Context(), &eventOrder)
		return c.JSON(http.StatusBadRequest, res)
	}
	eventTask = r.Handler.EventCreating("order creating", "done", order.CourierUUID, order.UUID)
	err = r.Handler.CreateEvent(c.Request().Context(), &eventOrder)

	err = r.Handler.CreateTasks(c.Request().Context(), order)
	if err != nil {
		res := logs.OutputRestError("failed to create event", err)
		eventTask = r.Handler.EventCreating("task creating", err.Error(), order.CourierUUID, order.UUID)
		err = r.Handler.CreateEvent(c.Request().Context(), &eventTask)
		return c.JSON(http.StatusBadRequest, res)
	}

	eventTask = r.Handler.EventCreating("task creating", "done", order.CourierUUID, order.UUID)
	err = r.Handler.CreateEvent(c.Request().Context(), &eventTask)

	return c.JSON(http.StatusCreated, order)
} */

func (r *Rest) CreateOrder(c echo.Context) error {
	order := new(proto.OrderFront)
	err := c.Bind(order)
	if err != nil {
		res := logs.OutputRestError("bind error on create order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	orderBack := proto.OrderFromFrontToCore(*order)
	orderBack, err = r.Handler.CreateOrder(c.Request().Context(), &orderBack)
	if err != nil {
		res := logs.OutputRestError("failed to create order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	*order = proto.OrderFromCoreToFront(orderBack)
	return c.JSON(http.StatusCreated, order)
}

func (r *Rest) UpdateOrder(c echo.Context) error {
	order := new(proto.OrderFront)
	err := c.Bind(order)
	if err != nil {
		res := logs.OutputRestError("bind error on update order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	orderBack := proto.OrderFromFrontToCore(*order)
	err = r.Handler.UpdateOrder(c.Request().Context(), &orderBack)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (r *Rest) DeleteOrder(c echo.Context) error {
	uuid := c.Param("uuid")
	err := r.Handler.DeleteOrder(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("deleted"))
}

func (r *Rest) UpdateStatusOrder(c echo.Context) error {
	var event models.Event
	order := new(proto.OrderFront)
	err := c.Bind(order)
	if err != nil {
		res := logs.OutputRestError("bind error on order status update", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	orderBack := proto.OrderFromFrontToCore(*order)
	err = r.Handler.UpdateOrder(c.Request().Context(), &orderBack)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	event = r.Handler.EventCreating("updating status order", "done", "", order.UUID)
	err = r.Handler.CreateEvent(c.Request().Context(), &event)
	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (r *Rest) GetOrderStates(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.OrderStateArray)
}

func (r *Rest) GetOrderTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.OrderTypeArray)
}

func (r *Rest) SetCourierToOrder(c echo.Context) error {
	ids := new(proto.OrderCourier)
	err := c.Bind(ids)
	if err != nil {
		res := logs.OutputRestError("bind error on setting courier to order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.SetCourierToOrder(c.Request().Context(), ids)
	if err != nil {
		res := logs.OutputRestError("courier set order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (r *Rest) RemoveCourierFromOrder(c echo.Context) error {
	ids := new(proto.OrderCourier)
	err := c.Bind(ids)
	if err != nil {
		res := logs.OutputRestError("bind error on removing courier from order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.RemoveCourierFromOrder(c.Request().Context(), ids)
	if err != nil {
		res := logs.OutputRestError("cant remove courier from order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, "removed")
}

func convertOrderArrayElements(orders []models.Order) (ordersFront []proto.OrderFront) {
	for _, v := range orders {
		ordersFront = append(ordersFront, proto.OrderFromCoreToFront(v))
	}
	return ordersFront
}
