package store

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"time"
)

// Token represents the stored code and associated data.
type Token struct {
	ID        string
	Recipient string
	CodeHash  []byte
	ExpiresAt time.Time
	CreatedAt time.Time
	Attempts  int // Track number of failed attempts
}

// TokenStore defines how tokens are saved, retrieved, verified, and deleted.
type TokenStore interface {
	// Store saves a Token. Implementation is responsible for hashing the code
	// or storing it already hashed, depending on design.
	Store(ctx context.Context, tok Token) error

	// Exists checks if a token with a given ID exists (and returns it if so).
	Exists(ctx context.Context, tokenID string) (*Token, error)

	// UpdateAttempts persists the latest failed-attempt count for the given token ID.
	// Implementations should not reset expiry or other fields when updating attempts.
	UpdateAttempts(ctx context.Context, tokenID string, attempts int) error

	// Verify checks if `code` matches the stored hash for tokenID, and
	// whether it's still valid. If valid, it may also consume or remove the token.
	Verify(ctx context.Context, tokenID, code string) (bool, error)

	// Delete permanently removes a token by ID (e.g. after verification).
	Delete(ctx context.Context, tokenID string) error
}

// Checks whether a given token is expired.
func IsTokenExpired(tok *Token) bool {
	return time.Now().After(tok.ExpiresAt)
}

// Verifies the provided code against the stored token's hash.
func VerifyToken(tok *Token, code string) bool {
	codeHash := sha256.Sum256([]byte(code))
	// Securely compares two hashes using constant-time comparison.
	if len(tok.CodeHash) != len(codeHash) {
		return false
	}
	return subtle.ConstantTimeCompare(codeHash[:], tok.CodeHash) == 1
}
