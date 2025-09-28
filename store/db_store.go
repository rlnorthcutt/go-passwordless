package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// DbStore represents a database-backed token store.
type DbStore struct {
	DB        *sql.DB // Reference to the database connection
	TableName string  // Name of the table used to store tokens
}

// NewDbStore initializes a new DbStore with a reference to *sql.DB and a table name.
func NewDbStore(db *sql.DB, tableName string) *DbStore {
	return &DbStore{
		DB:        db,
		TableName: tableName,
	}
}

// Store saves a new token in the database.
func (s *DbStore) Store(ctx context.Context, tok Token) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (id, recipient, code_hash, expires_at, created_at, attempts)
		VALUES (?, ?, ?, ?, ?, ?)`, s.TableName)

	_, err := s.DB.ExecContext(ctx, query,
		tok.ID,
		tok.Recipient,
		tok.CodeHash,
		tok.ExpiresAt,
		tok.CreatedAt,
		tok.Attempts,
	)
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}
	return nil
}

// Exists checks if a token exists in the database and returns it.
// If the token is expired, it is deleted automatically.
func (s *DbStore) Exists(ctx context.Context, tokenID string) (*Token, error) {
	query := fmt.Sprintf(`
                SELECT id, recipient, code_hash, expires_at, created_at, attempts
                FROM %s WHERE id = ?`, s.TableName)

	var tok Token
	err := s.DB.QueryRowContext(ctx, query, tokenID).Scan(
		&tok.ID,
		&tok.Recipient,
		&tok.CodeHash,
		&tok.ExpiresAt,
		&tok.CreatedAt,
		&tok.Attempts,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if the token has expired and delete it
	if IsTokenExpired(&tok) {
		_ = s.Delete(ctx, tokenID) // Purge expired token
		return nil, fmt.Errorf("token expired and was deleted")
	}

	return &tok, nil
}

// UpdateAttempts updates the failed-attempt counter for a token without altering other fields.
func (s *DbStore) UpdateAttempts(ctx context.Context, tokenID string, attempts int) error {
	query := fmt.Sprintf(`UPDATE %s SET attempts = ? WHERE id = ?`, s.TableName)

	res, err := s.DB.ExecContext(ctx, query, attempts, tokenID)
	if err != nil {
		return fmt.Errorf("failed to update token attempts: %w", err)
	}

	if rows, err := res.RowsAffected(); err == nil && rows == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

// Verify checks whether the provided code matches the stored hash.
func (s *DbStore) Verify(ctx context.Context, tokenID, code string) (bool, error) {
	tok, err := s.Exists(ctx, tokenID)
	if err != nil {
		return false, err // Token not found or expired
	}

	if !VerifyToken(tok, code) {
		return false, fmt.Errorf("invalid code")
	}

	// If verification succeeds, delete token (one-time use)
	if err := s.Delete(ctx, tokenID); err != nil {
		return false, fmt.Errorf("failed to delete token after verification: %w", err)
	}
	return true, nil
}

// Delete removes a token from the database.
func (s *DbStore) Delete(ctx context.Context, tokenID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, s.TableName)
	_, err := s.DB.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}
