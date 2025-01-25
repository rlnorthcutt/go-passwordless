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

	tokenID := "testid"
	code := "securecode"
	codeHash := sha256.Sum256([]byte(code))
	testToken := store.Token{
		ID:        tokenID,
		Recipient: "user@example.com",
		CodeHash:  codeHash[:],
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}

	// Store the token
	if err := memStore.Store(context.Background(), testToken); err != nil {
		t.Fatalf("Failed to store token: %v", err)
	}
	t.Logf("[DEBUG] Token stored with ID: %s", tokenID)

	// Verify token exists
	tok, err := memStore.Exists(context.Background(), tokenID)
	if err != nil {
		t.Fatalf("Token not found when it should exist: %v", err)
	}
	if tok.ID != tokenID {
		t.Errorf("Expected token ID %s, got %s", tokenID, tok.ID)
	}

	// Test token verification
	verified, err := memStore.Verify(context.Background(), tokenID, code)
	if err != nil {
		t.Fatalf("Error verifying token: %v", err)
	}
	if !verified {
		t.Fatalf("Expected token verification to succeed but it failed")
	}

	t.Logf("[DEBUG] Token ID %s successfully verified and deleted", tokenID)

	// Ensure the token does not exist anymore
	_, err = memStore.Exists(context.Background(), tokenID)
	if err == nil {
		t.Fatal("Expected token to be deleted but it still exists")
	}
}
