import * as React from 'react';
import { Avatar, Button, CssBaseline, TextField, FormControlLabel, Checkbox, Link, Box, Typography, Container, Snackbar, Alert, Paper, InputAdornment, IconButton, Card, CardContent, Divider, Stack } from '@mui/material';
import { LockOutlined as LockOutlinedIcon, Visibility as VisibilityIcon, VisibilityOff as VisibilityOffIcon, Email as EmailIcon } from '@mui/icons-material';
import { useNavigate, Navigate } from 'react-router-dom';
import { login as apiLogin, getProfile } from '../api/auth.ts';
import { useAuth } from '../context/AuthContext.tsx';

export default function Login() {
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [showPassword, setShowPassword] = React.useState(false);
  const [rememberMe, setRememberMe] = React.useState(false);
  const [isLoading, setIsLoading] = React.useState(false);
  const [snackbar, setSnackbar] = React.useState({
    open: false,
    message: '',
    severity: 'success' as 'success' | 'error',
  });

  const { login, token } = useAuth();
  const navigate = useNavigate();

  // Redirect if already logged in
  if (token) {
    return <Navigate to="/dashboard" replace />;
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    // Mark recent login timestamp to prevent transient 401 redirects
    localStorage.setItem('recent_login_ts', String(Date.now()));

    try {
      const res = await apiLogin({ username, password });
      console.log('Login response full:', res);
      console.log('Login response data:', res.data);
      console.log('Login response data.token:', res.data?.token);

      // Call login from context to save token and decode user
      // apiLogin already returns response.data, so res is the LoginResponse object
      const token = res.data?.token;
      if (!token) {
        throw new Error('No token received from server');
      }

      login(token);
      console.log('Token saved to context:', token);

      // Prime permission cache by fetching profile before navigation
      try {
        const profile = await getProfile();
        console.log('Profile fetched to warm permissions:', profile?.data?.username);
      } catch (warmErr) {
        console.warn('Profile warm-up failed, proceeding anyway:', warmErr);
      }

      if (rememberMe) {
        localStorage.setItem('rememberMe', 'true');
        localStorage.setItem('lastUsername', username);
      }

      setSnackbar({
        open: true,
        message: 'Login successful! 🎉',
        severity: 'success',
      });

      // Use setTimeout to ensure state updates propagate before navigation
      setTimeout(() => {
        console.log('Navigating to dashboard...');
        navigate('/dashboard', { replace: true });
      }, 100);
    } catch (err: any) {
      console.error('Login error:', err);
      setSnackbar({
        open: true,
        message: err.response?.status === 401 ? 'Invalid username or password' : 'Login failed. Please try again.',
        severity: 'error',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleMouseDownPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
  };

  return (
    <Box
      sx={{
        width: '100%',
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        py: 3,
        flexGrow: 1,
      }}
    >
      <Container component="main" maxWidth="sm">
        <CssBaseline />
        <Paper
          elevation={10}
          sx={{
            p: 4,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            borderRadius: 3,
            background: 'rgba(255, 255, 255, 0.95)',
            backdropFilter: 'blur(10px)',
          }}
        >
          {/* Header */}
          <Avatar
            sx={{
              m: 2,
              width: 60,
              height: 60,
              bgcolor: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              boxShadow: '0 4px 15px rgba(102, 126, 234, 0.4)',
            }}
          >
            <LockOutlinedIcon sx={{ fontSize: '2rem' }} />
          </Avatar>

          <Typography component="h1" variant="h4" sx={{ fontWeight: 'bold', mb: 1, color: 'text.primary' }}>
            Welcome Back
          </Typography>

          <Typography variant="body2" color="textSecondary" sx={{ mb: 3, textAlign: 'center' }}>
            Sign in to your account to continue
          </Typography>

          {/* Login Form */}
          <Box component="form" onSubmit={handleSubmit} sx={{ width: '100%' }}>
            <TextField
              margin="normal"
              required
              fullWidth
              label="Username or Email"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              disabled={isLoading}
              autoComplete="username"
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <EmailIcon sx={{ color: 'action.active' }} />
                  </InputAdornment>
                ),
              }}
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: 2,
                },
              }}
            />

            <TextField
              margin="normal"
              required
              fullWidth
              label="Password"
              type={showPassword ? 'text' : 'password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={isLoading}
              autoComplete="current-password"
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <LockOutlinedIcon sx={{ color: 'action.active' }} />
                  </InputAdornment>
                ),
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton aria-label="toggle password visibility" onClick={handleClickShowPassword} onMouseDown={handleMouseDownPassword} edge="end" disabled={isLoading}>
                      {showPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: 2,
                },
              }}
            />

            <FormControlLabel control={<Checkbox color="primary" checked={rememberMe} onChange={(e) => setRememberMe(e.target.checked)} />} label="Remember me" sx={{ mt: 1, mb: 2 }} />

            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{
                mt: 2,
                mb: 2,
                py: 1.5,
                fontSize: '1rem',
                fontWeight: 600,
                borderRadius: 2,
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                boxShadow: '0 4px 15px rgba(102, 126, 234, 0.4)',
                '&:hover': {
                  boxShadow: '0 6px 20px rgba(102, 126, 234, 0.6)',
                },
              }}
              disabled={isLoading}
            >
              {isLoading ? 'Signing in...' : 'Sign In'}
            </Button>

            {/* Links */}
            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1} justifyContent="space-between" sx={{ mb: 2 }}>
              <Link
                href="#"
                variant="body2"
                sx={{
                  color: 'primary.main',
                  textDecoration: 'none',
                  '&:hover': { textDecoration: 'underline' },
                }}
              >
                Forgot password?
              </Link>
              <Link
                component="button"
                type="button"
                variant="body2"
                onClick={(e) => {
                  e.preventDefault();
                  navigate('/register');
                }}
                sx={{
                  color: 'primary.main',
                  textDecoration: 'none',
                  '&:hover': { textDecoration: 'underline' },
                }}
              >
                Sign Up
              </Link>
            </Stack>
          </Box>

          <Divider sx={{ my: 2, width: '100%' }} />

          {/* Demo Credentials Card */}
          <Card
            sx={{
              width: '100%',
              backgroundColor: '#f5f7fa',
              border: '1px solid',
              borderColor: 'divider',
            }}
          >
            <CardContent sx={{ py: 2 }}>
              <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 1 }}>
                Demo Credentials (for testing):
              </Typography>
              <Stack spacing={0.5}>
                <Typography variant="caption">
                  <strong>User:</strong> user / secret123
                </Typography>
                <Typography variant="caption">
                  <strong>Admin:</strong> admin / secret123
                </Typography>
              </Stack>
            </CardContent>
          </Card>
        </Paper>

        {/* Footer */}
        <Box sx={{ mt: 4, textAlign: 'center', color: 'rgba(255, 255, 255, 0.8)' }}>
          <Typography variant="body2">© 2024 URL Shortener. All rights reserved.</Typography>
        </Box>
      </Container>

      {/* Snackbar */}
      <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
        <Alert severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
}
