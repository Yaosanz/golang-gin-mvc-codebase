# Go Starter App - MVC Architecture

A robust Go web application built with Gin framework, featuring MVC architecture, PostgreSQL database, JWT authentication, role-based permissions, and comprehensive seeding system.

## Features

- **MVC Architecture**: Clean separation of concerns with Models, Views (JSON responses), and Controllers
- **PostgreSQL Database**: Full database support with migrations and seeders
- **JWT Authentication**: Secure token-based authentication system
- **Role-Based Access Control**: Permission system with roles and permissions
- **RESTful API**: Well-structured REST endpoints
- **Database Seeding**: Automated data seeding for development and testing
- **Middleware Support**: CORS, authentication, authorization middleware
- **Swagger Documentation**: Auto-generated API documentation
- **Docker Support**: Containerized deployment ready

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository_url>
cd golang-gin-mvc-codebase
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables (create `.env` file based on `.env.example`):
```bash
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=corpu

# JWT configuration
JWT_SECRET=your_jwt_secret

# Other configurations...
```

## Database Setup

### Create Database
```bash
createdb -U postgres corpu
```

### Run Migrations
```bash
# Option 1: Using Go command
go run cmd/migration/main.go up

# Option 2: Using individual SQL files (if needed)
psql -U postgres -d corpu -f database/migrations/000002_users.up.sql
psql -U postgres -d corpu -f database/migrations/000003_create_shortenlink.up.sql
psql -U postgres -d corpu -f database/migrations/000004_permissions.up.sql
psql -U postgres -d corpu -f database/migrations/000005_role_permissions.up.sql
psql -U postgres -d corpu -f database/migrations/000006_create_roles_table.up.sql
psql -U postgres -d corpu -f database/migrations/000007_add_role_id_to_users.up.sql
psql -U postgres -d corpu -f database/migrations/000008_fix_role_permissions_fk.up.sql
```

### Run Seeders
```bash
# Run all seeders
go run ./cmd/seeder/main.go run:all

# Or run individual seeders
go run ./cmd/seeder/main.go run:one permission_seeder
go run ./cmd/seeder/main.go run:one role_seeder
go run ./cmd/seeder/main.go run:one role_permissions_seeder
go run ./cmd/seeder/main.go run:one user_seeder
```

### Database Rollback (if needed)
```bash
# Rollback migrations
go run cmd/migration/main.go down

# Or rollback individual migrations
psql -U postgres -d corpu -f database/migrations/000003_create_shortenlink.down.sql
psql -U postgres -d corpu -f database/migrations/000002_users.down.sql
psql -U postgres -d corpu -f database/migrations/000006_create_roles_table.down.sql
psql -U postgres -d corpu -f database/migrations/000007_add_role_id_to_users.down.sql
psql -U postgres -d corpu -f database/migrations/000008_fix_role_permissions_fk.down.sql
```

## Running the Application

### Development Mode
```bash
go run cmd/app/main.go
```

The server will start on `http://localhost:8080`

### Build and Run
```bash
go build -o app cmd/app/main.go
./app
```

### Clean Build
```bash
go clean -cache
go mod tidy
go run cmd/app/main.go
```

## API Documentation

### Swagger UI
Access the API documentation at: `http://localhost:8080/swagger/index.html`

### Authentication
The API uses JWT tokens for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Default Users
After running seeders, you can use these default credentials:

- **Admin User**: `admin` / `secret123`
- **Regular User**: `sandy` / `secret123`

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### Users Management (Admin only)
- `GET /api/v1/users` - List users
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create user
- `PATCH /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### URL Shortener
- `GET /api/v1/shorten-links` - List shortened links
- `POST /api/v1/shorten-links` - Create shortened link
- `GET /api/v1/shorten-links/{id}` - Get shortened link
- `PATCH /api/v1/shorten-links/{id}` - Update shortened link
- `DELETE /api/v1/shorten-links/{id}` - Delete shortened link
- `GET /api/v1/r/{code}` - Redirect to original URL

## Project Structure

```
├── app/
│   ├── http/
│   │   ├── controllers/     # HTTP controllers
│   │   ├── dto/            # Data transfer objects
│   │   ├── middleware/     # HTTP middleware
│   │   ├── routes/         # Route definitions
│   │   └── utils/          # HTTP utilities
│   ├── models/             # Database models
│   ├── repositories/       # Data access layer
│   └── services/           # Business logic layer
├── cmd/                    # Application entry points
│   ├── app/               # Main application
│   ├── migration/         # Database migration CLI
│   └── seeder/            # Database seeder CLI
├── config/                 # Configuration files
├── database/
│   ├── migrations/        # Database migration files
│   └── seeders/           # Database seeders
├── helpers/                # Utility functions
├── interfaces/             # Interface definitions
├── pkg/                    # Shared packages
└── scripts/                # Utility scripts
```

## Database Seeding

### Seeder CLI Usage
```bash
# List all available seeders
go run ./cmd/seeder/main.go list

# Run all seeders
go run ./cmd/seeder/main.go run:all

# Run specific seeder
go run ./cmd/seeder/main.go run:one <seeder_name>

# Run with transaction (rollback on failure)
go run ./cmd/seeder/main.go run:all-tx
```

### Available Seeders
- `permission_seeder` - Creates user permissions
- `role_seeder` - Creates user roles
- `role_permissions_seeder` - Assigns permissions to roles
- `user_seeder` - Creates default users

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./app/services/...
```

## Docker Support

Build and run with Docker:
```bash
# Build image
docker build -t go-starter-app .

# Run container
docker run -p 8080:8080 go-starter-app
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go naming conventions
- Write tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting PR

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [JWT](https://jwt.io/)

## Support

For support, please open an issue in the GitHub repository or contact the development team.
