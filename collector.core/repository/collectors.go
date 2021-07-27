package repository

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
	"time"
)

type CollectorsRepository interface {
	GetCollectorByUUID(ctx context.Context, uuid string) (*models.Collector, error)
	CreateCollector(ctx context.Context, collector *models.Collector) (*models.Collector, error)
	UpdateCollector(ctx context.Context, collector *models.Collector, event *models.Event) (*models.Collector, error)
	DeleteCollector(ctx context.Context, uuid string, event *models.Event) error
}

func (p *Pg) GetCollectorByUUID(ctx context.Context, uuid string) (*models.Collector, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	collector := new(models.Collector)
	err := p.Db.ModelContext(timeout, collector).
		Where("uuid = ?", uuid).
		Select()

	return collector, err
}

func (p *Pg) CreateCollector(ctx context.Context, collector *models.Collector) (*models.Collector, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier inserting",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {

		if _, err := tx.ModelContext(timeout, collector).Insert(); err != nil {
			log.WithField("reason", "failed to create collector").Error(err)
			return errors.Wrap(err, "failed to create collector")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return collector, errors.Wrap(err, "transaction failed")
	}
	return collector, nil
}

func (p *Pg) UpdateCollector(ctx context.Context, collector *models.Collector, event *models.Event) (*models.Collector, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	err := p.CheckExistsByUUID(ctx, collector, collector.UUID)
	if err != nil {
		return nil, err
	}

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier updating",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, collector).
			Where("uuid = ?", collector.UUID)

		if _, err := query.UpdateNotNull(); err != nil {
			log.WithField("reason", "failed to update collector").Error(err)
			return errors.Wrap(err, "failed to update collector")
		}

		if _, err := tx.ModelContext(timeout, event).Insert(); err != nil {
			log.WithField("reason", "failed to create events").Error(err)
			return errors.Wrap(err, "failed to create events")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return nil, errors.Wrap(err, "transaction failed")
	}

	return collector, nil
}

func (p *Pg) DeleteCollector(ctx context.Context, uuid string, event *models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	err := p.CheckExistsByUUID(ctx, &models.Collector{}, uuid)
	if err != nil {
		return err
	}

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier updating",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Collector{}).
			Set("updated_at = ?, deleted_at=?", time.Now(), time.Now()).
			Where("uuid = ?", uuid)

		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to delete courier").Error(err)
			return errors.Wrap(err, "failed to delete courier")
		}

		if _, err := tx.ModelContext(timeout, event).Insert(); err != nil {
			log.WithField("reason", "failed to create events").Error(err)
			return errors.Wrap(err, "failed to create events")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}
