package postgresql

import (
	"context"
	"errors"

	"github.com/aplavrov/wallet/internal/service"
	"github.com/aplavrov/wallet/internal/storage/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type WalletStorage struct {
	db db.DB
}

func NewWalletStorage(database db.DB) *WalletStorage {
	return &WalletStorage{db: database}
}

func (s *WalletStorage) CreateWallet(ctx context.Context, walletID uuid.UUID) error {
	_, err := s.db.Exec(ctx, "INSERT INTO wallets (id) VALUES ($1)", walletID)
	return err
}

func (s *WalletStorage) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	var balance int64
	err := s.db.Get(ctx, &balance, "SELECT balance FROM wallets WHERE id=$1", walletID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, service.ErrWalletNotFound
		}
		return 0, err
	}
	return balance, err
}

func (s *WalletStorage) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error {
	tag, err := s.db.Exec(ctx, "UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, walletID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return service.ErrWalletNotFound
	}

	return nil
}

func (s *WalletStorage) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error {
	tag, err := s.db.Exec(ctx, "UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, walletID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() > 0 {
		return nil
	}

	var exists bool
	err = s.db.Get(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1)`,
		walletID,
	)
	if err != nil {
		return err
	}

	if !exists {
		return service.ErrWalletNotFound
	}

	return service.ErrNotEnoughMoney
}
