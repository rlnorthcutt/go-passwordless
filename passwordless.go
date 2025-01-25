package passwordless

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/url"
	"time"

	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

type Manager struct {
	Store     store.TokenStore
	Transport transport.Transport
	Config    Config
}

// NewManager constructs a Manager using the default config (see config.go).
func NewManager(s store.TokenStore, t transport.Transport) *Manager {
	return &Manager{
		Store:     s,
		Transport: t,
		Config:    DefaultConfig(),
	}
}

// NewManagerWithConfig constructs a Manager using a custom Config.
func NewManagerWithConfig(s store.TokenStore, t transport.Transport, cfg Config) *Manager {
	// Fill in any zero-valued fields with defaults
	if cfg.CodeLength == 0 {
		cfg.CodeLength = 6
	}
	if cfg.TokenExpiry == 0 {
		cfg.TokenExpiry = 15 * time.Minute
	}
	if cfg.IDGenerator == nil {
		cfg.IDGenerator = defaultIDGenerator
	}
	if cfg.CodeCharset == "" {
		cfg.CodeCharset = "0123456789"
	}

	return &Manager{
		Store:     s,
		Transport: t,
		Config:    cfg,
	}
}

// StartLogin generates a code, stores it, and sends it to the recipient.
// Returns the generated token ID.
func (m *Manager) StartLogin(ctx context.Context, recipient string) (string, error) {
	// Generate code
	code, err := m.generateCode(m.Config.CodeLength, m.Config.CodeCharset)
	if err != nil {
		return "", err
	}

	// Hash it
	hash := sha256.Sum256([]byte(code))

	// Generate a token ID
	tokenID := m.Config.IDGenerator()

	// Build Token
	tok := store.Token{
		ID:        tokenID,
		Recipient: recipient,
		CodeHash:  hash[:],
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(m.Config.TokenExpiry),
	}

	// Store the token
	if err := m.Store.Store(ctx, tok); err != nil {
		return "", err
	}

	// Send the code to the user
	if err := m.Transport.Send(ctx, recipient, code); err != nil {
		// If sending fails, remove the token
		_ = m.Store.Delete(ctx, tokenID)
		return "", err
	}

	return tokenID, nil
}

// VerifyLogin checks the user-provided code against the stored token.
func (m *Manager) VerifyLogin(ctx context.Context, tokenID, code string) (bool, error) {
	ok, err := m.Store.Verify(ctx, tokenID, code)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// GenerateLoginLink generates a one-time login link containing a token.
func (m *Manager) GenerateLoginLink(ctx context.Context, recipient, baseURL string) (string, error) {
	// Start the login process and get a token ID
	tokenID, err := m.StartLogin(ctx, recipient)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Construct the login URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	// Add the token ID as a query parameter
	query := parsedURL.Query()
	query.Set("token", tokenID)
	parsedURL.RawQuery = query.Encode()

	// Return the final URL string
	return parsedURL.String(), nil
}

// generateCode produces a random code (numeric or alphanumeric) based on the config.
func (m *Manager) generateCode(length int, charset string) (string, error) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = charset[int(buf[i])%len(charset)]
	}
	return string(out), nil
}
