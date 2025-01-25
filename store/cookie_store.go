package store

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
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

// getRequestResponse extracts the http.Request and http.ResponseWriter from the context.
func getRequestResponse(ctx context.Context) (*http.Request, http.ResponseWriter) {
	req, _ := ctx.Value(ctxKeyRequest).(*http.Request)
	rsp, _ := ctx.Value(ctxKeyResponse).(http.ResponseWriter)
	return req, rsp
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
func (cs *CookieStore) Store(ctx context.Context, tok Token) error {
	req, rsp := getRequestResponse(ctx)
	if req == nil || rsp == nil {
		return errors.New("missing request or response in context")
	}

	session, _ := cs.store.Get(req, cs.CookieName)
	session.Values["tokenID"] = tok.ID
	session.Values["recipient"] = tok.Recipient
	session.Values["codeHash"] = tok.CodeHash
	session.Values["expiresAt"] = tok.ExpiresAt.Unix()

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(cs.DefaultExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true, // Set to false if testing locally without HTTPS
	}

	return session.Save(req, rsp)
}

// Exists checks if a token exists in the session.
func (cs *CookieStore) Exists(ctx context.Context, tokenID string) (*Token, error) {
	req, _ := getRequestResponse(ctx)
	if req == nil {
		return nil, errors.New("missing request in context")
	}

	session, err := cs.store.Get(req, cs.CookieName)
	if err != nil {
		return nil, err
	}

	if session.Values["tokenID"] != tokenID {
		return nil, errors.New("token not found")
	}

	expiresAtUnix, ok := session.Values["expiresAt"].(int64)
	if !ok || time.Now().After(time.Unix(expiresAtUnix, 0)) {
		return nil, errors.New("token expired")
	}

	return &Token{
		ID:        session.Values["tokenID"].(string),
		Recipient: session.Values["recipient"].(string),
		CodeHash:  session.Values["codeHash"].([]byte),
		ExpiresAt: time.Unix(expiresAtUnix, 0),
	}, nil
}

// Delete removes the session token from the store.
func (cs *CookieStore) Delete(ctx context.Context, tokenID string) error {
	req, rsp := getRequestResponse(ctx)
	if req == nil || rsp == nil {
		return errors.New("missing request or response in context")
	}

	session, _ := cs.store.Get(req, cs.CookieName)
	if session.Values["tokenID"] != tokenID {
		return nil // If the token doesn't exist, no need to proceed.
	}

	// Set MaxAge to -1 to delete the session immediately.
	session.Options.MaxAge = -1 // Invalidate the session immediately.
	session.Values = make(map[interface{}]interface{})
	return session.Save(req, rsp)
}
