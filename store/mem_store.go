package store

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MemStore struct {
	mu     sync.Mutex
	tokens map[string]Token
}

func NewMemStore() *MemStore {
	return &MemStore{
		tokens: make(map[string]Token),
	}
}

func (m *MemStore) Store(ctx context.Context, tok Token) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[tok.ID] = tok
	return nil
}

func (m *MemStore) Exists(ctx context.Context, tokenID string) (*Token, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tok, ok := m.tokens[tokenID]
	if !ok {
		return nil, fmt.Errorf("token not found")
	}

	if time.Now().After(tok.ExpiresAt) {
		delete(m.tokens, tokenID)
		return nil, fmt.Errorf("token expired")
	}

	return &tok, nil
}

func (m *MemStore) Verify(ctx context.Context, tokenID, code string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	tok, ok := m.tokens[tokenID]
	if !ok {
		return false, fmt.Errorf("token not found")
	}

	if IsTokenExpired(&tok) {
		delete(m.tokens, tokenID)
		return false, fmt.Errorf("token expired")
	}

	if !VerifyToken(&tok, code) {
		return false, nil
	}

	// If match, consume it (delete immediately):
	delete(m.tokens, tokenID)
	return true, nil
}

func (m *MemStore) Delete(ctx context.Context, tokenID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.tokens, tokenID)
	return nil
}
