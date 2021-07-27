package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
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
	barCode := c.Param("barCode")
	response, err := r.Handler.GetProductByBarCode(c.Request().Context(), barCode)
	if err != nil {
		res := logs.OutputRestError("failed to create product", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, response)
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

//func (r *Rest) CompareTwoProducts(c echo.Context) error {
//	barCodes := models.BarCodes{
//		BarCodeOrder:   c.QueryParam("barcode_order"),
//		BarCodeScanned: c.QueryParam("barcode_scanned"),
//	}
//	response, err := r.Handler.CompareTwoProducts(c.Request().Context(), barCodes)
//	if err != nil {
//		res := logs.OutputRestError("failed to compares product, try again", err)
//		return c.JSON(http.StatusBadRequest, res)
//	}
//	return c.JSON(http.StatusOK, response)
//}
