package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.uber.org/multierr"
)

func (c Courier) Validate() error {
	return multierr.Combine(
		validation.ValidateStruct(&c,
			validation.Field(&c.PhoneNumber, validation.Length(11, 11), validation.Match(regexp.MustCompile("^[0-9]"))),
		),
		validation.ValidateStruct(&c.CourierMeta,
			validation.Field(&c.CourierMeta.FullName, validation.Required, validation.Match(regexp.MustCompile("^[А-я]"))),
			validation.Field(&c.CourierMeta.Requisites, validation.Length(16, 16), validation.Match(regexp.MustCompile("^[0-9]"))),
		),
	)
}

func (o Order) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.PickupContactPhone, validation.Length(11, 11), validation.Match(regexp.MustCompile("^[0-9]"))),
	)
}
