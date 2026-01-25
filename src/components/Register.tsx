import * as React from 'react';
import {
  Avatar,
  Button,
  CssBaseline,
  TextField,
  FormControlLabel,
  Checkbox,
  Link,
  Grid,
  Box,
  Typography,
  Container,
  Paper,
  Snackbar,
  Alert,
  InputAdornment,
  IconButton,
  Stepper,
  Step,
  StepLabel,
  Card,
  CardContent,
  Divider,
  Stack,
  CircularProgress,
} from '@mui/material';
import { LockOutlined as LockOutlinedIcon, Visibility as VisibilityIcon, VisibilityOff as VisibilityOffIcon, Email as EmailIcon, Person as PersonIcon, AssignmentInd as AssignmentIndIcon } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { register as apiRegister } from '../api/auth.ts';
import { Snackbar as SnackbarComponent, Alert as AlertComponent } from '@mui/material';

export default function SignUp() {
  const [username, setUsername] = React.useState('');
  const [name, setName] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [confirmPassword, setConfirmPassword] = React.useState('');
  const [email, setEmail] = React.useState('');
  const [phone, setPhone] = React.useState('');
  const [showPassword, setShowPassword] = React.useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = React.useState(false);
  const [agreeTerms, setAgreeTerms] = React.useState(false);
  const [errors, setErrors] = React.useState<{ [key: string]: string }>({});
  const [isLoading, setIsLoading] = React.useState(false);
  const [snackbar, setSnackbar] = React.useState({
    open: false,
    message: '',
    severity: 'success' as 'success' | 'error',
  });
  const navigate = useNavigate();

  const validate = () => {
    const newErrors: { [key: string]: string } = {};
    if (!username.trim()) newErrors.username = 'Username is required';
    else if (username.length < 3) newErrors.username = 'Username must be at least 3 characters';
    if (!name.trim()) newErrors.name = 'Full name is required';
    if (!email.trim()) newErrors.email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(email)) newErrors.email = 'Email is invalid';
    if (!password.trim()) newErrors.password = 'Password is required';
    else if (password.length < 6) newErrors.password = 'Password must be at least 6 characters';
    if (password !== confirmPassword) newErrors.confirmPassword = 'Passwords do not match';
    if (!agreeTerms) newErrors.terms = 'You must agree to the terms';
    return newErrors;
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const validationErrors = validate();
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }
    setErrors({});
    setIsLoading(true);

    try {
      await apiRegister({ username, name, email, password, phone });
      setSnackbar({
        open: true,
        message: 'Registration successful! 🎉 Redirecting to login...',
        severity: 'success',
      });
      setTimeout(() => navigate('/login'), 2000);
    } catch (err: any) {
      const errorMessage = err.response?.status === 409 ? 'Username or email already exists.' : err.response?.data?.message || 'Registration failed. Please try again.';
      setSnackbar({ open: true, message: errorMessage, severity: 'error' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleClickShowConfirmPassword = () => {
    setShowConfirmPassword(!showConfirmPassword);
  };

  const handleMouseDownPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        py: 3,
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

          <Typography component="h1" variant="h4" sx={{ fontWeight: 'bold', mb: 1 }}>
            Create Account
          </Typography>

          <Typography variant="body2" color="textSecondary" sx={{ mb: 3, textAlign: 'center' }}>
            Join us today and start shortening your URLs
          </Typography>

          {/* Registration Form */}
          <Box component="form" noValidate onSubmit={handleSubmit} sx={{ width: '100%' }}>
            <Grid container spacing={2}>
              {/* Full Name */}
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Full Name"
                  name="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  disabled={isLoading}
                  error={!!errors.name}
                  helperText={errors.name}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <PersonIcon sx={{ color: 'action.active' }} />
                      </InputAdornment>
                    ),
                  }}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      borderRadius: 2,
                    },
                  }}
                />
              </Grid>

              {/* Username */}
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Username"
                  name="username"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  disabled={isLoading}
                  error={!!errors.username}
                  helperText={errors.username}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <AssignmentIndIcon sx={{ color: 'action.active' }} />
                      </InputAdornment>
                    ),
                  }}
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      borderRadius: 2,
                    },
                  }}
                />
              </Grid>

              {/* Email */}
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Email Address"
                  name="email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={isLoading}
                  error={!!errors.email}
                  helperText={errors.email}
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
              </Grid>

              {/* Phone (Optional) */}
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Phone (Optional)"
                  name="phone"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  disabled={isLoading}
                  placeholder="+62 812 3456 7890"
                  sx={{
                    '& .MuiOutlinedInput-root': {
                      borderRadius: 2,
                    },
                  }}
                />
              </Grid>

              {/* Password */}
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={isLoading}
                  error={!!errors.password}
                  helperText={errors.password || 'At least 6 characters'}
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
              </Grid>

              {/* Confirm Password */}
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Confirm Password"
                  type={showConfirmPassword ? 'text' : 'password'}
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  disabled={isLoading}
                  error={!!errors.confirmPassword}
                  helperText={errors.confirmPassword}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <LockOutlinedIcon sx={{ color: 'action.active' }} />
                      </InputAdornment>
                    ),
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton aria-label="toggle password visibility" onClick={handleClickShowConfirmPassword} onMouseDown={handleMouseDownPassword} edge="end" disabled={isLoading}>
                          {showConfirmPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
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
              </Grid>

              {/* Terms & Conditions */}
              <Grid item xs={12}>
                <FormControlLabel
                  control={<Checkbox color="primary" checked={agreeTerms} onChange={(e) => setAgreeTerms(e.target.checked)} disabled={isLoading} />}
                  label={
                    <Typography variant="body2">
                      I agree to the{' '}
                      <Link href="#" sx={{ color: 'primary.main' }}>
                        Terms of Service
                      </Link>{' '}
                      and{' '}
                      <Link href="#" sx={{ color: 'primary.main' }}>
                        Privacy Policy
                      </Link>
                    </Typography>
                  }
                />
                {errors.terms && (
                  <Typography variant="caption" sx={{ color: 'error.main', display: 'block', mt: -1 }}>
                    {errors.terms}
                  </Typography>
                )}
              </Grid>
            </Grid>

            {/* Sign Up Button */}
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{
                mt: 3,
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
              {isLoading ? (
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <CircularProgress size={20} color="inherit" />
                  Creating account...
                </Box>
              ) : (
                'Create Account'
              )}
            </Button>

            {/* Sign In Link */}
            <Box sx={{ textAlign: 'center' }}>
              <Typography variant="body2" color="textSecondary">
                Already have an account?{' '}
                <Link
                  component="button"
                  type="button"
                  variant="body2"
                  onClick={(e) => {
                    e.preventDefault();
                    navigate('/login');
                  }}
                  sx={{ color: 'primary.main', fontWeight: 600 }}
                >
                  Sign in
                </Link>
              </Typography>
            </Box>
          </Box>

          <Divider sx={{ my: 3, width: '100%' }} />

          {/* Features Card */}
          <Card sx={{ width: '100%', backgroundColor: '#f5f7fa', border: '1px solid', borderColor: 'divider' }}>
            <CardContent sx={{ py: 2 }}>
              <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 1 }}>
                ✨ What you get:
              </Typography>
              <Stack spacing={0.5}>
                <Typography variant="caption">✓ Free URL shortening service</Typography>
                <Typography variant="caption">✓ Unlimited link creation</Typography>
                <Typography variant="caption">✓ Link analytics & tracking</Typography>
                <Typography variant="caption">✓ Custom short codes (Pro)</Typography>
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
      <Snackbar open={snackbar.open} autoHideDuration={6000} onClose={() => setSnackbar({ ...snackbar, open: false })} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
        <Alert severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
}
