# **go-passwordless**

`go-passwordless` is a lightweight, extensible Go library that provides a secure, passwordless authentication system. It allows applications to verify users using one-time codes sent via email, SMS, or other messaging channels, eliminating the need for passwords. This approach improves security, simplifies the user experience, and reduces the risk of credential-based attacks.

## **🛠 Key Features**

- **Token Stores:** Choose from in-memory, cookie-based, file, or database storage options.
- **Flexible Transports:** Send tokens via log output (for testing), SMTP, or custom transports.
- **One-Time Login Links:** Automatically generate login URLs to simplify the authentication process.
- **Customizable Expiry:** Set expiration times to control token validity.
- **Stateless Authentication:** No need to manage sessions or passwords.
- **Secure by Default:** Supports encrypted token storage and best practices.

## **🔍 What Problem Does It Solve?**

Managing passwords is challenging and comes with security risks such as:

- **Security vulnerabilities:** Password leaks, brute-force attacks, and phishing.
- **User friction:** Users often forget passwords, leading to frequent resets.
- **Storage concerns:** Securely storing and hashing passwords requires careful implementation.

`go-passwordless` eliminates these concerns by providing a **passwordless authentication flow**, enabling users to log in with one-time codes or links via email, SMS, or other means.

## **🛠 How It Works**

1. **User Initiates Login:**
   - The application calls `StartLogin()` with the recipient's email/phone number.
   - A secure, time-limited token is generated and stored.
   - The token is sent via a transport method (email, SMS, etc.).

2. **User Clicks Login Link or Enters Code:**
   - If using codes, the user manually inputs it into the application.
   - If using links, they click the provided one-time login URL.

3. **Successful Verification:**
   - If valid, authentication succeeds, and the token is deleted.
   - If expired or incorrect, authentication fails.

## **🚀 Quick Start**

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

 // Generate a one-time login link
 loginURL, err := mgr.GenerateLoginLink(ctx, "user@example.com", "https://myapp.com/login")
 if err != nil {
  log.Fatalf("Error generating login link: %v", err)
 }

 log.Printf("Your one-time login link: %s", loginURL)

 // Simulate token verification (user submits the token from the URL)
 success, err := mgr.VerifyLogin(ctx, "token-from-url", "")
 if err != nil {
  log.Fatalf("Error verifying login: %v", err)
 }

 if success {
  log.Println("Login successful!")
 } else {
  log.Println("Invalid token!")
 }
}
```

## **🔗 Generating One-Time Login Links**

The `GenerateLoginLink()` helper simplifies the process of sending users a one-time login link, allowing them to authenticate by clicking the link.

### **How to Use It:**

```go
ctx := context.Background()

mgr := passwordless.NewManager(store.NewMemStore(), &transport.LogTransport{})
loginURL, err := mgr.GenerateLoginLink(ctx, "user@example.com", "https://myapp.com/login")
if err != nil {
    log.Fatalf("Error generating login link: %v", err)
}

log.Println("Login link:", loginURL)
```

### **Example Output:**

```bash
Login link: https://myapp.com/login?token=abc123xyz
```

### **How to Handle the Link in Your Frontend:**

When the user clicks the link, your frontend should extract the `token` parameter and send it to your backend for verification.

Example frontend handler in JavaScript:

```javascript
const params = new URLSearchParams(window.location.search);
const token = params.get('token');

fetch('https://api.myapp.com/verify', {
    method: 'POST',
    body: JSON.stringify({ token }),
    headers: { 'Content-Type': 'application/json' },
})
  .then(response => response.json())
  .then(data => {
      if (data.success) {
          console.log('Login successful!');
      } else {
          console.error('Invalid or expired token');
      }
  });
```

## **📖 How to Implement in Your Project**

### **Step 1: Install the package**

```bash
go get github.com/rlnorthcutt/go-passwordless
```

### **Step 2: Choose Storage and Transport**

Example using `MemStore` and `SMTPTransport`:

```go
memStore := store.NewMemStore()
smtpTransport := &transport.SMTPTransport{
    Host: "smtp.example.com",
    Port: "587",
    From: "noreply@example.com",
    Auth: smtp.PlainAuth("", "user", "pass", "smtp.example.com"),
}

mgr := passwordless.NewManager(memStore, smtpTransport)
```

### **Step 3: Start Login and Verify**

```go
tokenID, _ := mgr.StartLogin(context.Background(), "user@example.com")
success, _ := mgr.VerifyLogin(context.Background(), tokenID, "123456")
```

## **🔗 Dependencies**

`go-passwordless` has minimal dependencies to ensure lightweight performance:

- `github.com/gorilla/securecookie` (for secure cookie handling).
- `modernc.org/sqlite` (for database storage in `DbStore`).
- `net/smtp` (for email transport).

## **🧪 Running Tests**

Run tests to ensure the implementation works correctly:

```bash
go test -v ./...
```

To test specific modules:

```bash
go test -v ./store
go test -v ./transport
```

## **📦 Contributing**

We welcome contributions! If you'd like to contribute:

1. Fork the repository.
2. Create a feature branch.
3. Submit a pull request with a detailed description.

If you have any questions or suggestions, feel free to open an issue on [GitHub](https://github.com/rlnorthcutt/go-passwordless/issues).

## **📜 License**

`go-passwordless` is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
