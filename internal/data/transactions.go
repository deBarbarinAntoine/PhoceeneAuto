package data

import "time"

type Transaction struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Cars      []CarProduct
	Client    Client
	Status    string
	Leases    []float32
}
