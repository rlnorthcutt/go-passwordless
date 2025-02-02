package store_test

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/store"
)

func TestMemStore(t *testing.T) {
	memStore := store.NewMemStore()

	t.Run("PrimaryFlow", func(t *testing.T) {
		tokenID := "testid-primary"
		code := "securecode"
		codeHash := sha256.Sum256([]byte(code))

		testToken := store.Token{
			ID:        tokenID,
			Recipient: "user@example.com",
			CodeHash:  codeHash[:],
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		t.Run("StoreToken", func(t *testing.T) {
			if err := memStore.Store(context.Background(), testToken); err != nil {
				t.Fatalf("Failed to store token: %v", err)
			}
			t.Logf("[DEBUG] Token stored with ID: %s", tokenID)
		})

		t.Run("RetrieveToken", func(t *testing.T) {
			tok, err := memStore.Exists(context.Background(), tokenID)
			if err != nil {
				t.Fatalf("Token not found when it should exist: %v", err)
			}
			if tok.ID != tokenID {
				t.Errorf("Expected token ID %s, got %s", tokenID, tok.ID)
			}
			t.Logf("[DEBUG] Retrieved token with ID: %s", tok.ID)
		})

		t.Run("VerifyToken", func(t *testing.T) {
			verified, err := memStore.Verify(context.Background(), tokenID, code)
			if err != nil {
				t.Fatalf("Error verifying token: %v", err)
			}
			if !verified {
				t.Fatalf("Expected token verification to succeed but it failed")
			}
			t.Logf("[DEBUG] Token ID %s successfully verified and deleted", tokenID)
		})

		t.Run("VerifyDeletionAfterUse", func(t *testing.T) {
			_, err := memStore.Exists(context.Background(), tokenID)
			if err == nil {
				t.Fatal("Expected token to be deleted but it still exists")
			}
			t.Logf("[DEBUG] Token correctly deleted after verification")
		})
	})

	t.Run("ManualDeletion", func(t *testing.T) {
		tokenID := "testid-manual"
		code := "deletecode"
		codeHash := sha256.Sum256([]byte(code))

		testToken := store.Token{
			ID:        tokenID,
			Recipient: "delete@example.com",
			CodeHash:  codeHash[:],
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		t.Run("StoreToken", func(t *testing.T) {
			if err := memStore.Store(context.Background(), testToken); err != nil {
				t.Fatalf("Failed to store token for deletion: %v", err)
			}
			t.Logf("[DEBUG] Token stored with ID: %s for deletion", tokenID)
		})

		t.Run("DeleteToken", func(t *testing.T) {
			if err := memStore.Delete(context.Background(), tokenID); err != nil {
				t.Fatalf("Failed to delete token: %v", err)
			}
			t.Logf("[DEBUG] Token ID %s successfully deleted", tokenID)
		})

		t.Run("VerifyDeletion", func(t *testing.T) {
			_, err := memStore.Exists(context.Background(), tokenID)
			if err == nil {
				t.Fatal("Expected token to be deleted but it still exists")
			}
			t.Logf("[DEBUG] Token ID %s correctly deleted after manual deletion", tokenID)
		})
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		tokenID := "testid-expired"
		code := "expiredcode"
		codeHash := sha256.Sum256([]byte(code))

		expiredToken := store.Token{
			ID:        tokenID,
			Recipient: "expired@example.com",
			CodeHash:  codeHash[:],
			ExpiresAt: time.Now().Add(-1 * time.Minute), // Expired token
			CreatedAt: time.Now(),
		}

		t.Run("StoreExpiredToken", func(t *testing.T) {
			if err := memStore.Store(context.Background(), expiredToken); err != nil {
				t.Fatalf("Failed to store expired token: %v", err)
			}
			t.Logf("[DEBUG] Expired token stored with ID: %s", tokenID)
		})

		t.Run("VerifyExpiredToken", func(t *testing.T) {
			verified, err := memStore.Verify(context.Background(), tokenID, code)
			if err == nil || verified {
				t.Fatalf("Expected verification of expired token to fail, but it succeeded")
			}
			t.Logf("[DEBUG] Expired token verification correctly failed")
		})

		t.Run("VerifyDeletionAfterExpiration", func(t *testing.T) {
			_, err := memStore.Exists(context.Background(), tokenID)
			if err == nil {
				t.Fatal("Expected expired token to be deleted but it still exists")
			}
			t.Logf("[DEBUG] Expired token ID %s correctly deleted after expiration", tokenID)
		})
	})
}
