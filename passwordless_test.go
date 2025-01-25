package passwordless_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

func TestManager_GenerateLoginLink(t *testing.T) {
	t.Logf("[DEBUG] Starting TestManager_GenerateLoginLink...")

	ctx := context.Background()

	// Step 1: Setup memory store and log transport
	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	// Step 2: Define test parameters
	baseURL := "https://myapp.com/login"
	recipient := "user@example.com"

	t.Logf("[DEBUG] Generating login link for recipient: %s", recipient)

	// Step 3: Generate the login link
	loginLink, err := mgr.GenerateLoginLink(ctx, recipient, baseURL)
	if err != nil {
		t.Fatalf("GenerateLoginLink returned error: %v", err)
	}

	// Step 4: Validate the generated URL
	parsedURL, err := url.Parse(loginLink)
	if err != nil {
		t.Fatalf("Failed to parse generated login link: %v", err)
	}

	// Step 5: Extract token query parameter
	token := parsedURL.Query().Get("token")
	if token == "" {
		t.Fatal("Expected non-empty token in generated URL")
	}

	t.Logf("[DEBUG] Login link generated successfully: %s", loginLink)
	t.Logf("[DEBUG] Extracted token: %s", token)

	// Step 6: Verify token exists in the store
	tok, err := memStore.Exists(ctx, token)
	if err != nil {
		t.Fatalf("Failed to retrieve token from store: %v", err)
	}
	if tok.Recipient != recipient {
		t.Errorf("Expected recipient %q, got %q", recipient, tok.Recipient)
	}

	t.Logf("[DEBUG] Token retrieved successfully for recipient: %s", tok.Recipient)

	t.Logf("[DEBUG] TestManager_GenerateLoginLink completed successfully.")
}

func TestManager_GenerateLoginLink_InvalidBaseURL(t *testing.T) {
	t.Logf("[DEBUG] Starting TestManager_GenerateLoginLink_InvalidBaseURL...")

	ctx := context.Background()
	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	recipient := "user@example.com"
	invalidBaseURL := "://invalid-url"

	t.Logf("[DEBUG] Attempting to generate login link with invalid base URL: %s", invalidBaseURL)

	_, err := mgr.GenerateLoginLink(ctx, recipient, invalidBaseURL)
	if err == nil {
		t.Fatal("Expected error for invalid base URL, but got none")
	}

	t.Logf("[DEBUG] Correctly received error for invalid base URL: %v", err)

	t.Logf("[DEBUG] TestManager_GenerateLoginLink_InvalidBaseURL completed successfully.")
}
