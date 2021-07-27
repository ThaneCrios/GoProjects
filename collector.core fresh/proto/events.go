package proto

import (
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	"time"
)

type EventFilter struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	EventType string    `json:"event_type"`
}

type EventFront struct {
	UUID        string  `json:"uuid"`         //идентификатор эвента
	CreatedAt   int64   `json:"created_at"`   //время создания эвента
	EventType   string  `json:"event_type"`   //тип эвента
	Payload     Payload `json:"payload"`      //информация о ошибках в процессе(либо статус успешного выполнения) эвента
	CourierUUID string  `json:"courier_uuid"` //идентификатор курьера(если курьер был задействован в эвенте)
	OrderUUID   string  `json:"order_uuid"`   //идентификатор заказа(если заказ был задействован в эвенте)
}

type Payload struct {
	Information string
}

func EventFromCoreToFront(event models.Event) EventFront {
	return EventFront{
		UUID:      event.UUID,
		CreatedAt: event.CreatedAt.Unix(),
		EventType: event.EventType,
		Payload: Payload{
			Information: event.Payload.Information},
		CourierUUID: event.CourierUUID,
		OrderUUID:   event.OrderUUID,
	}
}
