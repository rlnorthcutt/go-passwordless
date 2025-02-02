# **REST API Authentication Example**

This example demonstrates how to implement passwordless authentication using the `go-passwordless` library in a REST API. The API allows users to request a login token via their email and verify it through query parameters.

## **How It Works**

1. **Start the server.**
   The server starts and displays a login URL for easy access.

2. **Initiate login by visiting the login endpoint.**
   The email is passed as a query parameter, and a login token is generated.

3. **Receive and use the token.**
   The API response includes the login code and a link to the verification endpoint.

4. **Verify the token.**
   The user submits the token ID and code via query parameters, and the API responds with success or failure.

## **Running the Example**

1. Start the server:

   ```bash
   go run -buildvcs=false main.go
   ```

   Expected output:

   ```bash
   Server running at http://localhost:8080
   To start login, visit: http://localhost:8080/login?email=me@here.com
   ```

2. Initiate login by opening the following URL in your browser or using `curl`:

   ```bash
   curl "http://localhost:8080/login?email=me@here.com"
   ```

   Example response:

   ```bash
   Login code sent to: me@here.com
   Token ID: abc123xyz
   Verify at: http://localhost:8080/verify?token_id=abc123xyz&code=YOUR_CODE
   ```

3. Verify the login by replacing `YOUR_CODE` with the actual code (printed in logs when using `LogTransport`):

   ```bash
   curl "http://localhost:8080/verify?token_id=abc123xyz&code=123456"
   ```

   Expected responses:

   - **Successful login:**

     ```bash
     Verification successful! You are now logged in.
     ```

   - **Failed verification (wrong code):**

     ```bash
     Verification failed! Invalid code.
     ```

## **Endpoints**

### **1. `/login` (Start login process)**

- **Query Parameters:**
  - `email` – The user's email address (e.g., `me@here.com`)

- **Example Request:**

  ```bash
  http://localhost:8080/login?email=me@here.com
  ```

- **Response:**

  ```bash
  Login code sent to: me@here.com
  Token ID: abc123xyz
  Verify at: http://localhost:8080/verify?token_id=abc123xyz&code=YOUR_CODE
  ```

### **2. `/verify` (Verify login token)**

- **Query Parameters:**
  - `token_id` – The token ID received during login.
  - `code` – The one-time verification code.

- **Example Request:**

  ```bash
  http://localhost:8080/verify?token_id=abc123xyz&code=123456
  ```

- **Response (Success):**

  ```bash
  Verification successful! You are now logged in.
  ```

- **Response (Failure):**

  ```bash
  Verification failed! Invalid code.
  ```

---

## **When to Use This Example**

This REST API authentication example is useful when:

- **You need a simple API-based login system:**
  Great for integrating with web applications, mobile apps, or automation scripts.

- **You're building a proof of concept:**
  Quick way to validate the passwordless authentication workflow.

- **You want to understand how passwordless login works:**
  Learn the key concepts of generating and verifying login tokens without passwords.

---

## **How to Extend This Example**

- **Persist Tokens:**
  Replace the in-memory store with `FileStore` or `DbStore` to persist tokens across restarts.

- **Send Real Emails:**
  Swap out `LogTransport` for `SMTPTransport` to send actual verification codes via email.

- **Enhance Security:**
  Implement rate limiting, request logging, and token expiration checks.

- **Add Frontend UI:**
  Create an HTML page to capture email input and display the verification prompt.

## **Conclusion**

This example provides a lightweight passwordless authentication system that can be easily integrated into various applications. It demonstrates how to:

- Generate and verify authentication tokens.
- Handle query parameters for user-friendly authentication.
- Respond with success/failure messages for a seamless login experience.
