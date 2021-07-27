package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

//GetCourierByUUID ...
func (r *Rest) GetCourierByUUID(c echo.Context) error {
	uuid := c.Param("uuid")
	courier, err := r.Handler.GetCourierByUUID(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	cor := proto.CourierFromCoreToFront(*courier)
	return c.JSON(http.StatusOK, cor)
}

//func (r *Rest) CreateCourier(c echo.Context) error {
//	var event models.Event
//	courier := new(models.Courier)
//	err := c.Bind(courier)
//	if err != nil {
//		res := logs.OutputRestError("bind error", err)
//		return c.JSON(http.StatusBadRequest, res)
//	}
//
//	err, courier = r.Handler.CreateCourier(c.Request().Context(), courier)
//	if err != nil {
//		res := logs.OutputRestError("", err)
//		event = r.Handler.EventCreating("create courier", err.Error(), "", "")
//		err = r.Handler.CreateEvent(c.Request().Context(), &event)
//		return c.JSON(http.StatusBadRequest, res)
//	}
//	event = r.Handler.EventCreating("create courier", "done", "", "")
//	err = r.Handler.CreateEvent(c.Request().Context(), &event)
//
//	err = r.Handler.CreateQueue(c.Request().Context(), courier.UUID)
//	if err != nil {
//		res := logs.OutputRestError("", err)
//		event = r.Handler.EventCreating("create queue for courier", err.Error(), courier.UUID, "")
//		err = r.Handler.CreateEvent(c.Request().Context(), &event)
//		return c.JSON(http.StatusBadRequest, res)
//	}
//	event = r.Handler.EventCreating("create queue for courier", "done", courier.UUID, "")
//	err = r.Handler.CreateEvent(c.Request().Context(), &event)
//
//	return c.JSON(http.StatusOK, courier)
//}

//CreateCourier ...
func (r *Rest) CreateCourier(c echo.Context) error {
	courierFront := new(proto.CourierFront)
	err := c.Bind(courierFront)
	if err != nil {
		res := logs.OutputRestError("bind error on create courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	courier, err := r.Handler.CreateCourier(c.Request().Context(), *courierFront)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = r.Handler.CreateQueue(c.Request().Context(), courier.UUID)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	cor := proto.CourierFromCoreToFront(courier)
	return c.JSON(http.StatusOK, cor)
}

//UpdateCourier ...
func (r *Rest) UpdateCourier(c echo.Context) error {
	courier := new(models.Courier)
	err := c.Bind(courier)
	if err != nil {
		res := logs.OutputRestError("bind error on update courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	courier, err = r.Handler.UpdateCourier(c.Request().Context(), courier)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, courier)
}

//CouriersFilter ...
func (r *Rest) CouriersFilter(c echo.Context) error {
	filter := new(proto.CouriersFilter)
	err := c.Bind(filter)
	if err != nil {
		res := logs.OutputRestError("bind error on couriers filtering", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	courier, err := r.Handler.CouriersFilter(c.Request().Context(), filter)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	couriersFront := convertCourierArrayElements(courier)
	return c.JSON(http.StatusOK, couriersFront)
}

//DeleteCourier ...
func (r *Rest) DeleteCourier(c echo.Context) error {
	uuid := c.Param("uuid")
	err := r.Handler.DeleteCourier(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, logs.OutputRestOK("deleted"))
}

//UpdateStatusCourier ...
func (r *Rest) UpdateStatusCourier(c echo.Context) error {
	courier := new(proto.CourierFront)
	err := c.Bind(courier)
	if err != nil {
		res := logs.OutputRestError("bind error on courier status updating", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	courierBack := proto.CourierFromFrontToCore(*courier)
	_, err = r.Handler.UpdateCourier(c.Request().Context(), &courierBack)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

//GetCourierByChatID ...
func (r *Rest) GetCourierByChatID(c echo.Context) error {
	chatID := c.FormValue("chatid")
	phoneNumber := c.FormValue("phone_number")
	err := r.Handler.GetCourierByChatID(c.Request().Context(), chatID, phoneNumber)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, "OK")
}

//GetCourierStates ...
func (r *Rest) GetCourierStates(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.CourierStateArray)
}

func (r *Rest) GetCourierTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.CourierTypeArray)
}

//UpdateCourierCoordinates ...
func (r *Rest) UpdateCourierCoordinates(c echo.Context) error {
	coor := new(proto.CouriersFilter)
	err := c.Bind(coor)
	if err != nil {
		res := logs.OutputRestError("bind error on update couriers coords", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = r.Handler.UpdateCourierCoordinates(c.Request().Context(), coor)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, "Coordinates updated")
}

//GetCourierTokenInfoByPhoneNumber ...
func (r *Rest) GetCourierTokenInfoByPhoneNumber(c echo.Context) error {
	phoneNumber := c.Param("phone_number")

	tokenInfo, err := r.Handler.GetCourierTokenInfoByPhoneNumber(c.Request().Context(), phoneNumber)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, tokenInfo)

}

//GetCourierUUIDByChatID ...
func (r *Rest) GetCourierUUIDByChatID(c echo.Context) error {
	chatID := c.Param("chat_id")
	courierUUID, err := r.Handler.GetCourierUUIDByChatID(c.Request().Context(), chatID)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, courierUUID)
}

//FillCoords  ...
func FillCoords(filter *proto.CouriersFilter) (coord models.CourierCoordinatesTable) {

	coord.CourierUUID = filter.CourierUUID
	coord.Lat = filter.CourierLat
	coord.Lon = filter.CourierLon

	return coord
}

func convertCourierArrayElements(couriers []models.Courier) (couriersFront []proto.CourierFront) {
	for _, v := range couriers {
		couriersFront = append(couriersFront, proto.CourierFromCoreToFront(v))
	}
	return couriersFront
}
