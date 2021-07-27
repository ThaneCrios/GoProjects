package handler

import (
	"gitlab.com/faemproject/backend/delivery/delivery.front/proto"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/core/shared/tools"
	"gitlab.com/faemproject/backend/delivery/delivery.front/models"
)

//TODO 1) Все входящие структуры лежат отдельно в папке proto
//TODO 2) Структуры отправляющиеся на коре сервисы изолированны от пользовательских структур
//TODO 3) Вызовы RPC обрамлены в функции, которые принимает А и отдает Б
func (h *Handler) GetCourierByUUID(c echo.Context) error {
	uuid := c.Param("uuid")

	courier := new(models.Courier)
	err := tools.RPC("GET", h.Config.Services.Delivery+"/couriers/"+uuid, nil, nil, courier)
	if err != nil {
		logs.Eloger.WithError(err).Error("getting courier by uuid")
		res := logs.OutputRestError("fail to get courier by uuid", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, courier)
}

func (h *Handler) GetFilteredCouriers(c echo.Context) error {
	couriers := new([]models.Courier)
	params := make(map[string]string)
	params["courier_state"] = c.QueryParam("courier_state")

	err := tools.RPC("GET", h.Config.Services.Delivery+"/couriers/filter", params, nil, couriers)
	if err != nil {
		logs.Eloger.WithError(err).Error("getting filtered couriers")
		res := logs.OutputRestError("fail to get filtered couriers", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, couriers)
}

func (h *Handler) CreateCourier(c echo.Context) error {
	courier := new(proto.CourierFront)
	err := c.Bind(courier)
	if err != nil {
		logs.Eloger.WithError(err).Error("courier bind error")
		res := logs.OutputRestError("fail to bind city", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = tools.RPC("POST", h.Config.Services.Delivery+"/courier", nil, courier, courier)
	if err != nil {
		logs.Eloger.WithError(err).Error("create courier")
		res := logs.OutputRestError("fail to create courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	logs.Eloger.Infof("courier with uuid=%s created", courier.UUID)

	return c.JSON(http.StatusOK, courier)
}

func (h *Handler) UpdateCourier(c echo.Context) error {
	uuid := c.Param("uuid")

	courier := new(models.Courier)
	err := c.Bind(courier)
	if err != nil {
		logs.Eloger.WithError(err).Error("courier bind error on update")
		res := logs.OutputRestError("fail to bind courier on update", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = tools.RPC("PUT", h.Config.Services.Delivery+"/courier/"+uuid, nil, courier, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("updating courier")
		res := logs.OutputRestError("fail to update courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	logs.Eloger.Infof("courier with uuid=%s updated", courier.UUID)

	return c.JSON(http.StatusOK, courier)
}

func (h *Handler) MarkCourierAsDeleted(c echo.Context) error {
	uuid := c.Param("uuid")

	err := tools.RPC("DELETE", h.Config.Services.Delivery+"/courier/delete"+uuid, nil, nil, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("deleting courier")
		res := logs.OutputRestError("fail to delete courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("deleted"))
}

func (h *Handler) UpdateStatusCourier(c echo.Context) error {
	meta := new(models.Courier)
	err := c.Bind(meta)

	if err != nil {
		res := logs.OutputRestError("bind error", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = tools.RPC("PUT", h.Config.Services.Delivery+"/courier/set_status", nil, meta, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("updating courier status")
		res := logs.OutputRestError("fail to update courier status", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (h *Handler) DeleteCourier(c echo.Context) error {
	uuid := c.Param("uuid")

	err := tools.RPC("PUT", h.Config.Services.Delivery+"/courier/delete/"+uuid, nil, nil, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("deleting courier")
		res := logs.OutputRestError("fail to delete courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, logs.OutputRestOK("deleted"))
}

func (h *Handler) CourierStates(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.States.Courier)
}

func (h *Handler) UpdateCoordinate(c echo.Context) error {
	coordinate := new(proto.CouriersFilter)
	err := c.Bind(coordinate)
	if err != nil {
		logs.Eloger.WithError(err).Error("courier bind error on update")
		res := logs.OutputRestError("fail to bind courier on update", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = tools.RPC("PUT", h.Config.Services.Delivery+"/courier/update_coor", nil, coordinate, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("deleting courier")
		res := logs.OutputRestError("fail to delete courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (h *Handler) GetCourierTokenInfoByPhoneNumber(c echo.Context) error {
	phoneNumber := c.Param("phoneNumber")
	kek, lol := c.ParamNames(), c.ParamValues()
	_, _ = kek, lol

	tokenInfo := new(models.AuthTokenDate)
	err := tools.RPC("GET", h.Config.Services.Delivery+"/courier/"+phoneNumber, nil, nil, tokenInfo)
	if err != nil {
		logs.Eloger.WithError(err).Error("getting courier token info")
		res := logs.OutputRestError("fail to getting courier token info", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, tokenInfo)
}
