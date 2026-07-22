package httpapi

import (
	"context"

	"github.com/ABAlmeida/pennyledger/internal/transfers"
)

type transferService interface {
	CreateTransfer(ctx context.Context, input transfers.CreateTransferInput) (*transfers.Transfer, error)
}
