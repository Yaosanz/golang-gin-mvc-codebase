import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Box, CircularProgress, Typography, Paper } from '@mui/material';
import { Error as ErrorIcon } from '@mui/icons-material';

const RedirectHandler: React.FC = () => {
  const { code } = useParams<{ code: string }>();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const redirect = async () => {
      if (!code) {
        setError('Invalid short code');
        return;
      }

      try {
        // Langsung redirect ke backend endpoint
        // Backend akan handle redirect ke original URL dengan HTTP 302
        const backendUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080';

        // Langsung redirect - browser akan follow 302 redirect dari backend
        window.location.href = `${backendUrl}/r/${code}`;
      } catch (err: any) {
        console.error('Redirect error:', err);
        setError('Failed to redirect. The link may be invalid or expired.');
      }
    };

    redirect();
  }, [code]);

  if (error) {
    return (
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          minHeight: '100vh',
          bgcolor: '#f5f5f5',
        }}
      >
        <Paper
          sx={{
            p: 4,
            maxWidth: 400,
            textAlign: 'center',
            borderRadius: 2,
          }}
        >
          <ErrorIcon sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
          <Typography variant="h6" gutterBottom>
            Redirect Failed
          </Typography>
          <Typography variant="body2" color="textSecondary">
            {error}
          </Typography>
        </Paper>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        bgcolor: '#f5f5f5',
      }}
    >
      <CircularProgress size={60} />
      <Typography variant="h6" sx={{ mt: 3, color: 'text.secondary' }}>
        Redirecting...
      </Typography>
    </Box>
  );
};

export default RedirectHandler;
