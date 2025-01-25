# **Token Storage Options in `go-passwordless`**

The `store` package provides several implementations for token storage, allowing you to choose the best option based on your application's needs. This document will help you understand the available storage options, their use cases, pros and cons, and how to implement your own custom store.

## **Available Storage Options**

### 1. **Memory Store (`MemStore`)**

**Description:**
Stores tokens in-memory using a simple map, making it suitable for short-lived sessions and development environments.

**Use Cases:**

- Suitable for local development and testing.
- Temporary in-memory authentication where persistence is not required.

**Pros:**

- Fast and lightweight.
- No external dependencies.
- Easy to use and set up.

**Cons:**

- Tokens are lost when the application restarts.
- Not suitable for distributed or production environments.

**Usage Example:**

```go
memStore := store.NewMemStore()
```

### 2. **Cookie Store (`CookieStore`)**

**Description:**
Stores encrypted tokens in the user's browser cookies, ensuring that authentication works without requiring server-side storage.

**Use Cases:**

- Stateless authentication where no server-side storage is preferred.
- Applications with minimal backend requirements.

**Pros:**

- No server-side storage required.
- Easy to integrate with web applications.
- Persistent across sessions until the cookie expires.

**Cons:**

- Limited by cookie size (~4KB).
- Prone to client-side attacks if not properly secured.
- Requires proper security settings (e.g., `HttpOnly`, `Secure`, `SameSite`).

**Usage Example:**

```go
secretKey := []byte("super-secret-key")
cookieStore := store.NewCookieStore(secretKey)
```

### 3. **Database Store (`DbStore`)**

**Description:**
Stores tokens securely in an SQL database, providing persistent and reliable storage.

**Use Cases:**

- Applications that require persistent authentication tokens.
- Distributed or cloud-based environments.

**Pros:**

- Tokens persist across server restarts.
- Can be scaled with distributed deployments.
- Centralized token management.

**Cons:**

- Slightly slower compared to memory-based solutions.
- Requires database setup and maintenance.
- Additional dependencies (e.g., SQLite, PostgreSQL).

**Usage Example:**

```go
db, err := sql.Open("sqlite", "tokens.db")
if err != nil {
    log.Fatal(err)
}
dbStore := store.NewDbStore(db, "tokens")
```

### 4. **File Store (`FileStore`)**

**Description:**
Stores tokens as encrypted session files on disk, allowing persistent authentication without needing a database.

**Use Cases:**

- Local or small-scale applications needing persistent sessions.
- Offline applications where cloud storage isn't an option.
- Lightweight alternative to databases.

**Pros:**

- Easy setup without an external database.
- Persistent across restarts.
- Simple and reliable for single-instance applications.

**Cons:**

- Slower compared to memory storage due to disk I/O.
- Not suitable for distributed or multi-node applications.
- Requires file management and cleanup.

**Usage Example:**

```go
secretKey := []byte("super-secret-key")
fileStore := store.NewFileStore("./session_data", secretKey)
```

## **Choosing the Right Storage Option**

| Feature          | MemStore      | CookieStore   | DbStore       | FileStore     |
|-----------------|---------------|---------------|---------------|---------------|
| Persistence     | ❌ (no)        | ✅ (limited)   | ✅ (permanent) | ✅ (persistent) |
| Performance     | ✅ (fastest)   | ✅ (fast)      | ⚠️ (depends on DB) | ⚠️ (slower than memory) |
| Scalability     | ❌ (single node) | ✅ (stateless) | ✅ (multi-node) | ❌ (single-node) |
| Setup Effort    | ✅ (none)       | ✅ (minimal)   | ⚠️ (moderate)  | ✅ (simple setup)  |
| Security        | ⚠️ (limited)    | ⚠️ (browser-based) | ✅ (secure storage) | ✅ (secure storage) |

## **How to Implement Your Own Token Store**

If the existing storage options don't fit your needs, you can create your own by implementing the `TokenStore` interface.

### **Interface Definition:**

```go
type TokenStore interface {
    Store(ctx context.Context, token Token) error
    Exists(ctx context.Context, tokenID string) (*Token, error)
    Verify(ctx context.Context, tokenID, code string) (bool, error)
    Delete(ctx context.Context, tokenID string) error
}
```

### **Steps to Create a Custom Store:**

1. **Define a struct that implements the `TokenStore` interface.**  
   Example:

   ```go
   type MyCustomStore struct {
       tokens map[string]store.Token
   }

   func NewMyCustomStore() *MyCustomStore {
       return &MyCustomStore{tokens: make(map[string]store.Token)}
   }
   ```

2. **Implement the `Store` method to save tokens.**

   ```go
   func (s *MyCustomStore) Store(ctx context.Context, tok store.Token) error {
       s.tokens[tok.ID] = tok
       return nil
   }
   ```

3. **Implement the `Exists` method to retrieve tokens.**

   ```go
   func (s *MyCustomStore) Exists(ctx context.Context, tokenID string) (*store.Token, error) {
       tok, exists := s.tokens[tokenID]
       if !exists {
           return nil, fmt.Errorf("token not found")
       }
       return &tok, nil
   }
   ```

4. **Implement the `Verify` method to validate tokens.**

   ```go
   func (s *MyCustomStore) Verify(ctx context.Context, tokenID, code string) (bool, error) {
       tok, exists := s.tokens[tokenID]
       if !exists {
           return false, fmt.Errorf("token not found")
       }
       return string(tok.CodeHash) == code, nil
   }
   ```

5. **Implement the `Delete` method to remove tokens.**

   ```go
   func (s *MyCustomStore) Delete(ctx context.Context, tokenID string) error {
       delete(s.tokens, tokenID)
       return nil
   }
   ```

6. **Use your custom store in your application.**

   ```go
   customStore := NewMyCustomStore()
   err := customStore.Store(context.Background(), store.Token{ID: "test", Recipient: "user@test.com"})
   ```

## **Security Considerations**

When choosing or implementing a token store, consider the following:

1. **Data Sensitivity:**
   - Use encryption to store tokens securely in databases and cookies.

2. **Token Expiry:**
   - Ensure tokens are expired and removed after their intended lifespan.

3. **Secure Cookie Flags:**
   - Always use `Secure`, `HttpOnly`, and `SameSite` flags when using cookies.

4. **Rate Limiting:**
   - Protect verification endpoints from brute-force attacks.

5. **Logging:**
   - Avoid logging sensitive token data.

## **Conclusion**

- Use **`MemStore`** for testing or short-lived tokens.
- Use **`CookieStore`** for lightweight, stateless authentication.
- Use **`DbStore`** for persistent, scalable solutions.
- Use **`FileStore`** for persistent, file-based storage.
- Implement a **custom store** if your requirements are unique.
