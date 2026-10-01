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

	// IncrementAttempts atomically increments the failed-attempt counter for the
	// given token ID and returns the new count. Implementations must perform the
	// read-modify-write as a single atomic operation (e.g. holding a lock, or
	// issuing an atomic `attempts = attempts + 1` update) so that concurrent
	// calls for the same token never lose an increment. Implementations should
	// not reset expiry or other fields when updating attempts.
	IncrementAttempts(ctx context.Context, tokenID string) (int, error)

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
	return SecureCompare(codeHash[:], tok.CodeHash)
}

// SecureCompare reports whether a and b are equal, using a constant-time
// comparison. subtle.ConstantTimeCompare already returns 0 for mismatched
// lengths, so no separate length check is needed before calling it.
func SecureCompare(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
