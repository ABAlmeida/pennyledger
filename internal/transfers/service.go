package transfers

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateTransfer(ctx context.Context, input CreateTransferInput) (*Transfer, error) {
	if input.FromAccountID == "" {
		return nil, ErrFromAccountRequired
	}
	if input.ToAccountID == "" {
		return nil, ErrToAccountRequired
	}
	if input.AmountPence <= 0 {
		return nil, ErrAmountMustBePositive
	}
	if input.FromAccountID == input.ToAccountID {
		return nil, ErrSameAccount
	}

	transfer, err := s.store.Create(ctx, uuid.NewString(), input)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}
