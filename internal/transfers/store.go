package transfers

import "context"

type Store interface {
	Create(ctx context.Context, id string, input CreateTransferInput) (Transfer, error)
}
