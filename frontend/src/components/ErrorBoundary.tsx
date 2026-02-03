import React from 'react';
import { Box, Button, Typography, Paper } from '@mui/material';

interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

export default class ErrorBoundary extends React.Component<React.PropsWithChildren<{}>, ErrorBoundaryState> {
  constructor(props: React.PropsWithChildren<{}>) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    // Log error to monitoring service if needed
    console.error('ErrorBoundary caught:', error, info);
  }

  handleReload = () => {
    this.setState({ hasError: false, error: undefined });
    window.location.reload();
  };

  render() {
    if (this.state.hasError) {
      return (
        <Box sx={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', p: 3 }}>
          <Paper sx={{ p: 4, borderRadius: 2, textAlign: 'center', maxWidth: 520 }}>
            <Typography variant="h5" sx={{ fontWeight: 'bold', mb: 1 }}>
              Something went wrong
            </Typography>
            <Typography variant="body2" color="textSecondary" sx={{ mb: 3 }}>
              An unexpected error occurred. Please try reloading the page.
            </Typography>
            <Button variant="contained" onClick={this.handleReload}>
              Reload Page
            </Button>
          </Paper>
        </Box>
      );
    }

    return this.props.children;
  }
}
