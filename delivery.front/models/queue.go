package models

/*type TasksQueue struct {
	UUID 		string 		`json:"uuid"`
	Address		string		`json:"address"`
	Action 		string		`json:"action"`
	TaskUUID	string		`json:"task_uuid"`
}*/

type Queue struct {
	TableName   struct{} `sql:"delivery_courier_queue"`
	UUID        string   `json:"uuid"`
	CourierUUID string   `json:"courier_uuid"`
	Tasks       []Task   `json:"tasks"`
}
