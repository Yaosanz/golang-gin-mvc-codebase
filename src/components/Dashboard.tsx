import React, { useState, useEffect, useCallback } from 'react';
import { Container, Typography, Button, List, ListItem, ListItemText, IconButton, Dialog, DialogTitle, DialogContent, DialogActions, TextField, Box } from '@mui/material';
import { Delete, Edit, Add } from '@mui/icons-material';
import axios from 'axios';
import { useAuth } from '../context/AuthContext.tsx';

interface ShortenLink {
  id: number;
  original_url: string;
  short_code: string;
  created_at: string;
}

const Dashboard: React.FC = () => {
  const [links, setLinks] = useState<ShortenLink[]>([]);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<ShortenLink | null>(null);
  const [originalUrl, setOriginalUrl] = useState('');
  const { token, logout } = useAuth();

  const fetchLinks = useCallback(async () => {
    try {
      const response = await axios.get('/api/v1/shorten-links', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setLinks(response.data);
    } catch (err) {
      console.error(err);
    }
  }, [token]);

  useEffect(() => {
    fetchLinks();
  }, [fetchLinks]);

  const handleCreate = () => {
    setEditing(null);
    setOriginalUrl('');
    setOpen(true);
  };

  const handleEdit = (link: ShortenLink) => {
    setEditing(link);
    setOriginalUrl(link.original_url);
    setOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await axios.delete(`/api/v1/shorten-links/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      fetchLinks();
    } catch (err) {
      console.error(err);
    }
  };

  const handleSave = async () => {
    try {
      if (editing) {
        await axios.patch(
          `/api/v1/shorten-links/${editing.id}`,
          { original_url: originalUrl },
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
      } else {
        await axios.post(
          '/api/v1/shorten-links',
          { original_url: originalUrl },
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
      }
      setOpen(false);
      fetchLinks();
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <Container>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mt: 4 }}>
        <Typography variant="h4">Shorten Links</Typography>
        <Box>
          <Button variant="contained" startIcon={<Add />} onClick={handleCreate}>
            Create
          </Button>
          <Button variant="outlined" onClick={logout} sx={{ ml: 2 }}>
            Logout
          </Button>
        </Box>
      </Box>
      <List>
        {links.map((link) => (
          <ListItem
            key={link.id}
            secondaryAction={
              <>
                <IconButton onClick={() => handleEdit(link)}>
                  <Edit />
                </IconButton>
                <IconButton onClick={() => handleDelete(link.id)}>
                  <Delete />
                </IconButton>
              </>
            }
          >
            <ListItemText primary={link.original_url} secondary={`Code: ${link.short_code}`} />
          </ListItem>
        ))}
      </List>
      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogTitle>{editing ? 'Edit' : 'Create'} Shorten Link</DialogTitle>
        <DialogContent>
          <TextField autoFocus margin="dense" label="Original URL" fullWidth value={originalUrl} onChange={(e) => setOriginalUrl(e.target.value)} />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
          <Button onClick={handleSave}>Save</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Dashboard;
