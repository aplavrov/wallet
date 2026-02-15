//go:build integration

package server_test

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperation_Deposit(t *testing.T) {
	db.SetUp(t, "wallets")
	defer db.TearDown(t)

	walletID, err := walletService.CreateWallet(context.Background())
	require.NoError(t, err)

	amount := int64(100)
	body := fmt.Sprintf(`{"walletId": "%s", "operationType": "DEPOSIT", "amount": %d}`, walletID, amount)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	srv.Operation(ctx)

	require.Equal(t, http.StatusOK, w.Code)

	balance, err := walletService.GetBalance(context.Background(), walletID)
	require.NoError(t, err)
	require.Equal(t, amount, balance)
}

func TestOperation_Withdraw(t *testing.T) {
	db.SetUp(t, "wallets")
	defer db.TearDown(t)

	walletID, err := walletService.CreateWallet(context.Background())
	require.NoError(t, err)

	amount := int64(100)
	body := fmt.Sprintf(`{"walletId": "%s", "operationType": "WITHDRAW", "amount": %d}`, walletID, amount)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	srv.Operation(ctx)

	require.Equal(t, http.StatusConflict, w.Code)
}

func TestLoad_1000RPS(t *testing.T) {
	walletID, err := walletService.CreateWallet(context.Background())
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()

			opType := "DEPOSIT"
			amount := int64(rand.Intn(100) + 1)
			if rand.Intn(2) == 0 {
				opType = "WITHDRAW"
				amount = int64(rand.Intn(50) + 1)
			}

			body := fmt.Sprintf(`{"walletId": "%s", "operationType": "%s", "amount": %d}`, walletID, opType, amount)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			srv.Operation(ctx)

			assert.NotEqual(t, w.Code, http.StatusInternalServerError)
		}()
	}

	wg.Wait()
}
