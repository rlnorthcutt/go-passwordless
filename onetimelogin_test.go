package passwordless_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

func TestGenerateLoginLink(t *testing.T) {
	ctx := context.Background()

	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	email := "user@example.com"
	baseURL := "https://myapp.com/login"

	// Generate login link
	loginLink, err := mgr.GenerateLoginLink(ctx, email, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate login link: %v", err)
	}
	if loginLink == "" {
		t.Fatal("Expected a non-empty login link")
	}
	t.Logf("Generated login link: %s", loginLink)
}

func TestVerifyLoginLink(t *testing.T) {
	ctx := context.Background()

	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	email := "user@example.com"
	baseURL := "https://myapp.com/login"

	// Generate login link
	loginLink, err := mgr.GenerateLoginLink(ctx, email, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate login link: %v", err)
	}

	// Extract token and hash from the generated URL
	parsedURL, err := url.Parse(loginLink)
	if err != nil {
		t.Fatalf("Failed to parse login link: %v", err)
	}
	queryParams := parsedURL.Query()
	token := queryParams.Get("token")
	hash := queryParams.Get("hash")

	if token == "" || hash == "" {
		t.Fatal("Missing token or hash in login link")
	}

	// Simulate failed verification with incorrect hash (before successful verification)
	success, err := mgr.VerifyLoginLink(ctx, token, "wronghash")
	if err == nil {
		t.Fatal("Expected an error for incorrect hash, got nil")
	}
	if success {
		t.Fatal("Expected verification to fail for incorrect hash")
	}

	// Simulate verification with correct hash
	success, err = mgr.VerifyLoginLink(ctx, token, hash)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}
	if !success {
		t.Fatal("Expected successful verification")
	}
}

func TestVerifyExpiredLoginLink(t *testing.T) {
	ctx := context.Background()

	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	email := "user@example.com"
	baseURL := "https://myapp.com/login"

	newTokenID, err := mgr.GenerateLoginLink(ctx, email, baseURL)
	if err != nil {
		t.Fatalf("Failed to generate new login link: %v", err)
	}

	parsedURL, _ := url.Parse(newTokenID)
	queryParams := parsedURL.Query()
	newToken := queryParams.Get("token")
	newHash := queryParams.Get("hash")

	tok, _ := memStore.Exists(ctx, newToken)
	tok.ExpiresAt = time.Now().Add(-1 * time.Minute)
	_ = memStore.Store(ctx, *tok)

	success, err := mgr.VerifyLoginLink(ctx, newToken, newHash)
	if err == nil {
		t.Fatal("Expected error missing when verifying an expired token")
	}
	if success {
		t.Fatalf("Expected verification to fail for expired token, got %v", success)
	}
}
