package data

import "time"

type CarProduct struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     string
	Kilometers float32
	OwnerNb    uint
	Color      string
	Price      float32
	Shop       string
	CarCatalog
}
