//go:build integration

package storage_test

import (
	"context"
	"log"

	"github.com/aplavrov/wallet/tests/integration/postgresql"
)

var (
	db *postgresql.TDB
)

func init() {
	var err error
	db, err = postgresql.NewTDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
