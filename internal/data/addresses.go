package data

import "PhoceeneAuto/internal/validator"

type Address struct {
	Street     string `json:"street"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city"`
	ZIP        string `json:"ZIP"`
	Country    string `json:"country"`
}

func (addr Address) validate(v *validator.Validator) {
	v.StringCheck(addr.Street, 2, 256, true, "street")
	v.StringCheck(addr.Complement, 2, 256, false, "complement")
	v.StringCheck(addr.City, 2, 40, true, "city")
	v.StringCheck(addr.ZIP, 2, 12, true, "zip")
	v.StringCheck(addr.Country, 2, 40, true, "country")
}

func (addr Address) toSQL(args []any) []any {
	return append(args, addr.Street, addr.Complement, addr.City, addr.ZIP, addr.Country)
}
