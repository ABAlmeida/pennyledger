package accounts

import (
	"errors"
	"time"
)

type Account struct {
	ID           string
	OwnerName    string
	Currency     string
	BalancePence int64
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateAccountInput struct {
	OwnerName string
}

var ErrOwnerNameRequired = errors.New("owner name is required")
var ErrNotFound = errors.New("account not found")
