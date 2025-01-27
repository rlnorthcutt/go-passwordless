package store

import (
	"bytes"
	"context"
	"crypto/sha256"
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
	req, rsp, err := getContextRequestResponse(ctx)
	if err != nil {
		return err
	}

	session, _ := fs.store.Get(req, fs.CookieName)
	setSessionValues(session, tok)

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(fs.DefaultExpiry.Seconds()),
		HttpOnly: true,
		Secure:   true,
	}

	return session.Save(req, rsp)
}

// Exists checks if a token exists and removes it if expired.
func (fs *FileStore) Exists(ctx context.Context, tokenID string) (*Token, error) {
	req, _, err := getContextRequestResponse(ctx)
	if err != nil {
		return nil, err
	}

	session, err := fs.store.Get(req, fs.CookieName)
	if err != nil {
		return nil, errors.New("failed to retrieve session")
	}

	tok, err := getSessionToken(session, tokenID)
	if err != nil {
		// Token not found or expired
		return nil, fs.Delete(ctx, tokenID)
	}
	return tok, nil
}

// Verify checks if the provided code matches the stored token's hash.
func (fs *FileStore) Verify(ctx context.Context, tokenID, code string) (bool, error) {
	tok, err := fs.Exists(ctx, tokenID)
	if err != nil {
		return false, err
	}

	// Hash the provided code and compare with stored hash
	codeHash := sha256.Sum256([]byte(code))
	if !bytes.Equal(tok.CodeHash, codeHash[:]) {
		tok.Attempts++
		_ = fs.Store(ctx, *tok) // Save attempts count
		return false, errors.New("invalid code")
	}

	// Delete token after successful verification (one-time use)
	return true, fs.Delete(ctx, tokenID)
}

// Delete removes the session token from the store.
func (fs *FileStore) Delete(ctx context.Context, tokenID string) error {
	return deleteSession(ctx, fs.store, fs.CookieName, tokenID)
}
