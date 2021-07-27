package models

type Task struct {
	TableName         struct{} `sql:"delivery_tasks"`
	UUID              string   `json:"uuid"`
	OrderNumber       string   `json:"order_number"`
	PhoneNumber       string   `json:"phone_number"`
	OrderUUID         string   `json:"order_uuid"`
	Type              string   `json:"type"`
	Route             Route    `json:"route"`
	State             string   `json:"state"`
	ParentTaskUUID    string   `json:"parent_task_uuid"`
	ProdComplTimeFrom string   `json:"prod_compl_time_from"`
	ProdComplTimeTo   string   `json:"prod_compl_time_to"`
	FinishedTime      string   `json:"finish_time"`
	CourierUUID       string   `json:"courier_uuid"`
}

type MetaTask struct {
}

type TaskType struct {
	PickUp                 string
	Deliver                string
	PickUpAndPay           string
	DeliverWaitDeliverBack string
}

var Type = TaskType{
	PickUp:                 "Забрать",
	Deliver:                "Отвезти",
	PickUpAndPay:           "Забрать и заплотить",
	DeliverWaitDeliverBack: "Доставить, подождать и отвезти обратно",
}
