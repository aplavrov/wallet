//go:generate mockgen -source ./server.go -destination=./mocks/server.go -package=mock_server

package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/aplavrov/wallet/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type walletService interface {
	Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
	CreateWallet(ctx context.Context) (uuid.UUID, error)
}

type WalletServer struct {
	walletService walletService
}

func NewWalletServer(s walletService) *WalletServer {
	return &WalletServer{
		walletService: s,
	}
}

func (s *WalletServer) Start(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	server.POST("/api/v1/wallet", s.Operation)
	server.GET("/api/v1/wallets/:id", s.Balance)
	server.POST("/api/v1/wallets", s.CreateWallet)

	return server.Run(addr)
}

func (s *WalletServer) Operation(ctx *gin.Context) {
	var operation OperationRequest
	err := ctx.ShouldBindJSON(&operation)
	if err != nil {
		log.Printf("Error: invalid JSON: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	walletID, err := uuid.Parse(operation.WalletId)
	if err != nil {
		log.Printf("Error: invalid walletID: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid walletID"})
		return
	}

	if operation.Amount <= 0 {
		log.Println("Error: amount must be positive")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	switch operation.OperationType {
	case "DEPOSIT":
		err = s.walletService.Deposit(ctx.Request.Context(), walletID, operation.Amount)
	case "WITHDRAW":
		err = s.walletService.Withdraw(ctx.Request.Context(), walletID, operation.Amount)
	default:
		log.Println("Error: invalid operation type")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation type"})
		return
	}

	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			log.Printf("Error: %v", err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}
		if errors.Is(err, service.ErrNotEnoughMoney) {
			log.Printf("Error: %v", err.Error())
			ctx.JSON(http.StatusConflict, gin.H{"error": "not enough money on a balance"})
			return
		}
		log.Printf("Error: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	ctx.Status(http.StatusOK)
}

func (s *WalletServer) Balance(ctx *gin.Context) {
	walletID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Printf("Error: invalid walletID: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid walletID"})
		return
	}

	balance, err := s.walletService.GetBalance(ctx.Request.Context(), walletID)
	if err != nil {
		if errors.Is(err, service.ErrWalletNotFound) {
			log.Printf("Error: %v", err.Error())
			ctx.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}
		log.Printf("Error: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (s *WalletServer) CreateWallet(ctx *gin.Context) {
	walletID, err := s.walletService.CreateWallet(ctx.Request.Context())
	if err != nil {
		log.Printf("Error: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"walletID": walletID})
}
