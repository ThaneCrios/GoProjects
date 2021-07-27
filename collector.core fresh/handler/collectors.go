package handler

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
)

func (h *Handler) GetCollectorByUUID(ctx context.Context, uuid string) (*models.User, error) {
	collector, err := h.DB.GetCollectorByUUID(ctx, uuid)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "getting courier by uuid",
				"uuid":  uuid,
			}).Error()
		return nil, errors.Wrap(err, "fail to get courier by uuid")
	}
	return collector, nil
}

func (h *Handler) CreateCollector(ctx context.Context, collector models.User) (models.User, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier creating",
	})

	collectorResult, err := h.DB.CreateCollector(ctx, &collector)
	if err != nil {
		log.WithField("reason", "failed to create collector in DB").Error(err)
		return *collectorResult, errors.Wrap(err, "fail to create courier")
	}

	logs.Eloger.Info("courier created")

	return *collectorResult, nil
}
