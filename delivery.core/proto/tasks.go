package proto

import (
	"time"

	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
)

type TasksFilter struct {
	TaskUUID       string `json:"task_uuid"`
	CourierUUID    string `json:"courier_uuid"`
	TaskType       string `json:"task_type"`
	TaskState      string `json:"task_state"`
	Status         string `json:"status"`
	ParentTaskUUID string `json:"parent_task_uuid"`
}

type QueueTasks struct {
	CourierUUID    string `json:"courier_uuid"`
	FirstTaskUUID  string `json:"first_task_uuid"`
	SecondTaskUUID string `json:"second_task_uuid"`
}

type TaskFront struct {
	UUID           string       `json:"uuid"`
	OrderNumber    string       `json:"order_number"`
	ClientData     ClientData   `json:"client_data`
	OrderUUID      string       `json:"order_uuid"`
	Type           string       `json:"type"`
	TypeHuman      string       `json:"type_human"`
	Route          models.Route `json:"route"`
	State          string       `json:"state"`
	StateHuman     string       `json:"state_human"`
	CreatedAt      int64        `json:"created_at,omitempty"`
	ExpectedTime   int64        `json:"expected_time,omitempty"`
	LastUpdateTime int64        `json:"last_update_time,omitempty"`
	FinishTime     int64        `json:"finish_time,omitempty"`
	CourierUUID    string       `json:"courier_uuid"`
}

type ClientData struct {
	ClientPhone string `json:"client_phone"`
	ClientName  string `json:"client_name"`
}

func TaskFromCoreToFront(task models.Task) TaskFront {
	return TaskFront{
		UUID:        task.UUID,
		OrderNumber: task.OrderNumber,
		ClientData: ClientData{
			ClientPhone: task.ClientData.ClientPhone,
			ClientName:  task.ClientData.ClientName,
		},
		OrderUUID: task.OrderUUID,
		Type:      task.Type,
		TypeHuman: Variable.Tasks[task.Type],
		Route: models.Route{
			UUID:              task.Route.UUID,
			PointType:         task.Route.PointType,
			UnrestrictedValue: task.Route.UnrestrictedValue,
			ValueForSearch:    task.Route.ValueForSearch,
			Value:             task.Route.Value,
			Country:           task.Route.Country,
			Region:            task.Route.Region,
			RegionType:        task.Route.RegionType,
			Type:              task.Route.Type,
			City:              task.Route.City,
			Category:          task.Route.Category,
			CityType:          task.Route.CityType,
			Street:            task.Route.Street,
			StreetType:        task.Route.StreetType,
			StreetWithType:    task.Route.StreetWithType,
			House:             task.Route.House,
			FrontDoor:         task.Route.FrontDoor,
			Comment:           task.Route.Comment,
			OutOfTown:         task.Route.OutOfTown,
			HouseType:         task.Route.HouseType,
			AccuracyLevel:     task.Route.AccuracyLevel,
			Radius:            task.Route.Radius,
			Lat:               task.Route.Lat,
			Lon:               task.Route.Lon,
		},
		State:       task.State,
		StateHuman:  Variable.Tasks[task.State],
		CourierUUID: task.CourierUUID,
	}
}

func TaskFromFrontToCore(task TaskFront) models.Task {
	return models.Task{
		UUID:        task.UUID,
		OrderNumber: task.OrderNumber,
		ClientData: models.ClientData{
			ClientPhone: task.ClientData.ClientPhone,
			ClientName:  task.ClientData.ClientName,
		},
		OrderUUID: task.OrderUUID,
		Type:      task.Type,
		Route: models.Route{
			UUID:              task.Route.UUID,
			PointType:         task.Route.PointType,
			UnrestrictedValue: task.Route.UnrestrictedValue,
			ValueForSearch:    task.Route.ValueForSearch,
			Value:             task.Route.Value,
			Country:           task.Route.Country,
			Region:            task.Route.Region,
			RegionType:        task.Route.RegionType,
			Type:              task.Route.Type,
			City:              task.Route.City,
			Category:          task.Route.Category,
			CityType:          task.Route.CityType,
			Street:            task.Route.Street,
			StreetType:        task.Route.StreetType,
			StreetWithType:    task.Route.StreetWithType,
			House:             task.Route.House,
			FrontDoor:         task.Route.FrontDoor,
			Comment:           task.Route.Comment,
			OutOfTown:         task.Route.OutOfTown,
			HouseType:         task.Route.HouseType,
			AccuracyLevel:     task.Route.AccuracyLevel,
			Radius:            task.Route.Radius,
			Lat:               task.Route.Lat,
			Lon:               task.Route.Lon,
		},
		State:          task.State,
		CreatedAt:      time.Unix(task.CreatedAt, 0),
		ExpectedTime:   time.Unix(task.ExpectedTime, 0),
		LastUpdateTime: time.Unix(task.LastUpdateTime, 0),
		FinishTime:     time.Unix(task.FinishTime, 0),
		CourierUUID:    task.CourierUUID,
	}
}
