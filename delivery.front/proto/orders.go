package proto

type OrdersFilter struct {
	CourierUUID string  `json:"courier_uuid"`
	OrderLat    float64 `json:"order_lat"`
	OrderLon    float64 `json:"order_lon"`
	OrderStatus string  `json:"order_status"`
}

type OrderCourier struct {
	OrderUUID   string `json:"order_uuid"`
	CourierUUID string `json:"courier_uuid"`
}

type OrderCreate struct {
	DeliveryItems       OrderItems `json:"delivery_items"`        //состав заказа, цена, вес и пр.
	Comment             string     `json:"comment"`               //комментарии по поводу доставки(в определённом положении, быстро и пр.)
	PickupRequestTime   string     `json:"pickup_request_time"`   //дэдлайн на забор(?) заказа
	DropOffRequestTime  string     `json:"drop_off_request_time"` //дэдлайн на доставку заказа
	PickupContactPhone  string     `json:"pickup_contact_phone"`  //номер телефона точки забора(?) заказа
	DropoffContactPhone string     `json:"dropoff_contact_phone"` //номер телефона точки доставки заказа
	PickupRoute         Route      `json:"pickup_route"`          //роут забора(?) заказа
	DropoffRoute        Route      `json:"dropoff_route"`         //роут доставки заказа
	Service             string     `json:"service"`               //описание алгоритма доставки(забрать, отвезти; забрать, оплатить, отвезти и пр.)
	PaymentType         string     `json:"payment_type"`          //способ оплаты(наличные, карта и пр.)
	PaymentStatus       string     `json:"payment_status"`        //статус оплаты(не оплачен, оплачен, ожидает оплаты и пр.)
	DeliveryPrice       float32    `json:"delivery_price"`
}

type Route struct {
	UUID              string  `json:"uuid"`
	PointType         string  `json:"point_type"`
	UnrestrictedValue string  `json:"unrestricted_value"`
	ValueForSearch    string  `json:"-"`
	Value             string  `json:"value"`
	Country           string  `json:"country"`
	Region            string  `json:"region"`
	RegionType        string  `json:"region_type"`
	Type              string  `json:"type"`
	City              string  `json:"city"`
	Category          string  `json:"category"`
	CityType          string  `json:"city_type"`
	Street            string  `json:"street"`
	StreetType        string  `json:"street_type"`
	StreetWithType    string  `json:"street_with_type"`
	House             string  `json:"house"`
	FrontDoor         int     `json:"front_door"`
	Comment           string  `json:"comment"`
	OutOfTown         bool    `json:"out_of_town"`
	HouseType         string  `json:"house_type"`
	AccuracyLevel     int     `json:"accuracy_level"`
	Radius            int     `json:"radius"`
	Lat               float32 `json:"lat"`
	Lon               float32 `json:"lon"`
}

//OrderItems состав заказа
type OrderItems struct {
	Items      []DeliveryItem `json:"items"`       //предметы доставки(наименования продкутов, товаров и пр.)
	ItemsPrice float64        `json:"items_price"` //цена заказа
	Weight     int            `json:"weight"`      //вес заказа
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
