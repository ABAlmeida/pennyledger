package transfers

import (
	"errors"
	"time"
)

type Transfer struct {
	ID            string
	FromAccountID string
	ToAccountID   string
	AmountPence   int64
	CreatedAt     time.Time
}

type CreateTransferInput struct {
	FromAccountID string
	ToAccountID   string
	AmountPence   int64
}

var ErrFromAccountRequired = errors.New("from account is required")
var ErrToAccountRequired = errors.New("to account is required")
var ErrAmountMustBePositive = errors.New("amount must be a positive value")
var ErrSameAccount = errors.New("from and to accounts cannot be the same")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrAccountNotFound = errors.New("account not found")
