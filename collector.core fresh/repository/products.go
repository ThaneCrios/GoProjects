package repository

import (
	"context"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
)

type ProductsRepository interface {
	CreateProduct(ctx context.Context, product models.Product) (models.Product, error)
	GetProductByBarCode(ctx context.Context, barCode string, storesUUID []string) (models.ProductWithBarCode, error)
	GetProductByUUID(ctx context.Context, uuid string) (models.Product, error)
	AppointBarCode(ctx context.Context, appointParams models.ProductWithBarCode) (models.ProductWithBarCode, error)
	GetProductUUIDsWithBarCodes(ctx context.Context, storesUUID []string) ([]models.ProductWithBarCode, error)
}

func (p *Pg) GetProductByUUID(ctx context.Context, uuid string) (models.Product, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get order by uuid",
	})

	product := new(models.Product)
	err := p.Db.ModelContext(timeout, product).
		Where("uuid = ?", uuid).
		Select()
	if err != nil {
		log.WithField("reason", "failed to get").Error(err)
		return *product, errors.Wrap(err, "failed to get order by uuid")
	}
	return *product, nil
}

func (p *Pg) AppointBarCode(ctx context.Context, appointParams models.ProductWithBarCode) (models.ProductWithBarCode, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "appoint bar code",
	})

	_, err := p.Db.ModelContext(timeout, &appointParams).Insert()

	if err != nil {
		log.WithField("reason", "failed to get").Error(err)
		return appointParams, errors.Wrap(err, "failed to get order by uuid")
	}
	return appointParams, nil
}

func (p *Pg) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "product inserting",
	})
	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.ModelContext(timeout, &product).Insert(); err != nil {
			log.WithField("reason", "failed to create product").Error(err)
			return errors.Wrap(err, "failed to create product")
		}

		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return product, errors.Wrap(err, "transaction failed")
	}
	return product, nil
}

func (p *Pg) GetProductByBarCode(ctx context.Context, barCode string, storesUUID []string) (models.ProductWithBarCode, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	var productRes models.ProductWithBarCode

	err := p.Db.ModelContext(timeout, &productRes).
		Where("bar_code = ? AND store_uuid IN (?)", barCode, pg.In(storesUUID)).
		Select()

	return productRes, err
}

func (p *Pg) GetProductUUIDsWithBarCodes(ctx context.Context, storesUUID []string) ([]models.ProductWithBarCode, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	var productRes []models.ProductWithBarCode

	err := p.Db.ModelContext(timeout, &productRes).
		Where("store_uuid IN (?)", pg.In(storesUUID)).
		Select()

	return productRes, err
}
