// file: transport/log_transport.go

package transport

import (
	"log"
)

// LogTransport simply logs the token code (for testing/dev).
type LogTransport struct{}

func (l *LogTransport) Send(recipient, tokenCode string) error {
	log.Printf("[LOG TRANSPORT] Recipient: %s, Code: %s\n", recipient, tokenCode)
	return nil
}
