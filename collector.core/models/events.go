package models

import (
	"time"
)

type Event struct {
	TableName   struct{}  `sql:"delivery_events"`
	UUID        string    `json:"uuid"`         //идентификатор эвента
	CreatedAt   time.Time `json:"created_at"`   //время создания эвента
	EventType   string    `json:"event_type"`   //тип эвента
	Payload     Payload   `json:"payload"`      //информация о ошибках в процессе(либо статус успешного выполнения) эвента
	CourierUUID string    `json:"courier_uuid"` //идентификатор курьера(если курьер был задействован в эвенте)
	OrderUUID   string    `json:"order_uuid"`   //идентификатор заказа(если заказ был задействован в эвенте)
}

type Payload struct {
	Information string
}
