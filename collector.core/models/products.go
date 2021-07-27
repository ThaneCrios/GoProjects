package models

import (
	"time"
)

type Product struct {
	tableName  struct{}    `sql:"collector_products,discard_unknown_columns"`
	UUID       string      `json:"uuid"`
	ExternalId string      `json:"external_id"`
	Name       string      `json:"name" required:"true"`
	StoreUUID  string      `json:"store_uuid" required:"true"`
	Comment    string      `json:"comment"`
	Url        string      `json:"url"`
	Deleted    bool        `json:"deleted"`
	Available  bool        `json:"available"`
	StopList   bool        `json:"stop_list"`
	DefaultSet bool        `json:"default_set"`
	Priority   float64     `json:"priority"`
	Type       ProductType `json:"type"`
	Leftover   int         `json:"leftover"`
	Price      float64     `json:"price"`
	Meta       ProductMeta `json:"meta"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type (
	ProductType string

	ProductMeta struct {
		ShortDescription  string   `json:"short_description"`
		Description       string   `json:"description"`
		Composition       string   `json:"composition"` // состав
		Weight            float64  `json:"weight"`
		WeightMeasurement string   `json:"weight_measurement"`
		Discount          float64  `json:"discount"`
		Images            []string `json:"images"`
	}
)

type ClientData struct {
	ClientPhone string `json:"client_phone"`
	ClientName  string `json:"client_name"`
}

type BarCodes struct {
	TableName      struct{} `sql:"collector_barcodes"`
	BarCodeOrder   string   `json:"bar_code_order"`
	BarCodeScanned string   `json:"bar_code_scanned"`
}

type ProductWithBarCode struct {
	tableName   struct{} `sql:"collector_barcodes"`
	UUID        string   `json:"uuid"`
	BarCode     string   `json:"bar_code"`
	ProductUUID string   `json:"product_uuid"`
}
