package proto

import (
	"gitlab.com/faemproject/backend/core/shared/structures"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	"time"
)

type OrderCollector struct {
	OrderUUID     string `json:"order_uuid"`
	CollectorUUID string `json:"collector_uuid"`
	CookingTime   int    `json:"cooking_time"`
}

type OrderCancel struct {
	OrderUUID string `json:"order_uuid"`
}

type ProductAdd struct {
	OrderUUID    string               `json:"order_uuid"`
	BarCode      string               `json:"bar_code"`
	ProductCount int                  `json:"product_count"`
	ProductUUID  string               `json:"product_uuid"`
	StoreUUID    string               `json:"store_uuid"`
	Products     []MassProductsAdding `json:"products"`
}

type MassProductsAdding struct {
	ProductUUID  string `json:"uuid"`
	ProductCount int    `json:"count"`
}

type OrderProductRemove struct {
	OrderUUID   string                   `json:"order_uuid"`
	ProductUUID int                      `json:"product_id"`
	Products    []MassOrderRemoveProduct `json:"products"`
}

type MassOrderRemoveProduct struct {
	ProductID int `json:"id"`
}

type OrderProductMark struct {
	OrderUUID    string               `json:"order_uuid"`
	ProductUUID  int                  `json:"product_id"`
	ProductCount int                  `json:"product_count"`
	Products     []MassProductMarking `json:"products"`
}

type MassProductMarking struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

type OrderProductChange struct {
	OrderUUID       string `json:"order_uuid"`
	ProductUUID     int    `json:"old_product_id"`
	NewProductUUID  string `json:"new_product_uuid"`
	NewProductCount int    `json:"new_product_count"`
}

type OrderForCollector struct {
	CollectorUUID string        `json:"collector_uuid"`
	OrderUUID     string        `json:"order_uuid"`
	ClientData    models.Client `json:"client_data"`
	CallbackPhone string        `json:"callback_phone"`
	Comment       string        `json:"comment"`
	CartItems     []CartItem    `json:"items"`
	TotalPrice    float64       `json:"total_price"`
	CollectTime   int           `json:"cooking_time"`
}

type Client struct {
	Phone string
	Name  string
}

type CollectorInfo struct {
	UUID        string `json:"uuid"`         //идентификатор курьера
	Name        string `json:"name"`         //ФИО
	PhoneNumber string `json:"phone_number"` //номер телефона
}

type (
	OrderInitiatorRole string

	Order struct {
		UUID                   string                 `json:"uuid"`
		ID                     string                 `json:"id"`
		StoreUUID              string                 `pg:",use_zero" json:"store_uuid" required:"true"`
		StoreData              Store                  `json:"store_data"`
		DeviceID               string                 `json:"device_id"`
		ClientUUID             string                 `json:"client_uuid" required:"true"`
		ClientData             Client                 `json:"client_data"`
		Application            string                 `json:"application" required:"true"`
		Source                 string                 `json:"source" required:"true"`
		State                  string                 `json:"state" required:"true"`
		CallbackPhone          string                 `json:"callback_phone"`
		Comment                string                 `json:"comment"`
		Deleted                bool                   `json:"-"`
		CartItems              []CartItem             `json:"items" sql:"items"`
		Promotion              Promotion              `json:"promotion"`
		PaymentType            string                 `json:"payment_type"`
		TotalPrice             float64                `json:"total_price" pg:",use_zero"`
		OwnDelivery            bool                   `json:"own_delivery"`
		WithoutDelivery        bool                   `json:"without_delivery"` // заберу сам
		EatInStore             bool                   `json:"eat_in_store"`
		DeliveryType           string                 `json:"delivery_type"`
		DeliveryPrice          float64                `json:"delivery_price"`
		DeliveryAddress        structures.Route       `json:"delivery_address"`
		DeliveryAddressDetails DeliveryAddressDetails `json:"delivery_address_details"`
		CookingTime            int                    `json:"cooking_time"`
		CookingTimeFinish      time.Time              `json:"cooking_time_finish"`
		LastUpdateUUID         string                 `json:"last_update_uuid"`
		LastUpdateRole         OrderInitiatorRole     `json:"last_update_role"` // кто совершил действие (client / user / system)
		CancelReason           string                 `json:"cancel_reason"`
		CancelComment          string                 `json:"cancel_comment"`
		CreatedAt              time.Time              `json:"created_at"`
		UpdatedAt              time.Time              `json:"-"`
	}

	Promotion struct {
		UUID               string  `json:"uuid"`
		BillingAccountUUID string  `json:"billing_account_uuid"`
		Amount             float64 `json:"amount"`
		PromoCode          string  `json:"promo_code"`
	}

	DeliveryAddressDetails struct {
		Entrance  string `json:"entrance"`  // подъезд
		Floor     string `json:"floor"`     // этаж
		Apartment string `json:"apartment"` // квартира
		Intercom  string `json:"intercom"`  // домофон
	}

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

type (
	PaymentType string

	Store struct {
		tableName       struct{}         `pg:"prd_stores,discard_unknown_columns"`
		UUID            string           `json:"uuid"`
		Name            string           `json:"name" required:"true"`
		PaymentTypes    []PaymentType    `json:"payment_types" pg:",type:text[]" required:"true"`
		CityUUID        string           `json:"city_uuid" required:"true"`
		LegalEntityUUID string           `json:"legal_entity_uuid" required:"true"`
		ParentUUID      string           `json:"parent_uuid"`
		Deleted         bool             `json:"-"`
		Available       StoreAvailable   `json:"available"`
		Validated       bool             `json:"-"`
		Type            string           `json:"type"`
		WorkSchedule    WorkSchedule     `json:"work_schedule"`
		Address         structures.Route `json:"address"`
		Contacts        []Contact        `json:"contacts"`
		Priority        float64          `json:"priority"`
		Lat             float32          `json:"lat"`
		Lon             float32          `json:"lon"`
		OwnDelivery     bool             `json:"own_delivery"`
		Url             string           `json:"url"`
		Meta            StoreMeta        `json:"meta"`
		CreatedAt       time.Time        `json:"-"`
		UpdatedAt       time.Time        `json:"-"`
	}

	StoreAvailable struct {
		Flag     bool   `json:"flag"`
		Reason   string `json:"reason"`
		Duration int    `json:"duration"`
	}

	StoreMeta struct {
		Images           []string `json:"images"`
		Rating           float64  `json:"rating"`
		DeliveryTime     string   `json:"delivery_time"`
		DeliveryPrice    string   `json:"delivery_price"`
		Description      string   `json:"description"`
		ConfirmationTime int      `json:"confirmation_time"`
		CookingTime      int      `json:"cooking_time"`
	}

	WorkSchedule struct {
		TimeZoneOffset int                        `json:"time_zone_offset"`
		Standard       []StandardWorkScheduleItem `json:"standard"`
		Holiday        []HolidayWorkScheduleItem  `json:"holiday"`
	}

	StandardWorkScheduleItem struct {
		BeginningTime string `json:"beginning_time"`
		EndingTime    string `json:"ending_time"`
		WeekDays      []bool `json:"week_days"`
	}

	HolidayWorkScheduleItem struct {
		BeginningDate string                     `json:"beginning_date"`
		EndingDate    string                     `json:"ending_date"`
		Items         []StandardWorkScheduleItem `json:"items"`
	}

	Contact struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}
)

const (
	// Статусы заказа
	OrderStateCart      = "cart"      // добавлен хотя бы 1 товар в корзину
	OrderStatePayment   = "payment"   // процесс оплаты заказа
	OrderStateCreated   = "created"   // заказ создан
	OrderStateCooking   = "cooking"   // заказ готовиться
	OrderStateReady     = "ready"     // заказ приготовлен
	OrderStateDelivery  = "delivery"  // заказ доставляется
	OrderStateFinished  = "finished"  // заказ выполнен
	OrderStateCancelled = "cancelled" // заказ отменен

	// Источники заказа
	OrderSourceMobile = "mobile"

	// Причины отмены заказа
	OrderCancelReasonNotConfirmed = "not_confirmed" // не принят заведением
)
