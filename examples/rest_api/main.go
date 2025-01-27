package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rlnorthcutt/go-passwordless"
	"github.com/rlnorthcutt/go-passwordless/store"
	"github.com/rlnorthcutt/go-passwordless/transport"
)

var mgr *passwordless.Manager

func init() {
	// Initialize passwordless manager with an in-memory token store and log transport.
	mgr = passwordless.NewManager(store.NewMemStore(), &transport.LogTransport{})
}

func main() {
	// Show login URL when server starts
	fmt.Println("Server running at http://localhost:8080")
	fmt.Println("To start login, visit: http://localhost:8080/login?email=me@here.com")

	http.HandleFunc("/login", startLoginHandler)
	http.HandleFunc("/verify", verifyLoginHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// startLoginHandler handles login requests and generates a token.
func startLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required. Use /login?email=me@here.com", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	tokenID, err := mgr.StartLogin(ctx, email)
	if err != nil {
		http.Error(w, "Failed to start login", http.StatusInternalServerError)
		return
	}

	// Display login token and verification link
	response := fmt.Sprintf(
		"Login code sent to: %s\nToken ID: %s\nVerify at: http://localhost:8080/verify?token_id=%s&code=YOUR_CODE",
		email, tokenID, tokenID)
	fmt.Fprintln(w, response)
}

// verifyLoginHandler handles the verification of the login token.
func verifyLoginHandler(w http.ResponseWriter, r *http.Request) {
	tokenID := r.URL.Query().Get("token_id")
	code := r.URL.Query().Get("code")

	if tokenID == "" || code == "" {
		http.Error(w, "token_id and code are required. Use /verify?token_id=123&code=456", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	success, err := mgr.VerifyLogin(ctx, tokenID, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verification failed: %v", err), http.StatusUnauthorized)
		return
	}

	if success {
		log.Printf("Verification successful for token ID: %s", tokenID)
		fmt.Fprintln(w, "Verification successful! You are now logged in.")
	} else {
		fmt.Fprintln(w, "Verification failed! Invalid code.")
	}
}
