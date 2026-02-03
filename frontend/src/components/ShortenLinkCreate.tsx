import React, { useState } from 'react';
import { Container, Typography, TextField, Button, Snackbar, Alert, Card, CardContent, Box, Grid, Paper, CircularProgress, Chip, IconButton, Tooltip, List, ListItem, ListItemText, Divider } from '@mui/material';
import { ContentCopy as ContentCopyIcon, CheckCircle as CheckCircleIcon, OpenInNew as OpenInNewIcon, Link as LinkIcon, Timer as TimerIcon } from '@mui/icons-material';
import { create, ShortenLinkResponse } from '../api/shortenlink.ts';
import AppLayout from './AppLayout.tsx';

const ShortenLinkCreate: React.FC = () => {
  const [url, setUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });
  const [createdLink, setCreatedLink] = useState<ShortenLinkResponse | null>(null);
  const [recentLinks, setRecentLinks] = useState<ShortenLinkResponse[]>([]);
  const [urlError, setUrlError] = useState('');

  const validateUrl = (urlString: string) => {
    try {
      new URL(urlString);
      setUrlError('');
      return true;
    } catch (e) {
      setUrlError('Please enter a valid URL (e.g., https://example.com)');
      return false;
    }
  };

  const handleUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setUrl(value);
    if (value) {
      validateUrl(value);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const trimmedUrl = url.trim();
    if (!trimmedUrl) {
      setSnackbar({
        open: true,
        message: 'Please enter a URL',
        severity: 'error',
      });
      return;
    }

    if (!validateUrl(trimmedUrl)) {
      return;
    }

    try {
      setIsLoading(true);
      // Backend expects 'url' per Postman collection
      const result = await create({ url: trimmedUrl });
      setCreatedLink(result);
      setRecentLinks([result, ...recentLinks].slice(0, 5));
      setSnackbar({ open: true, message: 'Link shortened successfully! 🎉', severity: 'success' });
      setUrl('');
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.response?.data?.error || err.message || 'Failed to shorten link';
      console.error('Create link error:', err.response?.data);
      setSnackbar({
        open: true,
        message: errorMsg,
        severity: 'error',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setSnackbar({ open: true, message: 'Copied to clipboard!', severity: 'success' });
  };

  const getShareLink = () => {
    if (!createdLink) return '';
    const baseUrl = process.env.REACT_APP_API_URL?.replace('/api', '') || 'http://localhost:8080';
    return `${baseUrl}/r/${createdLink.code}`;
  };

  return (
    <AppLayout>
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Grid container spacing={4}>
          {/* Left Column - Form */}
          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 4, borderRadius: 2 }}>
              <Box sx={{ mb: 3 }}>
                <Typography variant="h5" sx={{ fontWeight: 'bold', mb: 1 }}>
                  Shorten Your URL
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  Convert long URLs into short, shareable links in seconds.
                </Typography>
              </Box>

              <form onSubmit={handleSubmit}>
                <TextField
                  fullWidth
                  label="Enter your long URL"
                  value={url}
                  onChange={handleUrlChange}
                  margin="normal"
                  required
                  disabled={isLoading}
                  placeholder="https://example.com/very/long/url/that/needs/shortening"
                  helperText={urlError || 'Paste a long URL to shorten'}
                  error={!!urlError}
                  InputProps={{
                    startAdornment: <LinkIcon sx={{ mr: 1, color: 'action.active' }} />,
                  }}
                />

                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  fullWidth
                  size="large"
                  sx={{ mt: 3, py: 1.5 }}
                  disabled={isLoading || !url || !!urlError}
                  startIcon={isLoading ? <CircularProgress size={20} /> : <CheckCircleIcon />}
                >
                  {isLoading ? 'Shortening...' : 'Shorten URL'}
                </Button>
              </form>

              {/* Created Link Result */}
              {createdLink && (
                <Box sx={{ mt: 3 }}>
                  <Divider sx={{ mb: 3 }} />
                  <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2 }}>
                    ✓ Link Created Successfully!
                  </Typography>

                  <Card
                    sx={{
                      mb: 2,
                      border: '2px solid',
                      borderColor: 'success.main',
                      background: 'linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(16, 185, 129, 0.05) 100%)',
                    }}
                  >
                    <CardContent>
                      <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                        Short Link
                      </Typography>
                      <Box
                        sx={{
                          display: 'flex',
                          gap: 1,
                          alignItems: 'center',
                          p: 1.5,
                          backgroundColor: 'background.paper',
                          borderRadius: 1,
                          mb: 2,
                        }}
                      >
                        <Typography
                          variant="body2"
                          sx={{
                            flex: 1,
                            fontFamily: 'monospace',
                            fontWeight: 'bold',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                          }}
                        >
                          {getShareLink()}
                        </Typography>
                        <Tooltip title="Copy link">
                          <IconButton size="small" onClick={() => copyToClipboard(getShareLink())}>
                            <ContentCopyIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Open link">
                          <IconButton size="small" onClick={() => window.open(getShareLink(), '_blank')}>
                            <OpenInNewIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                      </Box>

                      <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                        Original URL
                      </Typography>
                      <Typography
                        variant="body2"
                        sx={{
                          p: 1,
                          backgroundColor: 'background.paper',
                          borderRadius: 1,
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          whiteSpace: 'nowrap',
                        }}
                      >
                        {createdLink.url}
                      </Typography>

                      <Box sx={{ display: 'flex', gap: 1, mt: 2, pt: 2, borderTop: '1px solid', borderColor: 'divider' }}>
                        <Chip label={`Code: ${createdLink.code}`} size="small" variant="outlined" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }} />
                        <Chip label={`Created: ${new Date(createdLink.created_at).toLocaleDateString('id-ID')}`} size="small" icon={<TimerIcon />} />
                      </Box>
                    </CardContent>
                  </Card>

                  <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                    💡 Share your short link:
                  </Typography>
                  <Button fullWidth variant="outlined" onClick={() => copyToClipboard(getShareLink())} sx={{ mb: 1 }}>
                    📋 Copy Short Link
                  </Button>
                </Box>
              )}
            </Paper>
          </Grid>

          {/* Right Column - Info & Stats */}
          <Grid item xs={12} md={6}>
            {/* Features Section */}
            <Paper sx={{ p: 3, mb: 3, borderRadius: 2 }}>
              <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2 }}>
                ✨ Features
              </Typography>
              <List disablePadding>
                <ListItem sx={{ py: 1 }}>
                  <ListItemText primary="Fast & Reliable" secondary="Create short links in milliseconds with 99.9% uptime" />
                </ListItem>
                <ListItem sx={{ py: 1 }}>
                  <ListItemText primary="Easy to Share" secondary="Copy, paste, and share across all platforms instantly" />
                </ListItem>
                <ListItem sx={{ py: 1 }}>
                  <ListItemText primary="Analytics Ready" secondary="Track your link performance and click statistics" />
                </ListItem>
                <ListItem sx={{ py: 1 }}>
                  <ListItemText primary="Secure & Private" secondary="Your links are encrypted and protected from unauthorized access" />
                </ListItem>
              </List>
            </Paper>

            {/* Recent Activity */}
            {recentLinks.length > 0 && (
              <Paper sx={{ p: 3, borderRadius: 2 }}>
                <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2 }}>
                  📊 Recent Links
                </Typography>
                <List disablePadding>
                  {recentLinks.map((link, idx) => (
                    <React.Fragment key={link.id || idx}>
                      <ListItem
                        sx={{
                          py: 1.5,
                          '&:hover': { backgroundColor: 'action.hover' },
                          borderRadius: 1,
                        }}
                      >
                        <ListItemText
                          primary={
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                              <Chip label={link.code} size="small" variant="outlined" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }} />
                            </Box>
                          }
                          secondary={
                            <Typography
                              variant="caption"
                              sx={{
                                display: 'block',
                                overflow: 'hidden',
                                textOverflow: 'ellipsis',
                                whiteSpace: 'nowrap',
                                mt: 0.5,
                              }}
                            >
                              {link.url}
                            </Typography>
                          }
                        />
                      </ListItem>
                      {idx < recentLinks.length - 1 && <Divider />}
                    </React.Fragment>
                  ))}
                </List>
              </Paper>
            )}

            {/* Quick Tips */}
            <Alert severity="info" sx={{ mt: 3 }}>
              <Typography variant="body2" sx={{ fontWeight: 'bold', mb: 1 }}>
                💡 Pro Tips:
              </Typography>
              <Typography variant="caption" component="div">
                • Use descriptive short links for better branding
              </Typography>
              <Typography variant="caption" component="div">
                • Test your links before sharing widely
              </Typography>
              <Typography variant="caption" component="div">
                • Monitor your links for performance insights
              </Typography>
            </Alert>
          </Grid>
        </Grid>

        {/* Snackbar */}
        <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
          <Alert severity={snackbar.severity} sx={{ width: '100%' }}>
            {snackbar.message}
          </Alert>
        </Snackbar>
      </Container>
    </AppLayout>
  );
};

export default ShortenLinkCreate;
