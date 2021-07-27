package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

func (r *Rest) GetPaymentTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.PaymentTypeArray)
}

func (r *Rest) GetPaymentStatuses(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.PaymentStatusArray)
}
