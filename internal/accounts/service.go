package accounts

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateAccount(ctx context.Context, input CreateAccountInput) (Account, error) {
	ownerName := strings.TrimSpace(input.OwnerName)
	if ownerName == "" {
		return Account{}, ErrOwnerNameRequired
	}

	return s.store.Create(ctx, uuid.NewString(), ownerName)
}

func (s *Service) GetAccountByID(ctx context.Context, id string) (Account, error) {
	return s.store.GetByID(ctx, id)
}
