package handler

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

func (h *Handler) CreateUser(ctx context.Context, user *models.User) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "user creating",
	})
	event := h.EventCreating("create user", "done", "", "")
	user.UUID = h.RAM.IDs.GenUUID()
	err := h.DB.CreateUser(ctx, user, &event)
	if err != nil {
		log.WithField("reason", "failed to create user in DB").Error(err)
		return errors.Wrap(err, "fail to create user")
	}

	logs.Eloger.Info("user created")

	return nil
}

func (h *Handler) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	user, err := h.DB.GetUserByLogin(ctx, login)
	if err != nil {
		logs.Eloger.WithError(err).
			WithFields(logrus.Fields{
				"event": "getting user login",
				"uuid":  login,
			}).Error()
		return nil, errors.Wrap(err, "fail to get user by login")
	}
	return user, nil
}

func (h *Handler) SetUserState(ctx context.Context, params *proto.UsersFilter) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "set user state",
	})
	event := h.EventCreating("set user state", "done", "", "")
	err := h.DB.SetUserState(ctx, params, &event)
	if err != nil {
		log.WithField("reason", "failed to set user state in DB").Error(err)
		return errors.Wrap(err, "fail to set user state")
	}

	logs.Eloger.Info("state updated")

	return nil
}

func (h *Handler) MarkUserAsDeleted(ctx context.Context, userId string) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "deleting user",
	})
	event := h.EventCreating("delete user", "done", "", "")
	err := h.DB.MarkUserAsDeleted(ctx, userId, &event)
	if err != nil {
		log.WithField("reason", "failed to create user in DB").Error(err)
		return errors.Wrap(err, "fail to create user")
	}
	logs.Eloger.Info("user deleted")

	return nil
}
