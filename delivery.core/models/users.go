package models

type User struct {
	tableName struct{} `sql:"delivery_users"`
	UUID      string   `json:"uuid"`
	Login     string   `json:"login"`
	State     string   `json:"state"`
	Meta      UserMeta `json:"meta"`
	Deleted   bool     `json:"deleted"`
}

type UserMeta struct {
}
