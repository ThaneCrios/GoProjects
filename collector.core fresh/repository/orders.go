package repository

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/proto"
	"time"
)

type OrdersRepository interface {
	GetOrderByUUID(ctx context.Context, uuid string) (models.Order, error)
	GetFreeOrders(ctx context.Context, storeUUID []string) ([]models.Order, error)
	GetMyOrders(ctx context.Context, collectorUUID string) ([]models.OrderForCollectorDuplicate, error)
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateOrderState(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order, event models.Event) error
	SetCollectorToOrder(ctx context.Context, order models.Order) error
	DuplicateOrder(ctx context.Context, orderCollector models.OrderForCollectorDuplicate) error
	GetDuplicateOrderByUUID(ctx context.Context, orderUUID string) (models.OrderForCollectorDuplicate, error)
	CancelOrder(ctx context.Context, uuid proto.OrderCancel) error
	UpdateDuplicateOrder(ctx context.Context, order models.OrderForCollectorDuplicate) error
	FinishCollectOrder(ctx context.Context, order *models.Order) (models.Order, error)
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
		log.WithField("reason", "failed to get order by uuid").Error(err)
		return *order, errors.Wrap(err, "Заказ с данным UUID не найден.")
	}
	return *order, nil
}

func (p *Pg) CancelOrder(ctx context.Context, uuid proto.OrderCancel) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "cancel order",
	})
	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.OrderForCollectorDuplicate{}).
			Set("state = ?, updated_at = ?", proto.OrderStateCancelled, time.Now()).
			Where("order_uuid = ?", uuid.OrderUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to cancel duplicate order").Error(err)
			return errors.Wrap(err, "Ошибка при отмене копии заказа.")
		}

		query = tx.ModelContext(timeout, &models.Order{}).
			Set("state = ?, updated_at = ?", "canceled", time.Now()).
			Where("uuid = ?", uuid.OrderUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to cancel duplicate order").Error(err)
			return errors.Wrap(err, "Ошибка при отмене оригинального заказа.")
		}

		return nil

	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (p *Pg) FinishCollectOrder(ctx context.Context, order *models.Order) (models.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "finish collect order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, order).
			Set("items = ?, state = ?, total_price = ?", order.CartItems, proto.OrderStateReady, order.TotalPrice+order.DeliveryData.Price).
			Where("uuid = ?", order.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to finish order").Error(err)
			return errors.Wrap(err, "Ошибка при завершении сборки оригинального заказа.")
		}

		query = tx.ModelContext(timeout, &models.OrderForCollectorDuplicate{}).
			Set("state = ?", proto.OrderStateReady).
			Where("order_uuid = ?", order.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to finish duplicate order").Error(err)
			return errors.Wrap(err, "Ошибка при завершении сборки копии заказа.")
		}
		return nil

	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return models.Order{}, errors.Wrap(err, "transaction failed")
	}

	return *order, nil
}

func (p *Pg) GetDuplicateOrderByUUID(ctx context.Context, orderUUID string) (models.OrderForCollectorDuplicate, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get order by uuid",
	})

	order := new(models.OrderForCollectorDuplicate)
	err := p.Db.ModelContext(timeout, order).
		Where("order_uuid = ?", orderUUID).
		Select()

	if err != nil {
		log.WithField("reason", "failed to get").Error(err)
		return *order, errors.Wrap(err, "Не могу найти копию заказа с указанным UUID.")
	}
	return *order, nil
}

func (p *Pg) GetFreeOrders(ctx context.Context, storeUUID []string) (orders []models.Order, err error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get all free orders",
	})
	err = p.Db.ModelContext(timeout, &orders).
		Where("state = ? AND store_uuid IN (?) AND store_data ->> 'type' = ?", proto.OrderStateCooking, pg.In(storeUUID), "grocery").
		Where("collector_uuid IS NULL").
		Select()

	if err != nil {
		log.WithField("reason", "failed to get all free orders").Error(err)
		return orders, errors.Wrap(err, "Не могу найти свободные заказы. Попробуйте позже.")
	}

	return orders, nil
}

func (p *Pg) GetMyOrders(ctx context.Context, collectorUUID string) ([]models.OrderForCollectorDuplicate, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get all free orders",
	})

	orders := make([]models.OrderForCollectorDuplicate, 0)

	err := p.Db.ModelContext(timeout, &orders).
		Where("state = ? AND collector_uuid = ?", "in progress", collectorUUID).
		Select()

	if err != nil {
		log.WithField("reason", "failed to get collector orders").Error(err)
		return orders, errors.Wrap(err, "Не могу найти Ваши заказы. Возможно, на Вас сейчас нет назначенных заказов.")
	}

	return orders, nil
}

func (p *Pg) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
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

		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return order, errors.Wrap(err, "transaction failed")
	}
	return order, nil
}

func (p *Pg) UpdateOrderState(ctx context.Context, order *models.Order) (*models.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "order state changing",
	})
	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, order).
			Set("state = ?, updated_at = ?", order.State, time.Now()).
			Where("uuid = ?", order.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to finish order").Error(err)
			return errors.Wrap(err, "failed to finish order")
		}

		query = tx.ModelContext(timeout, &models.OrderForCollectorDuplicate{}).
			Set("state = ?", order.State).
			Where("order_uuid = ?", order.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to finish duplicate order").Error(err)
			return errors.Wrap(err, "failed to finish duplicate order")
		}
		return nil

	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return &models.Order{}, errors.Wrap(err, "transaction failed")
	}
	return order, nil
}

func (p *Pg) DuplicateOrder(ctx context.Context, orderCollector models.OrderForCollectorDuplicate) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "order duplicate inserting",
	})
	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.ModelContext(timeout, &orderCollector).Insert(); err != nil {
			log.WithField("reason", "failed to duplicate order").Error(err)
			return errors.Wrap(err, "failed to duplicate order")
		}

		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
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

func (p *Pg) SetCollectorToOrder(ctx context.Context, order models.Order) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "setting courier on order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &order).
			Set("collector_uuid = ?, collector_data = ?", order.CollectorUUID, order.CollectorData).
			Where("uuid = ?", order.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update order").Error(err)
			return errors.Wrap(err, "Ошибка при назначении сборщика на заказ.")
		}

		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}

func (p *Pg) UpdateDuplicateOrder(ctx context.Context, order models.OrderForCollectorDuplicate) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "Update duplicate order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.OrderForCollectorDuplicate{}).
			Set("items = ?, total_price = ?", order.CartItems, order.TotalPrice).
			Where("order_uuid = ?", order.OrderUUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update order").Error(err)
			return errors.Wrap(err, "failed to update order")
		}

		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}
	return nil
}
