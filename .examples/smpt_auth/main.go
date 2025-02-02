package main

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/mocksmtp"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

func main() {
	// Start the mock SMTP server if using local testing
	// `USE_MOCK_SMTP=true go run examples/smtp_auth/main.go`
	if os.Getenv("USE_MOCK_SMTP") == "true" {
		go mocksmtp.StartMockSMTPServer("2525")
		time.Sleep(1 * time.Second) // Allow server to start
	}

	// Use the mock server by default if no real credentials are set
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "localhost"
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "2525" // Default to mock SMTP server port
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	smtpTransport := &transport.SMTPTransport{
		Host: smtpHost,
		Port: smtpPort,
		From: smtpUser,
		Auth: smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost),
	}

	mgr := passwordless.NewManager(store.NewMemStore(), smtpTransport)

	email := "user@example.com"
	loginURL, err := mgr.GenerateLoginLink(context.Background(), email, "https://myapp.com/login")
	if err != nil {
		log.Fatalf("Failed to generate login link: %v", err)
	}

	err = smtpTransport.Send(context.Background(), email, fmt.Sprintf("Click the link to login: %s", loginURL))
	if err != nil {
		log.Fatalf("Failed to send login email: %v", err)
	}

	log.Println("A login link has been sent to:", email)
}
