import React, { useEffect, useState, useCallback } from 'react';
import { Container, Paper, Box, Typography, Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, Snackbar, Alert, Card, CardContent, Chip, Tooltip, Grid, Skeleton } from '@mui/material';
import { DataGrid, GridColDef, GridActionsCellItem } from '@mui/x-data-grid';
import { Edit as EditIcon, Delete as DeleteIcon, Add as AddIcon, Refresh as RefreshIcon, Email as EmailIcon, Phone as PhoneIcon } from '@mui/icons-material';
import { getAll, update, remove, create, UserResponse } from '../api/user.ts';
import { useAuth } from '../context/AuthContext.tsx';
import AppLayout from './AppLayout.tsx';

interface CreateUserForm {
  name: string;
  username: string;
  email: string;
  phone: string;
  password: string;
}

interface UpdateUserForm {
  name: string;
  email: string;
  phone: string;
}

const UserList: React.FC = () => {
  const { user: authUser } = useAuth();
  const [users, setUsers] = useState<UserResponse[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [createDialog, setCreateDialog] = useState(false);
  const [editDialog, setEditDialog] = useState<{ open: boolean; user: UserResponse | null }>({
    open: false,
    user: null,
  });
  const [createForm, setCreateForm] = useState<CreateUserForm>({
    name: '',
    username: '',
    email: '',
    phone: '',
    password: '',
  });
  const [updateForm, setUpdateForm] = useState<UpdateUserForm>({ name: '', email: '', phone: '' });
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });
  const [isSaving, setIsSaving] = useState(false);
  const [paginationModel, setPaginationModel] = useState({ pageSize: 10, page: 0 });

  const fetchUsers = useCallback(async () => {
    try {
      setIsLoading(true);
      const result = await getAll({ limit: 100 });
      setUsers(result.data);
    } catch (err: any) {
      console.error('Failed to fetch users:', err);
      console.error('Error response:', err.response?.data);
      const errorMsg = err.response?.data?.message || err.message || 'Failed to fetch users';
      setSnackbar({
        open: true,
        message: `Backend error (500): ${errorMsg}. Check backend logs for details.`,
        severity: 'error',
      });
    } finally {
      setIsLoading(false);
    }
  }, []);

  const isUserAdmin = useCallback(() => {
    return authUser?.roles?.some((r) => (typeof r === 'string' ? r === 'admin' : r.name === 'admin'));
  }, [authUser]);

  useEffect(() => {
    if (isUserAdmin()) {
      fetchUsers();
    }
  }, [isUserAdmin, fetchUsers]);

  const handleCreateOpen = () => {
    setCreateForm({ name: '', username: '', email: '', phone: '', password: '' });
    setCreateDialog(true);
  };

  const handleCreateClose = () => {
    setCreateDialog(false);
  };

  const handleCreateSubmit = async () => {
    // Validate required fields
    if (!createForm.name?.trim() || !createForm.username?.trim() || !createForm.email?.trim() || !createForm.password?.trim()) {
      setSnackbar({
        open: true,
        message: 'Please fill in all required fields',
        severity: 'error',
      });
      return;
    }

    try {
      setIsSaving(true);
      // Only send defined fields to avoid backend validation errors
      const payload: any = {
        name: createForm.name.trim(),
        username: createForm.username.trim(),
        email: createForm.email.trim(),
        password: createForm.password,
      };

      // Only include phone if it has a value
      if (createForm.phone?.trim()) {
        payload.phone = createForm.phone.trim();
      }

      await create(payload);
      setSnackbar({ open: true, message: 'User created successfully', severity: 'success' });
      setCreateDialog(false);
      await fetchUsers();
    } catch (err: any) {
      const errorMsg = err.response?.data?.message || err.response?.data?.error || err.message || 'Failed to create user';
      setSnackbar({
        open: true,
        message: errorMsg,
        severity: 'error',
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleEdit = (user: UserResponse) => {
    setEditDialog({ open: true, user });
    setUpdateForm({
      name: user.name,
      email: user.email,
      phone: user.phone,
    });
  };

  const handleUpdate = async () => {
    if (!editDialog.user) return;
    try {
      setIsSaving(true);
      await update(editDialog.user.id, updateForm);
      setSnackbar({ open: true, message: 'User updated successfully', severity: 'success' });
      setEditDialog({ open: false, user: null });
      await fetchUsers();
    } catch (err: any) {
      setSnackbar({
        open: true,
        message: err.response?.data?.message || 'Failed to update user',
        severity: 'error',
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        await remove(id);
        setSnackbar({ open: true, message: 'User deleted successfully', severity: 'success' });
        await fetchUsers();
      } catch (err: any) {
        setSnackbar({
          open: true,
          message: err.response?.data?.message || 'Failed to delete user',
          severity: 'error',
        });
      }
    }
  };

  // Get row ID from response data structure
  const getRowId = (row: any) => {
    return row.id || row.ID || row.username;
  };

  const columns: GridColDef[] = [
    {
      field: 'username',
      headerName: 'Username',
      flex: 1,
      minWidth: 130,
      renderCell: (params) => <Chip label={params.value} variant="outlined" sx={{ fontFamily: 'monospace', fontWeight: 'bold' }} />,
    },
    {
      field: 'name',
      headerName: 'Full Name',
      flex: 1.5,
      minWidth: 180,
    },
    {
      field: 'email',
      headerName: 'Email',
      flex: 1.5,
      minWidth: 200,
      renderCell: (params) => (
        <Tooltip title={params.value}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, cursor: 'pointer' }}>
            <EmailIcon sx={{ fontSize: '1rem', color: 'action.active' }} />
            <Typography
              variant="body2"
              sx={{
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              }}
            >
              {params.value}
            </Typography>
          </Box>
        </Tooltip>
      ),
    },
    {
      field: 'phone',
      headerName: 'Phone',
      flex: 1,
      minWidth: 140,
      renderCell: (params) =>
        params.value ? (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
            <PhoneIcon sx={{ fontSize: '1rem', color: 'action.active' }} />
            {params.value}
          </Box>
        ) : (
          <Typography variant="body2" color="textSecondary">
            -
          </Typography>
        ),
    },
    {
      field: 'roles',
      headerName: 'Roles',
      flex: 1,
      minWidth: 130,
      renderCell: (params) => (
        <Box sx={{ display: 'flex', gap: 0.5 }}>
          {params.value?.map((role: any) => (
            <Chip key={role.id || role.name} label={role.name || role} size="small" color={role.name === 'admin' || role === 'admin' ? 'error' : 'default'} variant="outlined" />
          ))}
        </Box>
      ),
    },
    {
      field: 'created_at',
      headerName: 'Joined',
      flex: 1,
      minWidth: 140,
      renderCell: (params) =>
        new Date(params.value).toLocaleDateString('id-ID', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
        }),
    },
    {
      field: 'actions',
      type: 'actions',
      headerName: 'Actions',
      flex: 1,
      minWidth: 120,
      getActions: (params) => [
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
          disabled={params.row.id === authUser?.user_id}
        />,
      ],
    },
  ];

  const isAdmin = isUserAdmin();
  if (!isAdmin) {
    return (
      <AppLayout>
        <Container>
          <Alert severity="error">You do not have permission to access this page</Alert>
        </Container>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <Container maxWidth="lg">
        <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Box>
            <Typography variant="h4" sx={{ fontWeight: 'bold', mb: 1 }}>
              User Management
            </Typography>
            <Typography variant="body2" color="textSecondary">
              Manage all users in the system
            </Typography>
          </Box>
          <Box sx={{ display: 'flex', gap: 2 }}>
            <Button variant="contained" startIcon={<RefreshIcon />} onClick={fetchUsers} disabled={isLoading}>
              Refresh
            </Button>
            <Button variant="contained" color="success" startIcon={<AddIcon />} onClick={handleCreateOpen}>
              Add User
            </Button>
          </Box>
        </Box>

        {/* Stats Cards */}
        <Grid container spacing={2} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                {isLoading ? (
                  <Skeleton variant="text" width={100} />
                ) : (
                  <Typography color="textSecondary" gutterBottom>
                    Total Users
                  </Typography>
                )}
                {isLoading ? (
                  <Skeleton variant="rectangular" width={80} height={28} />
                ) : (
                  <Typography variant="h5" sx={{ fontWeight: 'bold', color: 'primary.main' }}>
                    {users.length}
                  </Typography>
                )}
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card>
              <CardContent sx={{ textAlign: 'center' }}>
                <Typography color="textSecondary" gutterBottom>
                  Admins
                </Typography>
                <Typography variant="h5" sx={{ fontWeight: 'bold', color: 'error.main' }}>
                  {users.filter((u) => u.roles?.some((r) => r.name === 'admin')).length}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Data Table */}
        <Paper sx={{ height: 'auto', width: '100%' }}>
          <DataGrid
            rows={users}
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
            getRowId={getRowId}
          />
        </Paper>

        {/* Create User Dialog */}
        <Dialog open={createDialog} onClose={handleCreateClose} maxWidth="sm" fullWidth>
          <DialogTitle>Create New User</DialogTitle>
          <DialogContent sx={{ pt: 2 }}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField fullWidth label="Full Name" value={createForm.name} onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })} required />
              </Grid>
              <Grid item xs={12}>
                <TextField fullWidth label="Username" value={createForm.username} onChange={(e) => setCreateForm({ ...createForm, username: e.target.value })} required />
              </Grid>
              <Grid item xs={12}>
                <TextField fullWidth label="Email" type="email" value={createForm.email} onChange={(e) => setCreateForm({ ...createForm, email: e.target.value })} required />
              </Grid>
              <Grid item xs={12}>
                <TextField fullWidth label="Phone" value={createForm.phone} onChange={(e) => setCreateForm({ ...createForm, phone: e.target.value })} />
              </Grid>
              <Grid item xs={12}>
                <TextField fullWidth label="Password" type="password" value={createForm.password} onChange={(e) => setCreateForm({ ...createForm, password: e.target.value })} required />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCreateClose}>Cancel</Button>
            <Button onClick={handleCreateSubmit} variant="contained" disabled={isSaving}>
              {isSaving ? 'Creating...' : 'Create'}
            </Button>
          </DialogActions>
        </Dialog>

        {/* Edit User Dialog */}
        <Dialog open={editDialog.open} onClose={() => setEditDialog({ open: false, user: null })} maxWidth="sm" fullWidth>
          <DialogTitle>Edit User</DialogTitle>
          <DialogContent sx={{ pt: 2 }}>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField fullWidth label="Username" value={editDialog.user?.username || ''} disabled variant="filled" />
              </Grid>
              <Grid item xs={12}>
                <TextField fullWidth label="Full Name" value={updateForm.name} onChange={(e) => setUpdateForm({ ...updateForm, name: e.target.value })} />
              </Grid>
              <Grid item xs={12}>
                <TextField fullWidth label="Email" type="email" value={updateForm.email} onChange={(e) => setUpdateForm({ ...updateForm, email: e.target.value })} />
              </Grid>
              <Grid item xs={12}>
                <TextField fullWidth label="Phone" value={updateForm.phone} onChange={(e) => setUpdateForm({ ...updateForm, phone: e.target.value })} />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setEditDialog({ open: false, user: null })}>Cancel</Button>
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

export default UserList;
