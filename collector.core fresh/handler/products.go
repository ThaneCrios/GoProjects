package handler

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
)

func (h *Handler) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "product creating",
	})

	product.UUID = h.Functions.IDs.GenUUID()
	product, err := h.DB.CreateProduct(ctx, product)
	if err != nil {
		log.WithField("reason", "failed to create product DB").Error(err)
		return product, errors.Wrap(err, "fail to create product")
	}

	return product, nil
}

func (h *Handler) GetProductByBarCode(ctx context.Context, barCode string) (models.Product, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get product by barcode and name",
	})
	params := make(map[string]string)
	params["barcode"] = barCode

	prodWithBar, err := h.RPC.GetProductUUIDByBarCode(ctx, params)
	if err != nil {
		log.WithField("reason", "failed to get product uuid by bar code").Error(err)
		return models.Product{}, errors.Wrap(err, "failed to get product uuid by bar code")
	}

	productRes, err := h.RPC.GetProductByUUID(ctx, prodWithBar[0].ProductUUID)
	if err != nil {
		log.WithField("reason", "failed to get product").Error(err)
		return *productRes, errors.Wrap(err, "failed to get product")
	}

	return *productRes, err
}

func (h *Handler) GetProductByUUID(ctx context.Context, uuid string) (models.Product, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":     "getting product",
		"with uuid": uuid,
	})
	product, err := h.RPC.GetProductByUUID(ctx, uuid)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return models.Product{}, errors.Wrap(err, "fail to get order by uuid")
	}
	return *product, nil
}

func (h *Handler) AppointBarCode(ctx context.Context, appointParams models.ProductWithBarCode) (models.ProductWithBarCode, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":     "appoint bar code",
		"with uuid": appointParams.ProductUUID,
	})
	appointParams.UUID = h.Functions.IDs.GenUUID()
	prodWithBar, err := h.DB.AppointBarCode(ctx, appointParams)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return prodWithBar, errors.Wrap(err, "fail to get order by uuid")
	}
	return prodWithBar, nil
}

func (h *Handler) GetProductUUIDsWithBarCodes(ctx context.Context, userUUID string) ([]models.ProductWithBarCode, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get product uuid with bar codes",
	})

	user, err := h.RPC.GetUserDataByUUID(ctx, userUUID)
	if err != nil {
		log.WithField("reason", "rpc error")
		return []models.ProductWithBarCode{}, errors.Wrap(err, "rpc error")
	}

	products, err := h.DB.GetProductUUIDsWithBarCodes(ctx, user.Meta.StoresUUID)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return []models.ProductWithBarCode{}, errors.Wrap(err, "fail to get order by uuid")
	}

	return products, nil
}
