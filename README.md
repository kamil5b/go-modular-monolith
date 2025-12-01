# Go Modular Monolith

A production-ready, modular monolithic application built with Go implementing clean architecture principles with pluggable components.

[![Go Version](https://img.shields.io/badge/Go-1.24.7-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- 🔄 **Switchable HTTP Frameworks** - Echo or Gin via configuration
- 🗄️ **Multiple Database Backends** - PostgreSQL and MongoDB support per module
- 🔐 **Complete Authentication** - JWT, Session-based, and Basic Auth
- 🛡️ **Middleware Support** - Authentication, authorization, and role-based access
- 📦 **Modular Architecture** - Independent versioning for handlers, services, and repositories
- 🎛️ **Feature Flags** - Enable/disable features through configuration
- 🔄 **Database Migrations** - Goose (SQL) and mongosh (MongoDB)

## Quick Start

### Prerequisites

- Go 1.24.7+
- PostgreSQL 14+
- MongoDB 6.0+ (optional)

### Installation

```bash
# Clone the repository
git clone https://github.com/kamil5b/go-modular-monolith.git
cd go-modular-monolith

# Install dependencies
go mod tidy

# Configure your database
# Edit config/config.yaml with your credentials

# Run the application
go run .
```

### CLI Commands

```bash
go run .                      # Run migrations and start server
go run . server               # Start server only
go run . migration sql up     # Apply SQL migrations
go run . migration sql down   # Rollback SQL migrations
go run . migration mongo up   # Apply MongoDB migrations
```

## Configuration

### Application Config (`config/config.yaml`)

```yaml
environment: development

app:
  server:
    port: "8080"
    grpc_port: "9090"

  database:
    sql:
      db_url: "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
    mongo:
      mongo_url: "mongodb://localhost:27017"
      mongo_db: "myapp_db"

  jwt:
    secret: "your-secret-key"
```

### Feature Flags (`config/featureflags.yaml`)

```yaml
http_handler: echo  # echo | gin

handler:
  authentication: v1  # v1 | disable
  product: v1
  user: v1

service:
  authentication: v1
  product: v1
  user: v1

repository:
  authentication: postgres  # postgres | mongo | disable
  product: postgres
  user: postgres
```

## Project Structure

```
go-modular-monolith/
├── cmd/bootstrap/          # Application bootstrapping
├── config/                 # Configuration files
├── docs/                   # Documentation
├── internal/
│   ├── app/               # Application core (DI, config, HTTP setup)
│   ├── domain/            # Domain models and interfaces
│   ├── infrastructure/    # Database connections, external services
│   ├── modules/           # Business modules (auth, product, user)
│   └── transports/        # HTTP framework adapters
└── pkg/                   # Shared utilities
```

## API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | User login |
| POST | `/auth/register` | User registration |
| POST | `/auth/refresh` | Refresh access token |
| POST | `/auth/validate` | Validate token |

### Authentication (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/logout` | User logout |
| GET | `/auth/profile` | Get user profile |
| PUT | `/auth/password` | Change password |
| GET | `/auth/sessions` | List active sessions |
| DELETE | `/auth/sessions/:id` | Revoke specific session |
| DELETE | `/auth/sessions` | Revoke all sessions |

### Products (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/product` | List products |
| POST | `/product` | Create product |

### Users (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/user` | List users |
| POST | `/user` | Create user |
| GET | `/user/:id` | Get user by ID |
| PUT | `/user/:id` | Update user |
| DELETE | `/user/:id` | Delete user |

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.24.7 |
| HTTP Frameworks | Echo v4, Gin v1 |
| SQL Database | PostgreSQL (sqlx) |
| NoSQL Database | MongoDB |
| Authentication | JWT (golang-jwt/jwt/v5) |
| Migrations | Goose, mongosh |
| Logging | Zerolog |

## Documentation

For detailed documentation, see [Technical Documentation](docs/TECHNICAL_DOCUMENTATION.md).

## Roadmap

- [x] Architecture & Infrastructure setup
- [x] Product & User CRUD
- [x] SQL & MongoDB repository support
- [x] Echo & Gin framework integration
- [x] Authentication (JWT, Session, Basic Auth)
- [x] Middleware integration
- [ ] Unit Tests
- [ ] Redis caching
- [ ] Worker support (Asynq, RabbitMQ)
- [ ] File storage (S3, GCS, MinIO)
- [ ] gRPC & Protocol Buffers
- [ ] WebSocket integration

## Contributing

1. Follow the module structure outlined in the documentation
2. Use feature flags for new components
3. Implement both PostgreSQL and MongoDB repositories when applicable
4. Add migrations for database schema changes
5. Update documentation for significant changes

## License

See [LICENSE](LICENSE) file for details.
