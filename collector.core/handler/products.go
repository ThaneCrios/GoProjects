package handler

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
)

func (h *Handler) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "product creating",
	})

	product.UUID = h.RAM.IDs.GenUUID()
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
	prodWithBar, err := h.DB.GetProductByBarCode(ctx, barCode)
	productRes, err := h.DB.GetProductByUUID(ctx, prodWithBar.ProductUUID)
	if err != nil {
		log.WithField("reason", "failed to get product").Error(err)
		return productRes, errors.Wrap(err, "failed to get product")
	}

	return productRes, err
}

func (h *Handler) GetProductByUUID(ctx context.Context, uuid string) (models.Product, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":     "getting product",
		"with uuid": uuid,
	})
	product, err := h.DB.GetProductByUUID(ctx, uuid)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return product, errors.Wrap(err, "fail to get order by uuid")
	}
	return product, nil
}

func (h *Handler) AppointBarCode(ctx context.Context, appointParams models.ProductWithBarCode) (models.ProductWithBarCode, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":     "appoint bar code",
		"with uuid": appointParams.ProductUUID,
	})
	appointParams.UUID = h.RAM.IDs.GenUUID()
	prodWithBar, err := h.DB.AppointBarCode(ctx, appointParams)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return prodWithBar, errors.Wrap(err, "fail to get order by uuid")
	}
	return prodWithBar, nil
}

//func (h *Handler) CompareTwoProducts(ctx context.Context, barCodes models.BarCodes) (string, error) {
//	log := logs.Eloger.WithFields(logrus.Fields{
//		"event": "get products by barcode and compare",
//	})
//
//	productOrder, err := h.DB.GetProductByBarCode(ctx, barCodes.BarCodeOrder)
//	if err != nil {
//		log.WithField("reason", "failed to get product from order").Error(err)
//		return "", errors.Wrap(err, "failed to get product")
//	}
//	productScanned, err := h.DB.GetProductByBarCode(ctx, barCodes.BarCodeScanned)
//	if err != nil {
//		log.WithField("reason", "failed to get product").Error(err)
//		return "", errors.Wrap(err, "failed to get product")
//	}
//	response := compareProducts(productOrder, productScanned)
//	return response, nil
//}

func compareProducts(productOrder, productScanned models.Product) (response string) {
	if productOrder.UUID == productScanned.UUID {
		response = "equal"
	} else {
		response = "unequal"
	}
	return response
}
