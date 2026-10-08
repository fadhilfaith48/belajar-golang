package store

import (
	"os"
	"testing"
)

func TestPostgresStore_Skip(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL tidak diset; skip tes Postgres (offline-safe)")
	}
}
