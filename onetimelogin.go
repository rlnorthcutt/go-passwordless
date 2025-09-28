package passwordless

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/url"
)

// GenerateLoginLink generates a one-time login link containing a token and hashed code.
func (m *Manager) GenerateLoginLink(ctx context.Context, recipient, baseURL string) (string, error) {
	tokenID, err := m.StartLogin(ctx, recipient)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Retrieve the stored token to get the code hash
	tok, err := m.Store.Exists(ctx, tokenID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token: %w", err)
	}

	// Generate a secure hash for URL encoding (do not expose raw hash)
	codeHash := sha256.Sum256(tok.CodeHash)
	encodedCodeHash := hex.EncodeToString(codeHash[:])

	// Construct the login URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("token", tokenID)
	query.Set("hash", encodedCodeHash)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

// VerifyLoginLink validates a one-time login link without requiring user input.
func (m *Manager) VerifyLoginLink(ctx context.Context, tokenID, providedHash string) (bool, error) {
	tok, err := m.Store.Exists(ctx, tokenID)
	if err != nil {
		return false, fmt.Errorf("token not found or expired")
	}

	// Hash the stored code and compare it with the hash from the URL
	expectedHash := sha256.Sum256(tok.CodeHash)
	expectedHashHex := make([]byte, hex.EncodedLen(len(expectedHash)))
	hex.Encode(expectedHashHex, expectedHash[:])

	providedHashBytes := []byte(providedHash)

	if len(providedHashBytes) != len(expectedHashHex) ||
		subtle.ConstantTimeCompare(expectedHashHex, providedHashBytes) != 1 {
		tok.Attempts++
		if tok.Attempts >= m.Config.MaxFailedAttempts {
			_ = m.Store.Delete(ctx, tokenID)
			return false, fmt.Errorf("too many failed attempts, token deleted")
		}
		if err := m.Store.UpdateAttempts(ctx, tokenID, tok.Attempts); err != nil {
			return false, fmt.Errorf("failed to persist attempt count: %w", err)
		}
		return false, fmt.Errorf("invalid login link")
	}

	// If verification succeeds, delete token (one-time use)
	_ = m.Store.Delete(ctx, tokenID)
	return true, nil
}
