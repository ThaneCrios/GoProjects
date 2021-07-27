package models

import (
	"gitlab.com/faemproject/backend/core/shared/structures"
	"gitlab.com/faemproject/backend/eda/eda.core/services/orders/proto"
	productModels "gitlab.com/faemproject/backend/eda/eda.core/services/products/models"
	"time"
)

type Order struct {
	tableName         struct{}            `sql:"collector_orders"`
	UUID              string              `json:"uuid"`
	ID                string              `json:"id"`
	StoreUUID         string              `pg:",use_zero" json:"store_uuid"`
	StoreData         productModels.Store `json:"store_data"`
	DeviceID          string              `json:"device_id"`
	ClientUUID        string              `json:"client_uuid"`
	ClientData        Client              `json:"client_data"`
	CollectorUUID     string              `json:"collector_uuid"`
	CollectorData     CollectorInfo       `json:"collector_data"`
	Application       string              `json:"application"`
	State             string              `json:"state"`
	Source            string              `json:"source"`
	CallbackPhone     string              `json:"callback_phone"`
	Promotion         Promotion           `json:"promotion"`
	Comment           string              `json:"comment"`
	Deleted           bool                `json:"-"`
	CartItems         []CartItem          `sql:"items" json:"items"`
	TotalPrice        float64             `json:"total_price"`
	CookingTime       int                 `json:"cooking_time"`
	DeliveryData      DeliveryData        `json:"delivery_data"`
	CookingTimeFinish time.Time           `json:"cooking_time_finish"`
	LastUpdateUUID    string              `json:"last_update_uuid"`
	LastUpdateRole    OrderInitiatorRole  `json:"last_update_role"`
	CancelReason      string              `json:"cancel_reason"`
	CancelComment     string              `json:"cancel_comment"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

type OrderInitiatorRole string

type DeliveryAddressDetails struct {
	Entrance  string `json:"entrance"`  // подъезд
	Floor     string `json:"floor"`     // этаж
	Apartment string `json:"apartment"` // квартира
	Intercom  string `json:"intercom"`  // домофон
}

type DeliveryData struct {
	Price          float64                `json:"price"`
	EstimatedTime  int                    `json:"estimated_time"` // в секундах
	Type           string                 `json:"type"`
	Address        structures.Route       `json:"address"`
	AddressDetails DeliveryAddressDetails `json:"address_details"`
}

type DeliveryTariff struct {
	Price         float64 `json:"price"`
	EstimatedTime int     `json:"estimated_time"` // в секундах
}

type Store struct {
	Type string `json:"type"`
}

type Client struct {
	UUID        string             `json:"uuid"`
	Name        string             `json:"name"`
	Application string             `json:"application" required:"true"`
	MainPhone   string             `json:"main_phone" required:"true"`
	DevicesID   []string           `json:"devices_id" pg:",type:text[]"`
	Deleted     bool               `json:"-"`
	Blocked     bool               `json:"blocked"`
	Addresses   []structures.Route `json:"addresses"`
	Meta        ClientMeta         `json:"meta"`
	CreatedAt   time.Time          `json:"-"`
	UpdatedAt   time.Time          `json:"-"`
}

type ClientMeta struct {
	InviteCount int `json:"invite_count"`
}

//CollectItems состав заказа
type (
	CartItem struct {
		ID              int                `json:"id"`
		Product         CartProduct        `json:"product"`
		VariantGroups   []CartVariantGroup `json:"variant_groups"`
		SingleItemPrice float64            `json:"single_item_price"` // суммарная одного итема
		TotalItemPrice  float64            `json:"total_item_price"`  // суммарная стоимость итема
		Count           int                `json:"count"`
		Hash            string             `json:"hash"` // TODO: удалить поле
	}

	CartProduct struct {
		UUID              string      `json:"uuid"`
		Name              string      `json:"name" required:"true"`
		StoreUUID         string      `json:"store_uuid" required:"true"`
		Type              ProductType `json:"type"`
		Price             float64     `json:"price"`
		Leftover          int         `json:"leftover"`
		Weight            float64     `json:"weight"`
		WeightMeasurement string      `json:"weight_measurement"`
		Meta              ProductMeta `json:"meta"`
	}

	CartVariantGroup struct {
		UUID        string        `json:"uuid"`
		Name        string        `json:"name"`
		ProductUUID string        `json:"product_uuid"`
		Required    bool          `json:"required"`
		Variants    []CartVariant `json:"variants"`
	}

	CartVariant struct {
		UUID             string  `json:"uuid"`
		Name             string  `json:"name"`
		ProductUUID      string  `json:"product_uuid"`
		VariantGroupUUID string  `json:"variant_group_uuid"`
		Price            float64 `json:"price"`
	}
)

//CollectorInfo - информация по курьру
type CollectorInfo struct {
	UUID        string `json:"uuid"`         //идентификатор курьера
	Name        string `json:"name"`         //ФИО
	PhoneNumber string `json:"phone_number"` //номер телефона
}

type OrderForCollectorDuplicate struct {
	tableName     struct{}            `sql:"collector_orders_duplicate"`
	UUID          string              `json:"uuid"`
	ID            string              `json:"id"`
	CollectorUUID string              `json:"collector_uuid"`
	OrderUUID     string              `json:"order_uuid"`
	ClientData    Client              `json:"client_data"`
	CallbackPhone string              `json:"callback_phone"`
	Comment       string              `json:"comment"`
	State         string              `json:"state"`
	CartItems     []CartItemDuplicate `sql:"items" json:"items"`
	TotalPrice    float64             `json:"total_price"`
	CollectTime   int                 `json:"cooking_time"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeliveryData  DeliveryData        `json:"delivery_data"`
}

type ClientDuplicate struct {
	Phone string
	Name  string
}

type (
	CartItemDuplicate struct {
		ID                  int                  `json:"id"`
		Product             CartProductDuplicate `json:"product"`
		VariantGroups       []CartVariantGroup   `json:"variant_groups"`
		SingleItemPrice     float64              `json:"single_item_price"` // суммарная одного итема
		TotalItemPrice      float64              `json:"total_item_price"`
		CollectionSign      string               `json:"collection_sign"`
		WasChanged          bool                 `json:"was_changed"`
		PreviousProductUUID string               `json:"previous_product_uuid"` // суммарная стоимость итема
		Count               int                  `json:"count"`
		Hash                string               `json:"hash"` // TODO: удалить поле
	}

	CartProductDuplicate struct {
		UUID              string      `json:"uuid"`
		Name              string      `json:"name" required:"true"`
		StoreUUID         string      `json:"store_uuid" required:"true"`
		Type              ProductType `json:"type"`
		Price             float64     `json:"price"`
		Leftover          int         `json:"leftover"`
		Weight            float64     `json:"weight"`
		WeightMeasurement string      `json:"weight_measurement"`
		Meta              ProductMeta `json:"meta"`
	}
)

type CollectorInfoCore struct {
	UUID        string `json:"uuid"`         //идентификатор курьера
	Name        string `json:"name"`         //ФИО
	PhoneNumber string `json:"phone_number"` //номер телефона
}

type OrderUpdated struct {
	UUID              string              `json:"uuid"`
	ID                string              `json:"id"`
	StoreUUID         string              `pg:",use_zero" json:"store_uuid"`
	StoreData         productModels.Store `json:"store_data"`
	DeviceID          string              `json:"device_id"`
	ClientUUID        string              `json:"client_uuid"`
	ClientData        Client              `json:"client_data"`
	Application       string              `json:"application"`
	State             string              `json:"state"`
	Source            string              `json:"source"`
	CallbackPhone     string              `json:"callback_phone"`
	Comment           string              `json:"comment"`
	Deleted           bool                `json:"-"`
	CartItems         []CartItem          `sql:"items" json:"items"`
	Promotion         Promotion           `json:"promotion"`
	TotalPrice        float64             `json:"total_price"`
	CookingTime       int                 `json:"cooking_time"`
	DeliveryData      DeliveryData        `json:"delivery_data"`
	CookingTimeFinish time.Time           `json:"cooking_time_finish"`
	LastUpdateUUID    string              `json:"last_update_uuid"`
	LastUpdateRole    OrderInitiatorRole  `json:"last_update_role"`
	CancelReason      string              `json:"cancel_reason"`
	CancelComment     string              `json:"cancel_comment"`
	OriginOrder       OrderOriginal       `json:"origin_order"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

type Promotion struct {
	UUID               string  `json:"uuid"`
	BillingAccountUUID string  `json:"billing_account_uuid"`
	DiscountPercentage float64 `json:"discount_percentage"`
	PromoCode          string  `json:"promo_code"`
}

type OrderOriginal struct {
	UUID              string             `json:"uuid"`
	ID                string             `json:"id"`
	StoreUUID         string             `pg:",use_zero" json:"store_uuid"`
	StoreData         Store              `json:"store_data"`
	DeviceID          string             `json:"device_id"`
	ClientUUID        string             `json:"client_uuid"`
	ClientData        Client             `json:"client_data"`
	Application       string             `json:"application"`
	State             string             `json:"state"`
	Source            string             `json:"source"`
	CallbackPhone     string             `json:"callback_phone"`
	Comment           string             `json:"comment"`
	Deleted           bool               `json:"-"`
	CartItems         []CartItem         `sql:"items" json:"items"`
	TotalPrice        float64            `json:"total_price"`
	CookingTime       int                `json:"cooking_time"`
	DeliveryData      DeliveryData       `json:"delivery_data"`
	CookingTimeFinish time.Time          `json:"cooking_time_finish"`
	LastUpdateUUID    string             `json:"last_update_uuid"`
	LastUpdateRole    OrderInitiatorRole `json:"last_update_role"`
	CancelReason      string             `json:"cancel_reason"`
	CancelComment     string             `json:"cancel_comment"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

func ConvertAllOrdersForSync(orderDuplicate OrderForCollectorDuplicate, order Order) OrderUpdated {
	return OrderUpdated{
		UUID: order.UUID,
		ID:   order.ID,
		Promotion: Promotion{
			UUID:               order.Promotion.UUID,
			BillingAccountUUID: order.Promotion.BillingAccountUUID,
			DiscountPercentage: order.Promotion.DiscountPercentage,
			PromoCode:          order.Promotion.PromoCode,
		},
		StoreUUID:  order.StoreUUID,
		StoreData:  order.StoreData,
		DeviceID:   order.DeviceID,
		ClientUUID: order.ClientUUID,
		ClientData: Client{
			UUID:        order.ClientData.UUID,
			Name:        order.ClientData.Name,
			Application: order.ClientData.Application,
			MainPhone:   order.ClientData.MainPhone,
			DevicesID:   order.ClientData.DevicesID,
			Deleted:     order.ClientData.Deleted,
			Blocked:     order.ClientData.Blocked,
			Addresses:   order.ClientData.Addresses,
			Meta: ClientMeta{
				InviteCount: order.ClientData.Meta.InviteCount},
			CreatedAt: order.ClientData.CreatedAt,
			UpdatedAt: order.ClientData.UpdatedAt,
		},
		Application:   order.Application,
		Source:        order.Source,
		CallbackPhone: order.CallbackPhone,
		Comment:       order.Comment,
		CartItems:     convertItemsToOriginal(orderDuplicate.CartItems),
		TotalPrice:    orderDuplicate.TotalPrice,
		DeliveryData: DeliveryData{
			Price:         order.DeliveryData.Price,
			EstimatedTime: order.DeliveryData.EstimatedTime,
			Type:          order.DeliveryData.Type,
			Address: structures.Route{
				UUID:              order.DeliveryData.Address.UUID,
				PointType:         order.DeliveryData.Address.PointType,
				UnrestrictedValue: order.DeliveryData.Address.UnrestrictedValue,
				ValueForSearch:    order.DeliveryData.Address.ValueForSearch,
				Value:             order.DeliveryData.Address.Value,
				Country:           order.DeliveryData.Address.Country,
				Region:            order.DeliveryData.Address.Region,
				RegionType:        order.DeliveryData.Address.RegionType,
				Type:              order.DeliveryData.Address.Type,
				City:              order.DeliveryData.Address.City,
				Category:          order.DeliveryData.Address.Category,
				CityType:          order.DeliveryData.Address.CityType,
				Street:            order.DeliveryData.Address.Street,
				StreetType:        order.DeliveryData.Address.StreetType,
				StreetWithType:    order.DeliveryData.Address.StreetWithType,
				House:             order.DeliveryData.Address.House,
				FrontDoor:         order.DeliveryData.Address.FrontDoor,
				Comment:           order.DeliveryData.Address.Comment,
				OutOfTown:         order.DeliveryData.Address.OutOfTown,
				HouseType:         order.DeliveryData.Address.HouseType,
				AccuracyLevel:     order.DeliveryData.Address.AccuracyLevel,
				Radius:            order.DeliveryData.Address.Radius,
				Lat:               order.DeliveryData.Address.Lat,
				Lon:               order.DeliveryData.Address.Lon,
			},
			AddressDetails: DeliveryAddressDetails{
				Entrance:  order.DeliveryData.AddressDetails.Entrance,
				Floor:     order.DeliveryData.AddressDetails.Floor,
				Apartment: order.DeliveryData.AddressDetails.Apartment,
				Intercom:  order.DeliveryData.AddressDetails.Intercom,
			},
		},
		OriginOrder: OrderOriginal{
			UUID:      order.UUID,
			ID:        order.ID,
			StoreUUID: order.StoreUUID,
			StoreData: Store{
				Type: order.StoreData.Type,
			},
			DeviceID:   order.DeviceID,
			ClientUUID: order.ClientUUID,
			ClientData: Client{
				UUID:        order.ClientData.UUID,
				Name:        order.ClientData.Name,
				Application: order.ClientData.Application,
				MainPhone:   order.ClientData.MainPhone,
				DevicesID:   order.ClientData.DevicesID,
				Deleted:     order.ClientData.Deleted,
				Blocked:     order.ClientData.Blocked,
				Addresses:   order.ClientData.Addresses,
				Meta: ClientMeta{
					InviteCount: order.ClientData.Meta.InviteCount},
				CreatedAt: order.ClientData.CreatedAt,
				UpdatedAt: order.ClientData.UpdatedAt,
			},
			Application:   order.Application,
			State:         order.State,
			Source:        order.Source,
			CallbackPhone: order.CallbackPhone,
			Comment:       order.Comment,
			Deleted:       order.Deleted,
			CartItems:     order.CartItems,
			TotalPrice:    order.TotalPrice,
			CookingTime:   order.CookingTime,
			DeliveryData: DeliveryData{
				Price:         order.DeliveryData.Price,
				EstimatedTime: order.DeliveryData.EstimatedTime,
				Type:          order.DeliveryData.Type,
				Address: structures.Route{
					UUID:              order.DeliveryData.Address.UUID,
					PointType:         order.DeliveryData.Address.PointType,
					UnrestrictedValue: order.DeliveryData.Address.UnrestrictedValue,
					ValueForSearch:    order.DeliveryData.Address.ValueForSearch,
					Value:             order.DeliveryData.Address.Value,
					Country:           order.DeliveryData.Address.Country,
					Region:            order.DeliveryData.Address.Region,
					RegionType:        order.DeliveryData.Address.RegionType,
					Type:              order.DeliveryData.Address.Type,
					City:              order.DeliveryData.Address.City,
					Category:          order.DeliveryData.Address.Category,
					CityType:          order.DeliveryData.Address.CityType,
					Street:            order.DeliveryData.Address.Street,
					StreetType:        order.DeliveryData.Address.StreetType,
					StreetWithType:    order.DeliveryData.Address.StreetWithType,
					House:             order.DeliveryData.Address.House,
					FrontDoor:         order.DeliveryData.Address.FrontDoor,
					Comment:           order.DeliveryData.Address.Comment,
					OutOfTown:         order.DeliveryData.Address.OutOfTown,
					HouseType:         order.DeliveryData.Address.HouseType,
					AccuracyLevel:     order.DeliveryData.Address.AccuracyLevel,
					Radius:            order.DeliveryData.Address.Radius,
					Lat:               order.DeliveryData.Address.Lat,
					Lon:               order.DeliveryData.Address.Lon,
				},
				AddressDetails: DeliveryAddressDetails{
					Entrance:  order.DeliveryData.AddressDetails.Entrance,
					Floor:     order.DeliveryData.AddressDetails.Floor,
					Apartment: order.DeliveryData.AddressDetails.Apartment,
					Intercom:  order.DeliveryData.AddressDetails.Intercom,
				},
			},
			CookingTimeFinish: order.CookingTimeFinish,
			LastUpdateUUID:    order.LastUpdateUUID,
			LastUpdateRole:    order.LastUpdateRole,
			CancelReason:      order.CancelReason,
			CancelComment:     order.CancelComment,
			CreatedAt:         order.CreatedAt,
			UpdatedAt:         order.UpdatedAt,
		},
	}
}

func CreateOrderDuplicate(order Order) (orderOutput OrderForCollectorDuplicate) {
	for i, v := range order.CartItems {
		orderOutput.CartItems = append(orderOutput.CartItems, convertOrderDuplicate(v))
		orderOutput.CartItems[i].CollectionSign = "uncollected"
	}
	orderOutput.TotalPrice = order.TotalPrice
	orderOutput.ID = order.ID
	orderOutput.State = proto.OrderStateCooking
	orderOutput.OrderUUID = order.UUID
	orderOutput.CollectorUUID = order.CollectorUUID
	orderOutput.CallbackPhone = order.CallbackPhone
	orderOutput.ClientData = order.ClientData
	orderOutput.DeliveryData = DeliveryData{
		Price:         order.DeliveryData.Price,
		EstimatedTime: order.DeliveryData.EstimatedTime,
		Type:          order.DeliveryData.Type,
		Address: structures.Route{
			UUID:              order.DeliveryData.Address.UUID,
			PointType:         order.DeliveryData.Address.PointType,
			UnrestrictedValue: order.DeliveryData.Address.UnrestrictedValue,
			ValueForSearch:    order.DeliveryData.Address.ValueForSearch,
			Value:             order.DeliveryData.Address.Value,
			Country:           order.DeliveryData.Address.Country,
			Region:            order.DeliveryData.Address.Region,
			RegionType:        order.DeliveryData.Address.RegionType,
			Type:              order.DeliveryData.Address.Type,
			City:              order.DeliveryData.Address.City,
			Category:          order.DeliveryData.Address.Category,
			CityType:          order.DeliveryData.Address.CityType,
			Street:            order.DeliveryData.Address.Street,
			StreetType:        order.DeliveryData.Address.StreetType,
			StreetWithType:    order.DeliveryData.Address.StreetWithType,
			House:             order.DeliveryData.Address.House,
			FrontDoor:         order.DeliveryData.Address.FrontDoor,
			Comment:           order.DeliveryData.Address.Comment,
			OutOfTown:         order.DeliveryData.Address.OutOfTown,
			HouseType:         order.DeliveryData.Address.HouseType,
			AccuracyLevel:     order.DeliveryData.Address.AccuracyLevel,
			Radius:            order.DeliveryData.Address.Radius,
			Lat:               order.DeliveryData.Address.Lat,
			Lon:               order.DeliveryData.Address.Lon,
		},
		AddressDetails: DeliveryAddressDetails{
			Entrance:  order.DeliveryData.AddressDetails.Entrance,
			Floor:     order.DeliveryData.AddressDetails.Floor,
			Apartment: order.DeliveryData.AddressDetails.Apartment,
			Intercom:  order.DeliveryData.AddressDetails.Intercom,
		},
	}
	return orderOutput
}

func convertOrderDuplicate(item CartItem) CartItemDuplicate {
	return CartItemDuplicate{
		ID:              item.ID,
		Product:         convertProducts(item.Product),
		VariantGroups:   item.VariantGroups,
		SingleItemPrice: item.SingleItemPrice,
		TotalItemPrice:  item.TotalItemPrice,
		Count:           item.Count,
		Hash:            item.Hash,
		WasChanged:      false,
	}
}

func convertProducts(item CartProduct) CartProductDuplicate {
	return CartProductDuplicate{
		UUID:              item.UUID,
		Name:              item.Name,
		StoreUUID:         item.StoreUUID,
		Type:              item.Type,
		Price:             item.Price,
		Weight:            item.Weight,
		WeightMeasurement: item.WeightMeasurement,
		Meta:              item.Meta,
	}
}

func convertItemsToOriginal(duplicateItem []CartItemDuplicate) []CartItem {
	var originalItem []CartItem

	for i, _ := range duplicateItem {
		var orig CartItem
		orig.TotalItemPrice = duplicateItem[i].TotalItemPrice
		orig.SingleItemPrice = duplicateItem[i].SingleItemPrice
		orig.Product = convertProductToOriginal(duplicateItem[i].Product)
		orig.Hash = duplicateItem[i].Hash
		orig.Count = duplicateItem[i].Count
		orig.ID = duplicateItem[i].ID
		orig.VariantGroups = duplicateItem[i].VariantGroups

		originalItem = append(originalItem, orig)
	}

	return originalItem
}

func convertProductToOriginal(duplicate CartProductDuplicate) (original CartProduct) {
	original.UUID = duplicate.UUID
	original.Meta = duplicate.Meta
	original.Type = duplicate.Type
	original.Name = duplicate.Name
	original.WeightMeasurement = duplicate.WeightMeasurement
	original.Weight = duplicate.Weight
	original.Price = duplicate.Price
	original.StoreUUID = duplicate.StoreUUID

	return original
}
