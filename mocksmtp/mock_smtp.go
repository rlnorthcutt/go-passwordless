package mocksmtp

import (
	"log"
	"net"
	"strings"
)

// StartMockSMTPServer runs a simple SMTP server on the given port.
func StartMockSMTPServer(port string) {
	go func() {
		listener, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("Failed to start mock SMTP server: %v", err)
		}
		defer listener.Close()

		log.Printf("Mock SMTP server is running on :%s", port)

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Connection error: %v", err)
				continue
			}
			go handleSMTPConnection(conn)
		}
	}()
}

func handleSMTPConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("220 mock-smtp-server ESMTP\r\n"))

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}

		cmd := strings.TrimSpace(string(buffer[:n]))
		log.Println("Received:", cmd)

		switch {
		case strings.HasPrefix(cmd, "HELO"):
			conn.Write([]byte("250 Hello\r\n"))
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			conn.Write([]byte("250 OK\r\n"))
		case strings.HasPrefix(cmd, "RCPT TO:"):
			conn.Write([]byte("250 OK\r\n"))
		case strings.HasPrefix(cmd, "DATA"):
			conn.Write([]byte("354 End data with <CR><LF>.<CR><LF>\r\n"))
		case strings.HasSuffix(cmd, "."):
			conn.Write([]byte("250 Message accepted\r\n"))
		case strings.HasPrefix(cmd, "QUIT"):
			conn.Write([]byte("221 Bye\r\n"))
			return
		default:
			conn.Write([]byte("500 Unrecognized command\r\n"))
		}
	}
}
