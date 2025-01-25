package store

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gorilla/sessions"
)

// FileStore manages passwordless tokens using Gorilla Sessions with file storage.
type FileStore struct {
	store         *sessions.FilesystemStore
	CookieName    string
	DefaultExpiry time.Duration
}

// NewFileStore initializes a new file-based session store.
func NewFileStore(path string, secretKey []byte) *FileStore {
	// Ensure the directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return nil
		}
	}

	fs := sessions.NewFilesystemStore(path, secretKey)
	fs.MaxAge(int(5 * time.Minute.Seconds())) // Default expiry

	return &FileStore{
		store:         fs,
		CookieName:    "pwdless_fsession",
		DefaultExpiry: 5 * time.Minute,
	}
}

// Store saves the token in the session.
func (fs *FileStore) Store(ctx context.Context, tok Token) error {
	req, rsp := getRequestResponse(ctx)
	if req == nil || rsp == nil {
		return errors.New("missing request or response in context")
	}

	session, _ := fs.store.Get(req, fs.CookieName)
	session.Values["tokenID"] = tok.ID
	session.Values["recipient"] = tok.Recipient
	session.Values["codeHash"] = tok.CodeHash
	session.Values["expiresAt"] = tok.ExpiresAt.Unix()

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(fs.DefaultExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true, // Set to false if testing locally without HTTPS
	}

	return session.Save(req, rsp)
}

// Exists checks if a token exists in the session.
func (fs *FileStore) Exists(ctx context.Context, tokenID string) (*Token, error) {
	req, _ := getRequestResponse(ctx)
	if req == nil {
		return nil, errors.New("missing request in context")
	}

	session, err := fs.store.Get(req, fs.CookieName)
	if err != nil {
		return nil, err
	}

	sessionTokenID, ok := session.Values["tokenID"].(string)
	if !ok || sessionTokenID != tokenID {
		return nil, errors.New("token not found")
	}

	expiresAtUnix, _ := session.Values["expiresAt"].(int64)
	if time.Now().After(time.Unix(expiresAtUnix, 0)) {
		return nil, errors.New("token expired")
	}

	return &Token{
		ID:        sessionTokenID,
		Recipient: session.Values["recipient"].(string),
		CodeHash:  session.Values["codeHash"].([]byte),
		ExpiresAt: time.Unix(expiresAtUnix, 0),
	}, nil
}

// Delete removes the session token from the store.
func (fs *FileStore) Delete(ctx context.Context, tokenID string) error {
	req, rsp := getRequestResponse(ctx)
	if req == nil || rsp == nil {
		return errors.New("missing request or response in context")
	}

	session, _ := fs.store.Get(req, fs.CookieName)
	if session.Values["tokenID"] != tokenID {
		return nil // No need to delete if the token isn't found
	}

	session.Options.MaxAge = -1 // Invalidate the session immediately
	return session.Save(req, rsp)
}
