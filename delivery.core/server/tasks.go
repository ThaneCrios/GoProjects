package server

import (
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

func (r *Rest) GetOrderTasks(c echo.Context) error {
	uuid := c.Param("uuid")

	tasks, err := r.Handler.GetOrderTasks(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	tasksFront := convertTasksArrayElements(tasks)
	return c.JSON(http.StatusOK, tasksFront)
}

//TasksListByStatus возвращает список тасков по определённому фильтеру
func (r *Rest) TasksListByStatus(c echo.Context) error {
	filter := new(proto.TasksFilter)
	err := c.Bind(filter)
	if err != nil {
		res := logs.OutputRestError("bind error on listing tasks by status", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	tasks, err := r.Handler.TasksListByStatus(c.Request().Context(), filter)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	taskFront := convertTasksArrayElements(tasks)
	return c.JSON(http.StatusOK, taskFront)
}

//GetCourierTasksByStatus возвращает список тасков курьера по определённому статусу
func (r *Rest) GetCourierTasksByStatus(c echo.Context) error {
	filter := new(proto.TasksFilter)
	err := c.Bind(filter)
	if err != nil {
		res := logs.OutputRestError("bind error in getting courier tasks by status", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	tasks, err := r.Handler.GetCourierTasksByStatus(c.Request().Context(), filter)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	taskFront := convertTasksArrayElements(tasks)
	return c.JSON(http.StatusOK, taskFront)

}

func (r *Rest) UpdateTask(c echo.Context) error {
	filter := new(proto.TasksFilter)
	err := c.Bind(filter)
	if err != nil {
		res := logs.OutputRestError("bind error on updating task", err)
		return c.JSON(http.StatusBadRequest, res)
	}

	task, err := r.Handler.UpdateTask(c.Request().Context(), filter)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	taskFront := proto.TaskFromCoreToFront(task)
	return c.JSON(http.StatusOK, taskFront)
}

func (r *Rest) NextTaskCourier(c echo.Context) error {
	uuid := c.Param("uuid")
	task, err := r.Handler.NextTaskCourier(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	taskFront := proto.TaskFromCoreToFront(task)
	return c.JSON(http.StatusOK, taskFront)
}

func (r *Rest) FinishTask(c echo.Context) error {
	courierUUID := c.Param("uuid")
	err := r.Handler.FinishTask(c.Request().Context(), courierUUID)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, logs.OutputRestOK("finished"))
}

func (r *Rest) ChangeQueue(c echo.Context) error {
	IDs := new(proto.QueueTasks)
	err := c.Bind(IDs)
	if err != nil {
		res := logs.OutputRestError("bind error on changing queue", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	err = r.Handler.ChangeQueue(c.Request().Context(), IDs)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	return c.JSON(http.StatusOK, logs.OutputRestOK("updated"))
}

func (r *Rest) CourierTasks(c echo.Context) error {
	uuid := c.Param("uuid")
	tasks, err := r.Handler.CourierTasks(c.Request().Context(), uuid)
	if err != nil {
		res := logs.OutputRestError("", err)
		return c.JSON(http.StatusBadRequest, res)
	}
	taskFront := convertTasksArrayElements(tasks)
	return c.JSON(http.StatusOK, taskFront)
}

func (r *Rest) GetTaskTypes(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.TaskTypeArray)
}

func (r *Rest) GetTaskStates(c echo.Context) error {
	return c.JSON(http.StatusOK, proto.TaskStateArray)
}

func convertTasksArrayElements(tasks []models.Task) (tasksFront []proto.TaskFront) {
	for _, v := range tasks {
		tasksFront = append(tasksFront, proto.TaskFromCoreToFront(v))
	}
	return tasksFront
}
