package handler

import (
	"gitlab.com/faemproject/backend/delivery/delivery.core/ram"
	"gitlab.com/faemproject/backend/delivery/delivery.core/repository"
)

type Repository interface {
	repository.VersionerRepository
	repository.OrdersRepository
	repository.CouriersRepository
	repository.TasksRepository
	repository.EventsRepository
	repository.UserRepository
}

type Memory interface {
	ram.CouriersRam
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
