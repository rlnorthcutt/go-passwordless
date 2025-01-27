# **SMTP Authentication Example**

This example demonstrates how to implement passwordless authentication using the `go-passwordless` library with an SMTP transport. It allows users to receive one-time login links via email for a secure, passwordless authentication experience.

## **How It Works**

1. **Start the server.**
   - The SMTP transport is initialized with credentials from environment variables or defaults to a local mock SMTP server.

2. **User initiates login.**
   - The system generates a secure one-time login link and sends it via email.

3. **User clicks the login link.**
   - The application verifies the token and grants access.

## **Running the Example**

### **With the Mock SMTP Server** (Default)

The example is set to run with a built-in mock SMTP server by default. Run the following command:

```bash
USE_MOCK_SMTP=true go run examples/smtp_auth/main.go
```

Expected output:

```bash
Mock SMTP server running on :2525
A login link has been sent to: user@example.com
```

### **With a Real SMTP Server**

To use a real SMTP server, export the necessary environment variables and run the example:

```bash
export SMTP_HOST="smtp.example.com"
export SMTP_PORT="587"
export SMTP_USER="your-email@example.com"
export SMTP_PASSWORD="yourpassword"
go run examples/smtp_auth/main.go
```

Expected output:

```bash
A login link has been sent to: user@example.com
```

## **Endpoints Overview**

This example does not include a web frontend, but it generates a login link that users can click to authenticate.

- **Login Link Example:**

  ```bash
  https://myapp.com/login?token=abc123xyz
  ```

Users should submit the token via the backend to verify authentication.

## **Customization**

You can customize the SMTP settings by modifying the environment variables or by editing the SMTP transport initialization in the code:

```go
smtpTransport := &transport.SMTPTransport{
    Host: os.Getenv("SMTP_HOST"),
    Port: os.Getenv("SMTP_PORT"),
    From: os.Getenv("SMTP_USER"),
    Auth: smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST")),
}
```

## **When to Use This Example**

This approach is ideal for:

- **Web Applications:** Providing passwordless login links via email.
- **Internal Tools:** Simple, secure authentication without passwords.
- **Development & Testing:** Easily test email flows with a mock server.

## **How to Extend This Example**

- **Switch to Database Storage:** Use `DbStore` instead of `MemStore` for persistence.
- **Integrate with Web Frameworks:** Connect it to a web frontend for user-friendly authentication.
- **Add Logging:** Monitor email sending activities with structured logging.
