package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/faemproject/backend/core/shared/logs"
)

const (
	errBind = "failed to parse the request"
)

func (r *Rest) Version(c echo.Context) error {
	out, err := r.Handler.GetCurrentVersion(c.Request().Context())
	if err != nil {
		logs.LoggerForContext(c.Request().Context()).
			Error(err) // you may add additional fields here
		res := logs.OutputRestError("can't get version", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, out)
}
