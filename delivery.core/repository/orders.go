package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

type OrdersRepository interface {
	GetOrderByUUID(ctx context.Context, uuid string) (models.Order, error)
	OrdersFilter(ctx context.Context, filter *proto.OrdersFilter) ([]models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order, tasks *[]models.Task, events *[]models.Event) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order, event models.Event) error
	DeleteOrder(ctx context.Context, order models.Order, queue models.Queue, events []models.Event) error
	FinishTasks(ctx context.Context, uuid string) error
	SetCourierToOrder(ctx context.Context, queue models.Queue, orderUUID string, event models.Event, courierData models.CourierData) error
	GetTasksByCourierUUID(ctx context.Context, uuid string) (*[]*models.Task, error)
	FinishTask(ctx context.Context, task models.Task, queue models.Queue, event models.Event) error
	GetNextTaskCourier(ctx context.Context, courierUUID string) (models.Task, error)
	UpdateTaskQueue(ctx context.Context, queue models.Queue, event models.Event) error
	RemoveCourierFromOrder(ctx context.Context, orderUUID string, queue models.Queue, event models.Event) error
}

func (p *Pg) GetOrderByUUID(ctx context.Context, uuid string) (models.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get order by uuid",
	})

	order := new(models.Order)
	err := p.Db.ModelContext(timeout, order).
		Where("uuid = ?", uuid).
		Select()
	if err != nil {
		log.WithField("reason", "failed to get").Error(err)
		return *order, errors.Wrap(err, "failed to get order by uuid")
	}
	return *order, nil
}

func (p *Pg) OrdersFilter(ctx context.Context, filter *proto.OrdersFilter) ([]models.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	fmt.Println(filter)
	order := new([]models.Order)
	query := p.Db.ModelContext(timeout, order)

	if filter.OrderService != "" {
		fmt.Println("AAAAAAAAAAAAAAA")
		query.Where("service = ?", filter.OrderService)
	}
	if filter.CourierUUID != "" {
		query.Where("courier_uuid = ?", filter.CourierUUID)
	}
	query.Where("deleted_at is null")

	if err := query.Select(); err != nil {
		return nil, err
	}

	return *order, nil
}

func (p *Pg) CreateOrder(ctx context.Context, order *models.Order, tasks *[]models.Task, event *[]models.Event) (*models.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "order inserting",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.ModelContext(timeout, order).Insert(); err != nil {
			log.WithField("reason", "failed to create order").Error(err)
			return errors.Wrap(err, "failed to create order")
		}

		if _, err := tx.ModelContext(timeout, tasks).Insert(); err != nil {
			log.WithField("reason", "failed to create tasks").Error(err)
			return errors.Wrap(err, "failed to create tasks")
		}

		if _, err := tx.ModelContext(timeout, event).Insert(); err != nil {
			log.WithField("reason", "failed to create events").Error(err)
			return errors.Wrap(err, "failed to create events")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return order, errors.Wrap(err, "transaction failed")
	}
	return order, nil
}

func (p *Pg) UpdateOrder(ctx context.Context, order *models.Order, event models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "order updating",
	})

	err := p.CheckExistsByUUID(ctx, &models.Order{}, order.UUID)
	if err != nil {
		return errors.Wrap(err, "failed check exists by uuid")
	}

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, order).
			Where("uuid = ?", order.UUID)
		if _, err := query.UpdateNotNull(); err != nil {
			log.WithField("reason", "failed to update order").Error(err)
			return errors.Wrap(err, "failed to update order")
		}

		if _, err = tx.ModelContext(timeout, &event).Insert(); err != nil {
			log.WithField("reason", "failed to create event").Error(err)
			return errors.Wrap(err, "failed to create event")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (p *Pg) DeleteOrder(ctx context.Context, order models.Order, queue models.Queue, events []models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "order deleting",
	})

	err := p.CheckExistsByUUID(ctx, &models.Order{}, order.UUID)
	if err != nil {
		return err
	}

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Order{}).
			Set("deleted_at = ?, updated_at = ?", time.Now(), time.Now()).
			Where("uuid = ?", order.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to delete order").Error(err)
			return errors.Wrap(err, "failed to delete order")
		}
		if order.CourierUUID != "" {
			query = tx.ModelContext(timeout, &models.Queue{}).
				Set("tasks = ?", queue.Tasks).
				Where("courier_uuid = ?", queue.CourierUUID)
			if _, err := query.Update(); err != nil {
				log.WithField("reason", "failed to delete tasks in queue").Error(err)
				return errors.Wrap(err, "failed to delete tasks in queue")
			}
		}

		query = tx.ModelContext(timeout, &models.Task{}).
			Set("finish_time = ?", time.Now()).
			Where("order_uuid = ?", order.UUID)
		if _, err = query.Update(); err != nil {
			log.WithField("reason", "failed to delete tasks").Error(err)
			return errors.Wrap(err, "failed to delete tasks")
		}

		if _, err = tx.ModelContext(timeout, &events).Insert(); err != nil {
			log.WithField("reason", "failed to create event").Error(err)
			return errors.Wrap(err, "failed to create event")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return err
}

func (p *Pg) FinishTasks(ctx context.Context, uuid string) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "delete tasks with order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Task{}).
			Set(" finish_time=?", true, time.Now()).
			Where("order_uuid = ?", uuid)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to delete tasks").Error(err)
			return errors.Wrap(err, "failed to delete tasks")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (p *Pg) SetCourierToOrder(ctx context.Context, queue models.Queue, orderUUID string, event models.Event, courierData models.CourierData) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "setting courier on order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Order{}).
			Set("courier_uuid = ?, state = ?, courier_data = ?", queue.CourierUUID, proto.Variables.Order.OrderState.InProgress, courierData).
			Where("uuid = ?", orderUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update order").Error(err)
			return errors.Wrap(err, "failed to update order")
		}

		query = tx.ModelContext(timeout, &models.Queue{}).
			Set("tasks = ?", queue.Tasks).
			Where("courier_uuid = ?", queue.CourierUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update tasks in queue").Error(err)
			return errors.Wrap(err, "failed to update tasks in queue")
		}

		query = tx.ModelContext(timeout, &models.Task{}).
			Set("courier_uuid = ?, state= ?", queue.CourierUUID, "in progress").
			Where("order_uuid = ?", orderUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update tasks").Error(err)
			return errors.Wrap(err, "failed to update tasks")
		}

		if _, err := tx.ModelContext(timeout, &event).Insert(); err != nil {
			log.WithField("reason", "failed to create event").Error(err)
			return errors.Wrap(err, "failed to create event")
		}

		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (p *Pg) GetTasksByCourierUUID(ctx context.Context, uuid string) (*[]*models.Task, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	tasks := new([]*models.Task)
	err := p.Db.ModelContext(timeout, tasks).
		Where("courier_uuid=?", uuid).
		Select()

	return tasks, err
}

func (p *Pg) FinishTask(ctx context.Context, task models.Task, queue models.Queue, event models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "finishing task",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Task{}).
			Set("state = ?, finish_time = ?", "done", time.Now()).
			Where("uuid = ?", task.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update order").Error(err)
			return errors.Wrap(err, "failed to update order")
		}

		query = tx.ModelContext(timeout, &models.Queue{}).
			Set("tasks = ?", queue.Tasks).
			Where("courier_uuid = ?", queue.CourierUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update tasks in queue").Error(err)
			return errors.Wrap(err, "failed to update tasks in queue")
		}

		if task.Type == "Отвезти" {
			query = tx.ModelContext(timeout, &models.Order{}).
				Set("state = ?,finish_at = ? ", "dose", time.Now()).
				Where("uuid = ?", task.OrderUUID)
			if _, err := query.Update(); err != nil {
				log.WithField("reason", "failed to update tasks").Error(err)
				return errors.Wrap(err, "failed to update tasks")
			}
		}
		if _, err := tx.ModelContext(timeout, &event).Insert(); err != nil {
			log.WithField("reason", "failed to create event").Error(err)
			return errors.Wrap(err, "failed to create event")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (p *Pg) GetNextTaskCourier(ctx context.Context, courierUUID string) (models.Task, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	task := new(models.Task)
	err := p.Db.ModelContext(timeout, task).
		Where("courier_uuid=? AND parent_task_uuid=?", courierUUID, "0").
		Select()

	return *task, err
}

func (p *Pg) UpdateTaskQueue(ctx context.Context, queue models.Queue, event models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "update queue",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Queue{}).
			Set("tasks = ?", queue.Tasks).
			Where("courier_uuid = ?", queue.CourierUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update tasks in queue").Error(err)
			return errors.Wrap(err, "failed to update tasks in queue")
		}

		if _, err := tx.ModelContext(timeout, &event).Insert(); err != nil {
			log.WithField("reason", "failed to create event").Error(err)
			return errors.Wrap(err, "failed to create event")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (p *Pg) RemoveCourierFromOrder(ctx context.Context, orderUUID string, queue models.Queue, event models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "setting courier on order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Order{}).
			Set("courier_uuid = ?, updated_at = ?, state = ?", nil, time.Now(), "свободен").
			Where("uuid = ?", orderUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update order").Error(err)
			return errors.Wrap(err, "failed to update order")
		}

		query = tx.ModelContext(timeout, &models.Queue{}).
			Set("tasks = ?", queue.Tasks).
			Where("courier_uuid = ?", queue.CourierUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update tasks in queue").Error(err)
			return errors.Wrap(err, "failed to update tasks in queue")
		}

		query = tx.ModelContext(timeout, &models.Task{}).
			Set("courier_uuid = ?, state = ?", nil, "wait courier").
			Where("order_uuid = ?", orderUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update tasks").Error(err)
			return errors.Wrap(err, "failed to update tasks")
		}

		if _, err := tx.ModelContext(timeout, &event).Insert(); err != nil {
			log.WithField("reason", "failed to create event").Error(err)
			return errors.Wrap(err, "failed to create event")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}
