import * as React from 'react';
import { Box, Typography, Button, TextField, List, ListItem, ListItemText, IconButton, Alert } from '@mui/material';
import { Add, Delete } from '@mui/icons-material';
import axios from 'axios';
import { useAuth } from '../context/AuthContext.tsx';
import { useNavigate } from 'react-router-dom';

interface ShortenLink {
  id: number;
  original_url: string;
  short_code: string;
  created_at: string;
}

export default function Dashboard() {
  const { token, logout } = useAuth();
  const navigate = useNavigate();
  const [links, setLinks] = React.useState<ShortenLink[]>([]);
  const [originalUrl, setOriginalUrl] = React.useState('');
  const [error, setError] = React.useState('');

  const fetchLinks = React.useCallback(async () => {
    try {
      const params: any = {};
      if (token) params.token = token;
      const response = await axios.get('/api/v1/shorten-links/', { params });
      setLinks(response.data.data);
      setError('');
    } catch (err: any) {
      setError('Failed to fetch links.');
      console.error(err);
    }
  }, [token]);

  React.useEffect(() => {
    fetchLinks();
  }, [fetchLinks]);

  const handleCreate = async () => {
    try {
      const params: any = {};
      if (token) params.token = token;
      await axios.post('/api/v1/shorten-links', { original_url: originalUrl }, { params });
      setOriginalUrl('');
      fetchLinks();
    } catch (err: any) {
      setError('Failed to create link.');
      console.error(err);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      const params: any = {};
      if (token) params.token = token;
      await axios.delete(`/api/v1/shorten-links/${id}`, { params });
      fetchLinks();
    } catch (err: any) {
      setError('Failed to delete link.');
      console.error(err);
    }
  };

  return (
    <Box sx={{ p: 4 }}>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      <Typography variant="body1" sx={{ mb: 3 }}>
        Welcome! Manage your shorten links here.
      </Typography>

      <Box sx={{ mb: 3 }}>
        <TextField label="Original URL" value={originalUrl} onChange={(e) => setOriginalUrl(e.target.value)} fullWidth sx={{ mb: 2 }} />
        <Button variant="contained" startIcon={<Add />} onClick={handleCreate}>
          Create Shorten Link
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Typography variant="h6" gutterBottom>
        Your Shorten Links
      </Typography>
      <List>
        {links.map((link) => (
          <ListItem
            key={link.id}
            secondaryAction={
              <IconButton onClick={() => handleDelete(link.id)}>
                <Delete />
              </IconButton>
            }
          >
            <ListItemText primary={link.original_url} secondary={`Short: http://localhost:8080/r/${link.short_code}`} />
          </ListItem>
        ))}
      </List>

      <Button
        variant="contained"
        color="error"
        onClick={() => {
          logout();
          navigate('/');
        }}
        sx={{ mt: 3 }}
      >
        Logout
      </Button>
    </Box>
  );
}
