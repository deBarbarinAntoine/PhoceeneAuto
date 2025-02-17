package data

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Cars       []CarProduct
	User       User
	Client     Client
	Status     string
	Leases     []float32
	TotalPrice float32
}

// TODO -> Transaction check fields

type TransactionModel struct {
	db *sql.DB
}

func (m TransactionModel) Insert(transaction *Transaction) error {
	
	// TODO -> implement create transaction
	
	return nil
}

func (m TransactionModel) Update(transaction *Transaction) error {
	
	// TODO -> implement update transaction
	
	return nil
}

func (m TransactionModel) Delete(id uint) error {
	
	// TODO -> implement delete transaction
	
	return nil
}

func (m TransactionModel) GetByID(id uint) (*Transaction, error) {
	
	// TODO -> implement getByID transaction
	
	return nil, nil
}

func (m TransactionModel) Search(search string) ([]*Transaction, error) {
	
	// TODO -> implement search transaction
	
	return nil, nil
}
