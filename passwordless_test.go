package passwordless_test

import (
	"context"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

// Custom LogTransport that stores the code for testing
type TestTransport struct {
	LastCode string
}

func (tt *TestTransport) Send(ctx context.Context, recipient, code string) error {
	tt.LastCode = code // Capture the code here
	return nil
}

func TestPasswordlessFlow(t *testing.T) {
	ctx := context.Background()
	memStore := store.NewMemStore()
	testTransport := &TestTransport{}

	mgr := passwordless.NewManager(memStore, testTransport)

	email := "user@example.com"

	t.Run("StartLogin", func(t *testing.T) {
		tokenID, err := mgr.StartLogin(ctx, email)
		if err != nil {
			t.Fatalf("StartLogin returned error: %v", err)
		}
		if tokenID == "" {
			t.Fatal("Expected a non-empty tokenID")
		}

		t.Run("RetrieveToken", func(t *testing.T) {
			tok, err := memStore.Exists(ctx, tokenID)
			if err != nil {
				t.Fatalf("Failed to retrieve token: %v", err)
			}
			if tok.Recipient != email {
				t.Errorf("Expected recipient %q, got %q", email, tok.Recipient)
			}

			t.Run("VerifyIncorrectCode", func(t *testing.T) {
				success, err := mgr.VerifyLogin(ctx, tokenID, "wrongcode")
				if err == nil {
					t.Fatal("Expected an error for incorrect code, got nil")
				}
				if success {
					t.Fatal("Expected verification to fail for incorrect code")
				}
			})

			t.Run("VerifyCorrectCode", func(t *testing.T) {
				correctCode := testTransport.LastCode // Use captured code

				success, err := mgr.VerifyLogin(ctx, tokenID, correctCode)
				if err != nil {
					t.Fatalf("Verification failed unexpectedly: %v", err)
				}
				if !success {
					t.Fatal("Expected verification to succeed")
				}
			})
			t.Run("VerifyExpiredToken", func(t *testing.T) {
				tok.ExpiresAt = time.Now().Add(-1 * time.Minute)
				_ = memStore.Store(ctx, *tok)

				success, err := mgr.VerifyLogin(ctx, tokenID, testTransport.LastCode)
				if err == nil {
					t.Fatal("Expected an error for expired token, got nil")
				}
				if success {
					t.Fatal("Expected verification to fail for expired token")
				}
			})
		})
	})

	t.Run("NewManagerWithConfig", func(t *testing.T) {
		memStore := store.NewMemStore()
		logTransport := &transport.LogTransport{}

		// Start with blanks to test the defaults
		customConfig := passwordless.Config{
			CodeLength:        0,
			TokenExpiry:       0,
			IDGenerator:       func() string { return "customID" },
			CodeCharset:       "",
			MaxFailedAttempts: 0,
		}

		mgr := passwordless.NewManagerWithConfig(memStore, logTransport, customConfig)

		if mgr.Config.CodeLength != 6 {
			t.Errorf("Expected CodeLength to be 6, got %d", mgr.Config.CodeLength)
		}
		if mgr.Config.TokenExpiry != 15*time.Minute {
			t.Errorf("Expected TokenExpiry to be 15 minutes, got %v", mgr.Config.TokenExpiry)
		}
		if mgr.Config.IDGenerator() != "customID" {
			t.Errorf("Expected IDGenerator to return 'customID', got %s", mgr.Config.IDGenerator())
		}
		if mgr.Config.CodeCharset != "0123456789" {
			t.Errorf("Expected CodeCharset to be '0123456789', got %s", mgr.Config.CodeCharset)
		}
		if mgr.Config.MaxFailedAttempts != 3 {
			t.Errorf("Expected MaxFailedAttempts to be 3, got %d", mgr.Config.MaxFailedAttempts)
		}
	})
}
