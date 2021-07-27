package models

import "time"

type (
	User struct {
		tableName interface{} `pg:"ord_users"`
		UUID      string      `json:"uuid"`
		Name      string      `json:"name"`
		Login     string      `json:"login"`
		Deleted   bool        `json:"-"`
		Blocked   bool        `json:"blocked"`
		Meta      UserMeta    `json:"meta"`
		CreatedAt time.Time   `json:"-"`
		UpdatedAt time.Time   `json:"-"`
	}

	UserMeta struct {
		StoresUUID []string `json:"partners_list"`
	}
)
