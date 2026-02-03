# Golang Gin MVC Codebase - Setup Guide

## ✅ Quick Start (Development)

### Prerequisites
- Go 1.20+
- Node.js 18+
- Redis server running on localhost:6379
- PostgreSQL (optional for database features)

### Backend Setup

1. **Navigate to backend directory**
   \`\`\`bash
   cd backend
   \`\`\`

2. **Create .env file** (copy from .env.example if needed)
   \`\`\`bash
   cp .env.example .env
   \`\`\`

3. **Run the application**
   \`\`\`bash
   go run cmd/app/main.go
   \`\`\`
   
   Backend will start on: **http://localhost:8080**

### Frontend Setup

1. **Navigate to frontend directory**
   \`\`\`bash
   cd frontend
   \`\`\`

2. **Install dependencies**
   \`\`\`bash
   npm install
   \`\`\`

3. **Create .env file**
   \`\`\`
   REACT_APP_API_URL=http://localhost:8080/api
   REACT_APP_ENV=development
   \`\`\`

4. **Start development server**
   \`\`\`bash
   npm start
   \`\`\`
   
   Frontend will start on: **http://localhost:3000** (or 3001 if 3000 is in use)

## ��� Docker Compose (Full Stack)

### Start all services
\`\`\`bash
docker compose up -d
\`\`\`

Services:
- **API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Redis**: localhost:6379
- **Swagger Docs**: http://localhost:8080/swagger/index.html

### Stop services
\`\`\`bash
docker compose down
\`\`\`

## ��� API Endpoints

Base URL: **http://localhost:8080/api**

### Authentication
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user
- `POST /auth/refresh` - Refresh JWT token
- `GET /auth/profile` - Get user profile

### Users
- `GET /users` - List all users (admin only)
- `GET /users/:id` - Get user by ID
- `POST /users` - Create user (admin only)
- `PATCH /users/:id` - Update user (admin only)
- `DELETE /users/:id` - Delete user (admin only)

### Short Links
- `GET /shorten-links` - List shortened links
- `GET /shorten-links/:id` - Get short link by ID
- `POST /shorten-links` - Create new shortened link
- `PUT /shortenlinks/:code` - Redirect to original URL

## ��� Authentication

The API uses JWT (JSON Web Tokens) for authentication.

1. Register or Login to get a token
2. Include token in request header:
   \`\`\`
   Authorization: Bearer <your_token>
   \`\`\`

## ��� Environment Variables

### Backend (.env in backend/)
- `APP_NAME` - Application name
- `APP_ENV` - Environment (development/production)
- `SERVER_PORT` - Server port (default: 8080)
- `JWT_SECRET` - Secret key for JWT
- `REDIS_HOST` - Redis host (default: 127.0.0.1)
- `REDIS_PORT` - Redis port (default: 6379)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password

### Frontend (.env in frontend/)
- `REACT_APP_API_URL` - Backend API URL
- `REACT_APP_ENV` - Environment (development/production)

## ��� Project Structure

\`\`\`
golang-gin-mvc-codebase/
├── backend/              # Golang Gin API
│   ├── app/             # Core application logic
│   ├── cmd/             # Entry points
│   ├── config/          # Configuration
│   ├── database/        # Migrations & setup
│   ├── helpers/         # Utilities
│   ├── .env             # Environment variables (local only)
│   └── .env.example     # Example configuration
│
├── frontend/            # React Web Application
│   ├── src/            # React components & logic
│   ├── public/         # Static files
│   ├── .env            # Environment variables (local only)
│   └── package.json    # Dependencies
│
├── docker-compose.yml  # Full stack configuration
└── .env                # Root environment variables
\`\`\`

## ��� Production Deployment

### Build Backend
\`\`\`bash
cd backend
go build -o app cmd/app/main.go
./app
\`\`\`

### Build Frontend
\`\`\`bash
cd frontend
npm run build
\`\`\`

## ���️ Troubleshooting

### Redis Connection Error
- Ensure Redis is running on localhost:6379
- Windows: Install Redis or use Docker: \`docker run -d -p 6379:6379 redis:latest\`
- Mac/Linux: \`redis-server\`

### Port Already in Use
- Backend: Change SERVER_PORT in .env
- Frontend: npm will ask to use another port

### Module Not Found (Frontend)
- Delete node_modules: \`rm -rf node_modules package-lock.json\`
- Reinstall: \`npm install\`

### Database Connection Error
- Ensure PostgreSQL is running
- Check DB_HOST, DB_USER, DB_PASSWORD in .env

## ��� Support

For API documentation, visit Swagger UI:
\`http://localhost:8080/swagger/index.html\`

