package session

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rlnorthcutt/go-passwordless/store"
)

// getContextRequestResponse extracts request and response from context with error handling.
func getContextRequestResponse(ctx context.Context) (*http.Request, http.ResponseWriter, error) {
	req, reqOk := ctx.Value(ctxKeyRequest).(*http.Request)
	rsp, rspOk := ctx.Value(ctxKeyResponse).(http.ResponseWriter)
	if !reqOk || !rspOk {
		return nil, nil, errors.New("missing request or response in context")
	}
	return req, rsp, nil
}

// setSessionValues is a helper to store common session data.
func setSessionValues(session *sessions.Session, tok store.Token) {
	session.Values["tokenID"] = tok.ID
	session.Values["recipient"] = tok.Recipient
	session.Values["codeHash"] = tok.CodeHash
	session.Values["expiresAt"] = tok.ExpiresAt.Unix()
	session.Values["createdAt"] = tok.CreatedAt.Unix()
	session.Values["attempts"] = tok.Attempts
}

// getSessionToken retrieves the token from session and checks expiration.
func getSessionToken(session *sessions.Session, tokenID string) (*store.Token, error) {
	if session.Values["tokenID"] != tokenID {
		return nil, errors.New("token not found")
	}

	expiresAtUnix, ok := session.Values["expiresAt"].(int64)
	if !ok {
		return nil, errors.New("invalid session expiration")
	}
	if time.Now().After(time.Unix(expiresAtUnix, 0)) { // Using the generic store method
		return nil, errors.New("token expired")
	}

	return &store.Token{
		ID:        session.Values["tokenID"].(string),
		Recipient: session.Values["recipient"].(string),
		CodeHash:  session.Values["codeHash"].([]byte),
		ExpiresAt: time.Unix(expiresAtUnix, 0),
		CreatedAt: time.Unix(session.Values["createdAt"].(int64), 0),
		Attempts:  session.Values["attempts"].(int),
	}, nil
}

// deleteSession invalidates the session and removes token data.
func deleteSession(ctx context.Context, store sessions.Store, cookieName, tokenID string) error {
	req, rsp, err := getContextRequestResponse(ctx)
	if err != nil {
		return err
	}

	session, _ := store.Get(req, cookieName)
	if session.Values["tokenID"] != tokenID {
		return nil // No need to delete if the token isn't found.
	}

	session.Options.MaxAge = -1 // Invalidate the session immediately.
	session.Values = make(map[interface{}]interface{})

	err = session.Save(req, rsp)
	return err
}
