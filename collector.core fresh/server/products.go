package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	"net/http"
)

func (r *Rest) CreateProduct(c echo.Context) error {
	product := new(models.Product)
	err := c.Bind(product)
	if err != nil {
		res := logs.OutputRestError("bind error on create product", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	*product, err = r.Handler.CreateProduct(c.Request().Context(), *product)
	if err != nil {
		res := logs.OutputRestError("failed to create product", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusCreated, product)
}

func (r *Rest) GetProductByBarCode(c echo.Context) error {
	barCode := c.Param("barcode")
	product, err := r.Handler.GetProductByBarCode(c.Request().Context(), barCode)
	if err != nil {
		res := logs.OutputRestError("failed to create product", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, product)
}

func (r *Rest) GetProductByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	order, err := r.Handler.GetProductByUUID(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, order)
}

func (r *Rest) AppointBarCode(c echo.Context) error {
	appointParams := new(models.ProductWithBarCode)
	err := c.Bind(appointParams)
	if err != nil {
		res := logs.OutputRestError("bind error on appoint bar code", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	prodWithBar, err := r.Handler.AppointBarCode(c.Request().Context(), *appointParams)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, prodWithBar)
}

func (r *Rest) GetProductUUIDsWithBarCodes(c echo.Context) error {
	userUUID := c.QueryParam("uuid")
	product, err := r.Handler.GetProductUUIDsWithBarCodes(c.Request().Context(), userUUID)
	if err != nil {
		res := logs.OutputRestError("fail to get product uuids with bar codes", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, product)
}
