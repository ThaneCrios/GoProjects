package models

import (
	"time"
)

type Order struct {
	tableName     struct{}      `sql:"collector_orders"`
	UUID          string        `json:"uuid"`
	ID            string        `json:"id"`
	StoreUUID     string        `pg:",use_zero" json:"store_uuid"`
	DeviceID      string        `json:"device_id"`
	ClientUUID    string        `json:"client_uuid"`
	ClientData    Client        `json:"client_data"`
	CollectorUUID string        `json:"collector_uuid"`
	CollectorData CollectorInfo `json:"collector_data"`
	Application   string        `json:"application"`
	State         string        `json:"state"`
	Source        string        `json:"source"`
	CallbackPhone string        `json:"callback_phone"`
	Comment       string        `json:"comment"`
	Deleted       bool          `json:"-"`
	CartItems     []CartItem    `pg:"items" json:"cart_items"`
	TotalPrice    float64       `json:"total_price"`
	CookingTime   int           `json:"cooking_time"`
	CancelReason  string        `json:"cancel_reason"`
	CancelComment string        `json:"cancel_comment"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Client struct {
	Phone string `json:"phone"`
	Name  string `json:"name"`
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
	tableName     struct{}            `sql:"collector_order_duplicate"`
	UUID          string              `json:"uuid"`
	CollectorUUID string              `json:"collector_uuid"`
	OrderUUID     string              `json:"order_uuid"`
	ClientData    ClientDuplicate     `json:"client_data"`
	CallbackPhone string              `json:"callback_phone"`
	Comment       string              `json:"comment"`
	CartItems     []CartItemDuplicate `json:"cart_items"`
	TotalPrice    float64             `json:"total_price"`
	CollectTime   int                 `json:"cooking_time"`
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

func CreateOrderDuplicate(order Order) (orderOutput OrderForCollectorDuplicate) {
	for i, v := range order.CartItems {
		orderOutput.CartItems = append(orderOutput.CartItems, convertOrderDuplicate(v))
		orderOutput.CartItems[i].CollectionSign = "uncollected"
	}
	orderOutput.TotalPrice = order.TotalPrice
	orderOutput.OrderUUID = order.UUID
	orderOutput.CollectorUUID = order.CollectorUUID
	orderOutput.CallbackPhone = order.CallbackPhone
	orderOutput.ClientData.Name = order.ClientData.Name
	orderOutput.ClientData.Phone = order.ClientData.Phone
	orderOutput.CollectTime = 1
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
