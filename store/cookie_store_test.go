package store_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/store"
)

func TestCookieStore(t *testing.T) {
	secretKey := []byte("super-secret-key")
	cs := store.NewCookieStore(secretKey)

	tokenID := "testid"
	testToken := store.Token{
		ID:        tokenID,
		Recipient: "cookie@test",
		CodeHash:  []byte("fake-hash"),
		ExpiresAt: time.Now().Add(1 * time.Minute),
		CreatedAt: time.Now(),
	}

	var cookies []*http.Cookie

	t.Run("StoreToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		err := cs.Store(ctx, testToken)
		if err != nil {
			t.Fatalf("Failed to store session: %v", err)
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

		tok, err := cs.Exists(ctx, tokenID)
		if err != nil {
			t.Fatalf("Failed to retrieve token: %v", err)
		}
		t.Logf("Retrieved token ID: %s", tok.ID)

		if tok.ID != tokenID {
			t.Fatalf("Expected token ID '%s', got '%s'", tokenID, tok.ID)
		}
	})

	t.Run("DeleteToken", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(cookies[0])
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		err := cs.Delete(ctx, tokenID)
		if err != nil {
			t.Fatalf("Failed to delete token: %v", err)
		}
		t.Logf("Token ID: %s successfully deleted", tokenID)
	})

	t.Run("VerifyDeletion", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		ctx := store.WithRequestResponse(context.Background(), req, w)

		token, _ := cs.Exists(ctx, tokenID)
		if token != nil {
			t.Fatal("Expected token to be deleted but it still exists")
		}
		t.Logf("Token correctly deleted")
	})
}
