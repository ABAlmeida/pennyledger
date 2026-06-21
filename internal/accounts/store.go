package accounts

import (
	"context"
)

type Store interface {
	Create(ctx context.Context, id string, ownerName string) (Account, error)
	GetByID(ctx context.Context, id string) (Account, error)
}
