import React, { useEffect, useState, useCallback } from 'react';
import { Container, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, IconButton, Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Snackbar, Alert } from '@mui/material';
import { Edit, Delete, ContentCopy, OpenInNew } from '@mui/icons-material';
import { getAll, update, remove } from '../api/shortenlink.ts';

interface ShortenLink {
  ID: string;
  ShortCode: string;
  OriginalURL: string;
  CreatedAt: string;
}

const ShortenLinkList: React.FC = () => {
  const [links, setLinks] = useState<ShortenLink[]>([]);
  const [editDialog, setEditDialog] = useState<{ open: boolean; link: ShortenLink | null }>({ open: false, link: null });
  const [newUrl, setNewUrl] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });

  const fetchLinks = useCallback(async () => {
    try {
      const data = await getAll();
      setLinks(data);
    } catch (err: any) {
      console.error('Failed to fetch links:', err);
      setSnackbar({ open: true, message: 'Failed to fetch links', severity: 'error' });
      setLinks([]);
    }
  }, []);

  useEffect(() => {
    fetchLinks();
  }, [fetchLinks]);

  const handleEdit = (link: ShortenLink) => {
    setEditDialog({ open: true, link });
    setNewUrl(link.OriginalURL);
  };

  const handleUpdate = async () => {
    if (!editDialog.link) return;
    try {
      await update(editDialog.link.ID, { original_url: newUrl });
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
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Original URL</TableCell>
              <TableCell>Short URL</TableCell>
              <TableCell>Created At</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {links.map((link) => (
              <TableRow key={link.ID}>
                <TableCell>{link.OriginalURL}</TableCell>
                <TableCell>
                  <a href={`${window.location.origin}/mydigilearn/${link.ShortCode}`} target="_blank" rel="noopener noreferrer">
                    {window.location.origin}/mydigilearn/{link.ShortCode}
                  </a>
                </TableCell>
                <TableCell>{new Date(link.CreatedAt).toLocaleDateString()}</TableCell>
                <TableCell>
                  <IconButton onClick={() => handleCopy(link.ShortCode)} title="Copy">
                    <ContentCopy />
                  </IconButton>
                  <IconButton onClick={() => window.open(`${window.location.origin}/r/${link.ShortCode}`, '_blank')} title="Open">
                    <OpenInNew />
                  </IconButton>
                  <IconButton onClick={() => handleEdit(link)} title="Edit">
                    <Edit />
                  </IconButton>
                  <IconButton onClick={() => handleDelete(link.ID)} title="Delete">
                    <Delete />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

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

      <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })}>
        <Alert severity={snackbar.severity}>{snackbar.message}</Alert>
      </Snackbar>
    </Container>
  );
};

export default ShortenLinkList;
