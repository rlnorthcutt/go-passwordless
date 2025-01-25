# **Transport Options in `go-passwordless`**

The `transport` package provides various ways to deliver authentication tokens to users. This document explains the available transport options, their use cases, pros and cons, and how to implement a custom transport.

## **Available Transport Options**

### 1. **Log Transport (`LogTransport`)**

**Description:**
A simple transport that logs token deliveries to the application console (stdout). It is primarily intended for development and debugging purposes.

**Use Cases:**

- Local development and testing.
- Applications where logging token codes can help with debugging.

**Pros:**

- Easy to set up with no external dependencies.
- Useful for debugging and development.
- No cost associated with email/SMS providers.

**Cons:**

- Not suitable for production environments.
- Tokens are only logged, not actually sent to users.

**Usage Example:**

```go
lt := &transport.LogTransport{}
err := lt.Send(context.Background(), "test@example.com", "123456")
if err != nil {
    log.Fatalf("Error sending token: %v", err)
}
```

---

### 2. **SMTP Transport (`SMTPTransport`)**

**Description:**
Sends tokens via email using an SMTP server. This transport is useful for production environments where email-based authentication is required.

**Use Cases:**

- Production applications requiring email-based authentication.
- Verification codes for user onboarding and passwordless login.

**Pros:**

- Works with any SMTP-compatible email service (e.g., Gmail, Mailgun, SendGrid).
- Provides a reliable way to deliver authentication codes to users.

**Cons:**

- Requires an SMTP server configuration.
- May introduce latency based on email provider performance.
- SMTP credentials must be securely managed.

**Usage Example:**

```go
tr := &transport.SMTPTransport{
    Host: "smtp.example.com",
    Port: "587",
    From: "noreply@example.com",
    Auth: smtp.PlainAuth("", "user", "pass", "smtp.example.com"),
}
err := tr.Send(context.Background(), "user@example.com", "123456")
if err != nil {
    log.Fatalf("Error sending token: %v", err)
}
```

## **Choosing the Right Transport Option**

| Feature          | LogTransport  | SMTPTransport  |
|-----------------|---------------|---------------|
| Ease of Setup   | ✅ (none)       | ⚠️ (moderate)  |
| Cost            | ✅ (free)       | ⚠️ (depends on provider) |
| Production Ready| ❌ (no)         | ✅ (yes)       |
| Latency         | ✅ (instant)    | ⚠️ (depends on SMTP) |
| Security        | ⚠️ (minimal)    | ✅ (secure with proper config) |

## **How to Implement Your Own Transport**

If the existing transport options do not meet your requirements, you can create a custom transport by implementing the `Transport` interface.

### **Interface Definition:**

```go
type Transport interface {
    Send(ctx context.Context, recipient, tokenCode string) error
}
```

### **Steps to Create a Custom Transport:**

1. **Define a struct that implements the `Transport` interface.**
   Example:

   ```go
   type CustomTransport struct {}
   ```

2. **Implement the `Send` method to handle token delivery.**
   Example:

   ```go
   func (c *CustomTransport) Send(ctx context.Context, recipient, tokenCode string) error {
       // Custom logic to send the token (e.g., SMS, push notifications)
       fmt.Printf("Sending token %s to %s\n", tokenCode, recipient)
       return nil
   }
   ```

3. **Use your custom transport in your application.**
   Example:

   ```go
   myTransport := &CustomTransport{}
   err := myTransport.Send(context.Background(), "user@example.com", "987654")
   if err != nil {
       log.Fatalf("Error sending token: %v", err)
   }
   ```

## **Security Considerations**

When choosing or implementing a transport, consider the following:

1. **Sensitive Data Handling:**
   - Ensure recipient details and tokens are not logged in production environments.

2. **Rate Limiting:**
   - Implement rate limiting to prevent abuse of email/SMS services.

3. **SMTP Security:**
   - Use secure authentication methods and avoid storing credentials in plain text.

4. **Monitoring:**
   - Consider logging transport errors for monitoring failed delivery attempts.

## **Testing Your Transport Implementation**

You can create automated tests to ensure your transport implementation works as expected. Here’s an example of a test structure:

```go
func TestCustomTransport_Send(t *testing.T) {
    ct := &CustomTransport{}

    err := ct.Send(context.Background(), "test@example.com", "654321")
    if err != nil {
        t.Fatalf("CustomTransport.Send failed: %v", err)
    }
}
```

For SMTP transports, consider using [MailHog](https://github.com/mailhog/MailHog) to test email deliveries locally without requiring an external SMTP service.

## **Conclusion**

- Use **`LogTransport`** for local development and debugging.  
- Use **`SMTPTransport`** for production email delivery.  
- Implement a **custom transport** for use cases like SMS, push notifications, or third-party services.
