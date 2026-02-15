package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_server "github.com/aplavrov/wallet/internal/server/mocks"
	"github.com/aplavrov/wallet/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestServer_Operation_JSON_validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Parallel()
	tests := []struct {
		name          string
		body          string
		wantCode      int
		wantJSONError string
	}{
		{
			name: "invalid walletID",
			body: `{
				"walletId": "12345",
				"amount": 100,
				"operationType": "DEPOSIT"
			}`,
			wantCode:      http.StatusBadRequest,
			wantJSONError: `{"error":"invalid walletID"}`,
		},
		{
			name: "negative deposit",
			body: `{
				"walletId": "550e8400-e29b-41d4-a716-446655440000",
				"amount": -10,
				"operationType": "DEPOSIT"
			}`,
			wantCode:      http.StatusBadRequest,
			wantJSONError: `{"error":"amount must be positive"}`,
		},
		{
			name: "negative withdraw",
			body: `{
				"walletId": "550e8400-e29b-41d4-a716-446655440000",
				"amount": -10,
				"operationType": "WITHDRAW"
			}`,
			wantCode:      http.StatusBadRequest,
			wantJSONError: `{"error":"amount must be positive"}`,
		},
		{
			name: "invalid operation type",
			body: `{
				"walletId": "550e8400-e29b-41d4-a716-446655440000",
				"amount": 100,
				"operationType": "UNKNOWN"
			}`,
			wantCode:      http.StatusBadRequest,
			wantJSONError: `{"error":"invalid operation type"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPost, "/operation", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			server := &WalletServer{}
			server.Operation(ctx)

			require.Equal(t, tt.wantCode, w.Code)
			require.JSONEq(t, tt.wantJSONError, w.Body.String())
		})
	}
}

func TestServer_Operation_withdraw_logic_validation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Parallel()
	tests := []struct {
		name          string
		body          string
		wantCode      int
		wantJSONError string
		wantError     error
	}{
		{
			name: "wallet not found",
			body: `{
				"walletId": "550e8400-e29b-41d4-a716-446655440000",
				"amount": 100,
				"operationType": "WITHDRAW"
			}`,
			wantCode:      http.StatusNotFound,
			wantError:     service.ErrWalletNotFound,
			wantJSONError: `{"error": "wallet not found"}`,
		},
		{
			name: "not enough balance",
			body: `{
				"walletId": "550e8400-e29b-41d4-a716-446655440000",
				"amount": 100,
				"operationType": "WITHDRAW"
			}`,
			wantCode:      http.StatusConflict,
			wantError:     service.ErrNotEnoughMoney,
			wantJSONError: `{"error": "not enough money on a balance"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodPost, "/operation", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req

			ctrl := gomock.NewController(t)
			mockService := mock_server.NewMockwalletService(ctrl)
			mockService.EXPECT().Withdraw(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.wantError)
			server := NewWalletServer(mockService)
			server.Operation(ctx)

			require.Equal(t, tt.wantCode, w.Code)
			require.JSONEq(t, tt.wantJSONError, w.Body.String())
		})
	}
}
