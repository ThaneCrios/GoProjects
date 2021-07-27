package models

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
