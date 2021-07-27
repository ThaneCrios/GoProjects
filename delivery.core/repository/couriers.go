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

type CouriersRepository interface {
	GetCourierByUUID(ctx context.Context, uuid string) (*models.Courier, error)
	CreateCourier(ctx context.Context, courier *models.Courier, event *models.Event) (*models.Courier, error)
	UpdateCourier(ctx context.Context, courier *models.Courier, event *models.Event) (*models.Courier, error)
	DeleteCourier(ctx context.Context, uuid string, event *models.Event) error
	CouriersFilter(ctx context.Context, filter *proto.CouriersFilter) ([]models.Courier, error)
	GetCourierByChatID(ctx context.Context, chatID string, phoneNumber string) (bool, error)
	UpdateCourierCoordinates(ctx context.Context, coordinates *proto.CouriersFilter, courCoordsTable *models.CourierCoordinatesTable) error
	UpdateCourierCoordinatesTable(ctx context.Context, coordinates *proto.CouriersFilter) error
	GetCourierByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Courier, error)
	GetCourierUUIDByChatID(ctx context.Context, chatID string) (string, error)
	CreateQueue(ctx context.Context, courier *models.Queue) error
	InsertCourierCoordinatesTable(ctx context.Context, courierCoordinatesTable *models.CourierCoordinatesTable) error
}

func (p *Pg) GetCourierByUUID(ctx context.Context, uuid string) (*models.Courier, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	courier := new(models.Courier)
	err := p.Db.ModelContext(timeout, courier).
		Where("uuid = ? AND deleted_at is null", uuid).
		Select()
	return courier, err
}

func (p *Pg) GetCourierQueue(ctx context.Context, status string) (*[]*models.Courier, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	courier := new([]*models.Courier)
	err := p.Db.ModelContext(timeout, courier).
		Where("status = ?", status).
		Select()
	return courier, err
}

func (p *Pg) GetCourierByChatID(ctx context.Context, ChatID string, phoneNumber string) (bool, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	courier := new(models.Courier)
	flag, err := p.Db.ModelContext(timeout, courier).
		Where("chat_id = ? AND phone_number=?", ChatID, phoneNumber).
		Exists()
	return flag, err
}

func (p *Pg) CouriersFilter(ctx context.Context, filter *proto.CouriersFilter) ([]models.Courier, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	courier := new([]models.Courier)
	query := p.Db.ModelContext(timeout, courier)
	if filter.CourierLat != 0 {
		query.Where("last_lat=?", filter.CourierLat)
	}
	if filter.CourierLon != 0 {
		query.Where("last_lon = ?", filter.CourierLon)
	}
	if filter.CourierStatus != "" {
		query.Where("status = ?", filter.CourierStatus)
	}
	if err := query.Select(); err != nil {
		return nil, err
	}
	return *courier, nil
}

func (p *Pg) CreateCourier(ctx context.Context, courier *models.Courier, event *models.Event) (*models.Courier, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier inserting",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, courier).
			Where("phone_number = ?", courier.PhoneNumber)

		if flag, err := query.Exists(); flag {
			return errors.New("this phone number is already exist in DB")
		} else {
			if err != nil {
				return errors.New("fail to validate phone number before insert courier")
			}
		}

		if _, err := tx.ModelContext(timeout, courier).Insert(); err != nil {
			log.WithField("reason", "failed to create courier").Error(err)
			return errors.Wrap(err, "failed to create courier")
		}
		if _, err := tx.ModelContext(timeout, event).Insert(); err != nil {
			log.WithField("reason", "failed to create events").Error(err)
			return errors.Wrap(err, "failed to create events")
		}
		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return courier, errors.Wrap(err, "transaction failed")
	}
	return courier, nil
}

func (p *Pg) UpdateCourier(ctx context.Context, courier *models.Courier, event *models.Event) (*models.Courier, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	err := p.CheckExistsByUUID(ctx, courier, courier.UUID)
	if err != nil {
		return nil, err
	}

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier updating",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, courier).
			Where("uuid = ?", courier.UUID)

		if _, err := query.UpdateNotNull(); err != nil {
			log.WithField("reason", "failed to update courier").Error(err)
			return errors.Wrap(err, "failed to update courier")
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

	return courier, nil
}

func (p *Pg) DeleteCourier(ctx context.Context, uuid string, event *models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()
	err := p.CheckExistsByUUID(ctx, &models.Courier{}, uuid)
	if err != nil {
		return err
	}

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier updating",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Courier{}).
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

func (p *Pg) UpdateCourierCoordinates(ctx context.Context, coordinates *proto.CouriersFilter, courCoordsTable *models.CourierCoordinatesTable) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "courier inserting",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.Courier{}).
			Set("last_lon = ?, last_lat = ?", coordinates.CourierLon, coordinates.CourierLat).
			Where("uuid = ?", coordinates.CourierUUID)

		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to update couriers coordinates").Error(err)
			return errors.Wrap(err, "failed to update couriers coordinates")
		}

		query = tx.ModelContext(timeout, courCoordsTable)
		if _, err := query.Insert(); err != nil {
			log.WithField("reason", "failed to insert last courier coordinates").Error(err)
			return errors.Wrap(err, "failed to insert last courier coordinates")
		}

		return nil
	}); err != nil {
		log.WithField("reason", "transaction failed").Error(err)
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}

func (p *Pg) UpdateCourierCoordinatesTable(ctx context.Context, coordinates *proto.CouriersFilter) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err := p.Db.ModelContext(timeout, &models.CourierCoordinatesTable{}).
		Set("lat = ?, lon = ?", coordinates.CourierLat, coordinates.CourierLon).
		Where("courier_uuid = ?", coordinates.CourierUUID).
		Update()
	return err
}

func (p *Pg) GetCourierByPhoneNumber(ctx context.Context, phoneNumber string) (*models.Courier, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	courier := new(models.Courier)
	err := p.Db.ModelContext(timeout, courier).
		Where("phone_number = ? AND deleted_at is null", phoneNumber, nil).
		Select()
	return courier, err
}

func (p *Pg) GetCourierUUIDByChatID(ctx context.Context, chatID string) (string, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	courier := new(models.Courier)
	err := p.Db.ModelContext(timeout, courier).
		Where("chat_id = ?", chatID).
		Select()
	return courier.UUID, err

}

func (p *Pg) CreateQueue(ctx context.Context, queue *models.Queue) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err := p.Db.ModelContext(timeout, queue).Insert()
	return err
}

func (p *Pg) InsertCourierCoordinatesTable(ctx context.Context, courierCoordinatesTable *models.CourierCoordinatesTable) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err := p.Db.ModelContext(timeout, courierCoordinatesTable).Insert()
	return err
}
