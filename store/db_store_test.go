package store_test

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/store"
	_ "modernc.org/sqlite"
)

func TestDbStore(t *testing.T) {
	dbFile := "test_tokens.db"
	sqlFile := "db_store_sample.sql"
	tableName := "passwordless_tokens"

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

	// Read table creation SQL from file
	sqlScript, err := os.ReadFile(sqlFile)
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}

	// Execute table creation script
	_, err = db.Exec(string(sqlScript))
	if err != nil {
		t.Fatalf("Failed to create %s table: %v", tableName, err)
	}
	t.Logf("[DEBUG] Tokens table created successfully using %s.", sqlFile)

	// Initialize DbStore
	dbStore := store.NewDbStore(db, tableName)
	if dbStore == nil {
		t.Fatal("Failed to initialize DbStore")
	}
	t.Logf("[DEBUG] DbStore initialized.")

	// Primary flow test
	t.Run("PrimaryFlow", func(t *testing.T) {
		tokenID := "testid-primary"
		code := "secure-code"
		codeHash := sha256.Sum256([]byte(code))

		testToken := store.Token{
			ID:        tokenID,
			Recipient: "dbstore@test",
			CodeHash:  codeHash[:],
			ExpiresAt: time.Now().Add(1 * time.Minute),
			CreatedAt: time.Now(),
		}

		t.Run("StoreToken", func(t *testing.T) {
			err := dbStore.Store(context.Background(), testToken)
			if err != nil {
				t.Fatalf("Failed to store token: %v", err)
			}
			t.Logf("Token stored successfully: %s", tokenID)
		})

		t.Run("RetrieveToken", func(t *testing.T) {
			tok, err := dbStore.Exists(context.Background(), tokenID)
			if err != nil {
				t.Fatalf("Failed to retrieve token: %v", err)
			}
			if tok.ID != tokenID {
				t.Fatalf("Expected token ID '%s', got '%s'", tokenID, tok.ID)
			}
			t.Logf("Token retrieved successfully: %s", tokenID)
		})

		t.Run("VerifyToken", func(t *testing.T) {
			valid, err := dbStore.Verify(context.Background(), tokenID, code)
			if err != nil {
				t.Fatalf("Failed to verify token: %v", err)
			}
			if !valid {
				t.Fatal("Expected token verification to succeed, but it failed")
			}
			t.Logf("Token verified successfully: %s", tokenID)
		})

		t.Run("VerifyDeletionAfterUse", func(t *testing.T) {
			_, err := dbStore.Exists(context.Background(), tokenID)
			if err == nil {
				t.Fatal("Expected token to be deleted after verification, but it still exists")
			}
			t.Logf("Token correctly deleted after use: %s", tokenID)
		})
	})

	// Manual deletion test
	t.Run("ManualDeletion", func(t *testing.T) {
		tokenID := "testid-delete"
		code := "delete-code"
		codeHash := sha256.Sum256([]byte(code))

		deleteToken := store.Token{
			ID:        tokenID,
			Recipient: "manual@test",
			CodeHash:  codeHash[:],
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		t.Run("StoreToken", func(t *testing.T) {
			err := dbStore.Store(context.Background(), deleteToken)
			if err != nil {
				t.Fatalf("Failed to store token for deletion test: %v", err)
			}
			t.Logf("Token stored successfully for deletion test: %s", tokenID)
		})

		t.Run("DeleteToken", func(t *testing.T) {
			err := dbStore.Delete(context.Background(), tokenID)
			if err != nil {
				t.Fatalf("Failed to delete token: %v", err)
			}
			t.Logf("Token deleted successfully: %s", tokenID)
		})

		t.Run("VerifyDeletion", func(t *testing.T) {
			_, err := dbStore.Exists(context.Background(), tokenID)
			if err == nil {
				t.Fatal("Expected token to be deleted but it still exists")
			}
			t.Logf("Token correctly deleted: %s", tokenID)
		})
	})
}
