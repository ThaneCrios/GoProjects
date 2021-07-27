package handler

import (
	"gitlab.com/faemproject/backend/delivery/collector.core/ram"
	"gitlab.com/faemproject/backend/delivery/collector.core/repository"
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

type Handler struct {
	DB  Repository
	Pub Publisher
	RAM ram.RAM
}
