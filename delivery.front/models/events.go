package models

import (
	"time"
)

type Event struct {
	TableName   struct{}  `sql:"delivery_events"`
	UUID        string    `json:"uuid"`
	CreatedAt   time.Time `json:"created_at"`
	EventType   string    `json:"event_type"`
	Payload     Payload   `json:"payload"`
	CourierUUID string    `json:"courier_uuid"`
	OrderUUID   string    `json:"order_uuid"`
}

type Payload struct {
	Information string
}
