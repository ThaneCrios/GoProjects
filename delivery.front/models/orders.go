package models

import "time"

type Order struct {
	TableName           struct{}    `sql:"delivery_orders"`
	UUID                string      `json:"uuid"`                  //идентификатор
	OrderNumber         string      `json:"order_number"`          //номер заказа(первые 5 символов идентификатора)
	DeliveryItems       OrderItems  `json:"delivery_items"`        //состав заказа, цена, вес и пр.
	Comment             string      `json:"comment"`               //комментарии по поводу доставки(в определённом положении, быстро и пр.)
	PickupRequestTime   string      `json:"pickup_request_time"`   //дэдлайн на забор(?) заказа
	DropOffRequestTime  string      `json:"drop_off_request_time"` //дэдлайн на доставку заказа
	PickupContactPhone  string      `json:"pickup_contact_phone"`  //номер телефона точки забора(?) заказа
	DropoffContactPhone string      `json:"dropoff_contact_phone"` //номер телефона точки доставки заказа
	PickupRoute         Route       `json:"pickup_route"`          //роут забора(?) заказа
	DropoffRoute        Route       `json:"dropoff_route"`         //роут доставки заказа
	Service             string      `json:"service"`               //описание алгоритма доставки(забрать, отвезти; забрать, оплатить, отвезти и пр.)
	CreatedAt           time.Time   `json:"created_at"`            //время создания заказа
	UpdatedAt           time.Time   `json:"updated_at"`            //время последнего обновления заказа
	FinishedAt          time.Time   `json:"finished_at"`           //время завершения заказа
	Canceled            bool        `json:"canceled"`              //
	CancelReason        string      `json:"cancel_reason"`         //причина отмены заказа
	PaymentType         string      `json:"payment_type"`          //способ оплаты(наличные, карта и пр.)
	PaymentStatus       string      `json:"payment_status"`        //статус оплаты(не оплачен, оплачен, ожидает оплаты и пр.)
	State               string      `json:"state"`                 //текущее состояние заказа(ожидает назначения, доставляется и пр.)
	DeletedAt           time.Time   `json:"deleted_at"`            //время "удаления" заказа
	DeliveryPrice       float32     `json:"delivery_price"`
	CourierUUID         string      `json:"courier_uuid"` //идентификатор курьера, назначенного на заказ
	CourierData         CourierData `json:"courier_data"` //мета-данные курьера, назначенного на заказ
}
type OrderItems struct {
	Item string `json:"item"`
}

type CourierData struct {
	Name string `json:"name"`
}
type OrderTypesConst struct {
	PickUpAndDeliver                   string `json:"pick_up_and_driver"`
	PickUpPayAndDeliver                string `json:"pick_up_pay_and_driver"`
	PickUpPayAndDeliverWaitDeliverBack string `json:"pick_up_pay_and_driver_wait_deliver_back"`
}

var OrderTypes = OrderTypesConst{
	PickUpAndDeliver:                   "забрать и доставить",
	PickUpPayAndDeliver:                "забрать, заплатить и доставить",
	PickUpPayAndDeliverWaitDeliverBack: "забрать, заплатить, доставить, подождать, отвезти обратно",
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
