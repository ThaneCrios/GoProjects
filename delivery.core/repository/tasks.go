package repository

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
	"time"
)

type TasksRepository interface {
	OrderTasks(ctx context.Context, uuid string) ([]models.Task, error)
	TasksListByStatus(ctx context.Context, filter *proto.TasksFilter) ([]models.Task, error)
	GetCourierTasksByStatus(ctx context.Context, filter *proto.TasksFilter) ([]models.Task, error)
	UpdateTask(ctx context.Context, filter *proto.TasksFilter, event models.Event) (models.Task, error)
	CourierTasks(ctx context.Context, uuid string) (models.Queue, error)
	CourierQueue(ctx context.Context, uuid string) (models.Queue, error)
	DeleteTasks(ctx context.Context, orderUUID string) error
}

func (p *Pg) OrderTasks(ctx context.Context, uuid string) ([]models.Task, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	tasks := new([]models.Task)
	err := p.Db.ModelContext(timeout, tasks).
		Where("order_uuid = ?", uuid).
		Select()
	return *tasks, err
}

//TasksListByStatus возвращает список тасков по определённому таску
func (p *Pg) TasksListByStatus(ctx context.Context, filter *proto.TasksFilter) ([]models.Task, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	var tasks []models.Task

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "tasks searching by status",
	})

	err := p.Db.ModelContext(timeout, tasks).
		Where("status = ?", filter.Status).
		Where("deleted = ", false).
		Select()

	if err != nil {
		if err != nil {
			log.WithField("reason", "failed to find tasks").Error(err)
			return nil, errors.Wrap(err, "failed to find tasks by status")
		}
	}

	return tasks, nil
}

//GetCourierTasksByStatus возвращает список тасков курьера по определённому статусу
func (p *Pg) GetCourierTasksByStatus(ctx context.Context, filter *proto.TasksFilter) ([]models.Task, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	var tasks []models.Task

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "tasks searching by status and courier uuid",
	})

	err := p.Db.ModelContext(timeout, tasks).
		Where("status = ?", filter.Status).
		Where("courier_uuid = ?", filter.CourierUUID).
		Where("deleted = ?", false).
		Select()

	if err != nil {
		if err != nil {
			log.WithField("reason", "failed to find tasks").Error(err)
			return nil, errors.Wrap(err, "failed to find tasks by status with currend courier uuid")
		}
	}

	return tasks, nil
}

func (p *Pg) UpdateTask(ctx context.Context, filter *proto.TasksFilter, event models.Event) (models.Task, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	var task models.Task
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "change status",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Task{}).Where("uuid = ?", filter.TaskUUID)
		if filter.CourierUUID != "" {
			query.Set("courier_uuid = ?", filter.CourierUUID)
		}
		if filter.Status != "" {
			query.Set("state = ?", filter.Status)
		}
		if filter.ParentTaskUUID != "" {
			query.Set("parent_task_uuid = ?", filter.ParentTaskUUID)
		}

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
		return models.Task{}, errors.Wrap(err, "transaction failed")
	}

	err := p.Db.ModelContext(timeout, &task).
		Where("uuid = ?", filter.TaskUUID).
		Select()
	if err != nil {
		log.WithField("reason", "failed to return updated task").Error(err)
		return models.Task{}, errors.Wrap(err, "failed to return updated task")
	}
	return task, nil
}

func (p *Pg) CourierTasks(ctx context.Context, uuid string) (models.Queue, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	var queue models.Queue
	err := p.Db.ModelContext(timeout, &queue).
		Where("courier_uuid = ?", uuid).
		Select()
	return queue, err
}

func (p *Pg) CourierQueue(ctx context.Context, uuid string) (models.Queue, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	queue := new(models.Queue)
	err := p.Db.ModelContext(timeout, queue).
		Where("courier_uuid = ?", uuid).
		Select()
	return *queue, err
}

func (p *Pg) DeleteTasks(ctx context.Context, orderUUID string) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err := p.Db.ModelContext(timeout, &models.Task{}).
		Set("deleted_at = ?", time.Now()).
		Where("order_uuid = ?", orderUUID).
		Update()

	return err
}
