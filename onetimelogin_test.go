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

func TestLoginLinkFlow(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	email := "user@example.com"
	baseURL := "https://myapp.com/login"

	t.Run("GenerateLoginLink", func(t *testing.T) {
		loginLink, err := mgr.GenerateLoginLink(ctx, email, baseURL)
		if err != nil {
			t.Fatalf("Failed to generate login link: %v", err)
		}
		if loginLink == "" {
			t.Fatal("Expected a non-empty login link")
		}
		t.Logf("Generated login link: %s", loginLink)
	})

	t.Run("VerifyLoginLink", func(t *testing.T) {
		loginLink, err := mgr.GenerateLoginLink(ctx, email, baseURL)
		if err != nil {
			t.Fatalf("Failed to generate login link: %v", err)
		}

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

		t.Run("FailIncorrectHash", func(t *testing.T) {
			success, err := mgr.VerifyLoginLink(ctx, token, "wronghash")
			if err == nil {
				t.Fatal("Expected an error for incorrect hash, got nil")
			}
			if success {
				t.Fatal("Expected verification to fail for incorrect hash")
			}
		})

		t.Run("SuccessCorrectHash", func(t *testing.T) {
			success, err := mgr.VerifyLoginLink(ctx, token, hash)
			if err != nil {
				t.Fatalf("Verification failed: %v", err)
			}
			if !success {
				t.Fatal("Expected successful verification")
			}
		})
	})

	t.Run("VerifyExpiredLoginLink", func(t *testing.T) {
		loginLink, err := mgr.GenerateLoginLink(ctx, email, baseURL)
		if err != nil {
			t.Fatalf("Failed to generate login link: %v", err)
		}

		parsedURL, _ := url.Parse(loginLink)
		queryParams := parsedURL.Query()
		token := queryParams.Get("token")
		hash := queryParams.Get("hash")

		tok, _ := memStore.Exists(ctx, token)
		tok.ExpiresAt = time.Now().Add(-1 * time.Minute) // Expire the token
		_ = memStore.Store(ctx, *tok)

		success, err := mgr.VerifyLoginLink(ctx, token, hash)
		if err == nil {
			t.Fatal("Expected an error for expired token, got nil")
		}
		if success {
			t.Fatalf("Expected verification to fail for expired token, got %v", success)
		}
	})
}
