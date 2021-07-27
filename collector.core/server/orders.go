package server

import (
	"gitlab.com/faemproject/backend/core/shared/structures"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/collector.core/proto"
)

func (r *Rest) GetOrderByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	order, err := r.Handler.GetOrderByUUID(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) GetProductsFromOrderByFilter(c echo.Context) error {
	param := c.QueryParam("param")
	uuid := c.QueryParam("uuid")
	order, err := r.Handler.GetProductsFromOrderByFilter(c.Request().Context(), uuid, param)
	if err != nil {
		res := logs.OutputRestError("failed to get order products", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) FinishCollectOrder(c echo.Context) error {
	id := new(proto.OrderCollector)
	err := c.Bind(id)
	if err != nil {
		res := logs.OutputRestError("bind error on finish order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	order, err := r.Handler.FinishCollectOrder(c.Request().Context(), id.OrderUUID)
	if err != nil {
		res := logs.OutputRestError("failed to finish collect order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, order)
}

func (r *Rest) CreateOrder(c echo.Context) error {
	order := new(models.Order)
	err := c.Bind(order)
	if err != nil {
		res := logs.OutputRestError("bind error on create order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	*order, err = r.Handler.CreateOrder(c.Request().Context(), order)
	if err != nil {
		res := logs.OutputRestError("failed to create order", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusCreated, order)
}

func (r *Rest) GetFreeOrders(c echo.Context) error {
	orders, err := r.Handler.GetFreeOrders(c.Request().Context())
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, orders)
}

func (r *Rest) GrabOrder(c echo.Context) error {
	ids := new(proto.OrderCollector)
	err := c.Bind(ids)
	if err != nil {
		res := logs.OutputRestError("bind error on setting collector to order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	orderOut, err := r.Handler.SetCollectorToOrder(c.Request().Context(), ids)
	if err != nil {
		res := logs.OutputRestError("fail to grab order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, orderOut)
}

func (r *Rest) MarkProduct(c echo.Context) error {
	params := new(proto.OrderProductMark)

	err := c.Bind(params)
	if err != nil {
		res := logs.OutputRestError("bind error on marking product as collected", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.MarkProduct(c.Request().Context(), *params)
	if err != nil {
		res := logs.OutputRestError("fail to mark product as collected", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, structures.ResponseStructOk)
}

func (r *Rest) ChangeProduct(c echo.Context) error {
	params := new(proto.OrderProductChange)

	err := c.Bind(params)
	if err != nil {
		res := logs.OutputRestError("bind error on changing product", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.ChangeProduct(c.Request().Context(), *params)
	if err != nil {
		res := logs.OutputRestError("fail to change product", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, structures.ResponseStructOk)
}
