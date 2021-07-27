package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
)

func (o Order) Validate() error {
	return validation.ValidateStruct(&o.ClientData,
		validation.Field(&o.ClientData.MainPhone, validation.Length(11, 11), validation.Match(regexp.MustCompile("^[0-9]"))),
	)
}
