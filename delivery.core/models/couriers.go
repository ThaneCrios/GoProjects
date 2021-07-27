package models

import (
	"time"
)

type Courier struct {
	TableName             struct{}    `sql:"delivery_couriers"`
	UUID                  string      `json:"uuid"`                    //идентификатор курьера
	CourierMeta           CourierMeta `json:"courier_meta"`            //мета-данные курьера(ФИО, дата рождения, реквизиты и пр.)
	CourierType           string      `json:"courier_type"`            //тип курьера(на машине, пеший и пр.)
	ChatID                string      `json:"chat_id"`                 //идентификатор для telegram-бота
	PhoneNumber           string      `json:"phone_number"`            //номер телефона
	Status                string      `json:"status"`                  //текущий статус курьера(ожидает заказ, на заказе, выходной и пр.)
	LastLat               float64     `json:"last_lat"`                //последняя координата широты курьера
	LastLon               float64     `json:"last_lon"`                //последняя координата долготы курьера
	StatusChangeTimestamp time.Time   `json:"status_change_timestamp"` //время последнего измения статуса курьера
	LatlonTimestamp       time.Time   `json:"latlon_timestamp"`        //время последнего изменения координат курьера
	CreatedAt             time.Time   `json:"created_at"`              //время создания курьера
	UpdatedAt             time.Time   `json:"updated_at"`              //время последнего обновления курьера
	DeletedAt             time.Time   `json:"deleted_at"`              //время "удаления" курьера
}

type AuthTokenDate struct {
	CourierUUID string   `json:"courier_uuid"` //идентификатор курьера
	Role        string   `json:"role"`         //роль
	MetaInfo    struct{} `json:"meta_info"`    //мета-данные
}

type CourierMeta struct {
	FullName   string `json:"full_name"`  //Полное имя курьера
	ShortName  string `json:"short_name"` //
	Requisites string `json:"requisites"` //реквизиты
	BirthDate  string `json:"birthdate"`  //дата рождения
	Params     Params `json:"params"`     //различные параметры курьера(наличие определённого оборудования, навыков и пр.)
}

type Params struct {
	ThermalBag bool
	Fridge     bool
}

type CourierCoordinatesTable struct {
	tableName       struct{}  `sql:"delivery_courier_coordinates"`
	UUID            string    `json:"uuid"`
	CourierUUID     string    `json:"courier_uuid"`
	Lat             float64   `json:"lat"`
	Lon             float64   `json:"lon"`
	LatlonTimestamp time.Time `json:"latlon_timestamp"`
}

type CourierResponse struct {
	UUID                  string      `json:"uuid"`                    //идентификатор курьера
	CourierMeta           CourierMeta `json:"courier_meta"`            //мета-данные курьера(ФИО, дата рождения, реквизиты и пр.)
	CourierType           string      `json:"courier_type"`            //тип курьера(на машине, пеший и пр.)
	PhoneNumber           string      `json:"phone_number"`            //номер телефона
	Status                string      `json:"status"`                  //текущий статус курьера(ожидает заказ, на заказе, выходной и пр.)
	LastLat               float64     `json:"last_lat"`                //последняя координата широты курьера
	LastLon               float64     `json:"last_lon"`                //последняя координата долготы курьера
	StatusChangeTimestamp time.Time   `json:"status_change_timestamp"` //время последнего измения статуса курьера
	LatlonTimestamp       time.Time   `json:"latlon_timestamp"`        //время последнего изменения координат курьера
	CreatedAt             time.Time   `json:"created_at"`              //время создания курьера
	UpdatedAt             time.Time   `json:"updated_at"`              //время последнего обновления курьера
	DeletedAt             time.Time   `json:"deleted_at"`              //время "удаления" курьера
}

// func (courierResponse *Courier) CourierConvertForResponse() CourierResponse {
// 	return CourierResponse{
// 		UUID: courierResponse.UUID,
// 		CourierMeta: CourierMeta{
// 			FullName:   courierResponse.CourierMeta.FullName,
// 			ShortName:  courierResponse.CourierMeta.ShortName,
// 			Requisites: courierResponse.CourierMeta.Requisites,
// 			BirthDate:  courierResponse.CourierMeta.BirthDate,
// 			Params: Params{
// 				ThermalBag: courierResponse.CourierMeta.Params.ThermalBag,
// 				Fridge:     courierResponse.CourierMeta.Params.Fridge,
// 			},
// 		},
// 		CourierType:           courierResponse.CourierType,
// 		PhoneNumber:           courierResponse.PhoneNumber,
// 		Status:                courierResponse.Status,
// 		LastLat:               courierResponse.LastLat,
// 		LastLon:               courierResponse.LastLon,
// 		StatusChangeTimestamp: courierResponse.StatusChangeTimestamp,
// 		LatlonTimestamp:       courierResponse.LatlonTimestamp,
// 		CreatedAt:             courierResponse.CreatedAt,
// 		UpdatedAt:             courierResponse.UpdatedAt,
// 		DeletedAt:             courierResponse.DeletedAt,
// 	}
// }
