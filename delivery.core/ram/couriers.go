package ram

import (
	"context"
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
	"gitlab.com/faemproject/backend/delivery/delivery.core/proto"
)

type CouriersRam interface {
	GetCourierByUUID(ctx context.Context, uuid string) (*models.Courier, error)
	GetCourierQueue(ctx context.Context, status string) (*[]*models.Courier, error)
	GetCouriersFromDb(ctx context.Context, courier models.Courier) error
	UpdateCourier(ctx context.Context, uuid string, city *models.Courier) error
	DeleteCourier(ctx context.Context, uuid string) error
	UpdateStatusCourier(ctx context.Context, value *proto.ChangeStat, courier *models.Courier) error
	CouriersFilter(ctx context.Context, filter *proto.CouriersFilter) (*[]*models.Courier, error)
	GetCourierByChatID(ctx context.Context, chatID string, phoneNumber string) (bool, error)
	CreateTasks(ctx context.Context, order models.Order) ([]models.Task, error)
}
