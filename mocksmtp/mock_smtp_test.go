package mocksmtp_test

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/rlnorthcutt/go-passwordless/mocksmtp"
)

func TestMockSMTPServer(t *testing.T) {
	port := "2526"
	mocksmtp.StartMockSMTPServer(port)

	// Allow the server to start
	time.Sleep(500 * time.Millisecond)

	t.Run("HELO Command", func(t *testing.T) {
		conn, err := net.Dial("tcp", ":"+port)
		if err != nil {
			t.Fatalf("Failed to connect to mock SMTP server: %v", err)
		}
		defer conn.Close()

		readResponse(t, conn) // Read initial 220 message
		sendCommand(t, conn, "HELO localhost\r\n", "250 Hello")
	})

	t.Run("MAIL FROM and RCPT TO Commands", func(t *testing.T) {
		conn, err := net.Dial("tcp", ":"+port)
		if err != nil {
			t.Fatalf("Failed to connect to mock SMTP server: %v", err)
		}
		defer conn.Close()

		readResponse(t, conn)
		sendCommand(t, conn, "MAIL FROM:<test@example.com>\r\n", "250 OK")
		sendCommand(t, conn, "RCPT TO:<recipient@example.com>\r\n", "250 OK")
	})

	t.Run("DATA Command", func(t *testing.T) {
		conn, err := net.Dial("tcp", ":"+port)
		if err != nil {
			t.Fatalf("Failed to connect to mock SMTP server: %v", err)
		}
		defer conn.Close()

		readResponse(t, conn)
		sendCommand(t, conn, "DATA\r\n", "354 End data with <CR><LF>.<CR><LF>")
		sendCommand(t, conn, "Test email body\r\n.\r\n", "250 Message accepted")
	})

	t.Run("QUIT Command", func(t *testing.T) {
		conn, err := net.Dial("tcp", ":"+port)
		if err != nil {
			t.Fatalf("Failed to connect to mock SMTP server: %v", err)
		}
		defer conn.Close()

		readResponse(t, conn)
		sendCommand(t, conn, "QUIT\r\n", "221 Bye")
	})
}

func sendCommand(t *testing.T, conn net.Conn, cmd string, expectedResponse string) {
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		t.Fatalf("Failed to send command %q: %v", cmd, err)
	}

	response := readResponse(t, conn)
	if !strings.Contains(response, expectedResponse) {
		t.Errorf("Expected response to contain %q, got %q", expectedResponse, response)
	}
}

func readResponse(t *testing.T, conn net.Conn) string {
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	fmt.Println("Server Response:", strings.TrimSpace(response))
	return strings.TrimSpace(response)
}
