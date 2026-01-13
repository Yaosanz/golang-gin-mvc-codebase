# URL Shortener Frontend

A modern, responsive React TypeScript frontend application for a URL shortener service with user authentication and management features.

## Description

This frontend application provides a user-friendly interface for creating, managing, and sharing shortened URLs. It includes user authentication, dashboard, user management, and link management functionalities. The application communicates with a Golang Gin MVC backend API to perform CRUD operations on users and shortened links.

## Features

- **User Authentication**: Secure login and registration system
- **Dashboard**: Overview of user activities and statistics
- **User Management**: Admin interface for managing users (create, read, update, delete)
- **URL Shortening**: Create and manage shortened links
- **Link Management**: View, edit, and delete shortened links
- **Responsive Design**: Mobile-friendly interface using Material-UI
- **Real-time Updates**: Dynamic updates after CRUD operations
- **Error Handling**: Comprehensive error handling with user-friendly messages

## Tech Stack

- **Frontend Framework**: React 18
- **Language**: TypeScript
- **Styling**: Material-UI (MUI)
- **HTTP Client**: Axios
- **Routing**: React Router DOM
- **State Management**: React Context API
- **Build Tool**: Create React App
- **Backend**: Golang Gin MVC (separate repository)

## Prerequisites

Before running this application, make sure you have the following installed:

- Node.js (version 16 or higher)
- npm or yarn
- Git

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/url-shortener-frontend.git
   cd url-shortener-frontend
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm start
   ```

The application will be available at `http://localhost:3000`.

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

### Authentication

1. **Register**: Create a new account using the registration form
2. **Login**: Log in with your credentials
3. **Logout**: Use the logout button in the navigation

### Dashboard

After logging in, you'll be redirected to the dashboard where you can:

- View your shortened links
- Create new shortened links
- Manage your account

### User Management (Admin)

If you have admin privileges:

- Navigate to the Users section
- View all users
- Create new users
- Edit existing users
- Delete users

### Link Management

- **Create Link**: Use the "Create Short Link" form to shorten a URL
- **View Links**: See all your shortened links in the list
- **Edit Link**: Modify the original URL of a shortened link
- **Delete Link**: Remove a shortened link
- **Copy Link**: Copy the shortened URL to clipboard
- **Open Link**: Test the shortened link

## Project Structure

```
src/
├── api/                    # API service functions
│   ├── auth.ts            # Authentication API calls
│   ├── config.ts          # Axios configuration and interceptors
│   ├── shortenlink.ts     # Shorten link API calls
│   └── user.ts            # User management API calls
├── components/            # React components
│   ├── Dashboard.tsx      # Main dashboard component
│   ├── Landing.tsx        # Landing page
│   ├── Login.tsx          # Login form
│   ├── Register.tsx       # Registration form
│   ├── ShortenLinkCreate.tsx  # Create shorten link form
│   ├── ShortenLinkList.tsx    # List of shorten links
│   ├── UserForm.tsx       # User creation/editing form
│   └── UserList.tsx       # List of users
├── context/               # React Context for state management
│   └── AuthContext.tsx    # Authentication context
├── App.tsx                # Main App component
├── index.tsx              # Application entry point
├── index.css              # Global styles
└── reportWebVitals.ts     # Performance monitoring
public/
├── index.html             # HTML template
└── ...                    # Static assets
```

## API Integration

The frontend communicates with the backend through RESTful API endpoints:

### Authentication Endpoints

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### User Management Endpoints

- `GET /api/v1/users` - Get all users (paginated)
- `POST /api/v1/users` - Create new user
- `GET /api/v1/users/:id` - Get user by ID
- `PATCH /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Shorten Link Endpoints

- `GET /api/v1/shorten-links` - Get all shorten links
- `POST /api/v1/shorten-links` - Create new shorten link
- `GET /api/v1/shorten-links/:id` - Get shorten link by ID
- `PATCH /api/v1/shorten-links/:id` - Update shorten link
- `DELETE /api/v1/shorten-links/:id` - Delete shorten link
- `GET /api/v1/r/:code` - Redirect to original URL

## Scripts

- `npm start` - Start the development server
- `npm run build` - Build the application for production
- `npm test` - Run tests
- `npm run eject` - Eject from Create React App (irreversible)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature/your-feature-name`
5. Open a pull request

### Development Guidelines

- Use TypeScript for type safety
- Follow React best practices
- Use Material-UI components for consistent UI
- Write descriptive commit messages
- Test your changes thoroughly

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-username/url-shortener-frontend/issues) page
2. Create a new issue with detailed information
3. Contact the maintainers

## Roadmap

- [ ] Add link analytics and statistics
- [ ] Implement link expiration dates
- [ ] Add bulk link creation
- [ ] Integrate with social media sharing
- [ ] Add dark mode theme
- [ ] Implement PWA features

## Authors

- **Sandy Budiwirawan** - _Initial work_ - [sandybudiwirawan](https://github.com/sandybudiwirawan)

## Acknowledgments

- Thanks to the React and Material-UI communities for excellent documentation
- Inspired by popular URL shortener services
- Built with love for the open-source community
