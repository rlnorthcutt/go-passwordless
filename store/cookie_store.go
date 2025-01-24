// file: store/cookie_store.go
package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

// The context keys we’ll use to retrieve http.Request and http.ResponseWriter
type ctxKey int

const (
	ctxKeyRequest ctxKey = iota
	ctxKeyResponse
)

func WithRequestResponse(ctx context.Context, r *http.Request, w http.ResponseWriter) context.Context {
	ctx = context.WithValue(ctx, ctxKeyRequest, r)
	ctx = context.WithValue(ctx, ctxKeyResponse, w)
	return ctx
}

// CookieStore implements the TokenStore interface using a signed (and optionally encrypted) cookie.
type CookieStore struct {
	// Name of the cookie to store the token data
	CookieName string
	// Securecookie instance for signing/encrypting
	SC *securecookie.SecureCookie
	// Default token expiry. If you also store it in the cookie data, rely on that as well.
	DefaultExpiry time.Duration
}

// TokenData is what we store in the cookie.
type TokenData struct {
	TokenID   string
	CodeHash  string
	ExpiresAt time.Time
}

// Store writes the token data into the user’s cookie.
func (cs *CookieStore) Store(ctx context.Context, tok Token) error {
	req, rsp := cs.getReqResp(ctx)
	if req == nil || rsp == nil {
		return fmt.Errorf("cookie store: no request/response in context")
	}

	data := TokenData{
		TokenID:   tok.ID,
		CodeHash:  hex.EncodeToString(tok.CodeHash),
		ExpiresAt: tok.ExpiresAt,
	}

	encoded, err := cs.SC.Encode(cs.CookieName, data)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  cs.CookieName,
		Value: encoded,
		Path:  "/",
		// Make sure to set Secure: true if using HTTPS
		HttpOnly: true,
		// If you have an explicit expiry or max-age:
		Expires: tok.ExpiresAt,
	}

	http.SetCookie(rsp, cookie)
	return nil
}

// Exists checks if the cookie is present and parses it.
// We won't strictly return a Token object in detail, but enough to confirm existence.
func (cs *CookieStore) Exists(ctx context.Context, tokenID string) (*Token, error) {
	data, err := cs.readCookie(ctx)
	if err != nil {
		return nil, fmt.Errorf("cookie store: %w", err)
	}
	if data == nil {
		return nil, fmt.Errorf("cookie not found or invalid")
	}
	if data.TokenID != tokenID {
		return nil, fmt.Errorf("tokenID mismatch in cookie")
	}

	codeHash, _ := hex.DecodeString(data.CodeHash)
	return &Token{
		ID:        data.TokenID,
		CodeHash:  codeHash,
		ExpiresAt: data.ExpiresAt,
	}, nil
}

// Verify compares the user-provided code with the hashed code in the cookie.
func (cs *CookieStore) Verify(ctx context.Context, tokenID, code string) (bool, error) {
	data, err := cs.readCookie(ctx)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, fmt.Errorf("no token data in cookie")
	}

	// Check the token ID
	if data.TokenID != tokenID {
		return false, fmt.Errorf("token ID does not match")
	}
	// Check expiry
	if time.Now().After(data.ExpiresAt) {
		// Remove the cookie
		_ = cs.Delete(ctx, tokenID)
		return false, fmt.Errorf("cookie token expired")
	}

	// Compare code hash
	rawHash := sha256.Sum256([]byte(code))
	if !strings.EqualFold(data.CodeHash, hex.EncodeToString(rawHash[:])) {
		return false, nil
	}

	// If everything matches, remove the cookie to consume it
	if err := cs.Delete(ctx, tokenID); err != nil {
		return false, err
	}
	return true, nil
}

// Delete removes the cookie, effectively invalidating the token.
func (cs *CookieStore) Delete(ctx context.Context, tokenID string) error {
	// We ignore tokenID here, because there's only 1 token stored in the cookie.
	_, rsp := cs.getReqResp(ctx)
	if rsp == nil {
		return fmt.Errorf("no response in context")
	}
	// Overwrite cookie with empty value & past expiry
	cookie := &http.Cookie{
		Name:    cs.CookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(rsp, cookie)
	return nil
}

// Helper to extract request/response from context
func (cs *CookieStore) getReqResp(ctx context.Context) (*http.Request, http.ResponseWriter) {
	req, _ := ctx.Value(ctxKeyRequest).(*http.Request)
	rsp, _ := ctx.Value(ctxKeyResponse).(http.ResponseWriter)
	return req, rsp
}

// readCookie decodes the stored cookie
func (cs *CookieStore) readCookie(ctx context.Context) (*TokenData, error) {
	req, _ := cs.getReqResp(ctx)
	if req == nil {
		return nil, fmt.Errorf("no request found in context")
	}
	c, err := req.Cookie(cs.CookieName)
	if err != nil {
		return nil, err
	}
	var data TokenData
	if err := cs.SC.Decode(cs.CookieName, c.Value, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
