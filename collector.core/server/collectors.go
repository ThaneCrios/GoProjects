package server

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
	"net/http"
)

func (r *Rest) GetCollectorByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	collector, err := r.Handler.GetCollectorByUUID(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, collector)
}

func (r *Rest) CreateCollector(c echo.Context) error {
	collector := new(models.Collector)
	err := c.Bind(collector)
	if err != nil {
		res := logs.OutputRestError("bind error on create collector", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	*collector, err = r.Handler.CreateCollector(c.Request().Context(), *collector)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, collector)
}
