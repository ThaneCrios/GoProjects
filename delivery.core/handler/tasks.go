package handler

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

//GetOrderTasks возврашает таски заказа по ID заказа
func (h *Handler) GetOrderTasks(ctx context.Context, uuid string) ([]models.Task, error) {
	tasks, err := h.DB.OrderTasks(ctx, uuid)
	if err != nil {
		logs.Eloger.WithError(err).Error("get task by id")
		return nil, errors.Wrap(err, "fail to get tasks by id")
	}
	return tasks, nil
}

//TasksListByStatus возвращает список тасков по определённому статусу
func (h *Handler) TasksListByStatus(ctx context.Context, filter *proto.TasksFilter) ([]models.Task, error) {
	tasks, err := h.DB.TasksListByStatus(ctx, filter)
	if err != nil {
		logs.Eloger.WithError(err).Error("tasks filtering")
		return nil, errors.Wrap(err, "fail to get filtered tasks")
	}
	return tasks, nil
}

//GetCourierTasksByStatus возвращает список тасков курьера по определённому статусу
func (h *Handler) GetCourierTasksByStatus(ctx context.Context, filter *proto.TasksFilter) ([]models.Task, error) {
	tasks, err := h.DB.GetCourierTasksByStatus(ctx, filter)
	if err != nil {
		logs.Eloger.WithError(err).Error("tasks for courier by status")
		return nil, errors.Wrap(err, "fail to get tasks for courier by status")
	}
	return tasks, nil
}

//UpdateTask меняет таск
func (h *Handler) UpdateTask(ctx context.Context, filter *proto.TasksFilter) (models.Task, error) {
	event := h.EventCreating("update tasks", "done", filter.CourierUUID, "")
	tasks, err := h.DB.UpdateTask(ctx, filter, event)
	if err != nil {
		logs.Eloger.WithError(err).Error("change task")
		return models.Task{}, errors.Wrap(err, "fail to change task")
	}
	return tasks, nil
}

//CreateTasks создает таски
func (h *Handler) CreateTasks(ctx context.Context, order *models.Order) []models.Task {
	var (
		firstTask  models.Task
		secondTask models.Task
	)
	switch order.Service {
	case proto.Variables.Order.OrderType.PickUpAndDeliver:
		firstTask = h.tasksCreating(order, proto.Variables.Tasks.TaskType.PickUp, order.PickupPersonContacts, order.PickupRoute)
		secondTask = h.tasksCreating(order, proto.Variables.Tasks.TaskType.Deliver, order.DropOffPersonContacts, order.DropoffRoute)

	case proto.Variables.Order.OrderType.PickUpPayAndDeliver:
		firstTask = h.tasksCreating(order, proto.Variables.Tasks.TaskType.PickUpAndPay, order.PickupPersonContacts, order.PickupRoute)
		secondTask = h.tasksCreating(order, proto.Variables.Tasks.TaskType.Deliver, order.DropOffPersonContacts, order.DropoffRoute)

	case proto.Variables.Order.OrderType.PickUpPayAndDeliverWaitDeliverBack:
		firstTask = h.tasksCreating(order, proto.Variables.Tasks.TaskType.PickUpAndPay, order.PickupPersonContacts, order.PickupRoute)
		secondTask = h.tasksCreating(order, proto.Variables.Tasks.TaskType.DeliverWaitDeliverBack, order.DropOffPersonContacts, order.DropoffRoute)
	}

	return []models.Task{firstTask, secondTask}
}

//NextTaskCourier возвращает следующий такс курьера по его ID
func (h *Handler) NextTaskCourier(ctx context.Context, courierUUID string) (models.Task, error) {
	queue, err := h.DB.CourierTasks(ctx, courierUUID)
	if err != nil {
		logs.Eloger.WithError(err).Error("next task")
		return models.Task{}, errors.Wrap(err, "fail to find next task")
	}

	if len(queue.Tasks) != 0 {
		return queue.Tasks[0], nil
	}
	return models.Task{}, errors.New("no tasks yet")
}

//FinishTask удаляет первый таск из очереди курьера по его ID,
//меняет поле finished и finish_time в БД
func (h *Handler) FinishTask(ctx context.Context, courierUUID string) error {
	queue, err := h.DB.CourierTasks(ctx, courierUUID)
	if err != nil {
		logs.Eloger.WithError(err).Error("finish task")
		return errors.Wrap(err, "fail to delete task from queue")
	}
	if len(queue.Tasks) != 0 {
		task := queue.Tasks[0]
		if len(queue.Tasks) > 1 {
			queue.Tasks = queue.Tasks[1:]
		} else {
			queue.Tasks = []models.Task{}
		}

		event := h.EventCreating("task finished", "done", queue.CourierUUID, "")
		err = h.DB.FinishTask(ctx, task, queue, event)
		if err != nil {
			logs.Eloger.WithError(err).Error("finish task")
			return errors.Wrap(err, "fail to delete task from db")
		}
		return nil
	}
	return errors.New("no more tasks")
}

//ChangeQueue меняет порядок очереди курьера по его ID
//на вход получает ID курьера и ID тасков которые надо поменять местами
func (h *Handler) ChangeQueue(ctx context.Context, IDs *proto.QueueTasks) error {
	queue, err := h.DB.CourierTasks(ctx, IDs.CourierUUID)
	if err != nil {
		logs.Eloger.WithError(err).Error("update couriers tasks queue")
		return errors.Wrap(err, "fail update couriers tasks queue")
	}
	var firstTaskIndex int
	var secondTaskIndex int
	for i, v := range queue.Tasks {
		if v.UUID == IDs.FirstTaskUUID {
			firstTaskIndex = i
		}
		if v.UUID == IDs.SecondTaskUUID {
			secondTaskIndex = i
		}
	}

	var task models.Task
	task = queue.Tasks[firstTaskIndex]
	queue.Tasks[firstTaskIndex] = queue.Tasks[secondTaskIndex]
	queue.Tasks[secondTaskIndex] = task
	event := h.EventCreating("queue updating", "done", queue.CourierUUID, "")
	err = h.DB.UpdateTaskQueue(ctx, queue, event)
	if err != nil {
		logs.Eloger.WithError(err).Error("update courier tasks")
		return errors.Wrap(err, "fail to update courier tasks")
	}
	return nil
}

//CourierTasks возврашает очерель курьера
func (h *Handler) CourierTasks(ctx context.Context, uuid string) ([]models.Task, error) {
	queue, err := h.DB.CourierTasks(ctx, uuid)
	if err != nil {
		logs.Eloger.WithError(err).Error("get courier tasks")
		return nil, errors.Wrap(err, "fail to get courier tasks")
	}

	return queue.Tasks, err
}

func (h *Handler) CourierQueue(ctx context.Context, uuid string) (models.Queue, error) {
	queue, err := h.DB.CourierQueue(ctx, uuid)
	if err != nil {
		logs.Eloger.WithError(err).Error("get courier tasks")
		return queue, errors.Wrap(err, "fail to get courier tasks")
	}
	return queue, nil
}

func (h *Handler) tasksCreating(order *models.Order, serviceType string, personData models.PersonsData, route models.Route) (task models.Task) {
	task.UUID = h.RAM.IDs.GenUUID()
	task.OrderNumber = order.OrderNumber
	task.ClientData.ClientPhone = personData.Phone
	task.ClientData.ClientName = personData.Name
	task.OrderUUID = order.UUID
	task.Route = route
	task.Type = serviceType

	return task
}
