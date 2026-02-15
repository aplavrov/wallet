//go:build integration

package server_test

import (
	"context"
	"log"

	"github.com/aplavrov/wallet/internal/server"
	"github.com/aplavrov/wallet/internal/service"
	repository "github.com/aplavrov/wallet/internal/storage/postgresql"
	"github.com/aplavrov/wallet/tests/integration/postgresql"
	"github.com/gin-gonic/gin"
)

var (
	db            *postgresql.TDB
	srv           *server.WalletServer
	walletService *service.WalletService
)

func init() {
	var err error
	db, err = postgresql.NewTDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	storage := repository.NewWalletStorage(db.DB)
	walletService = service.New(storage)

	gin.SetMode(gin.TestMode)
	srv = server.NewWalletServer(walletService)
}
