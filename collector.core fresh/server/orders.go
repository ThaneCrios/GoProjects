package server

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/structures"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/proto"
)

func (r *Rest) GetOrderByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event":      "getting order by uuid",
		"order-uuid": uuid,
	})

	order, err := r.Handler.GetOrderByUUID(c.Request().Context(), uuid)
	if err != nil {
		log.WithError(err).Error("failed to get order by uuid")
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) GetDuplicateOrderByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event":      "get duplicate order by uuid",
		"order-uuid": uuid,
	})

	order, err := r.Handler.GetDuplicateOrderByUUID(c.Request().Context(), uuid)
	if err != nil {
		log.WithError(err).Error("failed to get duplicate order by uuid")
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) GetProductsFromOrderByFilter(c echo.Context) error {
	param := c.QueryParam("param")
	uuid := c.QueryParam("uuid")

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event":      "getting products from order by filter",
		"filter":     param,
		"order-uuid": uuid,
	})

	order, err := r.Handler.GetProductsFromOrderByFilter(c.Request().Context(), uuid, param)
	if err != nil {
		log.WithError(err).Error("failed to get products from order by filter")
		res := logs.OutputRestError("failed to get products from order by filter", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) CancelOrder(c echo.Context) error {
	uuid := new(proto.OrderCancel)

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event": "cancel order",
	})

	err := c.Bind(uuid)
	if err != nil {
		log.WithError(err).Error("bind error on cancel order")
		res := logs.OutputRestError("bind error on cancel order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.CancelOrder(c.Request().Context(), *uuid)
	if err != nil {
		log.WithError(err).Error("failed to finish collecting order")
		res := logs.OutputRestError("failed to finish collect order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusCreated, structures.ResponseStructOk)
}

func (r *Rest) FinishCollectOrder(c echo.Context) error {
	id := new(proto.OrderCollector)

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event": "finish order collecting",
	})

	err := c.Bind(id)
	if err != nil {
		log.WithError(err).Error("bind error on finish order")
		res := logs.OutputRestError("bind error on finish order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	order, err := r.Handler.FinishCollectOrder(c.Request().Context(), id.OrderUUID)
	if err != nil {
		log.WithError(err).Error("failed to finish collecting order")
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
	collectorUUID := c.Param("uuid")

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event":          "getting order by uuid",
		"collector-uuid": collectorUUID,
	})

	orders, err := r.Handler.GetFreeOrders(c.Request().Context(), collectorUUID)
	if err != nil {
		log.WithError(err).Error("failed to get free orders")
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, orders)
}

func (r *Rest) GetMyOrders(c echo.Context) error {
	collectorUUID := c.Param("uuid")

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event":          "getting order by uuid",
		"collector-uuid": collectorUUID,
	})

	orders, err := r.Handler.GetMyOrders(c.Request().Context(), collectorUUID)
	if err != nil {
		log.WithError(err).Error("failed to get collector orders")
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, orders)
}

func (r *Rest) GrabOrder(c echo.Context) error {
	ids := new(proto.OrderCollector)

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event": "grab order",
	})

	err := c.Bind(ids)
	if err != nil {
		log.WithError(err).Error("bind error on grab order")
		res := logs.OutputRestError("bind error on setting collector to order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	orderOut, err := r.Handler.SetCollectorToOrder(c.Request().Context(), ids)
	if err != nil {
		log.WithError(err).Error("failed to set collector to order")
		res := logs.OutputRestError("fail to grab order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, orderOut)
}

func (r *Rest) AddProductToOrder(c echo.Context) error {
	params := new(proto.ProductAdd)

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event": "add product to order",
	})

	err := c.Bind(params)
	if err != nil {
		log.WithError(err).Error("bind error on add product")
		res := logs.OutputRestError("bind error on add product to order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	order, err := r.Handler.AddProductToOrder(c.Request().Context(), *params)
	if err != nil {
		log.WithError(err).Error("failed to add product to order")
		res := logs.OutputRestError("fail to product", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) MarkProduct(c echo.Context) error {
	params := new(proto.OrderProductMark)

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event": "grab order",
	})

	err := c.Bind(params)
	if err != nil {
		log.WithError(err).Error("bind error on marking product")
		res := logs.OutputRestError("bind error on marking product as collected", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.MarkProduct(c.Request().Context(), *params)
	if err != nil {
		log.WithError(err).Error("failed to mark product on order")
		res := logs.OutputRestError("fail to mark product as collected", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, structures.ResponseStructOk)
}

func (r *Rest) RemoveProductFromOrder(c echo.Context) error {
	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event": "remove product from order",
	})

	params := new(proto.OrderProductRemove)
	err := c.Bind(params)
	if err != nil {
		log.WithError(err).Error("bind error on removing product")
		res := logs.OutputRestError("bind error on removing product from order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	order, err := r.Handler.RemoveProductFromOrder(c.Request().Context(), *params)
	if err != nil {
		log.WithError(err).Error("failed to remove product from order")
		res := logs.OutputRestError("fail remove product from order", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) ChangeProduct(c echo.Context) error {
	params := new(proto.OrderProductChange)

	log := logs.LoggerForContext(c.Request().Context()).WithFields(logrus.Fields{
		"event":      "change product in order",
		"order-uuid": params.OrderUUID,
	})

	err := c.Bind(params)
	if err != nil {
		log.WithError(err).Error("bind error on changing product")
		res := logs.OutputRestError("bind error on changing product", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = r.Handler.ChangeProduct(c.Request().Context(), *params)
	if err != nil {
		log.WithError(err).Error("failed to change product on order")
		res := logs.OutputRestError("fail to change product", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, structures.ResponseStructOk)
}
