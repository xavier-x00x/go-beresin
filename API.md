# Go Beresin API Documentation

Base URL: `http://localhost:8080`

## Response Format

### Success
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error
```json
{
  "success": false,
  "message": "Error description"
}
```

### Validation Error (422)
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    { "field": "Email", "message": "invalid email format" },
    { "field": "Password", "message": "Password must be at least 8 characters" }
  ]
}
```

---

## Endpoints

### Health Check

`GET /health`

Check if the server is running.

**Response `200`**
```json
{
  "success": true,
  "message": "Server is healthy and running"
}
```

---

### Register

`POST /api/v1/auth/register`

Create a new user account.

**Request Body**
```json
{
  "email": "user@example.com",
  "phone": "6281234567890",
  "password": "Password123!",
  "full_name": "John Doe",
  "role": "user"
}
```

| Field | Type | Validation |
|---|---|---|
| `email` | string | required, valid email format |
| `phone` | string | required, Indonesian phone number (`+62`/`62`/`0` prefix, `8` followed by 7-11 digits) |
| `password` | string | required, min 8 characters |
| `full_name` | string | required |
| `role` | string | optional, must be `"user"` or `"talent"` (default: `"user"`) |

**Response `201`**
```json
{
  "success": true,
  "message": "Registration completed successfully. Please check your email to verify your account.",
  "data": {
    "user_id": "uuid-v7",
    "email": "user@example.com",
    "phone": "6281234567890",
    "full_name": "John Doe",
    "role": "user"
  }
}
```

**Response `409`** — Email or phone already registered

---

### Login

`POST /api/v1/auth/login`

Authenticate with email and password.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

| Field | Type | Validation |
|---|---|---|
| `email` | string | required |
| `password` | string | required |

**Response `200`**
```json
{
  "success": true,
  "message": "Successfully authenticated",
  "data": {
    "user_id": "uuid-v7",
    "full_name": "John Doe",
    "role": "user",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

**Response `401`** — Invalid email or password

**Response `429`** — IP temporarily blocked (after 5 failed attempts)

---

### Google Login

`POST /api/v1/auth/login/google`

Authenticate using Google OAuth ID token.

**Request Body**
```json
{
  "token": "google-oauth-id-token"
}
```

| Field | Type | Validation |
|---|---|---|
| `token` | string | required |

**Response `200`**
```json
{
  "success": true,
  "message": "Successfully authenticated via Google",
  "data": {
    "user_id": "uuid-v7",
    "full_name": "John Doe",
    "role": "user",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

### Send WhatsApp OTP

`POST /api/v1/auth/login/whatsapp/send-otp`

Send OTP code to phone number via WhatsApp.

**Request Body**
```json
{
  "phone": "6281234567890"
}
```

| Field | Type | Validation |
|---|---|---|
| `phone` | string | required, Indonesian phone number |

**Response `200`**
```json
{
  "success": true,
  "message": "OTP has been successfully dispatched to WhatsApp"
}
```

---

### Verify WhatsApp OTP

`POST /api/v1/auth/login/whatsapp/verify`

Verify OTP and authenticate via WhatsApp.

**Request Body**
```json
{
  "phone": "6281234567890",
  "otp": "123456"
}
```

| Field | Type | Validation |
|---|---|---|
| `phone` | string | required, Indonesian phone number |
| `otp` | string | required |

**Response `200`**
```json
{
  "success": true,
  "message": "Successfully authenticated via WhatsApp OTP",
  "data": {
    "user_id": "uuid-v7",
    "full_name": "John Doe",
    "role": "user",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

### Refresh Token

`POST /api/v1/auth/refresh-token`

Rotate access and refresh token pair.

**Request Body**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

| Field | Type | Validation |
|---|---|---|
| `refresh_token` | string | required |

**Response `200`**
```json
{
  "success": true,
  "message": "Tokens successfully rotated",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

### Logout

`POST /api/v1/auth/logout`

Invalidate refresh token and terminate session.

**Headers**
```
Authorization: Bearer <access_token>
```

**Response `200`**
```json
{
  "success": true,
  "message": "Successfully logged out, session terminated"
}
```

---

### Forgot Password

`POST /api/v1/auth/forgot-password`

Send password reset link to email.

**Request Body**
```json
{
  "email": "user@example.com"
}
```

| Field | Type | Validation |
|---|---|---|
| `email` | string | required, valid email format |

**Response `200`** (always returns success, even if email not found — security measure)
```json
{
  "success": true,
  "message": "If the email exists, password reset instructions have been dispatched."
}
```

---

### Reset Password

`POST /api/v1/auth/reset-password`

Reset password using email reset token.

**Request Body**
```json
{
  "token": "reset-token-from-email",
  "new_password": "NewPassword123!"
}
```

| Field | Type | Validation |
|---|---|---|
| `token` | string | required |
| `new_password` | string | required, min 8 characters |

**Response `200`**
```json
{
  "success": true,
  "message": "Password updated successfully. All other sessions have been logged out."
}
```

---

### Get Profile

`GET /api/v1/profile`

Fetch authenticated user identity from Bearer token claims.

**Headers**
```
Authorization: Bearer <access_token>
```

**Response `200`**
```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "user_id": "uuid-v7",
    "role": "user"
  }
}
```

---

### Upload File

`POST /api/v1/upload`

Upload a file (max 10 MB).

**Headers**
```
Authorization: Bearer <access_token>
```

**Request** — `multipart/form-data`

| Field | Type | Description |
|---|---|---|
| `file` | file | File to upload |

**Allowed MIME types:**
- `image/jpeg`, `image/png`, `image/webp`, `image/gif`
- `video/mp4`, `video/mpeg`, `video/quicktime`
- `application/pdf`, `application/msword`, `application/vnd.openxmlformats-officedocument.wordprocessingml.document`

**Response `201`**
```json
{
  "success": true,
  "message": "File uploaded and verified successfully",
  "data": {
    "file_name": "example.jpg",
    "file_size": 123456,
    "save_path": "uploads/example_sanitized.jpg"
  }
}
```

---

### Admin Dashboard

`GET /api/v1/admin/`

Access administrator dashboard (requires `admin` role).

**Headers**
```
Authorization: Bearer <access_token>
```

**Response `200`**
```json
{
  "success": true,
  "message": "Welcome to the Admin Dashboard",
  "data": {
    "admin_user_id": "uuid-v7",
    "privileges": "unlimited"
  }
}
```

---

## Authentication

### Token Format

All protected endpoints require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Token Lifetimes

| Token | Lifetime |
|---|---|
| Access Token | 15 minutes |
| Refresh Token | 30 days |

### Rate Limiting

| Scope | Limit |
|---|---|
| Global | 60 requests/minute per IP |
| Auth endpoints | 5 requests/minute per IP |
| Failed login block | 15 minutes after 5 failed attempts |
