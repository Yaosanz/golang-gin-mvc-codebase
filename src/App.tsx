import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Login from './components/Login.tsx';
import Register from './components/Register.tsx';
import Dashboard from './components/Dashboard.tsx';
import Profile from './components/Profile.tsx';
import ShortenLinksList from './components/ShortenLinksList.tsx';
import ShortenLinkCreate from './components/ShortenLinkCreate.tsx';
import UsersList from './components/UsersList.tsx';
import { AuthProvider, useAuth } from './context/AuthContext.tsx';

// Enhanced Material-UI theme with custom colors
const theme = createTheme({
  palette: {
    primary: {
      main: '#2563eb',
      light: '#60a5fa',
      dark: '#1e40af',
    },
    secondary: {
      main: '#8b5cf6',
      light: '#c4b5fd',
      dark: '#6d28d9',
    },
    success: {
      main: '#10b981',
      light: '#6ee7b7',
      dark: '#047857',
    },
    error: {
      main: '#ef4444',
      light: '#fca5a5',
      dark: '#dc2626',
    },
    warning: {
      main: '#f59e0b',
      light: '#fcd34d',
      dark: '#d97706',
    },
    info: {
      main: '#3b82f6',
      light: '#93c5fd',
      dark: '#1e3a8a',
    },
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 700,
      letterSpacing: '-0.5px',
    },
    h5: {
      fontWeight: 600,
    },
    button: {
      textTransform: 'none',
      fontWeight: 600,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: '8px',
          paddingTop: '10px',
          paddingBottom: '10px',
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
          },
        },
        contained: {
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: '12px',
          boxShadow: '0 2px 8px rgba(0, 0, 0, 0.08)',
          '&:hover': {
            boxShadow: '0 4px 16px rgba(0, 0, 0, 0.1)',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: '12px',
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          fontWeight: 600,
          fontSize: '0.85rem',
        },
      },
    },
  },
});

// Protected Route component
interface ProtectedRouteProps {
  element: React.ReactElement;
  requiredRole?: string;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ element, requiredRole }) => {
  const { token, user } = useAuth();

  if (!token) {
    return <Navigate to="/login" />;
  }

  if (requiredRole && !user?.roles?.some((r) => r.name === requiredRole || r === requiredRole)) {
    return <Navigate to="/dashboard" />;
  }

  return element;
};

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <Router>
          <Routes>
            {/* Public Routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/" element={<Navigate to="/dashboard" />} />

            {/* Protected Routes - Dashboard */}
            <Route path="/dashboard" element={<ProtectedRoute element={<Dashboard />} />} />

            {/* Protected Routes - Profile */}
            <Route path="/profile" element={<ProtectedRoute element={<Profile />} />} />

            {/* Protected Routes - Shortened Links */}
            <Route path="/links" element={<ProtectedRoute element={<ShortenLinksList />} />} />
            <Route path="/links/create" element={<ProtectedRoute element={<ShortenLinkCreate />} />} />

            {/* Protected Routes - Admin Only - User Management */}
            <Route path="/users" element={<ProtectedRoute element={<UsersList />} requiredRole="admin" />} />

            {/* Catch all */}
            <Route path="*" element={<Navigate to="/dashboard" />} />
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
