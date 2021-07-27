package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

//GetOrderByUUID возвращает заказ по указанному UUID
func (h *Handler) GetOrderByUUID(ctx context.Context, uuid string) (models.Order, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":        "getting order",
		"from address": uuid,
	})
	order, err := h.DB.GetOrderByUUID(ctx, uuid)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return order, errors.Wrap(err, "fail to get order by uuid")
	}
	return order, nil
}

//OrdersFilter возвращает заказы, удовлетворяющие указаннам параметрам в переменной filter
func (h *Handler) OrdersFilter(ctx context.Context, filter *proto.OrdersFilter) ([]models.Order, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "order filtration",
	})
	order, err := h.DB.OrdersFilter(ctx, filter)
	if err != nil {
		log.WithField("reason", "failed to filtration orders in DB").Error(err)
		return order, errors.Wrap(err, "fail to filtration orders")
	}
	return order, nil
}

//CreateOrder создаёт заказ
func (h *Handler) CreateOrder(ctx context.Context, order *models.Order) (models.Order, error) {
	var tasks []models.Task
	var events []models.Event
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":        "order creating",
		"from address": order.PickupRoute.Street,
	})
	err := order.Validate()
	if err != nil {
		log.WithField("reason", "failed to validate order data").Error(err)
		return *order, errors.Wrap(err, "failed to validate data")
	}

	order.UUID = h.RAM.IDs.GenUUID()
	order.OrderNumber = h.RAM.IDs.SliceUUID(order.UUID)
	tasks = h.CreateTasks(ctx, order)
	events = append(events, h.EventCreating("order created", "done", order.CourierUUID, order.UUID))
	events = append(events, h.EventCreating("task created", "done", order.CourierUUID, order.UUID))
	events = append(events, h.EventCreating("task created", "done", order.CourierUUID, order.UUID))

	order, err = h.DB.CreateOrder(ctx, order, &tasks, &events)
	if err != nil {
		log.WithField("reason", "failed to create order and tasks in DB").Error(err)
		return *order, errors.Wrap(err, "fail to create order")
	}
	return *order, nil
}

//UpdateOrder обновляет определённый заказ(внутри order приходит UUID заказа и параметры, которые нужно обновить)
func (h *Handler) UpdateOrder(ctx context.Context, order *models.Order) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":        "order update",
		"from address": order.UUID,
	})

	event := h.EventCreating("update order", "done", order.CourierUUID, order.UUID)
	err := h.DB.UpdateOrder(ctx, order, event)
	if err != nil {
		log.WithField("reason", "failed to order updated in DB").Error(err)
		return errors.Wrap(err, "failed to update order")
	}
	logs.Eloger.Info(fmt.Sprintf("order with uuid=%s updated", order.UUID))

	return nil
}

//DeleteOrder помечает заказ с определённым UUID как удалённый
func (h *Handler) DeleteOrder(ctx context.Context, orderUUID string) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":        "order delete",
		"from address": orderUUID,
	})
	order, err := h.DB.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		log.WithField("reason", "failed to getting tasks").Error(err)
		return errors.Wrap(err, "failed to getting tasks")
	}
	var queue models.Queue
	var tasks []models.Task
	if order.CourierUUID != "" {
		queue, err := h.DB.CourierTasks(ctx, order.CourierUUID)
		if err != nil {
			log.WithField("reason", "failed to getting tasks").Error(err)
			return errors.Wrap(err, "failed to getting tasks")
		}

		for _, v := range queue.Tasks {
			if v.OrderUUID != orderUUID {
				tasks = append(tasks, v)
			}
		}
		queue.Tasks = tasks
	}

	var events []models.Event
	events = append(events, h.EventCreating("delete order", "done", "", orderUUID))
	events = append(events, h.EventCreating("delete task", "done", "", orderUUID))
	events = append(events, h.EventCreating("delete task", "done", "", orderUUID))
	if order.CourierUUID != "" {
		events = append(events, h.EventCreating("update queue", "done", order.CourierUUID, ""))
	}

	err = h.DB.DeleteOrder(ctx, order, queue, events)
	if err != nil {
		log.WithField("reason", "failed to delete order in DB").Error(err)
		return errors.Wrap(err, "failed to delete order")
	}

	logs.Eloger.Info(fmt.Sprintf("order with uuid=%s deleted", orderUUID))
	return nil
}

//SetCourierToOrder назначает курьера на определённый заказ(в ids передаются UUID курьера и UUID заказа)
func (h *Handler) SetCourierToOrder(ctx context.Context, ids *proto.OrderCourier) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "set courier",
	})
	courierUUID := ids.CourierUUID
	orderUUID := ids.OrderUUID

	order, err := h.DB.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		log.WithField("reason", "failed to getting tasks").Error(err)
		return errors.Wrap(err, "failed to getting tasks")
	}
	if order.CourierUUID != "" {
		log.WithField("reason", "failed set courier").Error(err)
		return errors.New("courier has already been assigned")
	}
	tasks, err := h.DB.OrderTasks(ctx, orderUUID)
	if err != nil {
		log.WithField("reason", "failed to getting tasks").Error(err)
		return errors.Wrap(err, "failed to getting tasks")
	}
	queue, err := h.DB.CourierQueue(ctx, courierUUID)
	if err != nil {
		log.WithField("reason", "failed to getting queue").Error(err)
		return errors.Wrap(err, "failed to getting queue")
	}
	courier, err := h.DB.GetCourierByUUID(ctx, courierUUID)
	if err != nil {
		log.WithField("reason", "failed to getting courier").Error(err)
		return errors.Wrap(err, "failed to getting courier")
	}

	tasks[0].State = "in progress"
	tasks[1].State = "in progress"
	queue.Tasks = append(queue.Tasks, tasks[0])
	queue.Tasks = append(queue.Tasks, tasks[1])
	event := h.EventCreating("set courier", "done", courierUUID, orderUUID)

	var courierData models.CourierData
	courierData = courierDataConvert(*courier)
	err = h.DB.SetCourierToOrder(ctx, queue, orderUUID, event, courierData)
	if err != nil {
		log.WithField("reason", "failed to set courier to order in DB").Error(err)
		return errors.Wrap(err, "fail to set courier to order")
	}
	logs.Eloger.Info(fmt.Sprintf("order with uuid=%s updated", ids.OrderUUID))

	//if len(queue.Tasks) == 2 {
	//	courier, err := h.GetCourierByUUID(ctx, ids.CourierUUID)
	//	if err != nil {
	//		log.WithField("reason", "failed to set courier to order in DB").Error(err)
	//		return errors.Wrap(err, "fail to set courier to order")
	//	}
	//	chatId := courier.ChatID
	//
	//	var url = fmt.Sprintf("http://192.168.1.48:8000/api/v1/newTaskNotification/%s", chatId)
	//	err = tools.RPC(http.MethodPost, url, nil, nil, nil)
	//	if err != nil {
	//		log.WithField("reason", "failed to set courier to order in DB").Error(err)
	//		return errors.Wrap(err, "fail to set courier to order")
	//	}
	//}
	return nil
}

func (h *Handler) EventCreating(eventType string, payload string, courierUUID string, orderUUID string) models.Event {
	var event models.Event
	event.UUID = h.RAM.IDs.GenUUID()
	event.CreatedAt = time.Now()
	event.EventType = eventType
	event.Payload.Information = payload
	event.CourierUUID = courierUUID
	event.OrderUUID = orderUUID
	return event
}

func (h *Handler) RemoveCourierFromOrder(ctx context.Context, ids *proto.OrderCourier) error {

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "set courier to tasks",
	})

	queue, err := h.DB.CourierTasks(ctx, ids.CourierUUID)
	if err != nil {
		log.WithField("reason", "fail remove courier from order").Error(err)
		return errors.Wrap(err, "fail remove courier from order")
	}

	var tasks []models.Task
	for _, v := range queue.Tasks {
		if v.OrderUUID != ids.OrderUUID {
			tasks = append(tasks, v)
		}
	}
	queue.Tasks = tasks

	event := h.EventCreating("order canceled", "done", ids.CourierUUID, ids.OrderUUID)
	err = h.DB.RemoveCourierFromOrder(ctx, ids.OrderUUID, queue, event)
	if err != nil {
		log.WithField("reason", "fail remove courier from order").Error(err)
		return errors.Wrap(err, "fail remove courier from order")
	}

	logs.Eloger.Info(fmt.Sprintf("order with order_uuid=%s updated", ids.OrderUUID))
	return nil
}

func courierDataConvert(courier models.Courier) models.CourierData {
	return models.CourierData{
		UUID:        courier.UUID,
		Name:        courier.CourierMeta.FullName,
		PhoneNumber: courier.PhoneNumber,
	}
}
