package data

import (
	"time"
)

type Client struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Status    string
	Address   string
	Shop      string
}
