package data

import "PhoceeneAuto/internal/validator"

// Address represents an address with various fields.
type Address struct {
	Street     string `json:"street"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city"`
	ZIP        string `json:"ZIP"`
	Country    string `json:"country"`
}

// validate validates the Address instance using a Validator.
//
// Parameters:
//
//	v - The Validator instance to use for validation
func (addr Address) validate(v *validator.Validator) {

	// validating each field of the address
	v.StringCheck(addr.Street, 2, 256, true, "street")
	v.StringCheck(addr.Complement, 2, 256, false, "complement")
	v.StringCheck(addr.City, 2, 40, true, "city")
	v.StringCheck(addr.ZIP, 2, 12, true, "zip")
	v.StringCheck(addr.Country, 2, 40, true, "country")
}

// toSQL converts the Address instance into a slice of arguments suitable for SQL queries.
//
// Parameters:
//
//	args - The existing slice of arguments
//
// Returns:
//
//	[]any - A new slice containing all original arguments followed by the address fields
func (addr Address) toSQL(args []any) []any {

	// appending each field of the address to the provided slice of arguments
	return append(args, addr.Street, addr.Complement, addr.City, addr.ZIP, addr.Country)
}
