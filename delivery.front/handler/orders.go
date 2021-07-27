package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/core/shared/tools"
	"gitlab.com/faemproject/backend/delivery/delivery.front/models"
	"gitlab.com/faemproject/backend/delivery/delivery.front/proto"

	"net/http"
)

func (h *Handler) GetOrderByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	order := new(models.Order)

	err := tools.RPC("GET", h.Config.Services.Delivery+"/orders/"+uuid, nil, nil, order)
	if err != nil {
		logs.Eloger.WithError(err).WithField("uuid", uuid).Error("getting order by uuid")
		res := logs.OutputRestError("fail to get order by uuid", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (h *Handler) OrdersFilter(c echo.Context) error {
	var err error

	orders := new([]models.Order)

	filters := make(map[string]string)
	filters["courier_uuid"] = c.Param("courier_uuid")
	filters["order_lat"] = c.QueryParam("order_lat")
	filters["order_lon"] = c.QueryParam("order_lon")
	filters["order_status"] = c.QueryParam("order_status")

	err = tools.RPC("GET", h.Config.Services.Delivery+"/orders/filter", filters, nil, orders)
	if err != nil {
		logs.Eloger.WithError(err).WithField("filter", filters).Error("getting filtered orders")
		res := logs.OutputRestError("fail to get filtered orders", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) CreateOrder(c echo.Context) error {
	orderFront := new(proto.OrderCreate)
	orderResponse := new(models.Order)
	err := c.Bind(orderFront)
	if err != nil {
		logs.Eloger.WithError(err).Error("order binding error on create")
		res := logs.OutputRestError("order bind error on create", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = tools.RPC("POST", h.Config.Services.Delivery+"/orders", nil, orderFront, orderResponse)
	if err != nil {
		logs.Eloger.WithError(err).Error("deleting cart")
		res := logs.OutputRestError("fail to delete cart", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	logs.Eloger.Infof("order with uuid=%s created", orderResponse.UUID)

	return c.JSON(http.StatusOK, orderResponse)
}

func (h *Handler) UpdateOrder(c echo.Context) error {
	uuid := c.Param("uuid")

	order := new(models.Order)
	err := c.Bind(order)
	if err != nil {
		logs.Eloger.WithError(err).Error("order bind err on update")
		res := logs.OutputRestError("order bind error on update", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = tools.RPC("PUT", h.Config.Services.Delivery+"/orders/", nil, order, nil)
	if err != nil {
		logs.Eloger.WithError(err).
			WithField("uuid", uuid).
			Error("updating order by uuid")
		res := logs.OutputRestError("fail to update order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	logs.Eloger.Infof("order with uuid=%s updated", uuid)

	return c.JSON(http.StatusOK, order)
}

func (h *Handler) MarkOrderAsDeleted(c echo.Context) error {
	uuid := c.Param("uuid")

	err := tools.RPC("DELETE", h.Config.Services.Delivery+"/order/delete"+uuid, nil, nil, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("deleting state")
		res := logs.OutputRestError("fail to delete state", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	logs.Eloger.Infof("order with uuid=%s deleted", uuid)

	return c.JSON(http.StatusOK, logs.OutputRestOK("deleted"))
}

func (h *Handler) SetCourierToOrder(c echo.Context) error {
	ids := new(proto.OrderCourier)
	err := c.Bind(ids)
	fmt.Println(ids)
	if err != nil {
		res := logs.OutputRestError("bind error", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = tools.RPC("PUT", h.Config.Services.Delivery+"/order/set_courier", nil, ids, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("set courier to order")
		res := logs.OutputRestError("fail set courier to order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	logs.Eloger.Infof("courier with uuid=%s was set to order", ids.CourierUUID)

	return c.JSON(http.StatusOK, logs.OutputRestOK("courier was set"))
}

func (h *Handler) UpdateStatusOrder(c echo.Context) error {
	meta := new(models.Order)
	err := c.Bind(meta)
	fmt.Println(meta)
	if err != nil {
		res := logs.OutputRestError("bind error", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = tools.RPC("PUT", h.Config.Services.Delivery+"/orders", nil, meta, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("updating order status")
		res := logs.OutputRestError("fail to update order status", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (h *Handler) OrderStates(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.States.Order)
}

func (h *Handler) OrderRemoveCourier(c echo.Context) error {
	IDs := new(proto.OrderCourier)
	err := c.Bind(IDs)
	if err != nil {
		res := logs.OutputRestError("bind error", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = tools.RPC("PUT", h.Config.Services.Delivery+"/order/remove_courier", nil, IDs, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("updating order status")
		res := logs.OutputRestError("fail to update order status", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (h *Handler) OrderDeleted(c echo.Context) error {
	uuid := c.Param("uuid")
	err := tools.RPC("PUT", h.Config.Services.Delivery+"/order/delete/"+uuid, nil, nil, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("updating order status")
		res := logs.OutputRestError("fail to update order status", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("order deleted"))
}
