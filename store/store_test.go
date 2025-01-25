package store_test

import (
	"context"
	"crypto/sha256"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/store"
)

func TestTokenStore(t *testing.T) {
	// Initialize different store implementations
	mem := store.NewMemStore()
	// Add more stores as needed (e.g. DbStore, CookieStore)

	testers := map[string]store.TokenStore{
		"mem": mem,
	}

	for name, s := range testers {
		t.Run(name+"_store_verify", func(t *testing.T) {
			t.Logf("[DEBUG] Running store verification test for: %s", name)
			runStoreVerifyTest(t, s)
		})
		t.Run(name+"_expiry", func(t *testing.T) {
			t.Logf("[DEBUG] Running store expiry test for: %s", name)
			runStoreExpiryTest(t, s)
		})
	}
}

func runStoreVerifyTest(t *testing.T, s store.TokenStore) {
	ctx := context.Background()
	tokenID := "test-token"
	code := "123456"
	codeHash := sha256.Sum256([]byte(code))

	t.Logf("[DEBUG] Storing token ID: %s", tokenID)

	// Store the token
	err := s.Store(ctx, store.Token{
		ID:        tokenID,
		Recipient: "test@example.com",
		CodeHash:  codeHash[:],
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Store() error: %v", err)
	}

	// Verify correct code
	t.Logf("[DEBUG] Verifying token ID: %s", tokenID)
	ok, err := s.Verify(ctx, tokenID, code)
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}
	if !ok {
		t.Error("Verify() should succeed for correct code, got false")
	} else {
		t.Logf("[DEBUG] Token verification passed for ID: %s", tokenID)
	}

	// Verify again should fail (token consumed)
	ok2, err2 := s.Verify(ctx, tokenID, code)
	if err2 == nil && ok2 {
		t.Error("Verify() should fail after token has been consumed")
	} else {
		t.Logf("[DEBUG] Token ID: %s correctly consumed", tokenID)
	}
}

func runStoreExpiryTest(t *testing.T, s store.TokenStore) {
	ctx := context.Background()
	tokenID := "expire-token"
	code := "999999"
	codeHash := sha256.Sum256([]byte(code))

	t.Logf("[DEBUG] Storing expirable token ID: %s", tokenID)

	err := s.Store(ctx, store.Token{
		ID:        tokenID,
		Recipient: "expire@example.com",
		CodeHash:  codeHash[:],
		ExpiresAt: time.Now().Add(1 * time.Second),
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Store() error: %v", err)
	}

	t.Logf("[DEBUG] Waiting for token ID: %s to expire", tokenID)
	time.Sleep(2 * time.Second)

	// Try to verify expired token
	ok, err := s.Verify(ctx, tokenID, code)
	if err == nil {
		t.Error("Expected an error verifying expired token, got nil")
	} else {
		t.Logf("[DEBUG] Token ID: %s expired as expected", tokenID)
	}
	if ok {
		t.Error("Expired token should fail verification, got ok=true")
	}
}
