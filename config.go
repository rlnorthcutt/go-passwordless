package passwordless

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Config holds all configurable aspects of the passwordless flow.
type Config struct {
	// CodeLength is the number of digits (or characters) in the generated code.
	// For example, 6 => "123456".
	CodeLength int

	// TokenExpiry specifies how long a generated token remains valid (e.g., 15 minutes).
	TokenExpiry time.Duration

	// IDGenerator is a function that returns a unique string identifier
	// for each token stored in the TokenStore.
	//
	// By default, this creates a 16-byte random hex string (32 hex chars),
	// but you can plug in your own function for different formats.
	IDGenerator func() string

	// CodeCharset is an optional set of characters to use when generating the code.
	// If empty, the library might default to digits 0-9.
	// For example, "0123456789" for numeric codes
	// or "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" for alphanumeric codes.
	CodeCharset string

	// MaxFailedAttempts is the maximum number of failed attempts allowed before the token is invalidated.
	MaxFailedAttempts int
}

// DefaultConfig provides sensible defaults for a typical passwordless flow.
func DefaultConfig() Config {
	return Config{
		CodeLength:        6,
		TokenExpiry:       15 * time.Minute,
		IDGenerator:       defaultIDGenerator,
		CodeCharset:       "0123456789", // numeric-only
		MaxFailedAttempts: 3,
	}
}

// defaultIDGenerator returns a random 16-byte hex string.
// This is used if the user doesn't supply a custom IDGenerator.
func defaultIDGenerator() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	// In production code, you might handle err or panic if it fails
	if err != nil {
		panic("failed to generate random bytes for ID")
	}
	return hex.EncodeToString(b)
}
