package models

import "time"

type Order struct {
	TableName             struct{}    `sql:"delivery_orders"`
	UUID                  string      `json:"uuid"`                            //идентификатор
	OrderNumber           string      `json:"order_number"`                    //номер заказа(первые 5 символов идентификатора)
	DeliveryItems         OrderItems  `json:"delivery_items"`                  //состав заказа, цена, вес и пр.
	Comment               string      `json:"comment"`                         //комментарии по поводу доставки(в определённом положении, быстро и пр.)
	PickupRequestTime     time.Time   `json:"pickup_request_time,omitempty"`   //дэдлайн на забор(?) заказа
	DropOffRequestTime    time.Time   `json:"drop_off_request_time,omitempty"` //дэдлайн на доставку заказа
	PickupPersonContacts  PersonsData `json:"pickup_person_contacts"`          //номер телефона точки забора(?) заказа
	DropOffPersonContacts PersonsData `json:"drop_off_person_contacts"`        //номер телефона точки доставки заказа
	PickupRoute           Route       `json:"pickup_route"`                    //роут забора(?) заказа
	DropoffRoute          Route       `json:"dropoff_route"`                   //роут доставки заказа
	Service               string      `json:"service"`                         //описание алгоритма доставки(забрать, отвезти; забрать, оплатить, отвезти и пр.)
	CreatedAt             time.Time   `json:"created_at,omitempty"`            //время создания заказа
	UpdatedAt             time.Time   `json:"updated_at,omitempty"`            //время последнего обновления заказа
	FinishedAt            time.Time   `json:"finished_at,omitempty"`           //время завершения заказа
	CanceledAt            time.Time   `json:"canceled_at,omitempty"`           //
	CancelReason          string      `json:"cancel_reason"`                   //причина отмены заказа
	PaymentType           string      `json:"payment_type"`                    //способ оплаты(наличные, карта и пр.)
	PaymentStatus         string      `json:"payment_status"`                  //статус оплаты(не оплачен, оплачен, ожидает оплаты и пр.)
	State                 string      `json:"state"`                           //текущее состояние заказа(ожидает назначения, доставляется и пр.)
	DeletedAt             time.Time   `json:"deleted_at,omitempty"`            //время "удаления" заказа
	DeliveryPrice         float32     `json:"delivery_price"`
	CourierUUID           string      `json:"courier_uuid"` //идентификатор курьера, назначенного на заказ
	CourierData           CourierData `json:"courier_data"` //мета-данные курьера, назначенного на заказ
}

type PersonsData struct {
	Phone   string
	Name    string
	Comment string
}

//OrderItems состав заказа
type OrderItems struct {
	Items []DeliveryItem `json:"items"` //предметы доставки(наименования продкутов, товаров и пр.)
}

type DeliveryItem struct {
	Name  string  `json:"name"`  //наименование предмета доставки
	Price float64 `json:"price"` //цена предмета доставки
	Count int     `json:"count"` //количество
	SKU   string  `json:"sku"`   //идентификатор товарной позиции
}

//CourierData - информация по курьру
type CourierData struct {
	UUID        string `json:"uuid"`         //идентификатор курьера
	Name        string `json:"name"`         //ФИО
	PhoneNumber string `json:"phone_number"` //номер телефона
}

type OrderTest struct {
	TableName          struct{}  `sql:"delivery_full_order"`
	UUID               string    `json:"uuid"`
	Comment            string    `json:"comment"`
	PickupRequestTime  string    `json:"pickup_request_time"`
	PickupContactPhone string    `json:"pickup_contact_phone"`
	Status             string    `json:"status"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	FinishedAt         time.Time `json:"finished_at"`
	CancelReason       string    `json:"cancel_reason"`
	CancelComment      string    `json:"cancel_comment"`
	PaymentType        string    `json:"payment_type"`
	PaymentStatus      string    `json:"payment_status"`
	DeliveryType       string    `json:"delivery_type"`
	Deleted            bool      `json:"deleted"`
	DeliveryPrice      float64   `json:"delivery_price"`
	CourierUUID        string    `json:"courier_uuid"`
	OrderType          string    `json:"order_type"`
}
