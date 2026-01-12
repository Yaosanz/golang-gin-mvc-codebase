import * as React from 'react';
import { Avatar, Button, CssBaseline, TextField, FormControlLabel, Checkbox, Link, Grid, Box, Typography, Container, Snackbar, Alert } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { login as apiLogin } from '../api/auth.ts';
import { useAuth } from '../context/AuthContext.tsx';

const defaultTheme = createTheme();

export default function Login() {
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [snackbar, setSnackbar] = React.useState({
    open: false,
    message: '',
    severity: 'success' as 'success' | 'error',
  });

  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    try {
      const res = await apiLogin({ username, password });

      login(res.data.access_token);

      setSnackbar({
        open: true,
        message: 'Login successful! Redirecting...',
        severity: 'success',
      });

      setTimeout(() => {
        navigate('/dashboard');
      }, 2000);
    } catch (err: any) {
      setSnackbar({
        open: true,
        message: err.response?.status === 401 ? 'Invalid username or password' : 'Login failed. Please try again.',
        severity: 'error',
      });
    }
  };

  return (
    <ThemeProvider theme={defaultTheme}>
      <Container component="main" maxWidth="sm">
        <CssBaseline />
        <Box sx={{ mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component="h1" variant="h5">
            Sign in
          </Typography>

          <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
            <TextField margin="normal" required fullWidth label="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
            <TextField margin="normal" required fullWidth label="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />

            <FormControlLabel control={<Checkbox color="primary" />} label="Remember me" />

            <Button type="submit" fullWidth variant="contained" sx={{ mt: 3, mb: 2 }}>
              Sign In
            </Button>

            <Grid container justifyContent="space-between">
              <Link href="#" variant="body2">
                Forgot password?
              </Link>
              <Link component="button" variant="body2" onClick={() => navigate('/register')}>
                Don't have an account? Sign Up
              </Link>
            </Grid>
          </Box>
        </Box>
      </Container>

      <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
        <Alert severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </ThemeProvider>
  );
}
