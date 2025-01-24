package transport

import (
	"context"
	"fmt"
	"net/smtp"
)

// SMTPTransport sends token codes via an SMTP server.
type SMTPTransport struct {
	Host string // e.g. "smtp.example.com"
	Port string // e.g. "587"
	From string // e.g. "noreply@example.com"
	Auth smtp.Auth
}

func (t *SMTPTransport) Send(ctx context.Context, to, tokenCode string) error {
	// Check ctx for cancellation before doing work
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Construct message
	msg := []byte(fmt.Sprintf("Subject: Your Login Code\r\n\r\nYour code is: %s\r\n", tokenCode))
	addr := t.Host + ":" + t.Port

	// net/smtp.SendMail doesn't directly accept context, so you can't forcibly cancel it mid-flight.
	// But you can check ctx again before calling:
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return smtp.SendMail(addr, t.Auth, t.From, []string{to}, msg)
}
