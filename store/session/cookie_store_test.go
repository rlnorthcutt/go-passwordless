package session_test

import (
	"context"
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/store/session"
)

func TestCookieStore(t *testing.T) {
	secretKey := []byte("super-secret-key")
	cs := session.NewCookieStore(secretKey)

	t.Run("PrimaryFlow", func(t *testing.T) {
		tokenID := "testid-primary"
		code := "secure-code"
		codeHash := sha256.Sum256([]byte(code))

		testToken := store.Token{
			ID:        tokenID,
			Recipient: "cookie@test",
			CodeHash:  codeHash[:],
			ExpiresAt: time.Now().Add(1 * time.Minute),
			CreatedAt: time.Now(),
		}

		var cookies []*http.Cookie

		t.Run("StoreToken", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			ctx := session.WithRequestResponse(context.Background(), req, w)

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
			ctx := session.WithRequestResponse(context.Background(), req, w)

			tok, err := cs.Exists(ctx, tokenID)
			if err != nil {
				t.Fatalf("Failed to retrieve token: %v", err)
			}
			t.Logf("Retrieved token ID: %s", tok.ID)

			if tok.ID != tokenID {
				t.Fatalf("Expected token ID '%s', got '%s'", tokenID, tok.ID)
			}
		})

		t.Run("VerifyToken", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.AddCookie(cookies[0])
			w := httptest.NewRecorder()
			ctx := session.WithRequestResponse(context.Background(), req, w)

			valid, err := cs.Verify(ctx, tokenID, code)
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
			ctx := session.WithRequestResponse(context.Background(), req, w)

			token, err := cs.Exists(ctx, tokenID)
			if err == nil && token != nil {
				t.Fatal("Expected token to be deleted but it still exists")
			}
			t.Logf("Token correctly deleted")
		})
	})

	t.Run("ManualDeletion", func(t *testing.T) {
		tokenID := "testid-delete"
		code := "delete-code"
		codeHash := sha256.Sum256([]byte(code))

		deleteToken := store.Token{
			ID:        tokenID,
			Recipient: "delete@test",
			CodeHash:  codeHash[:],
			ExpiresAt: time.Now().Add(5 * time.Minute),
			CreatedAt: time.Now(),
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		ctx := session.WithRequestResponse(context.Background(), req, w)

		err := cs.Store(ctx, deleteToken)
		if err != nil {
			t.Fatalf("Failed to store token for deletion test: %v", err)
		}

		req.AddCookie(w.Result().Cookies()[0])
		w = httptest.NewRecorder()
		ctx = session.WithRequestResponse(context.Background(), req, w)

		err = cs.Delete(ctx, tokenID)
		if err != nil {
			t.Fatalf("Failed to delete token: %v", err)
		}
		t.Logf("Token ID: %s successfully deleted", tokenID)

		tok, _ := cs.Exists(ctx, tokenID)
		if tok != nil {
			t.Fatal("Expected token to be deleted but it still exists")
		}
		t.Logf("Token ID: %s correctly deleted", tokenID)
	})
}
