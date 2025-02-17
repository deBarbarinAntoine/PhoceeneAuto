package data

import "PhoceeneAuto/internal/validator"

type Address struct {
	Street     string `json:"street"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city"`
	ZIP        string `json:"ZIP"`
	State      string `json:"state"`
}

func (addr Address) validate(v *validator.Validator) {
	v.StringCheck(addr.Street, 2, 256, true, "street")
	v.StringCheck(addr.Complement, 2, 256, false, "complement")
	v.StringCheck(addr.City, 2, 40, true, "city")
	v.StringCheck(addr.ZIP, 2, 12, true, "zip")
	v.StringCheck(addr.State, 2, 40, true, "state")
}

func (addr Address) toSQL(args []any) {
	args = append(args, addr.Street, addr.Complement, addr.City, addr.ZIP, addr.State)
}
