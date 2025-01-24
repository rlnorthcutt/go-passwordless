package transport

import "context"

// Transport is responsible for delivering a token code to a user's "address."
// Example "address" might be an email address, phone number, etc.
type Transport interface {
	// Send delivers `tokenCode` to the user's `recipient`.
	// The context can handle cancellation or timeouts.
	Send(ctx context.Context, recipient, tokenCode string) error
}
