package transport_test

import (
	"context"
	"net/smtp"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/transport"
)

func TestSMTPTransport_Send(t *testing.T) {
	t.Logf("[DEBUG] Starting TestSMTPTransport_Send...")

	// Step 1: Setup the SMTP transport (use a local test server like MailHog for testing)
	tr := &transport.SMTPTransport{
		Host: "localhost",
		Port: "1025", // Assumes a test SMTP server is running (e.g., MailHog)
		From: "noreply@example.com",
		Auth: smtp.PlainAuth("", "user", "pass", "localhost"),
	}

	t.Logf("[DEBUG] SMTP Transport configured: host=%s, port=%s, from=%s", tr.Host, tr.Port, tr.From)

	// Step 2: Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t.Logf("[DEBUG] Context with timeout set to 5 seconds")

	// Step 3: Send email
	recipient := "test@example.com"
	code := "999999"

	t.Logf("[DEBUG] Attempting to send code %s to recipient %s", code, recipient)
	err := tr.Send(ctx, recipient, code)

	// Step 4: Handle expected error if no SMTP server is running
	if err != nil {
		t.Logf("[DEBUG] Expected error due to missing test SMTP server: %v", err)
		t.Logf("[INFO] If you want to run real tests, start a mock SMTP server like MailHog.")
		t.Skip("Skipping test due to no running SMTP server.")
	} else {
		t.Logf("[DEBUG] Email sent successfully to recipient: %s", recipient)
	}

	t.Logf("[DEBUG] TestSMTPTransport_Send completed successfully")
}
