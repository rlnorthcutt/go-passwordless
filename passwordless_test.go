package passwordless_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

func TestPasswordlessFlow(t *testing.T) {
	ctx := context.Background()

	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}

	mgr := passwordless.NewManager(memStore, logTransport)

	email := "user@example.com"
	tokenID, err := mgr.StartLogin(ctx, email)
	if err != nil {
		t.Fatalf("StartLogin returned error: %v", err)
	}
	if tokenID == "" {
		t.Fatal("Expected a non-empty tokenID")
	}

	// Retrieve the stored token
	tok, err := memStore.Exists(ctx, tokenID)
	if err != nil {
		t.Fatalf("Failed to retrieve token: %v", err)
	}
	if tok.Recipient != email {
		t.Errorf("Expected recipient %q, got %q", email, tok.Recipient)
	}

	// Simulate verifying with correct code
	correctCodeHash := sha256.Sum256(tok.CodeHash)
	correctCode := hex.EncodeToString(correctCodeHash[:])

	success, err := mgr.VerifyLogin(ctx, tokenID, correctCode)
	if err != nil {
		t.Fatalf("Verification failed unexpectedly: %v", err)
	}
	if !success {
		t.Fatal("Expected verification to succeed")
	}

	// Simulate verifying with an incorrect code
	success, err = mgr.VerifyLogin(ctx, tokenID, "wrongcode")
	if err == nil {
		t.Fatal("Expected an error for incorrect code, got nil")
	}
	if success {
		t.Fatal("Expected verification to fail for incorrect code")
	}

	// Simulate token expiration
	tok.ExpiresAt = time.Now().Add(-1 * time.Minute)
	_ = memStore.Store(ctx, *tok)
	success, err = mgr.VerifyLogin(ctx, tokenID, correctCode)
	if err == nil {
		t.Fatal("Expected an error for expired token, got nil")
	}
	if success {
		t.Fatal("Expected verification to fail for expired token")
	}
}
