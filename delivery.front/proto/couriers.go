package proto

import (
	"time"
)

type CouriersFilter struct {
	CourierUUID   string  `json:"courier_uuid"`
	CourierLat    float64 `json:"courier_lat"`
	CourierLon    float64 `json:"courier_lon"`
	CourierStatus string  `json:"courier_status"`
}

type CreateCourier struct {
	CourierMeta CourierMeta `json:"courier_meta"` //мета-данные курьера(ФИО, дата рождения, реквизиты и пр.)
	PhoneNumber string      `json:"phone_number"` //номер телефона
}

type CourierFront struct {
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

func (crCour *CreateCourier) CourierConvert() CourierFront {
	return CourierFront{
		CourierMeta: crCour.CourierMeta,
		PhoneNumber: crCour.PhoneNumber,
	}
}
