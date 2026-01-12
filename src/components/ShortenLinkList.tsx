import React, { useEffect, useState, useCallback } from 'react';
import { Container, Typography, List, ListItem, ListItemText, IconButton, Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Snackbar, Alert, Pagination, Box } from '@mui/material';
import { Edit, Delete, ContentCopy, Search } from '@mui/icons-material';
import { getAll, update, remove } from '../api/shortenlink.ts';

interface ShortenLink {
  id: string;
  short_code: string;
  original_url: string;
  created_at: string;
}

const ShortenLinkList: React.FC = () => {
  const [links, setLinks] = useState<ShortenLink[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [search, setSearch] = useState('');
  const [editDialog, setEditDialog] = useState<{ open: boolean; link: ShortenLink | null }>({ open: false, link: null });
  const [newUrl, setNewUrl] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });

  const fetchLinks = useCallback(async () => {
    try {
      const data = await getAll({ page, limit, search });
      if (data && data.data && Array.isArray(data.data)) {
        setLinks(data.data);
        setTotal(data.total || 0);
      } else if (Array.isArray(data)) {
        setLinks(data);
        setTotal(data.length);
      } else {
        setLinks([]);
        setTotal(0);
      }
    } catch (err: any) {
      setSnackbar({ open: true, message: 'Failed to fetch links', severity: 'error' });
      setLinks([]);
      setTotal(0);
    }
  }, [page, limit, search]);

  useEffect(() => {
    fetchLinks();
  }, [fetchLinks]);

  const handleEdit = (link: ShortenLink) => {
    setEditDialog({ open: true, link });
    setNewUrl(link.original_url);
  };

  const handleUpdate = async () => {
    if (!editDialog.link) return;
    try {
      await update(editDialog.link.id, { original_url: newUrl });
      setSnackbar({ open: true, message: 'Link updated successfully', severity: 'success' });
      setEditDialog({ open: false, link: null });
      fetchLinks();
    } catch (err: any) {
      setSnackbar({ open: true, message: 'Failed to update link', severity: 'error' });
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await remove(id);
      setSnackbar({ open: true, message: 'Link deleted successfully', severity: 'success' });
      fetchLinks();
    } catch (err: any) {
      setSnackbar({ open: true, message: 'Failed to delete link', severity: 'error' });
    }
  };

  const handleCopy = (code: string) => {
    navigator.clipboard.writeText(`${window.location.origin}/r/${code}`);
    setSnackbar({ open: true, message: 'Short URL copied to clipboard', severity: 'success' });
  };

  return (
    <Container>
      <Typography variant="h4" gutterBottom>
        My Shorten Links
      </Typography>
      <List>
        {links.map((link) => (
          <ListItem key={link.id} divider>
            <ListItemText primary={link.original_url} secondary={`Short: ${window.location.origin}/r/${link.short_code} | Created: ${new Date(link.created_at).toLocaleDateString()}`} />
            <IconButton onClick={() => handleCopy(link.short_code)}>
              <ContentCopy />
            </IconButton>
            <IconButton onClick={() => handleEdit(link)}>
              <Edit />
            </IconButton>
            <IconButton onClick={() => handleDelete(link.id)}>
              <Delete />
            </IconButton>
          </ListItem>
        ))}
      </List>

      <Dialog open={editDialog.open} onClose={() => setEditDialog({ open: false, link: null })}>
        <DialogTitle>Edit Link</DialogTitle>
        <DialogContent>
          <TextField fullWidth label="Original URL" value={newUrl} onChange={(e) => setNewUrl(e.target.value)} margin="normal" />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialog({ open: false, link: null })}>Cancel</Button>
          <Button onClick={handleUpdate} variant="contained">
            Update
          </Button>
        </DialogActions>
      </Dialog>

      <Box sx={{ mt: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <TextField
          label="Search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          InputProps={{
            endAdornment: <Search />,
          }}
          sx={{ width: 300 }}
        />
        <Pagination count={Math.ceil(total / limit)} page={page} onChange={(e, value) => setPage(value)} color="primary" />
      </Box>

      <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })}>
        <Alert severity={snackbar.severity}>{snackbar.message}</Alert>
      </Snackbar>
    </Container>
  );
};

export default ShortenLinkList;
