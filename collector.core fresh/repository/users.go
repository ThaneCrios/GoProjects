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

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User, event *models.Event) error
	GetUserByLogin(ctx context.Context, uuid string) (*models.User, error)
	SetUserState(ctx context.Context, params *proto.UsersFilter, event *models.Event) error
	MarkUserAsDeleted(ctx context.Context, userId string, event *models.Event) error
}

func (p *Pg) CreateUser(ctx context.Context, user *models.User, event *models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "create user",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		if _, err := tx.ModelContext(timeout, user).Insert(); err != nil {
			log.WithField("reason", "failed to create user").Error(err)
			return errors.Wrap(err, "failed to update courier")
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

func (p *Pg) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	user := new(models.User)
	err := p.Db.ModelContext(timeout, user).
		Where("uuid = ? AND deleted = ?", login, false).
		Select()
	return user, err
}

func (p *Pg) SetUserState(ctx context.Context, params *proto.UsersFilter, event *models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "create user",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.User{}).
			Set("state = ?, updated_at = ?", params.UserState, time.Now()).
			Where("uuid = ?", params.UserUUID)

		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to set user state").Error(err)
			return errors.Wrap(err, "failed to set user state")
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

func (p *Pg) MarkUserAsDeleted(ctx context.Context, userId string, event *models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "delete user",
	})

	if err := p.Db.RunInTransaction(func(tx *pg.Tx) error {
		query := tx.ModelContext(timeout, &models.User{}).
			Set("deleted = ?, updated_at = ?, deleted_at = ?", true, time.Now(), time.Now()).
			Where("uuid = ?", userId)

		if _, err := query.Update(); err != nil {
			log.WithField("reason", "failed to delete user").Error(err)
			return errors.Wrap(err, "failed to delete user")
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
