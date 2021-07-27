package models

import "time"

type Task struct {
	TableName      struct{}   `sql:"delivery_tasks"`
	UUID           string     `json:"uuid"`
	OrderNumber    string     `json:"order_number"`
	ClientData     ClientData `json:"client_data"`
	OrderUUID      string     `json:"order_uuid"`
	Type           string     `json:"type"`
	Route          Route      `json:"route"`
	State          string     `json:"state"`
	CreatedAt      time.Time  `json:"created_at"`
	ExpectedTime   time.Time  `json:"expected_time"`
	LastUpdateTime time.Time  `json:"last_update_time"`
	FinishTime     time.Time  `json:"finish_time"`
	CourierUUID    string     `json:"courier_uuid"`
}

type ClientData struct {
	ClientPhone string `json:"client_phone"`
	ClientName  string `json:"client_name"`
}
