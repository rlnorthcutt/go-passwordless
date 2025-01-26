package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

const (
	fileStorePath = "./session_data"
	tokenFilePath = "./tokens.json"
	secretKey     = "super-secret-key"
)

func main() {
	// Ensure cleanup of session data and tokens on exit
	defer cleanupFiles()

	// Setup mock HTTP request and response
	req := httptest.NewRequest("GET", "/", nil)
	rsp := httptest.NewRecorder()
	ctx := store.WithRequestResponse(context.Background(), req, rsp)

	// Initialize the FileStore with a secret key
	fileStore := store.NewFileStore(fileStorePath, []byte(secretKey))
	if fileStore == nil {
		log.Fatal("Failed to initialize FileStore")
	}

	// Initialize the passwordless manager
	mgr := passwordless.NewManager(fileStore, &transport.LogTransport{})

	// Prompt user for email
	fmt.Print("Enter your email: ")
	var email string
	fmt.Scanln(&email)

	// Start the login process and send a one-time code
	tokenID, err := mgr.StartLogin(ctx, email)
	if err != nil {
		log.Fatalf("Failed to start login: %v", err)
	}
	fmt.Println("A login code has been stored in session. Check your email (console log in this example).")

	// Simulate reading back the stored token via a new request with cookies
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("Cookie", rsp.Header().Get("Set-Cookie"))
	rsp2 := httptest.NewRecorder()
	ctx2 := store.WithRequestResponse(context.Background(), req2, rsp2)

	// Retrieve and verify token
	fmt.Print("Enter the verification code: ")
	var inputCode string
	fmt.Scanln(&inputCode)

	success, err := mgr.VerifyLogin(ctx2, tokenID, inputCode)
	if err != nil {
		log.Fatalf("Verification failed: %v", err)
	}

	if success {
		fmt.Println("Login successful!")
	} else {
		fmt.Println("Invalid code.")
	}
}

// cleanupFiles removes session and token data after the example finishes.
func cleanupFiles() {
	fmt.Println("Cleaning up session and token files...")
	if err := os.RemoveAll(fileStorePath); err != nil {
		log.Printf("Failed to remove session directory: %v", err)
	}
	if err := os.Remove(tokenFilePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove token file: %v", err)
	}
	fmt.Println("Cleanup complete.")
}
