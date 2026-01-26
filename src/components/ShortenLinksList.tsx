import React, { useEffect, useState, useCallback } from 'react';
import { Container, Paper, Box, Typography, Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Snackbar, Alert, Card, CardContent, Chip, Tooltip, Grid, Skeleton } from '@mui/material';
import { DataGrid, GridColDef, GridActionsCellItem } from '@mui/x-data-grid';
import { Edit as EditIcon, Delete as DeleteIcon, ContentCopy as ContentCopyIcon, OpenInNew as OpenInNewIcon, Refresh as RefreshIcon } from '@mui/icons-material';
import { getAll, update, remove, ShortenLinkResponse } from '../api/shortenlink.ts';
import AppLayout from './AppLayout.tsx';

const ShortenLinkList: React.FC = () => {
  const [links, setLinks] = useState<ShortenLinkResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [editDialog, setEditDialog] = useState<{ open: boolean; link: ShortenLinkResponse | null }>({
    open: false,
    link: null,
  });
  const [newUrl, setNewUrl] = useState('');
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });
  const [isSaving, setIsSaving] = useState(false);
  const [paginationModel, setPaginationModel] = useState({ pageSize: 10, page: 0 });

  const fetchLinks = useCallback(async () => {
    try {
      setIsLoading(true);
      const data = await getAll();
      setLinks(data);
    } catch (err: any) {
      console.error('Failed to fetch links:', err);
      setSnackbar({ open: true, message: 'Failed to fetch links', severity: 'error' });
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchLinks();
  }, [fetchLinks]);

  const handleEdit = (link: ShortenLinkResponse) => {
    setEditDialog({ open: true, link });
    setNewUrl(link.url);
  };

  const handleUpdate = async () => {
    if (!editDialog.link) return;
    try {
      setIsSaving(true);
      // Backend expects 'original_url'
      await update(editDialog.link.id, { original_url: newUrl });
      setSnackbar({ open: true, message: 'Link updated successfully', severity: 'success' });
      setEditDialog({ open: false, link: null });
      await fetchLinks();
    } catch (err: any) {
      setSnackbar({
        open: true,
        message: err.response?.data?.message || 'Failed to update link',
        severity: 'error',
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this link?')) {
      try {
        await remove(id);
        setSnackbar({ open: true, message: 'Link deleted successfully', severity: 'success' });
        await fetchLinks();
      } catch (err: any) {
        setSnackbar({
          open: true,
          message: err.response?.data?.message || 'Failed to delete link',
          severity: 'error',
        });
      }
    }
  };

  const handleCopy = (code: string) => {
    const shortUrl = `${window.location.origin}/mydigilearn/${code}`;
    navigator.clipboard.writeText(shortUrl);
    setSnackbar({ open: true, message: 'Short URL copied to clipboard', severity: 'success' });
  };

  const handleRefresh = () => {
    fetchLinks();
  };

  const columns: GridColDef[] = [
    {
      field: 'code',
      headerName: 'Short Code',
      flex: 1,
      minWidth: 120,
      renderCell: (params) => <Chip label={params.value} color="primary" variant="outlined" sx={{ fontWeight: 'bold', fontFamily: 'monospace' }} />,
    },
    {
      field: 'url',
      headerName: 'Original URL',
      flex: 2,
      minWidth: 250,
      renderCell: (params) => (
        <Tooltip title={params.value}>
          <Typography
            variant="body2"
            sx={{
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
              cursor: 'pointer',
              color: 'primary.main',
              '&:hover': { textDecoration: 'underline' },
            }}
            onClick={() => window.open(params.value, '_blank')}
          >
            {params.value}
          </Typography>
        </Tooltip>
      ),
    },
    {
      field: 'created_at',
      headerName: 'Created At',
      flex: 1,
      minWidth: 180,
      valueFormatter: (params) => {
        try {
          return new Date(params.value as string).toLocaleString();
        } catch {
          return String(params.value ?? '');
        }
      },
    },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Actions',
      flex: 1,
      minWidth: 160,
      getActions: (params) => [
        <GridActionsCellItem
          icon={
            <Tooltip title="Copy short URL">
              <ContentCopyIcon />
            </Tooltip>
          }
          label="Copy"
          onClick={() => handleCopy(params.row.code)}
          color="primary"
        />,
        <GridActionsCellItem
          icon={
            <Tooltip title="Open in New Tab">
              <OpenInNewIcon />
            </Tooltip>
          }
          label="Open"
          onClick={() => window.open(`${window.location.origin}/mydigilearn/${params.row.code}`, '_blank')}
          color="primary"
        />,
        <GridActionsCellItem
          icon={
            <Tooltip title="Edit">
              <EditIcon />
            </Tooltip>
          }
          label="Edit"
          onClick={() => handleEdit(params.row)}
          color="primary"
        />,
        <GridActionsCellItem
          icon={
            <Tooltip title="Delete">
              <DeleteIcon sx={{ color: 'error.main' }} />
            </Tooltip>
          }
          label="Delete"
          onClick={() => handleDelete(params.row.id)}
        />,
      ],
    },
  ];

  return (
    <AppLayout>
      <Container maxWidth="lg">
        <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Box>
            {isLoading ? (
              <React.Fragment>
                <Skeleton variant="text" width={220} height={40} />
                <Skeleton variant="text" width={260} />
              </React.Fragment>
            ) : (
              <React.Fragment>
                <Typography variant="h4" sx={{ fontWeight: 'bold', mb: 1 }}>
                  My Shortened Links
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  Manage and track all your shortened URLs
                </Typography>
              </React.Fragment>
            )}
          </Box>
          <Button variant="contained" startIcon={<RefreshIcon />} onClick={handleRefresh} disabled={isLoading}>
            Refresh
          </Button>
        </Box>

        {/* Stats Cards */}
        <Grid container spacing={2} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Typography color="textSecondary" gutterBottom>
                  Total Links
                </Typography>
                <Typography variant="h5" sx={{ fontWeight: 'bold', color: 'primary.main' }}>
                  {links.length}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Typography color="textSecondary" gutterBottom>
                  This Month
                </Typography>
                <Typography variant="h5" sx={{ fontWeight: 'bold', color: 'success.main' }}>
                  {links.filter((link) => new Date(link.created_at).getMonth() === new Date().getMonth() && new Date(link.created_at).getFullYear() === new Date().getFullYear()).length}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Data Table */}
        <Paper sx={{ height: 'auto', width: '100%' }}>
          <DataGrid
            rows={links}
            columns={columns}
            pageSizeOptions={[5, 10, 25]}
            paginationModel={paginationModel}
            onPaginationModelChange={setPaginationModel}
            loading={isLoading}
            sx={{
              '& .MuiDataGrid-cell': {
                borderColor: 'rgba(224, 224, 224, 0.5)',
              },
              '& .MuiDataGrid-row:hover': {
                backgroundColor: 'rgba(0, 0, 0, 0.02)',
              },
            }}
            getRowId={(row) => row.id}
          />
        </Paper>

        {/* Edit Dialog */}
        <Dialog open={editDialog.open} onClose={() => setEditDialog({ open: false, link: null })} maxWidth="sm" fullWidth>
          <DialogTitle>Edit Shortened Link</DialogTitle>
          <DialogContent sx={{ pt: 2 }}>
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="textSecondary">
                Short Code
              </Typography>
              <TextField fullWidth value={editDialog.link?.code || ''} disabled variant="filled" size="small" />
            </Box>
            <TextField fullWidth label="Original URL" value={newUrl} onChange={(e) => setNewUrl(e.target.value)} type="url" placeholder="https://example.com/very/long/url" />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setEditDialog({ open: false, link: null })}>Cancel</Button>
            <Button onClick={handleUpdate} variant="contained" disabled={isSaving}>
              {isSaving ? 'Saving...' : 'Save'}
            </Button>
          </DialogActions>
        </Dialog>

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

export default ShortenLinkList;
