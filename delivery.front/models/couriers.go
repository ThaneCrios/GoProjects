package models

import (
	"time"
)

type Courier struct {
	tableName             struct{}    `sql:"delivery_couriers"`
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
	Deleted               bool        `json:"deleted"`
}

type AuthTokenDate struct {
	CourierUUID string   `json:"courier_uuid"`
	Role        string   `json:"role"`
	MetaInfo    struct{} `json:"meta_info"`
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

//func (cour *Courier) CourierConvertToCore() proto.CourierFront {
//	return proto.CourierFront{
//		UUID: cour.UUID,
//		CourierMeta: proto.CourierMeta{
//			FullName:   cour.CourierMeta.FullName,
//			ShortName:  cour.CourierMeta.ShortName,
//			Requisites: cour.CourierMeta.Requisites,
//			BirthDate:  cour.CourierMeta.BirthDate,
//			Params: proto.Params{
//				ThermalBag: cour.CourierMeta.Params.ThermalBag,
//				Fridge:     cour.CourierMeta.Params.Fridge,
//			},
//		},
//		CourierType:           cour.CourierType,
//		PhoneNumber:           cour.PhoneNumber,
//		Status:                proto.Localisation[cour.Status],
//		LastLat:               cour.LastLat,
//		LastLon:               cour.LastLon,
//		StatusChangeTimestamp: cour.StatusChangeTimestamp,
//		LatlonTimestamp:       cour.LatlonTimestamp,
//		CreatedAt:             cour.CreatedAt,
//		UpdatedAt:             cour.UpdatedAt,
//		DeletedAt:             cour.DeletedAt,
//	}
//}
