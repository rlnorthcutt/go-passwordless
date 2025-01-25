package store_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/rlnorthcutt/go-passwordless/store"
	_ "modernc.org/sqlite"
)

func TestDbStore(t *testing.T) {
	dbFile := "test_tokens.db"

	// Cleanup before and after the test
	os.Remove(dbFile)
	defer os.Remove(dbFile)

	t.Logf("[DEBUG] Using database file: %s", dbFile)

	// Open the database using modernc.org/sqlite driver
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		t.Fatalf("Failed to open sqlite using modernc: %v", err)
	}
	defer db.Close()

	t.Logf("[DEBUG] Successfully opened database.")

	// Create tokens table
	_, err = db.Exec(`
        CREATE TABLE tokens (
            id TEXT PRIMARY KEY,
            recipient TEXT NOT NULL,
            code_hash BLOB NOT NULL,
            expires_at DATETIME NOT NULL,
            created_at DATETIME NOT NULL
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create tokens table: %v", err)
	}
	t.Logf("[DEBUG] Tokens table created successfully.")

	// Initialize DbStore
	dbStore := store.NewDbStore(db, "tokens")
	if dbStore == nil {
		t.Fatal("Failed to initialize DbStore")
	}
	t.Logf("[DEBUG] DbStore initialized.")

	// Run verification test
	t.Log("[DEBUG] Running verification test...")
	runStoreVerifyTest(t, dbStore)
	t.Log("[DEBUG] Verification test completed.")

	// Run expiry test
	t.Log("[DEBUG] Running expiry test...")
	runStoreExpiryTest(t, dbStore)
	t.Log("[DEBUG] Expiry test completed.")

	t.Logf("[DEBUG] TestDbStore completed successfully.")
}
