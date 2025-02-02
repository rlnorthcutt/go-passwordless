# **File-Based Authentication Example**

This example demonstrates how to implement passwordless authentication using the `go-passwordless` library with file-based storage. The authentication tokens are stored in local files, making it a simple and effective solution for CLI tools, local applications, or prototyping.

## **How It Works**

1. The user is prompted to enter their email.
2. A one-time login code is generated and stored in a local file.
3. The user enters the code they received (simulated by logging to the console).
4. If the code is correct, login is successful; otherwise, it fails.
5. Temporary session and token files are cleaned up after execution.

## **Running the Example**

To test the example, run the following command:

```bash
go run main.go
```

### **Example Run:**

```bash
Enter your email: me@here.com
A login code has been stored in session. Check your email (console log in this example).
Enter the verification code: 123456
Login successful!
Cleaning up session and token files...
Cleanup complete.
```

## **How to Extend This Example**

You can extend this example by:

- **Changing the storage method:**
  - Switch to `DbStore` for persistent storage instead of file-based sessions.

- **Replacing the transport:**
  - Currently, the example logs the token to the console using `LogTransport`.
  - Replace it with `SMTPTransport` to send real login codes via email.

  ```go
  smtpTransport := &transport.SMTPTransport{
      Host: "smtp.example.com",
      Port: "587",
      From: "noreply@example.com",
      Auth: transport.NewSMTPAuth("smtp-user", "smtp-pass", "smtp.example.com"),
  }
  mgr := passwordless.NewManager(fileStore, smtpTransport)
  ```

- **Running as a web service:**
  - Convert this example to an HTTP-based service to allow remote logins.

## **When to Use This**

This approach is ideal when:

- **Prototyping:**
  - Quickly testing passwordless authentication without a database or email setup.

- **Local CLI Applications:**
  - Storing temporary authentication data without relying on external services.

- **Small-Scale Applications:**
  - Projects that do not require centralized authentication but still need token-based login.

## **Cleanup Process**

To prevent leftover files, the example automatically removes:

- `./session_data` directory (stores session cookies)
- `tokens.json` file (stores token data)

If you'd like to persist tokens across sessions, remove the cleanup function from `main.go`:

```go
defer cleanupFiles()  // Comment out or remove this line to retain session files.
```

## **Conclusion**

This example provides a simple, easy-to-use implementation of passwordless authentication with local file storage. It's a great starting point for learning how passwordless login works before implementing it in larger systems.
