package passwordless_test

import (
	"context"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

func TestManager_BasicFlow(t *testing.T) {
	t.Logf("[DEBUG] Starting TestManager_BasicFlow...")

	ctx := context.Background()

	// Step 1: Setup memory store and log transport
	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	// Step 2: Start login process
	recipient := "user@example.com"
	tokenID, err := mgr.StartLogin(ctx, recipient)
	if err != nil {
		t.Fatalf("StartLogin returned error: %v", err)
	}
	if tokenID == "" {
		t.Fatal("Expected a non-empty tokenID")
	}
	t.Logf("[DEBUG] Login started successfully, tokenID: %s", tokenID)

	// Step 3: Verify token existence in store
	tok, err := memStore.Exists(ctx, tokenID)
	if err != nil {
		t.Fatalf("Failed to retrieve token from store: %v", err)
	}
	if tok.Recipient != recipient {
		t.Errorf("Expected recipient %q, got %q", recipient, tok.Recipient)
	}
	t.Logf("[DEBUG] Token retrieved successfully for recipient: %s", tok.Recipient)

	// Step 4: Attempt to verify login with an incorrect code
	fakeCode := "999999"
	ok, err := mgr.VerifyLogin(ctx, tokenID, fakeCode)
	if err == nil && ok {
		t.Error("VerifyLogin should have failed with an incorrect code, but succeeded.")
	}
	t.Logf("[DEBUG] VerifyLogin correctly failed with incorrect code.")

	// Step 5: Remove token after verification attempt
	if err := memStore.Delete(ctx, tokenID); err != nil {
		t.Fatalf("Failed to delete token: %v", err)
	}
	t.Logf("[DEBUG] Token deleted successfully, tokenID: %s", tokenID)

	t.Logf("[DEBUG] TestManager_BasicFlow completed successfully.")
}

func TestManager_Expiry(t *testing.T) {
	t.Logf("[DEBUG] Starting TestManager_Expiry...")

	ctx := context.Background()

	// Step 1: Setup memory store and log transport with short expiry config
	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}

	cfg := passwordless.Config{
		CodeLength:  4,
		TokenExpiry: 1 * time.Second, // expires quickly
	}
	mgr := passwordless.NewManagerWithConfig(memStore, logTransport, cfg)

	// Step 2: Start login with short expiry
	recipient := "user@example.com"
	tokenID, err := mgr.StartLogin(ctx, recipient)
	if err != nil {
		t.Fatalf("StartLogin error: %v", err)
	}
	t.Logf("[DEBUG] Login started successfully, tokenID: %s", tokenID)

	// Step 3: Wait for token expiry
	t.Logf("[DEBUG] Waiting for token to expire...")
	time.Sleep(2 * time.Second)

	// Step 4: Attempt verification after expiry
	ok, err := mgr.VerifyLogin(ctx, tokenID, "1234")
	if err == nil {
		t.Error("Expected an error due to expired token, got nil")
	}
	if ok {
		t.Error("Expected verify to fail due to expiry, but it succeeded.")
	} else {
		t.Logf("[DEBUG] Verification correctly failed due to expiry.")
	}

	t.Logf("[DEBUG] TestManager_Expiry completed successfully.")
}
