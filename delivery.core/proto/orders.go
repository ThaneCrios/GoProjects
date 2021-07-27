package proto

import (
	"time"

	"gitlab.com/faemproject/backend/delivery/delivery.core/models"
)

type OrdersFilter struct {
	OrderUUID    string  `json:"order_uuid"`
	CourierUUID  string  `json:"courier_uuid"`
	OrderLat     float64 `json:"order_lat"`
	OrderLon     float64 `json:"order_lon"`
	OrderService string  `json:"order_service"`
}

type OrderCourier struct {
	OrderUUID   string `json:"order_uuid"`
	CourierUUID string `json:"courier_uuid"`
}

type OrderFront struct {
	UUID                  string            `json:"uuid"`                   //идентификатор
	OrderNumber           string            `json:"order_number"`           //номер заказа(первые 5 символов идентификатора)
	DeliveryItems         models.OrderItems `json:"delivery_items"`         //состав заказа, цена, вес и пр.
	Comment               string            `json:"comment"`                //комментарии по поводу доставки(в определённом положении, быстро и пр.)
	PickupRequestTime     int64             `json:"pickup_request_time"`    //дэдлайн на забор(?) заказа
	DropOffRequestTime    int64             `json:"drop_off_request_time"`  //дэдлайн на доставку заказа
	PickupPersonContacts  PersonsData       `json:"pickup_person_contacts"` //номер телефона точки забора(?) заказа
	DropOffPersonContacts PersonsData       `json:"drop_off_person_contacts"`
	PickupRoute           models.Route      `json:"pickup_route"`  //роут забора(?) заказа
	DropoffRoute          models.Route      `json:"dropoff_route"` //роут доставки заказа
	Service               string            `json:"service"`       //описание алгоритма доставки(забрать, отвезти; забрать, оплатить, отвезти и пр.)
	ServiceHuman          string            `json:"service_human"`
	PaymentType           string            `json:"payment_type"` //способ оплаты(наличные, карта и пр.)
	PaymentTypeHuman      string            `json:"payment_type_human"`
	PaymentStatus         string            `json:"payment_status"` //статус оплаты(не оплачен, оплачен, ожидает оплаты и пр.)
	PaymentStatusHuman    string            `json:"payment_status_human"`
	State                 string            `json:"state"` //текущее состояние заказа(ожидает назначения, доставляется и пр.)
	StateHuman            string            `json:"state_human"`
	DeliveryPrice         float32           `json:"delivery_price"`
	CourierUUID           string            `json:"courier_uuid"` //идентификатор курьера, назначенного на заказ
	CourierData           CourierData       `json:"courier_data"` //мета-данные курьера, назначенного на заказ
}

type PersonsData struct {
	Phone   string `json:"phone"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type OrderItems struct {
	Items []DeliveryItem `json:"items"`
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

func OrderFromFrontToCore(orderFr OrderFront) models.Order {
	return models.Order{
		UUID:        orderFr.UUID,
		OrderNumber: orderFr.OrderNumber,
		DeliveryItems: models.OrderItems{
			Items: orderFr.DeliveryItems.Items,
		},
		Comment:            orderFr.Comment,
		PickupRequestTime:  time.Unix(orderFr.PickupRequestTime, 0),
		DropOffRequestTime: time.Unix(orderFr.DropOffRequestTime, 0),
		PickupPersonContacts: models.PersonsData{
			Phone:   orderFr.PickupPersonContacts.Phone,
			Name:    orderFr.PickupPersonContacts.Name,
			Comment: orderFr.PickupPersonContacts.Comment,
		},
		DropOffPersonContacts: models.PersonsData{
			Phone:   orderFr.DropOffPersonContacts.Phone,
			Name:    orderFr.DropOffPersonContacts.Name,
			Comment: orderFr.DropOffPersonContacts.Comment,
		},
		PickupRoute: models.Route{
			UUID:              orderFr.PickupRoute.UUID,
			PointType:         orderFr.PickupRoute.PointType,
			UnrestrictedValue: orderFr.PickupRoute.UnrestrictedValue,
			ValueForSearch:    orderFr.PickupRoute.ValueForSearch,
			Value:             orderFr.PickupRoute.Value,
			Country:           orderFr.PickupRoute.Country,
			Region:            orderFr.PickupRoute.Region,
			RegionType:        orderFr.PickupRoute.RegionType,
			Type:              orderFr.PickupRoute.Type,
			City:              orderFr.PickupRoute.City,
			Category:          orderFr.PickupRoute.Category,
			CityType:          orderFr.PickupRoute.CityType,
			Street:            orderFr.PickupRoute.Street,
			StreetType:        orderFr.PickupRoute.StreetType,
			StreetWithType:    orderFr.PickupRoute.StreetWithType,
			House:             orderFr.PickupRoute.House,
			FrontDoor:         orderFr.PickupRoute.FrontDoor,
			Comment:           orderFr.PickupRoute.Comment,
			OutOfTown:         orderFr.PickupRoute.OutOfTown,
			HouseType:         orderFr.PickupRoute.HouseType,
			AccuracyLevel:     orderFr.PickupRoute.AccuracyLevel,
			Radius:            orderFr.PickupRoute.Radius,
			Lat:               orderFr.PickupRoute.Lat,
			Lon:               orderFr.PickupRoute.Lon,
		},
		DropoffRoute: models.Route{
			UUID:              orderFr.DropoffRoute.UUID,
			PointType:         orderFr.DropoffRoute.PointType,
			UnrestrictedValue: orderFr.DropoffRoute.UnrestrictedValue,
			ValueForSearch:    orderFr.DropoffRoute.ValueForSearch,
			Value:             orderFr.DropoffRoute.Value,
			Country:           orderFr.DropoffRoute.Country,
			Region:            orderFr.DropoffRoute.Region,
			RegionType:        orderFr.DropoffRoute.RegionType,
			Type:              orderFr.DropoffRoute.Type,
			City:              orderFr.DropoffRoute.City,
			Category:          orderFr.DropoffRoute.Category,
			CityType:          orderFr.DropoffRoute.CityType,
			Street:            orderFr.DropoffRoute.Street,
			StreetType:        orderFr.DropoffRoute.StreetType,
			StreetWithType:    orderFr.DropoffRoute.StreetWithType,
			House:             orderFr.DropoffRoute.House,
			FrontDoor:         orderFr.DropoffRoute.FrontDoor,
			Comment:           orderFr.DropoffRoute.Comment,
			OutOfTown:         orderFr.DropoffRoute.OutOfTown,
			HouseType:         orderFr.DropoffRoute.HouseType,
			AccuracyLevel:     orderFr.DropoffRoute.AccuracyLevel,
			Radius:            orderFr.DropoffRoute.Radius,
			Lat:               orderFr.DropoffRoute.Lat,
			Lon:               orderFr.DropoffRoute.Lon,
		},
		Service:       orderFr.Service,
		PaymentType:   orderFr.PaymentType,
		PaymentStatus: orderFr.PaymentStatus,
		State:         orderFr.State,
		DeliveryPrice: orderFr.DeliveryPrice,
		CourierUUID:   orderFr.CourierUUID,
		CourierData: models.CourierData{
			UUID:        orderFr.CourierData.UUID,
			Name:        orderFr.CourierData.Name,
			PhoneNumber: orderFr.CourierData.PhoneNumber,
		},
	}
}

//FromCoreToFront ...
func OrderFromCoreToFront(orderCore models.Order) OrderFront {
	return OrderFront{
		UUID:        orderCore.UUID,
		OrderNumber: orderCore.OrderNumber,
		DeliveryItems: models.OrderItems{
			Items: orderCore.DeliveryItems.Items,
		},
		Comment:            orderCore.Comment,
		PickupRequestTime:  orderCore.PickupRequestTime.Unix(),
		DropOffRequestTime: orderCore.DropOffRequestTime.Unix(),
		PickupPersonContacts: PersonsData{
			Phone:   orderCore.PickupPersonContacts.Phone,
			Name:    orderCore.PickupPersonContacts.Name,
			Comment: orderCore.PickupPersonContacts.Comment,
		},
		DropOffPersonContacts: PersonsData{
			Phone:   orderCore.DropOffPersonContacts.Phone,
			Name:    orderCore.DropOffPersonContacts.Name,
			Comment: orderCore.DropOffPersonContacts.Comment,
		},
		PickupRoute: models.Route{
			UUID:              orderCore.PickupRoute.UUID,
			PointType:         orderCore.PickupRoute.PointType,
			UnrestrictedValue: orderCore.PickupRoute.UnrestrictedValue,
			ValueForSearch:    orderCore.PickupRoute.ValueForSearch,
			Value:             orderCore.PickupRoute.Value,
			Country:           orderCore.PickupRoute.Country,
			Region:            orderCore.PickupRoute.Region,
			RegionType:        orderCore.PickupRoute.RegionType,
			Type:              orderCore.PickupRoute.Type,
			City:              orderCore.PickupRoute.City,
			Category:          orderCore.PickupRoute.Category,
			CityType:          orderCore.PickupRoute.CityType,
			Street:            orderCore.PickupRoute.Street,
			StreetType:        orderCore.PickupRoute.StreetType,
			StreetWithType:    orderCore.PickupRoute.StreetWithType,
			House:             orderCore.PickupRoute.House,
			FrontDoor:         orderCore.PickupRoute.FrontDoor,
			Comment:           orderCore.PickupRoute.Comment,
			OutOfTown:         orderCore.PickupRoute.OutOfTown,
			HouseType:         orderCore.PickupRoute.HouseType,
			AccuracyLevel:     orderCore.PickupRoute.AccuracyLevel,
			Radius:            orderCore.PickupRoute.Radius,
			Lat:               orderCore.PickupRoute.Lat,
			Lon:               orderCore.PickupRoute.Lon,
		},
		DropoffRoute: models.Route{
			UUID:              orderCore.DropoffRoute.UUID,
			PointType:         orderCore.DropoffRoute.PointType,
			UnrestrictedValue: orderCore.DropoffRoute.UnrestrictedValue,
			ValueForSearch:    orderCore.DropoffRoute.ValueForSearch,
			Value:             orderCore.DropoffRoute.Value,
			Country:           orderCore.DropoffRoute.Country,
			Region:            orderCore.DropoffRoute.Region,
			RegionType:        orderCore.DropoffRoute.RegionType,
			Type:              orderCore.DropoffRoute.Type,
			City:              orderCore.DropoffRoute.City,
			Category:          orderCore.DropoffRoute.Category,
			CityType:          orderCore.DropoffRoute.CityType,
			Street:            orderCore.DropoffRoute.Street,
			StreetType:        orderCore.DropoffRoute.StreetType,
			StreetWithType:    orderCore.DropoffRoute.StreetWithType,
			House:             orderCore.DropoffRoute.House,
			FrontDoor:         orderCore.DropoffRoute.FrontDoor,
			Comment:           orderCore.DropoffRoute.Comment,
			OutOfTown:         orderCore.DropoffRoute.OutOfTown,
			HouseType:         orderCore.DropoffRoute.HouseType,
			AccuracyLevel:     orderCore.DropoffRoute.AccuracyLevel,
			Radius:            orderCore.DropoffRoute.Radius,
			Lat:               orderCore.DropoffRoute.Lat,
			Lon:               orderCore.DropoffRoute.Lon,
		},
		Service:            orderCore.Service,
		ServiceHuman:       Variable.Orders[orderCore.Service],
		PaymentType:        orderCore.PaymentType,
		PaymentTypeHuman:   Variable.Payments[orderCore.PaymentType],
		PaymentStatus:      orderCore.PaymentStatus,
		PaymentStatusHuman: Variable.Payments[orderCore.PaymentStatus],
		State:              orderCore.State,
		StateHuman:         Variable.Orders[orderCore.State],
		DeliveryPrice:      orderCore.DeliveryPrice,
		CourierUUID:        orderCore.CourierUUID,
		CourierData: CourierData{
			UUID:        orderCore.CourierData.UUID,
			Name:        orderCore.CourierData.Name,
			PhoneNumber: orderCore.CourierData.PhoneNumber,
		},
	}
}
