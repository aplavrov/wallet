//go:build integration

package postgresql

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aplavrov/wallet/internal/storage/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type TDB struct {
	DB db.DB
}

func NewTDB(ctx context.Context) (*TDB, error) {
	if err := godotenv.Load("../../../config.env"); err != nil {
		return nil, err
	}

	pool, err := pgxpool.Connect(ctx, generateTestDsn())
	if err != nil {
		return nil, err
	}

	return &TDB{DB: db.NewDatabase(pool)}, nil
}

func generateTestDsn() string {
	host := os.Getenv("DB_TEST_HOST")
	port := os.Getenv("DB_TEST_PORT")
	user := os.Getenv("DB_TEST_USER")
	password := os.Getenv("DB_TEST_PASSWORD")
	dbname := os.Getenv("DB_TEST_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func (d *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)
}

func (d *TDB) TearDown(t *testing.T) {
	t.Helper()
}

func (d *TDB) truncateTable(ctx context.Context, tableName ...string) {
	q := fmt.Sprintf("TRUNCATE %s", strings.Join(tableName, ","))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		log.Fatal(err)
	}
}
