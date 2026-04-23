# Ticket Engine - Clean Architecture & Domain-Driven Design in Go

A production-ready ticket engine API built with **Go**, implementing **Clean Architecture** and **Domain-Driven Design (DDD)** principles. The project demonstrates professional API design with secure authentication, comprehensive error handling, and modular architecture.

## Table of Contents

- [Overview](#overview)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Installation & Setup](#installation--setup)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
  - [Authentication Endpoints](#authentication-endpoints)
  - [Health Check](#health-check)
- [Authentication Flow](#authentication-flow)
- [Database Schema](#database-schema)
- [Development](#development)
- [Error Handling](#error-handling)

---

## Overview

The Ticket Engine is a backend service designed to manage ticket operations with a focus on security, scalability, and maintainability. It implements:

- **Session-Based Authentication** with Refresh Token rotation
- **Password Hashing** using bcrypt for enhanced security
- **Token Management** via SHA256-hashed refresh tokens stored in PostgreSQL
- **Role-Based Access Control** (User, Organizer, Admin)
- **Comprehensive Validation** with detailed error responses
- **Redis Caching** for performance optimization
- **Structured Logging** with Zerolog for debugging and monitoring
- **Clean Architecture** separation of concerns (Delivery, Domain, Infrastructure)

### Key Features

✅ Stateless session authentication with secure refresh tokens
✅ User registration with email validation
✅ Secure login with bcrypt password verification
✅ Token refresh mechanism for long-lived sessions
✅ Logout with token revocation
✅ User profile retrieval (protected endpoint)
✅ PostgreSQL for persistent storage
✅ Redis for caching and performance
✅ Comprehensive input validation with field-specific error messages
✅ Professional error handling and HTTP status codes

---

## Technology Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.26.1 |
| **Web Framework** | Gin v1.12.0 |
| **Database** | PostgreSQL 5.9.1 |
| **Cache** | Redis v9.18.0 |
| **Hashing** | bcrypt (via `golang.org/x/crypto`) |
| **Validation** | go-playground/validator/v10 |
| **Logging** | Zerolog v1.35.0 |
| **JWT/Token** | golang-jwt v5.3.1 |
| **Database Migration** | sqlc (type-safe SQL) |
| **ID Generation** | Google UUID v1.6.0 |
| **Environment** | godotenv v1.5.1, caarlos0/env v11.4.0 |

---

## Project Structure

```
go-ticket-engine/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── core/
│   │   ├── config/
│   │   │   └── config.go           # Configuration management
│   │   ├── database/
│   │   │   ├── postgres.go         # PostgreSQL connection pool
│   │   │   └── redis.go            # Redis client initialization
│   │   ├── delivery/
│   │   │   └── http/
│   │   │       ├── response.go     # HTTP response structures
│   │   │       └── validation.go   # Input validation utilities
│   │   └── middleware/
│   │       ├── auth.go             # JWT authentication middleware
│   │       └── logger.go           # Request logging middleware
│   ├── infra/                      # Infrastructure layer adapters
│   │   ├── crypto/
│   │   │   ├── bcrypt.go           # Password hashing
│   │   │   └── sha256.go           # Token hashing
│   │   ├── identifier/
│   │   │   └── uuid.go             # UUID generation
│   │   ├── sqlc/
│   │   │   └── *.go                # Auto-generated database queries
│   │   └── token/
│   │       └── jwt.go              # JWT token engine
│   └── modules/                    # Feature modules (DDD)
│       ├── auth/
│       │   ├── router.go           # Route registration
│       │   ├── delivery/           # HTTP request/response handling
│       │   ├── domain/             # Business logic & entities
│       │   │   ├── user.go
│       │   │   ├── refresh_token.go
│       │   │   ├── ports.go        # Interfaces (repository, hasher)
│       │   │   └── errors.go       # Domain-specific errors
│       │   ├── features/           # Use cases
│       │   │   ├── register/       # User registration
│       │   │   ├── login/          # User authentication
│       │   │   ├── logout/         # Session termination
│       │   │   ├── refresh/        # Token refresh
│       │   │   └── me/             # User profile retrieval
│       │   └── infra/              # Repository implementations
│       │       ├── user_repository.go
│       │       ├── refresh_token_repository.go
│       │       ├── bcrypt_hasher.go
│       │       └── jwt_token_generator.go
│       └── health/
│           └── handler.go          # Health check endpoint
├── db/
│   ├── migrations/                 # Database migration files
│   │   ├── 000001_init_scheme.up.sql
│   │   └── 000001_init_scheme.down.sql
│   └── queries/                    # SQL query definitions
│       ├── user.sql
│       └── auth.sql
├── go.mod                          # Go module definition
├── go.sum                          # Go module checksums
├── Makefile                        # Build and development commands
└── README.md                       # This file
```

---

## Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│         Delivery Layer (HTTP)           │  ← HTTP Handlers, Request/Response DTOs
├─────────────────────────────────────────┤
│         Domain Layer (Business Logic)   │  ← Entities, Value Objects, Errors, Interfaces
├─────────────────────────────────────────┤
│   Application Layer (Use Cases)         │  ← Feature handlers, Business rules
├─────────────────────────────────────────┤
│    Infrastructure Layer (Adapters)      │  ← Repository, Crypto, Token, Database
└─────────────────────────────────────────┘
```

### Data Flow Example: User Registration

```
HTTP Request
    ↓
[register.HTTPHandler] → Validates Input
    ↓
[register.Handler] → Use Case Logic
    ↓
[UserRepository] → Database Operation
    ↓
[BcryptHasher] → Password Hashing
    ↓
[UUIDGenerator] → ID Generation
    ↓
Database (PostgreSQL)
    ↓
[User Domain Model]
    ↓
HTTP Response (201 Created)
```

---

## Installation & Setup

### Prerequisites

- **Go** 1.26.1 or higher
- **PostgreSQL** 12.0 or higher
- **Redis** 6.0 or higher
- **Make** (optional, for build automation)

### Clone the Repository

```bash
git clone https://github.com/hogiabao7725/gin-auth-playground.git
cd gin-auth-playground
```

### Install Dependencies

```bash
go mod download
go mod tidy
```

### Setup PostgreSQL

```bash
# Create database
createdb ticket_engine

# Run migrations
# Migration files are located in db/migrations/
# Apply them using your preferred migration tool or psql
psql -U postgres -d ticket_engine -f db/migrations/000001_init_scheme.up.sql
```

### Setup Redis

```bash
# Start Redis server (if not already running)
redis-server

# Or using Docker
docker run -d -p 6379:6379 redis:latest
```

### Environment Configuration

Create a `.env` file in the root directory:

```bash
# Server Configuration
SERVER_PORT=8080
ENV=development

# PostgreSQL Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ticket_engine
DB_SSLMODE=disable
DB_MAX_CONNS=25
DB_MIN_CONNS=5
DB_CONNECT_TIMEOUT=10s
DB_CONN_LIFETIME=1h
DB_CONN_IDLE_TIME=30m

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_CONNECT_TIMEOUT=5s

# JWT Configuration (Session-based, NOT using JWT for tokens)
JWT_ACCESS_SECRET=your_access_secret_key_min_32_chars_recommended
JWT_REFRESH_SECRET=your_refresh_secret_key_min_32_chars_recommended
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=7d
```

### Build & Run

```bash
# Using Go directly
go run cmd/server/main.go

# Using Makefile (if available)
make build
make run

# For development with auto-reload, use a tool like air or CompileDaemon
air
```

The server will start on `http://localhost:8080` by default.

---

## Configuration

Configuration is managed through environment variables using the `caarlos0/env` package. The application loads configuration from:

1. `.env` file (if present)
2. Environment variables
3. Default values (where applicable)

### Configuration Structure

```go
type Config struct {
    Server ServerConfig    // HTTP server settings
    DB     DBConfig        // PostgreSQL connection settings
    Redis  RedisConfig     // Redis connection settings
    JWT    JWTConfig       // Token secrets and TTLs
}
```

### Important Configuration Notes

- **JWT_ACCESS_SECRET**: Minimum 32 characters recommended for security
- **JWT_REFRESH_SECRET**: Separate secret for refresh tokens
- **JWT_ACCESS_TTL**: Access token expiration (default: 15m)
- **JWT_REFRESH_TTL**: Refresh token expiration (default: 7d)
- **Database Pool Settings**: Tuned for production use (25 max, 5 min connections)

---

## API Endpoints

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### 1. User Registration

Register a new user account with email and password.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Request Parameters:**

| Field | Type | Required | Validation | Description |
|-------|------|----------|-----------|-------------|
| `name` | string | Yes | Not blank, max 255 chars | User's full name |
| `email` | string | Yes | Valid email format | User's email address (case-insensitive) |
| `password` | string | Yes | Min 6 characters | User's password (bcrypt hashed) |

**Success Response (201 Created):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "role": "user",
    "created_at": "2026-04-23T10:30:00Z"
  }
}
```

**Error Responses:**

| Status | Error Code | Message | Cause |
|--------|-----------|---------|-------|
| 400 | VALIDATION_ERROR | See `fields` object | Invalid request format or missing required fields |
| 409 | CONFLICT | `user with this email already exists` | Email already registered |
| 500 | SERVER_ERROR | `internal server error` | Unexpected server error |

**Validation Error Response Example:**
```json
{
  "message": "validation error",
  "fields": {
    "email": "must be a valid email address",
    "password": "must be at least 6 characters long",
    "name": "cannot be blank"
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securePassword123"
  }'
```

---

#### 2. User Login

Authenticate user and receive access token with refresh token in secure cookie.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Request Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address |
| `password` | string | Yes | User's password (not hashed in request) |

**Success Response (200 OK):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "user"
    }
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `access_token` | string | JWT token for authenticating subsequent requests (15 min TTL) |
| `expires_in` | integer | Seconds until access token expires |
| `user` | object | Authenticated user information |
| `user.id` | string | User's UUID |
| `user.name` | string | User's full name |
| `user.email` | string | User's email |
| `user.role` | string | User's role: `user`, `organizer`, or `admin` |

**Response Cookies:**

| Cookie | Value | Attributes |
|--------|-------|-----------|
| `refresh_token` | JWT token | HttpOnly, Secure, Path=/, Max-Age=604800 (7 days) |

**Error Responses:**

| Status | Message | Cause |
|--------|---------|-------|
| 400 | `validation error` | Missing email or password |
| 401 | `invalid email or password` | Incorrect credentials |
| 404 | `user not found` | User doesn't exist |
| 500 | `internal server error` | Unexpected error |

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securePassword123"
  }' \
  -c cookies.txt
```

---

#### 3. Refresh Access Token

Obtain a new access token using the refresh token (stored in secure cookie).

**Endpoint:** `POST /auth/refresh`

**Request:**
- No request body required
- Requires `refresh_token` cookie from login response

**Success Response (200 OK):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `access_token` | string | New JWT access token (15 min TTL) |
| `expires_in` | integer | Seconds until new token expires |

**Response Cookies:**

| Cookie | Description |
|--------|-------------|
| `refresh_token` | New refresh token (rotated for security) |

**Error Responses:**

| Status | Message | Cause |
|--------|---------|-------|
| 401 | `refresh token is missing` | Cookie not provided |
| 401 | `invalid refresh token` | Token tampered or invalid format |
| 401 | `refresh token has expired` | Token TTL exceeded |
| 401 | `refresh token has been revoked` | Token was logged out or manually revoked |
| 500 | `internal server error` | Unexpected error |

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -b cookies.txt
```

---

#### 4. User Logout

Revoke the refresh token and end the session.

**Endpoint:** `POST /auth/logout`

**Request:**
- No request body required
- Requires `refresh_token` cookie

**Success Response (200 OK):**
```json
{
  "data": {
    "message": "logged out successfully"
  }
}
```

**Response Cookies:**

| Cookie | Action |
|--------|--------|
| `refresh_token` | Deleted (Max-Age: -1) |

**Error Responses:**

| Status | Message | Cause |
|--------|---------|-------|
| 401 | `invalid refresh token` | Token invalid or tampered |
| 401 | `refresh token has expired` | Token already expired |
| 401 | `refresh token has been revoked` | Already logged out |
| 500 | `internal server error` | Unexpected error |

**Notes:**
- After logout, the refresh token is permanently revoked in the database
- Subsequent login requests will be required

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -b cookies.txt
```

---

#### 5. Get Current User Profile

Retrieve authenticated user's information (protected endpoint).

**Endpoint:** `GET /auth/me`

**Request Headers:**

| Header | Value | Required | Description |
|--------|-------|----------|-------------|
| `Authorization` | `Bearer <access_token>` | Yes | JWT access token from login/refresh |

**Success Response (200 OK):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "role": "user",
    "created_at": "2026-04-23T10:30:00Z"
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | User's UUID |
| `name` | string | User's full name |
| `email` | string | User's email address |
| `role` | string | User's role: `user`, `organizer`, or `admin` |
| `created_at` | string | Account creation timestamp (RFC3339 format) |

**Error Responses:**

| Status | Message | Cause |
|--------|---------|-------|
| 401 | `unauthorized` | Missing or invalid access token |
| 404 | `user not found` | User doesn't exist (rare) |
| 500 | `internal server error` | Unexpected error |

**cURL Example:**
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### Health Check

#### Health Status

Check if the API is running and healthy.

**Endpoint:** `GET /healthz`

**Request:**
- No authentication required
- No request body

**Success Response (200 OK):**
```json
{
  "status": "ok",
  "message": "Ticket Engine is healthy"
}
```

**cURL Example:**
```bash
curl http://localhost:8080/api/v1/healthz
```

---

## Authentication Flow

### Session-Based Authentication with Refresh Tokens

The API uses **session-based authentication** (NOT JWT-based) with the following flow:

```
┌─────────────────────────────────────────────────────────────┐
│                    AUTHENTICATION FLOW                      │
└─────────────────────────────────────────────────────────────┘

1. REGISTRATION
   User → POST /auth/register → Server
   └─ Create bcrypt-hashed password
   └─ Store user in PostgreSQL
   └─ Return 201 Created with user info

2. LOGIN
   User → POST /auth/login → Server
   ├─ Verify credentials (bcrypt password comparison)
   ├─ Generate access token (JWT, 15 min expiry)
   ├─ Generate refresh token (JWT, 7 day expiry)
   ├─ Hash refresh token with SHA256
   ├─ Store hashed refresh token in PostgreSQL
   ├─ Set refresh_token in HttpOnly, Secure cookie
   └─ Return 200 OK with access_token

3. PROTECTED REQUEST (e.g., GET /auth/me)
   User → Request with Authorization: Bearer <access_token>
   ├─ AuthMiddleware validates JWT signature
   ├─ Verify token not expired
   ├─ Extract user ID from claims
   └─ Proceed to handler if valid, else 401 Unauthorized

4. TOKEN REFRESH
   User → POST /auth/refresh with refresh_token cookie
   ├─ Extract and hash refresh token
   ├─ Lookup in PostgreSQL
   ├─ Verify not revoked and not expired
   ├─ Generate new access token
   ├─ Rotate refresh token (new hash stored in DB)
   ├─ Set new refresh_token in cookie
   └─ Return 200 OK with new access_token

5. LOGOUT
   User → POST /auth/logout with refresh_token cookie
   ├─ Extract and hash refresh token
   ├─ Mark as revoked in PostgreSQL (delete from table)
   ├─ Clear refresh_token cookie
   └─ Return 200 OK

6. SUBSEQUENT LOGIN AFTER LOGOUT
   User → POST /auth/login → Server
   ├─ Old refresh_token in cookie is ignored
   ├─ New session created with new tokens
   └─ Return 200 OK with new access_token
```

### Key Security Features

| Feature | Implementation |
|---------|------------------|
| **Password Security** | bcrypt hashing with salt |
| **Token Storage** | Refresh tokens stored as SHA256 hashes in DB |
| **Cookie Security** | HttpOnly, Secure, Path=/ attributes |
| **Token Expiration** | Access (15 min), Refresh (7 days) |
| **Token Revocation** | Refresh tokens can be revoked via logout |
| **CORS** | Can be configured in middleware layer |
| **HTTPS** | Recommended in production (Secure cookie flag) |

### Token Validation Middleware

The `AuthMiddleware` validates incoming access tokens:

```go
// Request with valid token → ✅ Allowed
GET /auth/me
Header: Authorization: Bearer <valid_jwt>
Result: 200 OK

// Request with expired token → ❌ Rejected
GET /auth/me
Header: Authorization: Bearer <expired_jwt>
Result: 401 Unauthorized

// Request without token → ❌ Rejected
GET /auth/me
Result: 401 Unauthorized

// Request with invalid signature → ❌ Rejected
GET /auth/me
Header: Authorization: Bearer <tampered_jwt>
Result: 401 Unauthorized
```

---

## Database Schema

### User Table

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL DEFAULT 'user',
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

**Indexes:**
- `UNIQUE INDEX ON email` - Enforce email uniqueness

### Refresh Token Table

```sql
CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash VARCHAR(64) NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
```

**Indexes:**
- `UNIQUE INDEX ON user_id` - One active refresh token per user
- `INDEX ON expires_at` - For efficient cleanup of expired tokens

### User Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| `user` | Regular user (default) | View own profile, use tickets |
| `organizer` | Event organizer | Manage events, create tickets |
| `admin` | System administrator | Full system access |

---

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/modules/auth/...
```

### Code Generation

The project uses `sqlc` for type-safe database queries:

```bash
# Regenerate SQL code
sqlc generate
```

### Database Migrations

```bash
# Run migrations (using your preferred migration tool)
# Migrations are in db/migrations/

# Up
psql -d ticket_engine -f db/migrations/000001_init_scheme.up.sql

# Down
psql -d ticket_engine -f db/migrations/000001_init_scheme.down.sql
```

### Project Commands (Makefile)

```bash
make build       # Build the application
make run         # Run the application
make test        # Run tests
make fmt         # Format code
make lint        # Run linter
```

---

## Error Handling

### Error Response Format

All error responses follow a consistent format:

```json
{
  "message": "Error description",
  "fields": {
    "field_name": "Specific error message"
  }
}
```

### HTTP Status Codes

| Code | Usage | Example |
|------|-------|---------|
| **200 OK** | Successful request | Login, logout, token refresh |
| **201 Created** | Resource created | User registration |
| **400 Bad Request** | Invalid input | Missing fields, validation errors |
| **401 Unauthorized** | Authentication failed | Invalid credentials, expired token |
| **404 Not Found** | Resource not found | User doesn't exist |
| **409 Conflict** | Resource conflict | Email already exists |
| **500 Internal Server Error** | Server error | Database connection failure |

### Common Error Messages

| Message | Status | Cause |
|---------|--------|-------|
| `validation error` | 400 | See `fields` for details |
| `invalid email or password` | 401 | Wrong credentials |
| `user with this email already exists` | 409 | Email already registered |
| `refresh token is missing` | 401 | Cookie not provided |
| `refresh token has expired` | 401 | Token TTL exceeded |
| `refresh token has been revoked` | 401 | User logged out |
| `unauthorized` | 401 | Invalid/missing access token |
| `user not found` | 404 | User doesn't exist |
| `internal server error` | 500 | Unexpected server error |

### Validation Error Details

When validation fails, the `fields` object contains specific error messages:

```json
{
  "message": "validation error",
  "fields": {
    "email": "must be a valid email address",
    "password": "must be at least 6 characters long",
    "name": "cannot be blank"
  }
}
```

**Validation Messages:**

| Message | Field | Meaning |
|---------|-------|---------|
| `is required` | Any | Field is missing |
| `cannot be blank` | `name` | Field is empty or whitespace only |
| `must be a valid email address` | `email` | Invalid email format |
| `must be at least X characters long` | `password` | Password too short |
| `user with this email already exists` | `email` | Email already registered |

**Last Updated:** April 2026
