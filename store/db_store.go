// file: store/db_store.go
package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type DbStore struct {
	DB        *sql.DB
	TableName string
}

// NewDbStore constructs a DbStore with a reference to *sql.DB and a table name
func NewDbStore(db *sql.DB, tableName string) *DbStore {
	return &DbStore{
		DB:        db,
		TableName: tableName,
	}
}

func (s *DbStore) Store(ctx context.Context, tok Token) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, recipient, code_hash, expires_at, created_at)
                          VALUES (?, ?, ?, ?, ?)`, s.TableName)
	_, err := s.DB.ExecContext(ctx, query,
		tok.ID,
		tok.Recipient,
		tok.CodeHash,
		tok.ExpiresAt,
		tok.CreatedAt,
	)
	return err
}

func (s *DbStore) Exists(ctx context.Context, tokenID string) (*Token, error) {
	query := fmt.Sprintf(`SELECT id, recipient, code_hash, expires_at, created_at
                          FROM %s WHERE id = ?`, s.TableName)
	row := s.DB.QueryRowContext(ctx, query, tokenID)

	var tok Token
	err := row.Scan(&tok.ID, &tok.Recipient, &tok.CodeHash, &tok.ExpiresAt, &tok.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, err
	}
	return &tok, nil
}

func (s *DbStore) Verify(ctx context.Context, tokenID, code string) (bool, error) {
	tok, err := s.Exists(ctx, tokenID)
	if err != nil {
		// If not found or DB error, bail out
		return false, err
	}

	// Check expiry
	if time.Now().After(tok.ExpiresAt) {
		// If expired, delete it
		_ = s.Delete(ctx, tokenID)
		return false, fmt.Errorf("token expired")
	}

	// Compare hash
	codeHash := sha256.Sum256([]byte(code))
	if string(codeHash[:]) != string(tok.CodeHash) {
		return false, nil
	}

	// If match, consume it by deleting
	if err := s.Delete(ctx, tokenID); err != nil {
		return false, err
	}
	return true, nil
}

func (s *DbStore) Delete(ctx context.Context, tokenID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, s.TableName)
	_, err := s.DB.ExecContext(ctx, query, tokenID)
	return err
}
