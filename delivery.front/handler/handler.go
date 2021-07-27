package handler

import (
	"gitlab.com/faemproject/backend/delivery/delivery.front/config"
)

type Repository interface {
}

type Publisher interface {
}

type Handler struct {
	DB     Repository
	Pub    Publisher
	Config *config.Config
}
