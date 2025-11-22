# EINFRA CRM Backend

Enterprise Infrastructure CRM Backend built with Go, Clean Architecture, and PostgreSQL.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Make (optional)

### Development Setup

1. **Clone the repository**
```bash
git clone https://github.com/unitechio/einfra-be.git
cd einfra-crm-be
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup environment**
```bash
# Copy development environment file
cp .env.development .env

# Or set APP_ENV to automatically load the right file
export APP_ENV=development
```

4. **Run database**
```bash
# Using Docker
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:14

# Create database
docker exec -it postgres psql -U postgres -c "CREATE DATABASE einfra_crm_dev;"
```

5. **Run the application**
```bash
# Development mode (auto-loads .env.development)
APP_ENV=development go run cmd/api/main.go

# Or using Make
make run-dev
```

The server will start on `http://localhost:8080`

## 📁 Project Structure

```
einfra-crm-be/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── adapter/
│   │   └── http/
│   │       ├── handler/            # HTTP handlers
│   │       └── middleware/         # HTTP middleware
│   ├── auth/
│   │   └── jwt.go                  # JWT service
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── domain/                     # Domain models & interfaces
│   │   ├── user.go
│   │   ├── auth.go
│   │   ├── audit.go
│   │   └── notification.go
│   ├── infrastructure/
│   │   ├── database/
│   │   │   └── postgres.go         # Database connection & migrations
│   │   └── repository/             # Repository implementations
│   │       ├── auth_repository.go
│   │       ├── user_repository.go
│   │       ├── session_repository.go
│   │       ├── audit_repository.go
│   │       └── notification_repository.go
│   └── usecase/                    # Business logic
│       ├── auth_usecase.go
│       └── user_usecase.go
├── .env.development                # Development environment
├── .env.production                 # Production environment
├── .env.example                    # Environment template
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🔧 Environment Configuration

The application uses environment-based configuration. Set `APP_ENV` to load the appropriate `.env` file:

```bash
# Development
APP_ENV=development go run cmd/api/main.go

# Production
APP_ENV=production go run cmd/api/main.go
```

### Environment Files

- `.env.development` - Development settings (debug mode, local database)
- `.env.production` - Production settings (optimized, secure)
- `.env.example` - Template for all available options

### Key Environment Variables

```bash
# Application
APP_NAME=EINFRA-CRM-BE
APP_ENV=development|production
APP_PORT=8080
APP_DEBUG=true|false

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=einfra_crm_dev
DB_SSLMODE=disable|require

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600
REFRESH_TOKEN_EXPIRY=604800
```

## 🗄️ Database Migrations

The application uses **Code-First** approach with GORM AutoMigrate:

```go
// Migrations run automatically on startup
database.AutoMigrate(db)
```

### Manual Migration

```bash
# Run migrations
make migrate

# Seed default data
make seed
```

### Default Data

On first run, the application seeds:
- **Admin Role** with all permissions
- **User Role** with basic permissions
- **Default Permissions** (user.*, role.*)

## 🔐 Authentication

The system uses JWT-based authentication with refresh tokens:

1. **Login** → Get access token + refresh token
2. **Access Token** → Short-lived (1 hour)
3. **Refresh Token** → Long-lived (7 days)
4. **Token Refresh** → Get new access token

### Security Features

- ✅ Password hashing (bcrypt)
- ✅ JWT access & refresh tokens
- ✅ Account locking after failed attempts
- ✅ Session management
- ✅ Login attempt tracking
- ✅ Email verification
- ✅ Password reset

## 📝 API Endpoints

### Health Check
```bash
GET /health
```

### Authentication
```bash
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
POST   /api/v1/auth/verify-email
```

### Users
```bash
GET    /api/v1/users
GET    /api/v1/users/:id
POST   /api/v1/users
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
```

## 🛠️ Development

### Run Tests
```bash
go test ./...
```

### Build
```bash
# Development
go build -o bin/api cmd/api/main.go

# Production
make build
```

### Docker
```bash
# Build image
docker build -t einfra-crm-be .

# Run container
docker run -p 8080:8080 --env-file .env.production einfra-crm-be
```

## 📦 Dependencies

- **Gin** - HTTP web framework
- **GORM** - ORM library
- **JWT** - JSON Web Tokens
- **Bcrypt** - Password hashing
- **Godotenv** - Environment management
- **UUID** - Unique identifiers

## 🏗️ Architecture

The project follows **Clean Architecture** principles:

1. **Domain Layer** - Business entities & interfaces
2. **Use Case Layer** - Business logic
3. **Infrastructure Layer** - External dependencies (DB, APIs)
4. **Adapter Layer** - HTTP handlers, middleware

### Key Principles

- ✅ Dependency Injection
- ✅ Interface-based design
- ✅ Separation of concerns
- ✅ Testability
- ✅ SOLID principles

## 📄 License

MIT License

## 👥 Contributors

- Your Team

---

**Happy Coding! 🚀**
