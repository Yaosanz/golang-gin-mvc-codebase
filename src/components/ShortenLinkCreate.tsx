import React, { useState } from 'react';
import { Container, Typography, TextField, Button, Snackbar, Alert, Card, CardContent } from '@mui/material';
import { create } from '../api/shortenlink.ts';

const ShortenLinkCreate: React.FC = () => {
  const [url, setUrl] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });
  const [createdLink, setCreatedLink] = useState<{ short_code: string; original_url: string } | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const result = await create({ original_url: url });
      setCreatedLink(result);
      setSnackbar({ open: true, message: 'Link created successfully', severity: 'success' });
      setUrl('');
    } catch (err: any) {
      setSnackbar({ open: true, message: 'Failed to create link', severity: 'error' });
    }
  };

  return (
    <Container maxWidth="sm">
      <Typography variant="h4" gutterBottom>
        Create Shorten Link
      </Typography>
      <form onSubmit={handleSubmit}>
        <TextField fullWidth label="Original URL" value={url} onChange={(e) => setUrl(e.target.value)} margin="normal" required type="url" />
        <Button type="submit" variant="contained" fullWidth sx={{ mt: 2 }}>
          Create Link
        </Button>
      </form>

      {createdLink && (
        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6">Created Link</Typography>
            <Typography>Original: {createdLink.original_url}</Typography>
            <Typography>
              Short: {window.location.origin}/r/{createdLink.short_code}
            </Typography>
          </CardContent>
        </Card>
      )}

      <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })}>
        <Alert severity={snackbar.severity}>{snackbar.message}</Alert>
      </Snackbar>
    </Container>
  );
};

export default ShortenLinkCreate;
