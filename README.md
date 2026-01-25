# URL Shortener - Full Stack Application

A modern, production-ready full-stack application for URL shortening with comprehensive authentication, role-based access control, and Redis caching. Built with Golang Gin MVC backend and React TypeScript Material-UI frontend.

## Description

This application provides a complete URL shortening service with enterprise-grade features including multi-role JWT authentication, RBAC (Role-Based Access Control), Redis caching for high performance, transaction management with rollback capabilities, and a beautiful Material-UI frontend. The system supports both regular users and administrators with distinct capabilities and interfaces.

## Features

### Backend Features

- **JWT Multi-Role Authentication**: Secure token-based authentication with support for multiple roles per user
- **RBAC (Role-Based Access Control)**: Fine-grained permission system for admin and user roles
- **Redis Caching**: High-performance caching with TTL management (5min-48hrs) for optimized response times (1-3ms)
- **Transaction Management**: Atomic operations with automatic rollback on failures
- **RESTful API**: Clean API design following REST principles
- **Middleware Protection**: JWT validation and role-based route protection
- **Database Migrations**: Automated schema management with GORM

### Frontend Features

- **Modern React UI**: Built with React 18 and TypeScript for type safety
- **Material-UI Design**: Beautiful, responsive components with custom theming
- **Protected Routes**: Client-side route protection with JWT expiration checks
- **Role-Based UI**: Dynamic navigation and features based on user roles
- **Error Boundary**: Global error handling for graceful error recovery
- **Loading States**: Skeleton loaders and loading indicators for better UX
- **Real-time Validation**: Form validation with instant feedback
- **Responsive Design**: Mobile-first design that works on all screen sizes
- **DataGrid Tables**: Advanced tables with sorting, filtering, and pagination
- **Profile Management**: Self-service profile updates for users
- **Admin Dashboard**: Comprehensive user management for administrators
- **Settings Panel**: Configurable user preferences and notifications

## Tech Stack

### Backend

- **Framework**: Golang with Gin (MVC architecture)
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis 7+ for high-performance caching
- **Authentication**: JWT (JSON Web Tokens) with multi-role support
- **Validation**: Go Validator v10
- **API Documentation**: RESTful design

### Frontend

- **Framework**: React 18 with TypeScript
- **UI Library**: Material-UI (MUI) v5
- **Data Grid**: MUI X-DataGrid
- **HTTP Client**: Axios with interceptors
- **Routing**: React Router DOM v6
- **State Management**: React Context API
- **Build Tool**: Create React App
- **Icons**: Material-UI Icons

## Prerequisites

Before running this application, make sure you have:

### Backend Requirements

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Redis 7 or higher
- Git

### Frontend Requirements

- Node.js 16 or higher
- npm or yarn
- Git

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/golang-gin-mvc-codebase.git
cd golang-gin-mvc-codebase
```

### 2. Backend Setup

1. Install Go dependencies:

   ```bash
   go mod download
   ```

2. Set up PostgreSQL database:

   ```sql
   CREATE DATABASE url_shortener;
   ```

3. Set up Redis (make sure Redis server is running):

   ```bash
   redis-server
   ```

4. Configure environment variables (create `.env` file in root):

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=url_shortener

   REDIS_HOST=localhost
   REDIS_PORT=6379

   JWT_SECRET=your-secret-key-here
   PORT=8080
   ```

5. Run database migrations:

   ```bash
   go run main.go migrate
   ```

6. Start the backend server:
   ```bash
   go run main.go
   ```

The backend API will be available at `http://localhost:8080`

### 3. Frontend Setup

1. Navigate to frontend directory (if separated) or install in root:

   ```bash
   npm install
   ```

2. Configure API endpoint (create `.env` in root or use default proxy):

   ```env
   REACT_APP_API_URL=http://localhost:8080/api/v1
   ```

3. Start the development server:
   ```bash
   npm start
   ```

The frontend application will be available at `http://localhost:3000`

## Configuration

### Environment Variables

Create a `.env` file in the root directory for production configuration:

```env
REACT_APP_API_URL=https://your-backend-api-url.com/api/v1
```

If `REACT_APP_API_URL` is not set, the application will use the proxy configuration for development.

### Proxy Configuration

For development, the application uses a proxy to avoid CORS issues. The proxy is configured in `package.json`:

```json
"proxy": "http://localhost:8080"
```

This proxies API requests from `/api/v1/*` to `http://localhost:8080/api/v1/*`.

## Usage

### Getting Started

1. **Register a New Account**
   - Navigate to `/register`
   - Fill in: username, name, email, phone, password
   - Click "Sign Up"
   - Default role: "user"

2. **Login**
   - Navigate to `/login`
   - Enter credentials
   - Demo accounts for testing:
     - Admin: `admin` / `admin123`
     - User: `user` / `user123`

3. **Access Dashboard**
   - After login, you'll see the main dashboard
   - View statistics, recent links, and quick actions
   - Role badges display your current permissions

### User Features

#### Dashboard (`/dashboard`)

- View total links count
- See monthly statistics
- Quick access to create new links
- Recent links overview
- Profile summary with role display

#### Create Short Links (`/links/create`)

- Enter long URL to shorten
- Automatic code generation
- URL validation
- Recent links sidebar
- Copy short URL to clipboard
- Preview and test links

#### Manage Links (`/links`)

- View all your shortened links in DataGrid
- Edit original URLs
- Delete links
- Copy short codes
- Open links in new tab
- Filter and search capabilities
- Pagination support

#### Profile Management (`/profile`)

- View and edit personal information
- Update name, email, phone
- View account creation date
- See assigned roles
- Change password (coming soon)

#### Settings (`/settings`)

- Configure notifications (email, push)
- Appearance preferences
- Language selection (EN, ID, ES, FR)
- Timezone configuration
- Security settings
- Auto-save preferences

### Admin Features (Requires "admin" Role)

#### User Management (`/users`)

- View all users in the system
- Create new users with roles
- Edit user information
- Delete users
- View user statistics
- Search and filter users
- Bulk operations support

#### Admin Dashboard Extras

- Total users count
- System-wide statistics
- All users' links visibility

### API Features Utilized

#### Authentication & Authorization

- JWT token stored in localStorage
- Automatic token injection in API calls
- Token expiration validation
- 401 redirect to login
- Role-based route protection

#### Caching & Performance

- Redis cache for link lookups (1-3ms response)
- Cache invalidation on CRUD operations
- Refresh button to clear cache
- TTL management (5min-48hrs)

#### Transaction Safety

- Atomic create/update/delete operations
- Automatic rollback on errors
- Duplicate detection
- Constraint validation

## Project Structure

```
golang-gin-mvc-codebase/
├── backend/                    # Backend Go application (if separated)
│   ├── controllers/           # HTTP request handlers
│   ├── models/                # Data models and database schemas
│   ├── routes/                # API route definitions
│   ├── middleware/            # JWT auth, RBAC, logging
│   ├── services/              # Business logic layer
│   ├── config/                # Configuration and database setup
│   └── main.go                # Application entry point
│
├── src/                       # Frontend React application
│   ├── api/                   # API service layer
│   │   ├── auth.ts           # Authentication API (login, register, profile)
│   │   ├── config.ts         # Axios config with JWT interceptor
│   │   ├── shortenlink.ts    # Link CRUD operations
│   │   └── user.ts           # User management (admin)
│   │
│   ├── components/            # React components
│   │   ├── AppLayout.tsx     # Main layout with navigation
│   │   ├── Dashboard.tsx     # Dashboard with stats
│   │   ├── Login.tsx         # Login form with validation
│   │   ├── Register.tsx      # Registration form
│   │   ├── Profile.tsx       # User profile management
│   │   ├── ShortenLinkCreate.tsx  # Create short links
│   │   ├── ShortenLinksList.tsx   # DataGrid of links
│   │   ├── UsersList.tsx     # Admin user management
│   │   ├── Settings.tsx      # User settings panel
│   │   └── ErrorBoundary.tsx # Global error handler
│   │
│   ├── context/               # React Context
│   │   └── AuthContext.tsx   # Auth state & JWT decode
│   │
│   ├── App.tsx                # Routes & theme config
│   ├── index.tsx              # Application entry
│   └── index.css              # Global styles
│
├── public/                    # Static assets
│   └── index.html            # HTML template
│
├── build/                     # Production build output
├── package.json              # npm dependencies
├── tsconfig.json             # TypeScript configuration
├── .env                      # Environment variables
└── README.md                 # This file
```

## API Integration

The frontend seamlessly integrates with all backend endpoints:

### Authentication Endpoints

| Method | Endpoint                | Description                 | Auth Required |
| ------ | ----------------------- | --------------------------- | ------------- |
| POST   | `/api/v1/auth/register` | Register new user           | No            |
| POST   | `/api/v1/auth/login`    | Login and get JWT token     | No            |
| GET    | `/api/v1/auth/profile`  | Get current user profile    | Yes           |
| PATCH  | `/api/v1/auth/profile`  | Update current user profile | Yes           |

### User Management Endpoints (Admin Only)

| Method | Endpoint            | Description               | Role Required |
| ------ | ------------------- | ------------------------- | ------------- |
| GET    | `/api/v1/users`     | Get all users (paginated) | admin         |
| POST   | `/api/v1/users`     | Create new user           | admin         |
| GET    | `/api/v1/users/:id` | Get user by ID            | admin         |
| PATCH  | `/api/v1/users/:id` | Update user               | admin         |
| DELETE | `/api/v1/users/:id` | Delete user               | admin         |

### Shorten Link Endpoints

| Method | Endpoint                     | Description               | Auth Required |
| ------ | ---------------------------- | ------------------------- | ------------- |
| GET    | `/api/v1/shorten-links`      | Get all user's links      | Yes           |
| POST   | `/api/v1/shorten-links`      | Create new short link     | Yes           |
| GET    | `/api/v1/shorten-links/:id`  | Get link by ID            | Yes           |
| PATCH  | `/api/v1/shorten-links/:id`  | Update link               | Yes           |
| DELETE | `/api/v1/shorten-links/:id`  | Delete link               | Yes           |
| GET    | `/api/v1/shortenlinks/:code` | Get link by code (cached) | No            |
| GET    | `/r/:code`                   | Redirect to original URL  | No            |

### Response Format

All API responses follow this structure:

```typescript
{
  success: boolean;
  message: string;
  data: T; // Generic data type
}
```

### Error Handling

- `400 Bad Request`: Validation errors
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate resource
- `500 Internal Server Error`: Server-side errors

## Scripts

### Backend Scripts

```bash
go run main.go          # Start backend server
go run main.go migrate  # Run database migrations
go test ./...           # Run backend tests
go build                # Build production binary
```

### Frontend Scripts

```bash
npm start              # Start development server (http://localhost:3000)
npm run build          # Build for production (outputs to build/)
npm test               # Run tests
npm run eject          # Eject from CRA (irreversible)
```

### Production Deployment

```bash
# Backend
go build -o app main.go
./app

# Frontend
npm run build
# Serve build/ directory with nginx or static server
npx serve -s build
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature/your-feature-name`
5. Open a pull request

### Development Guidelines

#### Backend

- Follow MVC architecture pattern
- Use GORM for database operations
- Implement proper error handling
- Use middleware for cross-cutting concerns
- Write clear API documentation
- Follow Go naming conventions

#### Frontend

- Use TypeScript for type safety
- Follow React best practices and hooks patterns
- Use Material-UI components for consistency
- Implement proper error boundaries
- Write descriptive commit messages
- Keep components focused and reusable
- Use Context API for global state
- Implement loading states for async operations

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-username/url-shortener-frontend/issues) page
2. Create a new issue with detailed information
3. Contact the maintainers

## Roadmap

### Completed ✅

- [x] JWT Multi-Role Authentication
- [x] RBAC with admin/user roles
- [x] Redis caching with TTL
- [x] Transaction management with rollback
- [x] Full CRUD for users and links
- [x] Material-UI frontend with DataGrid
- [x] Protected routes and role-based UI
- [x] Profile management
- [x] Error boundary and global error handling
- [x] Loading states and skeletons
- [x] Settings page

### Planned 🚀

- [ ] Link analytics and click tracking
- [ ] Custom short codes
- [ ] Link expiration dates
- [ ] Bulk link creation
- [ ] QR code generation
- [ ] Social media sharing integration
- [ ] Dark mode theme
- [ ] PWA features (offline support)
- [ ] Email notifications
- [ ] Two-factor authentication
- [ ] API rate limiting
- [ ] Link categories/tags
- [ ] Export links to CSV
- [ ] Password change functionality
- [ ] Unit and integration tests

## Authors

- **Sandy Budi Wirawan** - _Initial Work_ - [sandybudiwirawan](https://github.com/Yaosanz)

## Acknowledgments

- Thanks to the React and Material-UI communities for excellent documentation
- Inspired by popular URL shortener services
- Built with love for the open-source community
