package proto

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

type AppointParams struct {
	BarCode     string `json:"bar_code"`
	ProductUUID string `json:"product_uuid"`
}

type ProductBarcodesFilter struct {
	ProductUUID string `query:"product_uuid" json:"product_uuid"`
	Barcode     string `query:"barcode" json:"barcode"`
}
