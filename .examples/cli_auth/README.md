# **CLI Authentication Example**

This example demonstrates a simple passwordless authentication flow using the `go-passwordless` library in a command-line interface (CLI).

## **How It Works**

1. The program prompts the user to enter their email address.
2. A one-time login token is generated and displayed in the console (simulating an email being sent).
3. The user enters the token to complete the login.
4. If the correct token is entered, login is successful; otherwise, it fails.

Note that this example doesn't actually send an email. It uses the `LogTransport` to output to the CLI for testing.

**Example Run:**

```plaintext
Enter your email: user@example.com
A verification code has been sent to: user@example.com
Enter the verification code: 123456
Login successful!
```

## **Running the Example**

To test the CLI example, run the following command:

```bash
go run main.go
```

## **When to Use This**

- **For quick prototyping:** Easily test passwordless authentication without building a full web interface.
- **For internal tools:** Secure access to CLI-based applications without requiring passwords.
- **For learning purposes:** Understand the core concepts of passwordless authentication before integrating into larger systems.
