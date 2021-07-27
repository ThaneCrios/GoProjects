package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
	"net/http"
)

func (r *Rest) EventsFilter(c echo.Context) error {
	filter := new(proto.EventFilter)
	events, err := r.Handler.EventsFilter(c.Request().Context(), filter)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	eventsFront := convertEventArrayElements(events)
	return c.JSON(http.StatusOK, eventsFront)
}

func (r *Rest) GetCourierEvents(c echo.Context) error {
	courierUUID := c.Param("uuid")
	fmt.Println(courierUUID)
	events, err := r.Handler.GetCourierEvents(c.Request().Context(), courierUUID)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	eventsFront := convertEventArrayElements(events)
	return c.JSON(http.StatusOK, eventsFront)
}

func (r *Rest) GetOrderEvents(c echo.Context) error {
	orderUUID := c.Param("uuid")

	events, err := r.Handler.GetOrderEvents(c.Request().Context(), orderUUID)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	eventsFront := convertEventArrayElements(events)
	return c.JSON(http.StatusOK, eventsFront)
}

func convertEventArrayElements(events []models.Event) (eventsFront []proto.EventFront) {
	for _, v := range events {
		eventsFront = append(eventsFront, proto.EventFromCoreToFront(v))
	}
	return eventsFront
}
