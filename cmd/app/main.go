package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aplavrov/wallet/internal/server"
	"github.com/aplavrov/wallet/internal/service"
	"github.com/aplavrov/wallet/internal/storage/db"
	"github.com/aplavrov/wallet/internal/storage/postgresql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbPool, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.GetPool().Close()

	sqlDB, err := sql.Open("pgx", db.GenerateDsn())
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	for i := 0; i < 10; i++ {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		log.Println("Waiting for PostgreSQL...")
		time.Sleep(2 * time.Second)
	}

	if err := goose.Up(sqlDB, "internal/storage/db/migrations"); err != nil {
		log.Fatal(err)
	}

	storage := postgresql.NewWalletStorage(dbPool)
	walletService := service.New(storage)
	handler := server.NewWalletServer(walletService)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-exit
		log.Println("Работа утилиты завершена")
		cancel()
	}()

	go func() {
		log.Println("Сервер запущен на :9000")
		if err := handler.Start(":9000"); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
