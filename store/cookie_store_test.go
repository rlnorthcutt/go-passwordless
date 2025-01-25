// file: store/cookie_store_test.go
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

	// Step 1: Store a token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	ctx := store.WithRequestResponse(context.Background(), req, w)

	err := cs.Store(ctx, store.Token{
		ID:        "testid",
		Recipient: "cookie@test",
		CodeHash:  []byte("fake-hash"),
		ExpiresAt: time.Now().Add(1 * time.Minute),
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to store session: %v", err)
	}

	// Print cookies after storing
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected cookie to be set, but none found")
	}

	// Step 2: Retrieve token using the stored cookie
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.AddCookie(cookies[0])
	w2 := httptest.NewRecorder()
	ctx2 := store.WithRequestResponse(context.Background(), req2, w2)

	tok, err := cs.Exists(ctx2, "testid")
	if err != nil {
		t.Fatalf("failed to retrieve token: %v", err)
	}
	t.Logf("Retrieved token ID: %s", tok.ID)

	// Step 3: Delete the token
	err = cs.Delete(ctx2, "testid")
	if err != nil {
		t.Fatalf("failed to delete token: %v", err)
	}

	// Step 4: Inspect cookies after deletion
	respCookies := w2.Result().Cookies()
	t.Logf("Cookies after deletion: %+v", respCookies)

	// Step 5: Make a new request without the cookie to verify deletion
	req3 := httptest.NewRequest(http.MethodGet, "/", nil)
	w3 := httptest.NewRecorder()
	ctx3 := store.WithRequestResponse(context.Background(), req3, w3)

	_, err = cs.Exists(ctx3, "testid")
	if err == nil {
		t.Fatal("expected token to be deleted but it still exists")
	}
}
