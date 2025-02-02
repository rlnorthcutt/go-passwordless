package passwordless_test

import (
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("Default Configuration", func(t *testing.T) {
		dc := passwordless.DefaultConfig()

		t.Run("Code Length", func(t *testing.T) {
			if dc.CodeLength != 6 {
				t.Errorf("Expected CodeLength to be 6, got %d", dc.CodeLength)
			}
		})

		t.Run("Token Expiry", func(t *testing.T) {
			expected := 15 * time.Minute
			if dc.TokenExpiry != expected {
				t.Errorf("Expected TokenExpiry to be %v, got %v", expected, dc.TokenExpiry)
			}
		})

		t.Run("ID Generator", func(t *testing.T) {
			id := dc.IDGenerator()
			if len(id) != 32 { // 16 bytes in hex = 32 characters
				t.Errorf("Expected ID length to be 32, got %d", len(id))
			}
		})

		t.Run("Code Charset", func(t *testing.T) {
			expected := "0123456789"
			if dc.CodeCharset != expected {
				t.Errorf("Expected CodeCharset to be %q, got %q", expected, dc.CodeCharset)
			}
		})

		t.Run("Max Failed Attempts", func(t *testing.T) {
			if dc.MaxFailedAttempts != 3 {
				t.Errorf("Expected MaxFailedAttempts to be 3, got %d", dc.MaxFailedAttempts)
			}
		})
	})
}

func TestCustomIDGenerator(t *testing.T) {
	t.Run("Custom ID Generator", func(t *testing.T) {
		customGenerator := func() string {
			return "custom-id-123"
		}

		cfg := passwordless.Config{
			IDGenerator: customGenerator,
		}

		id := cfg.IDGenerator()
		if id != "custom-id-123" {
			t.Errorf("Expected custom ID generator to return 'custom-id-123', got %q", id)
		}
	})
}

func TestDefaultIDGeneratorUniqueness(t *testing.T) {
	t.Run("ID Generator Uniqueness", func(t *testing.T) {
		generatedIDs := make(map[string]bool)
		for i := 0; i < 1000; i++ {
			id := passwordless.DefaultConfig().IDGenerator()
			if generatedIDs[id] {
				t.Errorf("Duplicate ID generated: %s", id)
			}
			generatedIDs[id] = true
		}
	})
}
