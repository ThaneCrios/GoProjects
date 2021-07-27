package handler

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

func (h *Handler) CreateEvent(ctx context.Context, event *models.Event) error {
	event.UUID = h.RAM.IDs.GenUUID()
	err := h.DB.CreateEvent(ctx, event)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "create event",
			}).Error()
		return errors.Wrap(err, "fail to create event")
	}
	logs.Eloger.Info("event created")

	return nil
}

func (h *Handler) GetCourierEvents(ctx context.Context, courierUUID string) ([]models.Event, error) {
	events, err := h.DB.GetCourierEvents(ctx, courierUUID)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "get courier events",
			}).Error()
		return events, errors.Wrap(err, "fail to get courier events")
	}
	logs.Eloger.Info("events returned")

	return events, nil
}

func (h *Handler) GetOrderEvents(ctx context.Context, courierUUID string) ([]models.Event, error) {
	events, err := h.DB.GetOrderEvents(ctx, courierUUID)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "get order events",
			}).Error()
		return events, errors.Wrap(err, "fail to get order events")
	}
	logs.Eloger.Info("events returned")

	return events, nil
}

func (h *Handler) EventsFilter(ctx context.Context, filter *proto.EventFilter) ([]models.Event, error) {
	event, err := h.DB.EventsFilter(ctx, filter)
	if err != nil {
		logs.Eloger.WithError(err).Error("events filtering")
		return nil, errors.Wrap(err, "fail to get filtered events")
	}
	return event, nil
}
