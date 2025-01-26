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

	// Step 1: Store a token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := store.WithRequestResponse(context.Background(), req, w)

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

	err := fs.Store(ctx, testToken)
	if err != nil {
		t.Fatalf("Failed to store session: %v", err)
	}
	t.Logf("Stored token with ID: %s", testToken.ID)

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected cookie to be set, but none found")
	}
	t.Logf("Stored Cookie: %+v", cookies[0])

	// Step 2: Retrieve token using the stored cookie
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.AddCookie(cookies[0])
	w2 := httptest.NewRecorder()
	ctx2 := store.WithRequestResponse(context.Background(), req2, w2)

	t.Logf("Retrieving token with ID: %s", tokenID)
	tok, err := fs.Exists(ctx2, "testid")
	if err != nil {
		t.Fatalf("Failed to retrieve token: %v", err)
	}
	t.Logf("Retrieved token ID: %s, Recipient: %s", tok.ID, tok.Recipient)

	if tok.ID != "testid" {
		t.Fatalf("Expected token ID 'testid', got '%s'", tok.ID)
	}

	// Step 3: Verify the token
	t.Logf("Verifying token with ID: %s", tokenID)
	valid, err := fs.Verify(ctx2, tokenID, code)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}
	if !valid {
		t.Fatal("Expected token verification to succeed, but it failed")
	}
	t.Logf("Token ID: %s successfully verified", tokenID)

	// Step 4: Ensure token does not exist after verification (one-time use)
	t.Logf("Ensuring token does not exist after verification...")
	_, err = fs.Exists(ctx2, tokenID)
	if err == nil {
		t.Fatal("Expected token to be deleted after verification, but it still exists")
	}
	t.Logf("Token ID: %s correctly deleted after verification", tokenID)

	// Step 5: Delete the token manually (should gracefully handle non-existing token)
	t.Logf("Deleting token with ID: %s", tokenID)
	err = fs.Delete(ctx2, tokenID)
	if err != nil {
		t.Fatalf("Failed to delete token: %v", err)
	}
	t.Logf("Token ID: %s successfully deleted (or already deleted)", tokenID)

	t.Logf("TestFileStore completed successfully")
}
