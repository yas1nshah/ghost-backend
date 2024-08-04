package data

import (
	"database/sql"
	"errors"
)

type Models struct {
	Users    UserModel
	Tokens   TokenModel
	Listings ListingsModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Tokens:   TokenModel{DB: db},
		Listings: ListingsModel{DB: db},
	}
}

// Errors
var (
	ErrRecordNotFound      = errors.New("No Record Found")
	ErrEditConflict        = errors.New("Conflict in Edit")
	ErrListingLimitReached = errors.New("You have reached your Listings Limit")
)
