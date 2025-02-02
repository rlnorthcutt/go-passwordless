package main

import (
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store/session"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

const (
	fileStorePath = "./session_data"
	secretKey     = "super-secret-key"
)

func main() {
	// Ensure cleanup of session data after execution
	defer cleanupFiles()

	// Setup mock HTTP request and response for session handling
	req := httptest.NewRequest("GET", "/", nil)
	rsp := httptest.NewRecorder()
	ctx := session.WithRequestResponse(context.Background(), req, rsp)

	// Initialize the FileStore
	fileStore := session.NewFileStore(fileStorePath, []byte(secretKey))
	if fileStore == nil {
		log.Fatal("Failed to initialize FileStore")
	}

	// Initialize the passwordless manager with FileStore and LogTransport
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
	fmt.Println("A login code has been logged. Check the console output in this example.")

	// Simulate reading the stored token via a new request with session cookies
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("Cookie", rsp.Header().Get("Set-Cookie"))
	rsp2 := httptest.NewRecorder()
	ctx2 := session.WithRequestResponse(context.Background(), req2, rsp2)

	// Prompt the user to enter the verification code
	fmt.Print("Enter the verification code: ")
	var inputCode string
	fmt.Scanln(&inputCode)

	// Verify the entered code
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

// cleanupFiles removes session data after the example finishes.
func cleanupFiles() {
	fmt.Println("Cleaning up session data...")
	if err := os.RemoveAll(fileStorePath); err != nil {
		log.Printf("Failed to remove session directory: %v", err)
	}
	fmt.Println("Cleanup complete.")
}
