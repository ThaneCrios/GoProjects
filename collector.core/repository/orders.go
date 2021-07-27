package repository

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
)

type OrdersRepository interface {
	GetOrderByUUID(ctx context.Context, uuid string) (models.Order, error)
	GetFreeOrders(ctx context.Context) ([]models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order, event models.Event) error
	SetCollectorToOrder(ctx context.Context, order models.Order) error
	DuplicateOrder(ctx context.Context, orderCollector models.OrderForCollectorDuplicate) error
	GetDuplicateOrderByUUID(ctx context.Context, orderUUID string) (models.OrderForCollectorDuplicate, error)
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
		return *order, errors.Wrap(err, "failed to get order by uuid")
	}
	return *order, nil
}

func (p *Pg) FinishCollectOrder(ctx context.Context, order *models.Order) (models.Order, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "finish collect order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, order).
			Set("cart_items = ?, state = ?", order.CartItems, "finished").
			Where("uuid = ?", order.UUID)
		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to finish order").Error(err)
			return errors.Wrap(err, "failed to finish order")
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
		return *order, errors.Wrap(err, "failed to get order by uuid")
	}
	return *order, nil
}

func (p *Pg) GetFreeOrders(ctx context.Context) (orders []models.Order, err error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get all free orders",
	})

	err = p.Db.ModelContext(timeout, &orders).
		Where("state = ?", "created").
		Select()
	if err != nil {
		log.WithField("reason", "failed to get all free orders").Error(err)
		return orders, errors.Wrap(err, "failed to get all free orders")
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
			Set("collector_uuid = ?, collector_data = ?, state = ?", order.CollectorUUID, order.CollectorData, "in progress").
			Where("uuid = ?", order.UUID)
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

func (p *Pg) UpdateDuplicateOrder(ctx context.Context, order models.OrderForCollectorDuplicate) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "Update duplicate order",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.OrderForCollectorDuplicate{}).
			Set("cart_items = ?", order.CartItems).
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
