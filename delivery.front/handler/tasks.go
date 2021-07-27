package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/core/shared/tools"
	"gitlab.com/faemproject/backend/delivery/delivery.front/models"
	"gitlab.com/faemproject/backend/delivery/delivery.front/proto"
)

func (h *Handler) GetCourierTasks(c echo.Context) error {
	uuid := c.Param("uuid")

	tasks := new([]models.Task)
	err := tools.RPC("GET", h.Config.Services.Delivery+"/tasks/courier_tasks/"+uuid, nil, nil, tasks)
	if err != nil {
		logs.Eloger.WithError(err).Error("getting courier tasks by courier uuid")
		res := logs.OutputRestError("fail to get courier tasks by courier uuid", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *Handler) ChangeQueue(c echo.Context) error {
	IDs := new(proto.QueueTasks)
	err := tools.RPC("PUT", h.Config.Services.Delivery+"/tasks/change_queue", nil, IDs, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("change courier`s queue")
		res := logs.OutputRestError("fail to change courier`s queue", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, logs.OutputRestOK("deleted"))
}

func (h *Handler) GetOrderTasks(c echo.Context) error {
	uuid := c.Param("uuid")

	tasks := new([]models.Task)
	err := tools.RPC("GET", h.Config.Services.Delivery+"/tasks/order/"+uuid, nil, nil, tasks)
	if err != nil {
		logs.Eloger.WithError(err).Error("getting courier tasks by courier uuid")
		res := logs.OutputRestError("fail to get courier tasks by courier uuid", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *Handler) UpdateTask(c echo.Context) error {
	IDs := new(proto.TasksFilter)
	task := new(models.Task)
	err := c.Bind(IDs)
	if err != nil {
		logs.Eloger.WithError(err).Error("task bind error on update")
		res := logs.OutputRestError("fail to bind task on update", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	err = tools.RPC("PUT", h.Config.Services.Delivery+"/task/update", nil, IDs, task)
	if err != nil {
		logs.Eloger.WithError(err).Error("updating courier")
		res := logs.OutputRestError("fail to update courier", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	logs.Eloger.Infof("courier with uuid=%s updated", IDs.TaskUUID)

	return c.JSON(http.StatusOK, task)
}

func (h *Handler) FinishTask(c echo.Context) error {
	courierUUID := c.Param("uuid")

	err := tools.RPC("GET", h.Config.Services.Delivery+"/tasks/courier_tasks/"+courierUUID, nil, nil, nil)
	if err != nil {
		logs.Eloger.WithError(err).Error("finishing courier task by courier uuid")
		res := logs.OutputRestError("fail to finished courier task by courier uuid", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, "finished")
}

func (h *Handler) NextTaskCourier(c echo.Context) error {
	courierUUID := c.Param("uuid")

	task := new(models.Task)
	err := tools.RPC("GET", h.Config.Services.Delivery+"/next_task/"+courierUUID, nil, nil, task)
	if err != nil {
		logs.Eloger.WithError(err).Error("getting courier task by courier uuid")
		res := logs.OutputRestError("fail to getting courier task by courier uuid", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	return c.JSON(http.StatusOK, task)
}
