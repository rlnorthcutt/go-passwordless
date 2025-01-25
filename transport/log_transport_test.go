package transport_test

import (
	"context"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/transport"
)

func TestLogTransport_Send(t *testing.T) {
	t.Logf("[DEBUG] Starting TestLogTransport_Send...")

	// Step 1: Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	t.Logf("[DEBUG] Context with timeout set to 3 seconds")

	lt := &transport.LogTransport{}

	// Step 2: Test sending a log transport message
	recipient := "test@example.com"
	code := "123456"

	t.Logf("[DEBUG] Sending code %s to recipient %s", code, recipient)

	err := lt.Send(ctx, recipient, code)
	if err != nil {
		t.Fatalf("LogTransport.Send returned error: %v", err)
	}
	t.Logf("[DEBUG] LogTransport.Send executed successfully for recipient: %s", recipient)

	// Step 3: Test handling of a canceled context
	cancel() // Cancel the context

	t.Logf("[DEBUG] Context canceled, attempting to send with canceled context")

	err = lt.Send(ctx, "test-cancel@example.com", "abcdef")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, but got: %v", err)
	} else {
		t.Logf("[DEBUG] LogTransport correctly returned context.Canceled error")
	}

	t.Logf("[DEBUG] TestLogTransport_Send completed successfully")
}
