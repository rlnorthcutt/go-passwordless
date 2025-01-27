package store_test

import (
	"context"
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/store"
)

func TestFileStore(t *testing.T) {
	secretKey := []byte("super-secret-key")
	filePath := "./session_data"

	t.Logf("Starting TestFileStore...")

	// Ensure clean test environment
	os.RemoveAll(filePath)
	defer func() {
		t.Logf("Cleaning up test directory: %s", filePath)
		os.RemoveAll(filePath)
	}()

	fs := store.NewFileStore(filePath, secretKey)

	tokenID := "testid"
	code := "securecode"
	codeHash := sha256.Sum256([]byte(code))
	testToken := store.Token{
		ID:        tokenID,
		Recipient: "file@test",
		CodeHash:  codeHash[:],
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}

	var cookies []*http.Cookie

	t.Run("StoreToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		err := fs.Store(ctx, testToken)
		if err != nil {
			t.Fatalf("Failed to store token: %v", err)
		}

		cookies = w.Result().Cookies()
		if len(cookies) == 0 {
			t.Fatal("Expected cookie to be set, but none found")
		}
		t.Logf("Stored Cookie: %+v", cookies[0])
	})

	t.Run("RetrieveToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookies[0])
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		tok, err := fs.Exists(ctx, tokenID)
		if err != nil {
			t.Fatalf("Failed to retrieve token: %v", err)
		}
		t.Logf("Retrieved token ID: %s, Recipient: %s", tok.ID, tok.Recipient)

		if tok.ID != tokenID {
			t.Fatalf("Expected token ID '%s', got '%s'", tokenID, tok.ID)
		}
	})

	t.Run("VerifyToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookies[0])
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		valid, err := fs.Verify(ctx, tokenID, code)
		if err != nil {
			t.Fatalf("Failed to verify token: %v", err)
		}
		if !valid {
			t.Fatal("Expected token verification to succeed, but it failed")
		}
		t.Logf("Token ID: %s successfully verified", tokenID)
	})

	t.Run("VerifyDeletionAfterUse", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		token, _ := fs.Exists(ctx, tokenID)
		if token != nil {
			t.Fatal("Expected token to be deleted after verification, but it still exists")
		}
		t.Logf("Token ID: %s correctly deleted after verification", tokenID)
	})

	t.Run("DeleteTokenManually", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		err := fs.Delete(ctx, tokenID)
		if err != nil {
			t.Fatalf("Failed to delete token: %v", err)
		}
		t.Logf("Token ID: %s successfully deleted (or already deleted)", tokenID)
	})

	t.Logf("TestFileStore completed successfully")
}
