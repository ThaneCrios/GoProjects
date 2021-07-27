package proto

import (
	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
)

type CouriersFilter struct {
	CourierUUID   string  `json:"courier_uuid"`
	CourierLat    float64 `json:"courier_lat"`
	CourierLon    float64 `json:"courier_lon"`
	CourierStatus string  `json:"courier_status"`
}

type CourierFront struct {
	UUID                  string      `json:"uuid"`         //идентификатор курьера
	CourierMeta           CourierMeta `json:"courier_meta"` //мета-данные курьера(ФИО, дата рождения, реквизиты и пр.)
	CourierType           string      `json:"courier_type"` //тип курьера(на машине, пеший и пр.)
	PhoneNumber           string      `json:"phone_number"` //номер телефона
	StatusBack            string      `json:"status_back"`  //текущий статус курьера(ожидает заказ, на заказе, выходной и пр.)
	StatusFront           string      `json:"status_front"`
	LastLat               float64     `json:"last_lat"`                //последняя координата широты курьера
	LastLon               float64     `json:"last_lon"`                //последняя координата долготы курьера
	StatusChangeTimestamp int64       `json:"status_change_timestamp"` //время последнего измения статуса курьера
	LatlonTimestamp       int64       `json:"latlon_timestamp"`        //время последнего изменения координат курьера
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

//FromFrontToCore преобразует модель из фронта в модель из кор
func CourierFromFrontToCore(courFr CourierFront) models.Courier {
	return models.Courier{
		UUID: courFr.UUID,
		CourierMeta: models.CourierMeta{
			FullName:   courFr.CourierMeta.FullName,
			Requisites: courFr.CourierMeta.Requisites,
			BirthDate:  courFr.CourierMeta.BirthDate,
			Params: models.Params{
				ThermalBag: courFr.CourierMeta.Params.ThermalBag,
				Fridge:     courFr.CourierMeta.Params.Fridge,
			},
		},
		CourierType: courFr.CourierType,
		PhoneNumber: courFr.PhoneNumber,
		Status:      courFr.StatusBack,
	}
}

//FromCoreToFront ...
func CourierFromCoreToFront(courCore models.Courier) CourierFront {
	return CourierFront{
		UUID: courCore.UUID,
		CourierMeta: CourierMeta{
			FullName:   courCore.CourierMeta.FullName,
			ShortName:  courCore.CourierMeta.ShortName,
			Requisites: courCore.CourierMeta.Requisites,
			BirthDate:  courCore.CourierMeta.BirthDate,
			Params: Params{
				ThermalBag: courCore.CourierMeta.Params.ThermalBag,
				Fridge:     courCore.CourierMeta.Params.Fridge,
			},
		},
		CourierType:           courCore.CourierType,
		PhoneNumber:           courCore.PhoneNumber,
		StatusBack:            courCore.Status,
		StatusFront:           Variable.Couriers[courCore.Status],
		LastLat:               courCore.LastLat,
		LastLon:               courCore.LastLon,
		StatusChangeTimestamp: courCore.StatusChangeTimestamp.Unix(),
		LatlonTimestamp:       courCore.LatlonTimestamp.Unix(),
	}
}
