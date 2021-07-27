package handler

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

//GetCourierByUUID возвращает курьера по ID
func (h *Handler) GetCourierByUUID(ctx context.Context, uuid string) (*models.Courier, error) {
	courier, err := h.DB.GetCourierByUUID(ctx, uuid)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "getting courier by uuid",
				"uuid":  uuid,
			}).Error()
		return nil, errors.Wrap(err, "fail to get courier by uuid")
	}
	return courier, nil
}

func (h *Handler) GetCourierTokenInfoByPhoneNumber(ctx context.Context, phoneNumber string) (*models.AuthTokenDate, error) {
	var tokenDate models.AuthTokenDate
	courier, err := h.DB.GetCourierByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "getting courier by phone number",
				"uuid":  phoneNumber,
			}).Error()
		return nil, errors.Wrap(err, "fail to get courier by phone number")
	}
	tokenDate.CourierUUID = courier.UUID
	tokenDate.Role = "Курьер"

	return &tokenDate, nil
}

//CreateCourier создает курьера по полученным данным и возврашает его обратно
func (h *Handler) CreateCourier(ctx context.Context, courFront proto.CourierFront) (models.Courier, error) {
	var event models.Event
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier creating",
	})
	courierBack := proto.CourierFromFrontToCore(courFront)
	err := courierBack.Validate()
	if err != nil {
		log.WithField("reason", "failed to validate").Error(err)
		return courierBack, errors.Wrap(err, "failed to validate data")
	}
	courierBack.UUID = h.RAM.IDs.GenUUID()
	event = h.EventCreating("courier creating", "done", courierBack.UUID, "")
	courierResult, err := h.DB.CreateCourier(ctx, &courierBack, &event)
	if err != nil {
		log.WithField("reason", "failed to create courier in DB").Error(err)
		return *courierResult, errors.Wrap(err, "fail to create courier")
	}

	logs.Eloger.Info("courier created")

	return *courierResult, nil
}

//UpdateCourier обновляет курьера по его ID, изменяя нужные поля на новые
func (h *Handler) UpdateCourier(ctx context.Context, courier *models.Courier) (*models.Courier, error) {
	event := h.EventCreating("update courier", "done", courier.UUID, "")
	courier, err := h.DB.UpdateCourier(ctx, courier, &event)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "update courier by uuid",
				"uuid":  courier.UUID,
			}).Error()
		return nil, errors.Wrap(err, "fail to update courier")
	}

	logs.Eloger.Info(fmt.Sprintf("courier with uuid=%s updated", courier.UUID))

	return courier, nil
}

//CouriersFilter возвращает список курьеров отфильтрованных по полученным параметрам
func (h *Handler) CouriersFilter(ctx context.Context, filter *proto.CouriersFilter) ([]models.Courier, error) {
	courier, err := h.DB.CouriersFilter(ctx, filter)
	if err != nil {
		logs.Eloger.WithError(err).Error("couriers filtering")
		return nil, errors.Wrap(err, "fail to get filtered couriers")
	}
	return courier, nil
}

//DeleteCourier меняет, для курьера по его ID, поле deleted на true в БД
func (h *Handler) DeleteCourier(ctx context.Context, uuid string) error {
	event := h.EventCreating("delete courier", "done", uuid, "")
	err := h.DB.DeleteCourier(ctx, uuid, &event)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "delete courier by uuid",
				"uuid":  uuid,
			}).Error()
		return errors.Wrap(err, "fail to delete courier")
	}

	logs.Eloger.Info(fmt.Sprintf("city with uuid=%s deleted", uuid))

	return nil
}

//GetCourierByChatID возвращает курьера по номеру телефона и chat_ID
func (h *Handler) GetCourierByChatID(ctx context.Context, chatID string, phoneNumber string) error {
	flag, err := h.DB.GetCourierByChatID(ctx, chatID, phoneNumber)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event":  "getting courier by chatID",
				"chatID": chatID,
			}).Error()
		return errors.Wrap(err, "cant find courier, try again.")
	}
	if flag == false {
		return errors.New("cant find courier, try again.")
	}
	return nil
}

//UpdateCourierCoordinates обновляет координаты курьера по его ID
func (h *Handler) UpdateCourierCoordinates(ctx context.Context, coordinates *proto.CouriersFilter) error {
	courierCoordsTable := h.CreateCourierCoordinatesForTable(coordinates)
	err := h.DB.UpdateCourierCoordinates(ctx, coordinates, &courierCoordsTable)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event":       "update courier coordinates",
				"coordinates": coordinates,
			}).Error()
		return errors.Wrap(err, "cant update courier coordinates, try again.")
	}
	return nil
}

//UpdateCourierCoordinatesTable обновляет координаты курьера
func (h *Handler) UpdateCourierCoordinatesTable(ctx context.Context, coordinates *proto.CouriersFilter) error {
	err := h.DB.UpdateCourierCoordinatesTable(ctx, coordinates)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event":       "update courier coordinates",
				"coordinates": coordinates,
			}).Error()
		return errors.Wrap(err, "cant update courier coordinates, try again.")
	}
	return nil
}

//GetCourierUUIDByChatID возвращает ID курьера по его chatID
func (h *Handler) GetCourierUUIDByChatID(ctx context.Context, chatID string) (string, error) {
	uuid, err := h.DB.GetCourierUUIDByChatID(ctx, chatID)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event":   "get courier uuid by chat_id",
				"chat_id": chatID,
			}).Error()
		return uuid, errors.Wrap(err, "cant update courier coordinates, try again.")
	}
	return uuid, nil
}

//CreateQueue создает очередь для курьера
func (h *Handler) CreateQueue(ctx context.Context, courierUUID string) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "queue creating",
	})

	var queue models.Queue
	queue.CourierUUID = courierUUID
	queue.UUID = h.RAM.IDs.GenUUID()

	err := h.DB.CreateQueue(ctx, &queue)
	if err != nil {
		log.WithField("reason", "failed to create queue in DB").Error(err)
		return errors.Wrap(err, "fail to create queue")
	}

	logs.Eloger.Info("queue created")

	return nil
}

//CreateCourierCoordinatesForTable ....
func (h *Handler) CreateCourierCoordinatesForTable(coord *proto.CouriersFilter) models.CourierCoordinatesTable {
	return models.CourierCoordinatesTable{
		UUID:        h.RAM.IDs.GenUUID(),
		CourierUUID: coord.CourierUUID,
		Lat:         coord.CourierLat,
		Lon:         coord.CourierLon,
	}
}
