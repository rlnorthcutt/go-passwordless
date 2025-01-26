package main

import (
	"context"
	"log"
	"os"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

func main() {
	ctx := context.Background()

	// Load SMTP credentials from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	smtpTransport := &transport.SMTPTransport{
		Host: smtpHost,
		Port: smtpPort,
		From: smtpUser,
		Auth: transport.NewSMTPAuth(smtpUser, smtpPassword, smtpHost),
	}

	mgr := passwordless.NewManager(store.NewMemStore(), smtpTransport)

	email := "user@example.com"
	tokenID, err := mgr.StartLogin(ctx, email)
	if err != nil {
		log.Fatalf("Failed to start login: %v", err)
	}

	log.Println("A login code has been sent to:", email)

	// Simulate verification (in a real app, get the code from email)
	success, err := mgr.VerifyLogin(ctx, tokenID, "123456")
	if err != nil {
		log.Fatalf("Verification failed: %v", err)
	}
	if success {
		log.Println("Login successful!")
	} else {
		log.Println("Invalid code.")
	}
}
