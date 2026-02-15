package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func NewDB(ctx context.Context) (*Database, error) {
	if err := godotenv.Load("config.env"); err != nil {
		return nil, err
	}

	pool, err := pgxpool.Connect(ctx, GenerateDsn())
	if err != nil {
		return nil, err
	}

	return NewDatabase(pool), nil
}

func GenerateDsn() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
