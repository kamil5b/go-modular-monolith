# Go Modular Monolith - Technical Documentation

> **Version:** 1.1.0  
> **Last Updated:** December 1, 2025  
> **Go Version:** 1.24.7

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Project Structure](#project-structure)
4. [Getting Started](#getting-started)
5. [Configuration](#configuration)
6. [Feature Flags](#feature-flags)
7. [Modules](#modules)
8. [Domain Layer](#domain-layer)
9. [Infrastructure Layer](#infrastructure-layer)
10. [Transport Layer](#transport-layer)
11. [Authentication & Middleware](#authentication--middleware)
12. [Database & Migrations](#database--migrations)
13. [API Reference](#api-reference)
14. [Development Guide](#development-guide)
15. [Roadmap](#roadmap)

---

## Overview

Go Modular Monolith is a production-ready, modular monolithic application built with Go. It implements clean architecture principles with pluggable components, allowing teams to:

- **Switch HTTP frameworks** (Echo/Gin) via configuration
- **Swap database backends** (PostgreSQL/MongoDB) per module
- **Version handlers, services, and repositories** independently
- **Enable/disable features** through feature flags

### Key Technologies

| Component | Technology |
|-----------|------------|
| Language | Go 1.24.7 |
| HTTP Frameworks | Echo v4, Gin v1 |
| SQL Database | PostgreSQL (via sqlx) |
| NoSQL Database | MongoDB |
| Migrations | Goose (SQL), mongosh (MongoDB) |
| Authentication | JWT (golang-jwt/jwt/v5) |
| Logging | Zerolog |
| Password Hashing | golang.org/x/crypto/bcrypt |
| UUID Generation | github.com/google/uuid |

---

## Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                      Transport Layer                        │
│                   (Echo/Gin Adapters)                       │
├─────────────────────────────────────────────────────────────┤
│                      Handler Layer                          │
│              (HTTP Request/Response Handling)               │
├─────────────────────────────────────────────────────────────┤
│                      Service Layer                          │
│                  (Business Logic)                           │
├─────────────────────────────────────────────────────────────┤
│                    Repository Layer                         │
│              (Data Access - SQL/MongoDB)                    │
├─────────────────────────────────────────────────────────────┤
│                   Infrastructure Layer                      │
│           (Database Connections, External Services)         │
└─────────────────────────────────────────────────────────────┘
```

### Dependency Flow

```
main.go
   │
   ▼
bootstrap/
   │
   ├── LoadConfig()
   ├── LoadFeatureFlags()
   │
   ▼
Container (DI)
   │
   ├── Repository (SQL/MongoDB)
   ├── Service (v1/noop)
   └── Handler (v1/noop)
   │
   ▼
HTTP Server (Echo/Gin)
   │
   ▼
Routes → Handlers → Services → Repositories → Database
```

---

## Project Structure

```
go-modular-monolith/
├── main.go                          # Application entry point
├── go.mod                           # Go module definition
├── config/
│   ├── config.yaml                  # Application configuration
│   └── featureflags.yaml            # Feature flag configuration
├── cmd/
│   └── bootstrap/
│       ├── bootstrap.server.go      # Server initialization
│       └── bootstrap.migration.go   # Migration runner
├── internal/
│   ├── app/
│   │   ├── core/
│   │   │   ├── config.go            # Config structs & loader
│   │   │   ├── container.go         # Dependency injection container
│   │   │   └── feature_flag.go      # Feature flag structs & loader
│   │   └── http/
│   │       ├── echo.go              # Echo server setup
│   │       ├── gin.go               # Gin server setup
│   │       └── routes.go            # Route definitions
│   ├── domain/
│   │   ├── auth/
│   │   │   ├── model.auth.go        # Auth entities (Session, Credential)
│   │   │   ├── interface.auth.go    # Handler/Service/Repo/Middleware interfaces
│   │   │   ├── request.auth.go      # Auth request DTOs
│   │   │   └── response.auth.go     # Auth response DTOs
│   │   ├── product/
│   │   │   ├── model.product.go     # Product entity
│   │   │   ├── interface.product.go # Handler/Service/Repo interfaces
│   │   │   ├── request.product.go   # Request DTOs
│   │   │   └── response.product.go  # Response DTOs
│   │   └── user/
│   │       ├── model.user.go        # User entity
│   │       ├── interface.user.go    # Handler/Service/Repo interfaces
│   │       ├── request.user.go      # Request DTOs
│   │       └── response.user.go     # Response DTOs
│   ├── infrastructure/
│   │   ├── db/
│   │   │   ├── sql/
│   │   │   │   ├── sql.go           # PostgreSQL connection
│   │   │   │   └── migration/       # SQL migration files
│   │   │   └── mongo/
│   │   │       ├── mongo.go         # MongoDB connection
│   │   │       └── migration/       # MongoDB migration scripts
│   │   ├── cache/                   # Redis (planned)
│   │   ├── logger/                  # Logger infrastructure
│   │   └── storage/                 # File storage (planned)
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── handler/
│   │   │   │   ├── v1/              # v1 handler implementation
│   │   │   │   └── noop/            # No-op/disabled handler
│   │   │   ├── middleware/          # Auth middleware (JWT, Session, Basic)
│   │   │   ├── service/
│   │   │   │   ├── v1/              # v1 service implementation
│   │   │   │   └── noop/            # No-op/disabled service
│   │   │   └── repository/
│   │   │       ├── sql/             # PostgreSQL repository
│   │   │       ├── mongo/           # MongoDB repository
│   │   │       └── noop/            # No-op repository
│   │   ├── product/
│   │   │   ├── handler/
│   │   │   │   ├── v1/              # v1 handler implementation
│   │   │   │   └── noop/            # No-op/disabled handler
│   │   │   ├── service/
│   │   │   │   ├── v1/              # v1 service implementation
│   │   │   │   └── noop/            # No-op/disabled service
│   │   │   └── repository/
│   │   │       ├── sql/             # PostgreSQL repository
│   │   │       ├── mongo/           # MongoDB repository
│   │   │       └── noop/            # No-op repository
│   │   └── user/
│   │       ├── handler/v1/          # v1 handler implementation
│   │       ├── service/v1/          # v1 service implementation
│   │       └── repository/sql/      # PostgreSQL repository
│   ├── transports/
│   │   └── http/
│   │       ├── echo/
│   │       │   ├── adapter.echo.go  # Echo route adapter
│   │       │   └── context.echo.go  # Echo context wrapper
│   │       └── gin/
│   │           ├── adapter.gin.go   # Gin route adapter
│   │           └── context.gin.go   # Gin context wrapper
│   └── proto/                       # gRPC protobuf definitions (planned)
└── pkg/
    ├── constant/                    # Shared constants
    ├── logger/
    │   └── logger.go                # Shared logger utilities
    ├── model/
    │   ├── request.go               # Common request models
    │   └── response.go              # Common response models
    └── routes/
        └── route.go                 # Route struct definition
```

---

## Getting Started

### Prerequisites

- Go 1.24.7+
- PostgreSQL 14+
- MongoDB 6.0+ (optional)
- mongosh CLI (for MongoDB migrations)

### Installation

```bash
# Clone the repository
git clone https://github.com/kamil5b/go-modular-monolith.git
cd go-modular-monolith

# Install dependencies
go mod tidy

# Configure database connection
# Edit config/config.yaml with your database credentials
```

### Running the Application

```bash
# Default: Run migrations and start server
go run .

# Run only the server (skip migrations)
go run . server

# Run SQL migrations
go run . migration sql up
go run . migration sql down

# Run MongoDB migrations
go run . migration mongo up
```

### CLI Commands

| Command | Description |
|---------|-------------|
| `go run .` | Run SQL migrations (up) then start server |
| `go run . server` | Start server only |
| `go run . migration sql up` | Apply SQL migrations |
| `go run . migration sql down` | Rollback SQL migrations |
| `go run . migration mongo up` | Apply MongoDB migrations |

---

## Configuration

### config/config.yaml

```yaml
environment: development  # development | production

app:
  server:
    port: "8080"          # HTTP server port
    grpc_port: "9090"     # gRPC server port (planned)

  database:
    sql:
      db_url: "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
    mongo:
      mongo_url: "mongodb://localhost:27017"
      mongo_db: "myapp_db"

  jwt:
    secret: "supersecretkey"  # JWT signing secret
```

### Configuration Struct

```go
type Config struct {
    Environment string    `yaml:"environment"`
    App         AppConfig `yaml:"app"`
}

type AppConfig struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
    Port     string `yaml:"port"`
    GRPCPort string `yaml:"grpc_port"`
}

type DatabaseConfig struct {
    SQL   SQLConfig   `yaml:"sql"`
    Mongo MongoConfig `yaml:"mongo"`
}
```

---

## Feature Flags

Feature flags allow dynamic component selection without code changes.

### config/featureflags.yaml

```yaml
http_handler: echo  # echo | gin

handler:
  authentication: v1   # v1 | disable
  product: v1          # v1 | disable
  user: v1             # v1 | disable

service:
  authentication: v1   # v1 | disable
  product: v1          # v1 | disable
  user: v1             # v1 | disable

repository:
  authentication: postgres  # postgres | mongo | disable
  product: postgres         # postgres | mongo | disable
  user: postgres            # postgres | mongo | disable
```

### Feature Flag Options

| Component | Options | Description |
|-----------|---------|-------------|
| `http_handler` | `echo`, `gin` | HTTP framework selection |
| `handler.*` | `v1`, `disable` | Handler version or disabled |
| `service.*` | `v1`, `disable` | Service version or disabled |
| `repository.*` | `postgres`, `mongo`, `disable` | Database backend |

### How It Works

The `Container` in `internal/app/core/container.go` reads feature flags and instantiates the appropriate implementations:

```go
// Repository selection example
switch featureFlag.Repository.Product {
case "mongo":
    productRepository = repoMongo.NewMongoRepository(mongo, "appdb")
case "postgres":
    productRepository = repoSQL.NewSQLRepository(db)
default:
    // No-op or disabled
}
```

---

## Modules

### Module Structure

Each business module follows a consistent structure:

```
modules/<module>/
├── handler/
│   ├── v1/handler_v1.<module>.go    # Version 1 implementation
│   └── noop/handler_noop.<module>.go # No-op implementation
├── service/
│   ├── v1/service_v1.<module>.go    # Version 1 implementation
│   └── noop/service_noop.<module>.go # No-op implementation
└── repository/
    ├── sql/repository_sql.<module>.go   # PostgreSQL implementation
    ├── mongo/repository_mongo.<module>.go # MongoDB implementation
    └── noop/repository_noop.<module>.go  # No-op implementation
```

### Current Modules

#### Product Module
- **Status:** ✅ Complete
- **Features:** CRUD operations
- **Repository:** PostgreSQL, MongoDB

#### User Module
- **Status:** ✅ Complete  
- **Features:** CRUD operations
- **Repository:** PostgreSQL

#### Auth Module
- **Status:** ✅ Complete (untested)
- **Features:** JWT authentication, session management, Basic Auth, middleware
- **Repository:** PostgreSQL, MongoDB

---

## Domain Layer

### Context Interface

Each module defines a `Context` interface that abstracts HTTP framework specifics:

```go
// Product/User Context
type Context interface {
    BindJSON(obj any) error      // Bind JSON body
    BindURI(obj any) error       // Bind URI parameters
    BindQuery(obj any) error     // Bind query parameters
    BindHeader(obj any) error    // Bind headers
    Bind(obj any) error          // Auto-bind based on content type
    JSON(code int, v any) error  // Send JSON response
    Param(name string) string    // Get URL parameter
    GetUserID() string           // Get authenticated user ID
    Get(key string) any          // Get context value
    GetContext() context.Context // Get Go context
}

// Auth Context (extended with additional methods)
type AuthContext interface {
    Context
    Set(key string, value any)                                                  // Set context value
    GetHeader(key string) string                                                // Get request header
    SetHeader(key, value string)                                                // Set response header
    GetCookie(name string) (string, error)                                      // Get cookie value
    SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) // Set cookie
    RemoveCookie(name string)                                                   // Remove cookie
    GetClientIP() string                                                        // Get client IP address
    GetUserAgent() string                                                       // Get user agent
}
```

### Product Domain

#### Entity

```go
type Product struct {
    ID          string     `db:"id" json:"id" bson:"id"`
    Name        string     `db:"name" json:"name" bson:"name"`
    Description string     `db:"description" json:"description" bson:"description"`
    CreatedAt   time.Time  `db:"created_at" json:"created_at" bson:"created_at"`
    CreatedBy   string     `db:"created_by" json:"created_by" bson:"created_by"`
    UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
    UpdatedBy   *string    `db:"updated_by" json:"updated_by,omitempty"`
    DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
    DeletedBy   *string    `db:"deleted_by" json:"deleted_by,omitempty"`
}
```

#### Interfaces

```go
type ProductHandler interface {
    Create(c Context) error
    Get(c Context) error
    List(c Context) error
    Update(c Context) error
    Delete(c Context) error
}

type ProductService interface {
    Create(ctx context.Context, req *CreateProductRequest, createdBy string) (*Product, error)
    Get(ctx context.Context, id string) (*Product, error)
    List(ctx context.Context) ([]Product, error)
    Update(ctx context.Context, req *UpdateProductRequest, updatedBy string) (*Product, error)
    Delete(ctx context.Context, id, deletedBy string) error
}

type ProductRepository interface {
    StartContext(ctx context.Context) context.Context
    DeferErrorContext(ctx context.Context, err error)
    Create(ctx context.Context, p *Product) error
    GetByID(ctx context.Context, id string) (*Product, error)
    List(ctx context.Context) ([]Product, error)
    Update(ctx context.Context, p *Product) error
    SoftDelete(ctx context.Context, id, deletedBy string) error
}
```

### User Domain

#### Entity

```go
type User struct {
    ID        string     `db:"id" json:"id" bson:"id"`
    Name      string     `db:"name" json:"name" bson:"name"`
    Email     string     `db:"email" json:"email" bson:"email"`
    CreatedAt time.Time  `db:"created_at" json:"created_at" bson:"created_at"`
    CreatedBy string     `db:"created_by" json:"created_by" bson:"created_by"`
    UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
    UpdatedBy *string    `db:"updated_by" json:"updated_by,omitempty"`
    DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
    DeletedBy *string    `db:"deleted_by" json:"deleted_by,omitempty"`
}
```

### Auth Domain

#### Session Entity

```go
type Session struct {
    ID        string     `db:"id" json:"id" bson:"id"`
    UserID    string     `db:"user_id" json:"user_id" bson:"user_id"`
    Token     string     `db:"token" json:"token" bson:"token"`
    ExpiresAt time.Time  `db:"expires_at" json:"expires_at" bson:"expires_at"`
    CreatedAt time.Time  `db:"created_at" json:"created_at" bson:"created_at"`
    UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty" bson:"updated_at,omitempty"`
    RevokedAt *time.Time `db:"revoked_at" json:"revoked_at,omitempty" bson:"revoked_at,omitempty"`
    UserAgent string     `db:"user_agent" json:"user_agent" bson:"user_agent"`
    IPAddress string     `db:"ip_address" json:"ip_address" bson:"ip_address"`
}
```

#### Credential Entity

```go
type Credential struct {
    ID           string     `db:"id" json:"id" bson:"id"`
    UserID       string     `db:"user_id" json:"user_id" bson:"user_id"`
    Username     string     `db:"username" json:"username" bson:"username"`
    Email        string     `db:"email" json:"email" bson:"email"`
    PasswordHash string     `db:"password_hash" json:"-" bson:"password_hash"`
    IsActive     bool       `db:"is_active" json:"is_active" bson:"is_active"`
    LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty" bson:"last_login_at,omitempty"`
    CreatedAt    time.Time  `db:"created_at" json:"created_at" bson:"created_at"`
    UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty" bson:"updated_at,omitempty"`
    DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
```

#### Token Claims

```go
type TokenClaims struct {
    UserID   string   `json:"user_id"`
    Username string   `json:"username"`
    Email    string   `json:"email"`
    Roles    []string `json:"roles,omitempty"`
}
```

#### Auth Interfaces

```go
type AuthHandler interface {
    Login(c Context) error
    Register(c Context) error
    Logout(c Context) error
    RefreshToken(c Context) error
    ValidateToken(c Context) error
    ChangePassword(c Context) error
    GetProfile(c Context) error
    GetSessions(c Context) error
    RevokeSession(c Context) error
    RevokeAllSessions(c Context) error
}

type AuthService interface {
    // Authentication
    Login(ctx context.Context, req *LoginRequest, userAgent, ipAddress string) (*LoginResponse, error)
    Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
    Logout(ctx context.Context, userID string, req *LogoutRequest) error
    RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
    ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error)

    // Password management
    ChangePassword(ctx context.Context, userID string, req *ChangePasswordRequest) error
    ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
    ConfirmResetPassword(ctx context.Context, req *ConfirmResetPasswordRequest) error

    // Session management
    GetSessions(ctx context.Context, userID string) (*SessionListResponse, error)
    RevokeSession(ctx context.Context, userID, sessionID string) error
    RevokeAllSessions(ctx context.Context, userID string) error

    // Token utilities
    GenerateAccessToken(claims *TokenClaims) (string, error)
    GenerateRefreshToken(userID string) (string, error)
    ParseToken(token string) (*TokenClaims, error)

    // Password utilities
    HashPassword(password string) (string, error)
    VerifyPassword(hashedPassword, password string) error
}

type AuthRepository interface {
    StartContext(ctx context.Context) context.Context
    DeferErrorContext(ctx context.Context, err error)

    // Credential operations
    CreateCredential(ctx context.Context, cred *Credential) error
    GetCredentialByUsername(ctx context.Context, username string) (*Credential, error)
    GetCredentialByEmail(ctx context.Context, email string) (*Credential, error)
    GetCredentialByUserID(ctx context.Context, userID string) (*Credential, error)
    UpdateCredential(ctx context.Context, cred *Credential) error
    UpdatePassword(ctx context.Context, userID, passwordHash string) error
    UpdateLastLogin(ctx context.Context, userID string) error

    // Session operations
    CreateSession(ctx context.Context, session *Session) error
    GetSessionByToken(ctx context.Context, token string) (*Session, error)
    GetSessionByID(ctx context.Context, id string) (*Session, error)
    GetSessionsByUserID(ctx context.Context, userID string) ([]Session, error)
    RevokeSession(ctx context.Context, sessionID string) error
    RevokeAllUserSessions(ctx context.Context, userID string) error
    DeleteExpiredSessions(ctx context.Context) error
}
```

---

## Infrastructure Layer

### PostgreSQL Connection

```go
// internal/infrastructure/db/sql/sql.go
func Open(dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    return db, nil
}
```

**Connection Pool Settings:**
- Max Open Connections: 25
- Max Idle Connections: 25
- Connection Max Lifetime: 5 minutes

### MongoDB Connection

```go
// internal/infrastructure/db/mongo/mongo.go
func OpenMongo(uri string) (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    return client, nil
}

func CloseMongo(client *mongo.Client) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return client.Disconnect(ctx)
}
```

---

## Transport Layer

### HTTP Framework Abstraction

The transport layer provides adapters that convert framework-specific contexts to domain contexts:

#### Echo Adapter

```go
func AdapterToEchoRoutes[T any](
    e *echo.Group,
    route *routes.Route,
    domainContext func(echo.Context) T,
) *echo.Group {
    handler := func(ctx echo.Context) error {
        return route.Handler.(func(T) error)(domainContext(ctx))
    }
    // Register route based on method
    switch route.Method {
    case echo.GET:
        e.GET(route.Path, handler)
    // ... other methods
    }
    return e
}
```

#### Gin Adapter

```go
func AdapterToGinRoutes[T any](
    r *gin.RouterGroup,
    route *routes.Route,
    domainContext func(*gin.Context) T,
) *gin.RouterGroup {
    handler := func(ctx *gin.Context) {
        _ = route.Handler.(func(T) error)(domainContext(ctx))
    }
    // Register route based on method
    switch route.Method {
    case "GET":
        r.GET(route.Path, handler)
    // ... other methods
    }
    return r
}
```

### Route Definition

```go
type Route struct {
    Method      string
    Path        string
    Handler     any
    Middlewares []any
    Flags       []string
}
```

---

## Authentication & Middleware

### Authentication Types

The system supports multiple authentication methods:

| Type | Description |
|------|-------------|
| `jwt` | JSON Web Token authentication via Bearer token |
| `session` | Session-based authentication via cookies |
| `basic` | HTTP Basic Authentication |
| `none` | No authentication required |

### Middleware Configuration

```go
type MiddlewareConfig struct {
    AuthType       AuthType   // jwt, session, basic, none
    SkipPaths      []string   // Paths to skip authentication
    SessionCookie  string     // Cookie name for session auth
    BasicAuthRealm string     // Realm for Basic Auth
}

func DefaultMiddlewareConfig() MiddlewareConfig {
    return MiddlewareConfig{
        AuthType:       AuthTypeJWT,
        SkipPaths:      []string{"/auth/login", "/auth/register", "/health"},
        SessionCookie:  "session_token",
        BasicAuthRealm: "Restricted",
    }
}
```

### Middleware Interface

```go
type Middleware interface {
    // Authenticate validates the request and sets auth context
    Authenticate() func(next func(Context) error) func(Context) error

    // RequireAuth ensures the request is authenticated
    RequireAuth() func(next func(Context) error) func(Context) error

    // OptionalAuth tries to authenticate but allows unauthenticated requests
    OptionalAuth() func(next func(Context) error) func(Context) error

    // RequireRoles ensures the authenticated user has specific roles
    RequireRoles(roles ...string) func(next func(Context) error) func(Context) error
}
```

### JWT Configuration

```go
type AuthConfig struct {
    JWTSecret            string        // JWT signing secret
    AccessTokenDuration  time.Duration // Default: 15 minutes
    RefreshTokenDuration time.Duration // Default: 7 days
    SessionDuration      time.Duration // Default: 24 hours
    BcryptCost           int           // Default: bcrypt.DefaultCost
}
```

### Protected Routes

Routes can be protected using middleware chains:

```go
{
    Method:      "GET",
    Path:        "/user",
    Handler:     userHandler.List,
    Middlewares: []any{authMiddleware.Authenticate(), authMiddleware.RequireAuth()},
    Flags:       []string{"protected"},
}
```

---

## Database & Migrations

### SQL Migrations (Goose)

Migrations are located in `internal/infrastructure/db/sql/migration/`

#### Users Table

```sql
-- 00001_create_users_table.sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  active BOOLEAN NOT NULL DEFAULT false,
  activation_token TEXT,
  reset_token TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

#### Products Table

```sql
-- 00002_create_products_table.sql
-- +goose Up
CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  created_by UUID,
  updated_at TIMESTAMP WITH TIME ZONE,
  updated_by UUID,
  deleted_at TIMESTAMP WITH TIME ZONE,
  deleted_by UUID
);

-- +goose Down
DROP TABLE IF EXISTS products;
```

#### Auth Tables

```sql
-- 00003_create_auth_tables.sql
-- Auth credentials table
CREATE TABLE IF NOT EXISTS auth_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Auth sessions table
CREATE TABLE IF NOT EXISTS auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    revoked_at TIMESTAMP WITH TIME ZONE,
    user_agent TEXT,
    ip_address VARCHAR(45)
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_auth_credentials_username ON auth_credentials(username) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_auth_credentials_email ON auth_credentials(email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_auth_credentials_user_id ON auth_credentials(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_auth_sessions_token ON auth_sessions(token) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id ON auth_sessions(user_id) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_auth_sessions_expires_at ON auth_sessions(expires_at);
```

### MongoDB Migrations

MongoDB migrations are JavaScript files executed via `mongosh` CLI.  
Located in `internal/infrastructure/db/mongo/migration/`

### Running Migrations

```bash
# SQL migrations
go run . migration sql up    # Apply all pending
go run . migration sql down  # Rollback last migration

# MongoDB migrations
go run . migration mongo up  # Apply all .js files in order
```

---

## API Reference

### Base URL

```
http://localhost:8080/v1
```

### Authentication

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "username": "johndoe",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "base64-encoded-token",
  "token_type": "Bearer",
  "expires_in": 900,
  "expires_at": "2025-12-01T00:15:00Z",
  "user": {
    "id": "uuid",
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

#### Register
```http
POST /auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Validation:**
- `username`: Required, min 3, max 50 characters
- `email`: Required, valid email format
- `password`: Required, min 8 characters
- `name`: Required

**Response:** `201 Created`
```json
{
  "user": {
    "id": "uuid",
    "username": "johndoe",
    "email": "john@example.com",
    "name": "John Doe"
  },
  "message": "Registration successful"
}
```

#### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "base64-encoded-token"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "base64-encoded-token",
  "token_type": "Bearer",
  "expires_in": 900,
  "expires_at": "2025-12-01T00:15:00Z"
}
```

#### Validate Token
```http
POST /auth/validate
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:** `200 OK`
```json
{
  "valid": true,
  "user": {
    "id": "uuid",
    "username": "johndoe",
    "email": "john@example.com"
  }
}
```

#### Logout (Protected)
```http
POST /auth/logout
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "refresh_token": "base64-encoded-token",
  "all_devices": false
}
```

**Response:** `200 OK`
```json
{
  "message": "Logged out successfully",
  "success": true
}
```

#### Get Profile (Protected)
```http
GET /auth/profile
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "username": "johndoe",
  "email": "john@example.com"
}
```

#### Change Password (Protected)
```http
PUT /auth/password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Response:** `200 OK`
```json
{
  "message": "Password changed successfully",
  "success": true
}
```

#### Get Sessions (Protected)
```http
GET /auth/sessions
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "sessions": [
    {
      "id": "session-uuid",
      "user_agent": "Mozilla/5.0...",
      "ip_address": "192.168.1.1",
      "created_at": "2025-12-01T00:00:00Z",
      "expires_at": "2025-12-08T00:00:00Z",
      "current": true
    }
  ]
}
```

#### Revoke Session (Protected)
```http
DELETE /auth/sessions/:id
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "message": "Session revoked successfully",
  "success": true
}
```

#### Revoke All Sessions (Protected)
```http
DELETE /auth/sessions
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "message": "All sessions revoked successfully",
  "success": true
}
```

### Products (Protected)

All product endpoints require authentication.

#### List Products
```http
GET /product
Authorization: Bearer <access_token>
```

**Response:**
```json
[
  {
    "id": "uuid",
    "name": "Product Name",
    "description": "Product description",
    "created_at": "2025-12-01T00:00:00Z",
    "created_by": "user-uuid"
  }
]
```

#### Create Product
```http
POST /product
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product description"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "name": "Product Name",
  "description": "Product description",
  "created_at": "2025-12-01T00:00:00Z",
  "created_by": "user-uuid"
}
```

### Users (Protected)

All user endpoints require authentication.

#### List Users
```http
GET /user
Authorization: Bearer <access_token>
```

#### Get User
```http
GET /user/:id
Authorization: Bearer <access_token>
```

#### Create User
```http
POST /user
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Validation:**
- `name`: Required
- `email`: Required, must be valid email format

#### Update User
```http
PUT /user/:id
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

#### Delete User
```http
DELETE /user/:id
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "status": "deleted"
}
```

---

## Development Guide

### Adding a New Module

1. **Create Domain Types**
   ```
   internal/domain/<module>/
   ├── model.<module>.go      # Entity definition
   ├── interface.<module>.go  # Handler/Service/Repo interfaces
   ├── request.<module>.go    # Request DTOs
   └── response.<module>.go   # Response DTOs
   ```

2. **Implement Repository**
   ```
   internal/modules/<module>/repository/
   ├── sql/repository_sql.<module>.go
   └── mongo/repository_mongo.<module>.go
   ```

3. **Implement Service**
   ```
   internal/modules/<module>/service/
   └── v1/service_v1.<module>.go
   ```

4. **Implement Handler**
   ```
   internal/modules/<module>/handler/
   └── v1/handler_v1.<module>.go
   ```

5. **Register in Container**
   Update `internal/app/core/container.go`

6. **Add Routes**
   Update `internal/app/http/routes.go`

7. **Add Feature Flags**
   Update `config/featureflags.yaml`

### Transaction Handling

Repositories support transactional operations:

```go
func (r *SQLRepository) StartContext(ctx context.Context) context.Context {
    tx := r.db.MustBeginTx(ctx, nil)
    return context.WithValue(ctx, driverName, tx)
}

func (r *SQLRepository) DeferErrorContext(ctx context.Context, err error) {
    tx := r.getTxFromContext(ctx)
    if tx != nil {
        if err != nil {
            tx.Rollback()
        } else {
            tx.Commit()
        }
    }
}
```

### Soft Delete Pattern

All entities support soft delete with audit fields:

```go
func (r *SQLRepository) SoftDelete(ctx context.Context, id, deletedBy string) error {
    now := time.Now().UTC()
    query := `UPDATE products SET deleted_at=$1, deleted_by=$2 WHERE id=$3`
    _, err := r.db.Exec(query, now, deletedBy, id)
    return err
}
```

### Common Response Models

```go
// Pagination Request
type PaginationRequest struct {
    Page  int `json:"page" binding:"required,min=1"`
    Limit int `json:"limit" binding:"required,min=5,max=100"`
}

func (p *PaginationRequest) Offset() int {
    return (p.Page - 1) * p.Limit
}

// Common Response
type CommonResponse struct {
    RequestID string `json:"requestId"`
}

// Paginated Response
type PaginatedResponse[T any] struct {
    CommonResponse
    Metadata PaginationMetadata `json:"metadata"`
    Data     []T                `json:"data"`
}

type PaginationMetadata struct {
    TotalItems int `json:"totalItems"`
    TotalPages int `json:"totalPages"`
    Page       int `json:"page"`
    Limit      int `json:"limit"`
}

func NewPaginatedResponse[T any](requestID string, data []T, totalItems int, meta PaginationMetadata) *PaginatedResponse[T] {
    totalPages := (totalItems + meta.Limit - 1) / meta.Limit

    return &PaginatedResponse[T]{
        CommonResponse: CommonResponse{
            RequestID: requestID,
        },
        Metadata: PaginationMetadata{
            TotalItems: totalItems,
            TotalPages: totalPages,
            Page:       meta.Page,
            Limit:      meta.Limit,
        },
        Data: data,
    }
}
```

---

## Roadmap

### Completed ✅
- [x] Architecture & Infrastructure setup
- [x] Product CRUD: HTTP Echo → v1 → SQL repository
- [x] SQL & MongoDB repository support with migrations
- [x] Gin framework integration
- [x] Request & Response models
- [x] User CRUD implementation
- [x] Authentication system: JWT, Basic Auth, Session-based (untested)
- [x] Middleware integration (untested)

### Planned 📋
- [ ] Unit Tests
- [ ] Redis integration (caching)
- [ ] Worker support: Asynq, RabbitMQ, Redpanda
- [ ] Storage support: S3-Compatible, GCS, MinIO, Local
- [ ] gRPC & Protocol Buffers support
- [ ] WebSocket integration

---

## Contributing

1. Follow the module structure outlined above
2. Use feature flags for new components
3. Implement both PostgreSQL and MongoDB repositories when applicable
4. Add migrations for database schema changes
5. Update this documentation for significant changes

---

## License

See [LICENSE](../LICENSE) file for details.
