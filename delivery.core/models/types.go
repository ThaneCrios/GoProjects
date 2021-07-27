package models

const (
	TypeMachine = "machine"
	TypeFoot = "on foot"
)

type CourierType struct {
	CourierMachine string
	CourierFoot	   string
}

