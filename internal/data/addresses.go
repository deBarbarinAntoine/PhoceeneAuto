package data

import (
	"database/sql"

	"PhoceeneAuto/internal/validator"
)

// Address represents an address with various fields.
type Address struct {
	Street     string `json:"street"`
	Complement string `json:"complement,omitempty"`
	City       string `json:"city"`
	ZIP        string `json:"ZIP"`
	Country    string `json:"country"`
}

type AddressSql struct {
	Street     sql.NullString `json:"street"`
	Complement sql.NullString `json:"complement,omitempty"`
	City       sql.NullString `json:"city"`
	ZIP        sql.NullString `json:"ZIP"`
	Country    sql.NullString `json:"country"`
}

func (addr AddressSql) toAddress() Address {
	address := Address{}
	if addr.Street.Valid {
		address.Street = addr.Street.String
	}
	if addr.Complement.Valid {
		address.Complement = addr.Complement.String
	}
	if addr.City.Valid {
		address.City = addr.City.String
	}
	if addr.ZIP.Valid {
		address.ZIP = addr.ZIP.String
	}
	if addr.Country.Valid {
		address.Country = addr.Country.String
	}
	return address
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
