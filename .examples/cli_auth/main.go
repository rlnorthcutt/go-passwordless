package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

func main() {
	ctx := context.Background()
	memStore := store.NewMemStore()
	logTransport := &transport.LogTransport{}
	mgr := passwordless.NewManager(memStore, logTransport)

	fmt.Print("Enter your email: ")
	var email string
	fmt.Scanln(&email)

	tokenID, err := mgr.StartLogin(ctx, email)
	if err != nil {
		log.Fatalf("Failed to start login: %v", err)
	}
	fmt.Println("Check your email for the verification code.")

	fmt.Print("Enter the verification code: ")
	var code string
	fmt.Scanln(&code)

	success, err := mgr.VerifyLogin(ctx, tokenID, code)
	if err != nil {
		log.Fatalf("Verification failed: %v", err)
	}
	if success {
		fmt.Println("Login successful!")
	} else {
		fmt.Println("Invalid code.")
	}
}
