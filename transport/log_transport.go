// file: transport/log_transport.go

package transport

import (
	"context"
	"log"
)

// LogTransport simply logs the token code (for testing/dev).
type LogTransport struct{}

// Send checks for cancellation via ctx, then logs the token code.
func (l *LogTransport) Send(ctx context.Context, recipient, tokenCode string) error {
	// Optionally check for cancellation before doing anything.
	select {
	case <-ctx.Done():
		// If the context is already canceled, return immediately.
		return ctx.Err()
	default:
	}

	log.Printf("[LOG TRANSPORT] Recipient: %s, Code: %s\n", recipient, tokenCode)
	return nil
}
