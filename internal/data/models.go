package data

import (
	"database/sql"
	"errors"
)

const (
	UserToActivate = "to-activate"
	UserActivated  = "activated"
	
	TokenActivation = "activation"
	TokenReset      = "reset"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	TokenModel       *TokenModel
	UserModel        *UserModel
	ClientModel      *ClientModel
	CarProductModel  *CarProductModel
	CarCatalogModel  *CarCatalogModel
	TransactionModel *TransactionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		TokenModel:       &TokenModel{db},
		UserModel:        &UserModel{db},
		ClientModel:      &ClientModel{db},
		CarProductModel:  &CarProductModel{db},
		CarCatalogModel:  &CarCatalogModel{db},
		TransactionModel: &TransactionModel{db},
	}
}
