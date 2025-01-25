# **go-passwordless**

`go-passwordless` is a lightweight, extensible Go library that provides a secure, passwordless authentication system. It allows applications to verify users using one-time codes sent via email, SMS, or other messaging channels, eliminating the need for passwords. This approach improves security, simplifies the user experience, and reduces the risk of credential-based attacks.

## **🔍 What Problem Does It Solve?**

Managing passwords is challenging and comes with security risks such as:

- **Security vulnerabilities:** Password leaks, brute-force attacks, and phishing.
- **User friction:** Users often forget passwords, leading to frequent resets.
- **Storage concerns:** Securely storing and hashing passwords requires careful implementation.

`go-passwordless` eliminates these concerns by providing a **passwordless authentication flow**. Instead of passwords, users receive a secure one-time code via email, SMS, or other methods to verify their identity.

## **🚀 TL;DR (Quick Start)**

### **Installation**

```bash
go get github.com/rlnorthcutt/go-passwordless
```

### **Basic Usage Example**

```go
package main

import (
 "context"
 "log"

 "github.com/rlnorthcutt/go-passwordless"
 "github.com/rlnorthcutt/go-passwordless/store"
 "github.com/rlnorthcutt/go-passwordless/transport"
)

func main() {
 ctx := context.Background()

 // Initialize token store (choose MemStore for ephemeral storage)
 memStore := store.NewMemStore()

 // Initialize transport (LogTransport for development)
 logTransport := &transport.LogTransport{}

 // Create the passwordless manager
 mgr := passwordless.NewManager(memStore, logTransport)

 // Start login process
 recipient := "user@example.com"
 tokenID, err := mgr.StartLogin(ctx, recipient)
 if err != nil {
  log.Fatalf("Error starting login: %v", err)
 }
 log.Printf("Token ID: %s", tokenID)

 // Simulate token verification
 success, err := mgr.VerifyLogin(ctx, tokenID, "123456") // Replace with actual code
 if err != nil {
  log.Fatalf("Error verifying login: %v", err)
 }

 if success {
  log.Println("Login successful!")
 } else {
  log.Println("Invalid code!")
 }
}
```

## **🛠 How It Works**

1. **User Initiates Login:**
   - The application calls `StartLogin()` with the recipient's email/phone number.
   - A secure, time-limited token is generated and stored.
   - The token is sent to the recipient via the configured transport.

2. **User Provides Token:**
   - The user enters the received code in the application.
   - The application calls `VerifyLogin()` to validate the code.

3. **Successful Verification:**
   - If valid, the token is marked as used, and authentication is successful.
   - If expired or incorrect, authentication fails.

## **🛠 Deep Dive into Components**

### **1. Stores (Token Management)**

`go-passwordless` offers various storage options to suit different needs:

| Store Type   | Description                           | Use Case                      |
|--------------|---------------------------------------|-------------------------------|
| `MemStore`   | In-memory storage for tokens.         | Development, short-lived apps.|
| `CookieStore`| Stores tokens in secure cookies.      | Stateless web applications.   |
| `DbStore`    | SQL-based token persistence.          | Scalable, persistent storage. |
| `FileStore`  | File-based storage for small apps.    | Local, non-distributed usage. |

### **2. Transports (Token Delivery)**

Transports are responsible for delivering authentication tokens to users:

| Transport Type  | Description                         | Use Case                       |
|-----------------|------------------------------------|-------------------------------|
| `LogTransport`  | Logs token to stdout (debug mode). | Development and testing.      |
| `SMTPTransport` | Sends tokens via email.           | Production authentication.    |
| `Custom`        | Implement your own delivery method | e.g., SMS, push notifications |

### **3. Manager (Core Logic)**

The `Manager` component handles the entire login and verification lifecycle:

```go
type Manager struct {
 store     store.TokenStore
 transport transport.Transport
 config    Config
}
```

- **`StartLogin(ctx, recipient)`** – Generates and sends a token.
- **`VerifyLogin(ctx, tokenID, code)`** – Verifies the provided token.
- **`Config`** – Customize code length, expiration, and other settings.

## **📖 How to Use in Your Project**

### **Step 1: Install the package**

```bash
go get github.com/rlnorthcutt/go-passwordless
```

### **Step 2: Choose Storage and Transport**

Example using `MemStore` and `SMTPTransport`:

```go
import (
 "context"
 "github.com/rlnorthcutt/go-passwordless"
 "github.com/rlnorthcutt/go-passwordless/store"
 "github.com/rlnorthcutt/go-passwordless/transport"
)

// Setup store and transport
memStore := store.NewMemStore()
smtpTransport := &transport.SMTPTransport{
    Host: "smtp.example.com",
    Port: "587",
    From: "noreply@example.com",
    Auth: smtp.PlainAuth("", "user", "pass", "smtp.example.com"),
}

// Initialize manager
mgr := passwordless.NewManager(memStore, smtpTransport)

// Start login
tokenID, err := mgr.StartLogin(context.Background(), "user@example.com")
```

### **Step 3: Verify Login**

```go
success, err := mgr.VerifyLogin(context.Background(), tokenID, "123456")
if err != nil {
    log.Fatal("Verification failed:", err)
}
if success {
    log.Println("Login successful")
}
```

## **🔗 Dependencies**

`go-passwordless` has minimal dependencies to ensure lightweight, fast performance. Key dependencies include:

- `github.com/gorilla/securecookie` (for secure cookie handling).
- `modernc.org/sqlite` (for database storage in `DbStore`).
- `net/smtp` (for email transport).

Install dependencies using:

```bash
go mod tidy
```

## **🧪 Running Tests**

Unit tests ensure the correctness of token storage, transport mechanisms, and authentication flows.

```bash
go test -v ./...
```

To test with verbose logging for debugging:

```bash
go test -v ./store
go test -v ./transport
```

## **📦 Contributing**

We welcome contributions to improve `go-passwordless`. If you’d like to contribute:

1. Fork the repository.
2. Create a feature branch.
3. Submit a pull request with a detailed description.

Please follow the [contribution guidelines](CONTRIBUTING.md).

## **❓ FAQ**

**Q: Is passwordless authentication secure?**
A: Yes, when combined with secure transports and short-lived token expiration, it provides a strong authentication method.

**Q: Can I customize token expiration times?**
A: Yes, using the `Config` struct to set custom expiry durations.

**Q: How do I integrate with SMS providers?**
A: You can create a custom transport by implementing the `Transport` interface.

## **📜 License**

`go-passwordless` is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details

---

Let us know your thoughts or questions via [GitHub Issues](https://github.com/rlnorthcutt/go-passwordless/issues). Happy coding! 🚀
