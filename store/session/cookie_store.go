package session

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rlnorthcutt/go-passwordless/store"
)

// contextKey is a custom type to avoid collisions with other context values.
type contextKey string

const (
	ctxKeyRequest  contextKey = "request"
	ctxKeyResponse contextKey = "response"
)

// WithRequestResponse adds the http.Request and http.ResponseWriter to the context.
func WithRequestResponse(ctx context.Context, r *http.Request, w http.ResponseWriter) context.Context {
	ctx = context.WithValue(ctx, ctxKeyRequest, r)
	ctx = context.WithValue(ctx, ctxKeyResponse, w)
	return ctx
}

// CookieStore manages passwordless tokens using Gorilla Sessions.
type CookieStore struct {
	store         *sessions.CookieStore
	CookieName    string
	DefaultExpiry time.Duration
}

// NewCookieStore initializes a new cookie store with encryption keys.
func NewCookieStore(secretKey []byte) *CookieStore {
	return &CookieStore{
		store:         sessions.NewCookieStore(secretKey),
		CookieName:    "pwdless_session",
		DefaultExpiry: 5 * time.Minute,
	}
}

// Store saves the token in the session.
func (cs *CookieStore) Store(ctx context.Context, tok store.Token) error {
	req, rsp, err := getContextRequestResponse(ctx)
	if err != nil {
		return err
	}

	session, _ := cs.store.Get(req, cs.CookieName)
	setSessionValues(session, tok)

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(cs.DefaultExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true,
	}

	return session.Save(req, rsp)
}

// Exists checks if a token exists in the session and removes it if expired.
func (cs *CookieStore) Exists(ctx context.Context, tokenID string) (*store.Token, error) {
	req, _, err := getContextRequestResponse(ctx)
	if err != nil {
		return nil, err
	}

	session, err := cs.store.Get(req, cs.CookieName)
	if err != nil {
		return nil, errors.New("failed to retrieve session")
	}

	tok, err := getSessionToken(session, tokenID)
	if err != nil {
		// Token not found or expired
		return nil, cs.Delete(ctx, tokenID)
	}
	return tok, nil
}

// Verify checks if the provided code matches the stored token's hash.
func (cs *CookieStore) Verify(ctx context.Context, tokenID, code string) (bool, error) {
	tok, err := cs.Exists(ctx, tokenID)
	if err != nil {
		return false, err
	}

	// Hash the provided code and compare with stored hash
	codeHash := sha256.Sum256([]byte(code))
	if !bytes.Equal(tok.CodeHash, codeHash[:]) {
		tok.Attempts++
		_ = cs.Store(ctx, *tok) // Save attempts count
		return false, errors.New("invalid code")
	}

	// Token is verified; remove it for one-time use.
	return true, cs.Delete(ctx, tokenID)
}

// Delete removes the session token from the store.
func (cs *CookieStore) Delete(ctx context.Context, tokenID string) error {
	return deleteSession(ctx, cs.store, cs.CookieName, tokenID)
}
