package service

import (
	"context"

	"github.com/google/uuid"
)

type storage interface {
	Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
	CreateWallet(ctx context.Context, walletID uuid.UUID) error
}

type WalletService struct {
	storage storage
}

func New(s storage) *WalletService {
	return &WalletService{
		storage: s,
	}
}

func (s *WalletService) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error {
	if err := s.storage.Deposit(ctx, walletID, amount); err != nil {
		return err
	}
	return nil
}

func (s *WalletService) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error {
	if err := s.storage.Withdraw(ctx, walletID, amount); err != nil {
		return err
	}
	return nil
}

func (s *WalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	balance, err := s.storage.GetBalance(ctx, walletID)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (s *WalletService) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	walletID := uuid.New()
	if err := s.storage.CreateWallet(ctx, walletID); err != nil {
		return uuid.Nil, err
	}
	return walletID, nil
}
