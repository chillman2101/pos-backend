# POS Backend API

Backend API untuk Point of Sale (POS) application dengan support offline sync.

## Tech Stack

- **Language**: Golang 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Validation**: go-playground/validator

## Features

- ✅ User authentication & authorization (JWT)
- ✅ Role-based access control (Admin, Manager, Cashier)
- ✅ Product management with categories
- ✅ Transaction processing
- ✅ Inventory tracking
- ✅ Offline sync support
- ✅ RESTful API architecture
- ✅ Auto migrations with GORM
- ✅ Standardized API responses

## Project Structure

```
pos-backend/
├── cmd/api/                    # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── database/               # Database connection & migrations
│   ├── domain/                 # Business entities & interfaces
│   ├── repository/             # Data access layer (GORM)
│   ├── service/                # Business logic layer
│   ├── handler/                # HTTP handlers
│   ├── middleware/             # Middlewares (auth, cors, etc)
│   ├── dto/                    # Data transfer objects
│   └── router/                 # Route definitions
└── pkg/                        # Reusable packages
    ├── response/               # Standardized responses
    ├── jwt/                    # JWT utilities
    └── utils/                  # Helper functions
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+

## Installation

1. Clone repository
```bash
git clone <repository-url>
cd pos-backend
```

2. Install dependencies
```bash
make deps
```

3. Setup environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Create database
```bash
createdb pos_db
```

5. Run application (migrations run automatically)
```bash
make run
```

Server will start on `http://localhost:8080`

## Available Commands

```bash
make run              # Run application
make build            # Build binary
make test             # Run tests
make docker-build     # Build Docker image
make docker-run       # Run Docker container
make clean            # Clean build artifacts
make deps             # Install dependencies
```

## API Endpoints

### Health Check
- `GET /health` - Check API health

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

### Protected Routes (Require JWT Token)
- `GET /api/v1/profile` - Get current user profile

### Coming Soon
- Users management
- Products management
- Categories management
- Transactions management
- Inventory tracking

## API Documentation

### Register User
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe",
  "role": "cashier"  // optional: admin, manager, cashier (default: cashier)
}

Response:
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "cashier",
      "is_active": true
    }
  }
}
```

### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "password123"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "cashier",
      "is_active": true
    }
  }
}
```

### Get Profile (Protected)
```bash
GET /api/v1/profile
Authorization: Bearer <token>

Response:
{
  "user_id": "uuid",
  "username": "john_doe",
  "role": "cashier"
}
```

## Database Schema

### Tables (Auto-migrated with GORM)
- **users** - User accounts with roles
- **categories** - Product categories
- **products** - Product catalog with inventory
- **transactions** - Sales transactions
- **transaction_items** - Transaction line items
- **inventory_movements** - Inventory tracking history

## Environment Variables

```env
SERVER_PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgres://user:password@localhost:5432/pos_db?sslmode=disable
JWT_SECRET=your-secret-key
```

## User Roles

- **admin** - Full access to all features
- **manager** - Can manage products, categories, and view reports
- **cashier** - Can process transactions and view products

## Development Roadmap

- [x] Base project structure
- [x] GORM integration
- [x] Auto migrations
- [x] Middleware setup (auth, cors, logger)
- [x] Response standardization
- [x] Auth module (login, register) ✅
- [x] JWT authentication ✅
- [x] User repository ✅
- [ ] User management endpoints
- [ ] Product management
- [ ] Category management
- [ ] Transaction processing
- [ ] Offline sync mechanism
- [ ] Inventory tracking
- [ ] Reporting & analytics
- [ ] Unit tests
- [ ] API documentation (Swagger)

## Testing with cURL

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@pos.com",
    "password": "admin123",
    "full_name": "Admin User",
    "role": "admin"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## License

MIT

## Author

Adit - Backend Engineer
```
