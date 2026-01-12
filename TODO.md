# TODO: Integrate Backend API with React + Material UI Frontend

## Overview

This TODO list guides the implementation of a React frontend using Material UI to consume the Go Gin backend APIs. Ensure all API calls are made correctly with proper authentication, error handling, and data validation. Base URL: `http://localhost:8080/api/v1` (adjust as needed).

## Prerequisites

- [ ] Set up React project: `npx create-react-app frontend`
- [ ] Install dependencies: `npm install @mui/material @emotion/react @emotion/styled axios react-router-dom`
- [ ] Set up Material UI theme provider in `App.js`
- [ ] Configure React Router for navigation

## API Configuration

- [ ] Create `src/api/config.js` with base URL and axios instance
- [ ] Add request interceptor to include JWT token in headers for protected routes
- [ ] Add response interceptor to handle 401 (token expired) by redirecting to login

## Authentication

- [ ] Create `src/api/auth.js` with functions:
  - `login(username, password)`: POST /auth/login, store token in localStorage
  - `register(data)`: POST /auth/register
  - `logout()`: Clear token from localStorage
- [ ] Create Login component (`src/components/Auth/Login.js`) with form validation
- [ ] Create Register component (`src/components/Auth/Register.js`)
- [ ] Implement token persistence: Check token on app load, redirect if invalid

## Shorten Link CRUD

- [ ] Create `src/api/shortenlink.js` with functions (include Authorization header):
  - `getAll()`: GET /shorten-links
  - `create(data)`: POST /shorten-links, data: {original_url}
  - `getById(id)`: GET /shorten-links/:id
  - `update(id, data)`: PATCH /shorten-links/:id, data: {original_url}
  - `delete(id)`: DELETE /shorten-links/:id
  - `redirect(code)`: GET /r/:code (public, no auth)
- [ ] Create ShortenLinkList component with Material UI Table
- [ ] Create ShortenLinkForm component for create/edit with validation
- [ ] Add delete confirmation dialog
- [ ] Handle API errors: Show Snackbar with error messages

## User Management (Admin Only)

- [ ] Create `src/api/user.js` with functions (include Authorization header):
  - `getAll()`: GET /users (requires user:read permission)
  - `create(data)`: POST /users, data: {name, username, email, phone?, is_active?, password}
  - `getById(id)`: GET /users/:id
  - `update(id, data)`: PATCH /users/:id, data: partial update
  - `delete(id)`: DELETE /users/:id
- [ ] Create UserList component with Material UI Table
- [ ] Create UserForm component for create/edit
- [ ] Add role-based UI: Show user management only if user has admin role (from JWT claims)

## UI/UX Enhancements

- [ ] Add loading states for all API calls (CircularProgress)
- [ ] Implement error handling: Display errors in Snackbar or Alert components
- [ ] Add success messages for create/update/delete operations
- [ ] Implement pagination for list views if backend supports it
- [ ] Add search/filter functionality for lists

## Testing

- [ ] Test login/register: Ensure token is stored and used in subsequent requests
- [ ] Test protected routes: Verify 401 handling when token is missing/expired
- [ ] Test CRUD operations: Create, read, update, delete for shortenlinks and users
- [ ] Test error scenarios: Invalid data, network errors, permission denied
- [ ] Test public redirect: Ensure /r/:code works without auth

## Deployment

- [ ] Build React app: `npm run build`
- [ ] Serve static files or integrate with backend
- [ ] Ensure CORS is configured in backend for frontend origin

## Notes

- All protected routes require `Authorization: Bearer <token>` header
- User permissions are checked server-side; frontend should hide UI based on user role
- Handle JWT expiration: Refresh token or re-login
- Validate form data client-side before API calls
