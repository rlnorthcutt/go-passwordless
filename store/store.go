package store

import (
	"context"
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

	// Verify checks if `code` matches the stored hash for tokenID, and
	// whether it's still valid. If valid, it may also consume or remove the token.
	Verify(ctx context.Context, tokenID, code string) (bool, error)

	// Delete permanently removes a token by ID (e.g. after verification).
	Delete(ctx context.Context, tokenID string) error
}
