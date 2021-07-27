package proto

import (
	"gitlab.com/faemproject/backend/delivery/delivery.front/models"
	"time"
)

type TasksFilter struct {
	TaskUUID    string `json:"task_uuid"`
	CourierUUID string `json:"courier_uuid"`
	TaskType    string `json:"task_type"`
	TaskState   string `json:"task_state"`
	Status      string `json:"status"`
}

type QueueTasks struct {
	CourierUUID    string `json:"courier_uuid"`
	FirstTaskUUID  string `json:"first_task_uuid"`
	SecondTaskUUID string `json:"second_task_uuid"`
}

type Task struct {
	UUID           string       `json:"uuid"`
	OrderNumber    string       `json:"order_number"`
	ClientPhone    string       `json:"client_phone"`
	ClientName     string       `json:"client_name"`
	OrderUUID      string       `json:"order_uuid"`
	Type           string       `json:"type"`
	Route          models.Route `json:"route"`
	State          string       `json:"state"`
	CreatedAt      time.Time    `json:"created_at"`
	ExpectedTime   time.Time    `json:"expected_time"`
	LastUpdateTime time.Time    `json:"last_update_time"`
	FinishTime     string       `json:"finish_time"`
	CourierUUID    string       `json:"courier_uuid"`
}
