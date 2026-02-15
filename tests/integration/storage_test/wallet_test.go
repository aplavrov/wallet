//go:build integration

package storage_test

import (
	"context"
	"testing"

	"github.com/aplavrov/wallet/internal/storage/postgresql"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateWallet(t *testing.T) {
	db.SetUp(t, "wallets")
	defer db.TearDown(t)
	storage := postgresql.NewWalletStorage(db.DB)
	err := storage.CreateWallet(context.Background(), uuid.New())
	require.NoError(t, err)
}

func TestGetBalance(t *testing.T) {
	db.SetUp(t, "wallets")
	defer db.TearDown(t)

	walletID := uuid.New()
	balance := int64(100)
	fillDb(t, walletID, balance)
	storage := postgresql.NewWalletStorage(db.DB)

	resp, err := storage.GetBalance(context.Background(), walletID)
	require.NoError(t, err)
	assert.Equal(t, resp, balance)
}

func TestDeposit(t *testing.T) {
	db.SetUp(t, "wallets")
	defer db.TearDown(t)

	walletID := uuid.New()
	amount := int64(100)
	fillDb(t, walletID, 0)
	storage := postgresql.NewWalletStorage(db.DB)
	err := storage.Deposit(context.Background(), walletID, amount)
	require.NoError(t, err)

	resp, _ := storage.GetBalance(context.Background(), walletID)
	assert.Equal(t, resp, amount)
}

func TestWithdraw(t *testing.T) {
	db.SetUp(t, "wallets")
	defer db.TearDown(t)

	walletID := uuid.New()
	balance := int64(1000)
	amount := int64(100)
	fillDb(t, walletID, balance)
	storage := postgresql.NewWalletStorage(db.DB)
	err := storage.Withdraw(context.Background(), walletID, amount)
	require.NoError(t, err)

	resp, _ := storage.GetBalance(context.Background(), walletID)
	assert.Equal(t, resp, balance-amount)
}

func fillDb(t *testing.T, walletID uuid.UUID, balance int64) {
	t.Helper()
	_, err := db.DB.Exec(context.Background(), "INSERT INTO wallets(id, balance) VALUES ($1, $2)", walletID, balance)
	require.NoError(t, err)
}
