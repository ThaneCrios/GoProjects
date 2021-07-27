package handler

import (
	"context"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/config"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/functions"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/repository"
	orderModels "gitlab.com/faemproject/backend/eda/eda.core/services/orders/proto"
)

type Repository interface {
	repository.VersionerRepository
	repository.OrdersRepository
	repository.CollectorsRepository
	repository.EventsRepository
	repository.UserRepository
	repository.ProductsRepository
}

type Publisher interface {
	BrokerPublisher
}

type BrokerPublisher interface {
}

type RPC interface {
	UsersRPC
	ProductsRPC
	OrdersRPC
}

type Handler struct {
	DB        Repository
	Pub       Publisher
	Functions functions.Functions
	Config    config.Config
	RPC       RPC
}

type UsersRPC interface {
	GetUserDataByUUID(ctx context.Context, uuid string) (*models.User, error)
}

type ProductsRPC interface {
	GetProductByUUID(ctx context.Context, uuid string) (*models.Product, error)
	GetProductUUIDByBarCode(ctx context.Context, params map[string]string) ([]models.ProductBarcode, error)
}

type OrdersRPC interface {
	SetOrderState(ctx context.Context, params orderModels.OrdersStateOptions, uuid string) error
	UpdateOrder(ctx context.Context, order models.OrderUpdated) error
}
