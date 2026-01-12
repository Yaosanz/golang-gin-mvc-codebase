import React, { useEffect, useState, useCallback } from 'react';
import {
  Container,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Snackbar,
  Alert,
  Fab,
  Pagination,
  Box,
  TextField,
} from '@mui/material';
import { Edit, Delete, Add, Search } from '@mui/icons-material';
import { getAll, remove } from '../api/user.ts';
import UserForm from './UserForm.tsx';

interface User {
  id: number;
  name: string;
  username: string;
  email: string;
}

const UserList: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [search, setSearch] = useState('');
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; user: User | null }>({ open: false, user: null });
  const [formDialog, setFormDialog] = useState<{ open: boolean; user: User | null }>({ open: false, user: null });
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });

  const fetchUsers = useCallback(async () => {
    try {
      const data = await getAll({ page, limit, search });
      if (data && data.data && Array.isArray(data.data)) {
        setUsers(data.data);
        setTotal(data.total || 0);
      } else if (Array.isArray(data)) {
        setUsers(data);
        setTotal(data.length);
      } else {
        setUsers([]);
        setTotal(0);
      }
    } catch (err: any) {
      setSnackbar({ open: true, message: 'Failed to fetch users', severity: 'error' });
      setUsers([]);
      setTotal(0);
    }
  }, [page, limit, search]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const handleDelete = async () => {
    if (!deleteDialog.user) return;
    try {
      await remove(deleteDialog.user.id.toString());
      setSnackbar({ open: true, message: 'User deleted successfully', severity: 'success' });
      setDeleteDialog({ open: false, user: null });
      fetchUsers();
    } catch (err: any) {
      setSnackbar({ open: true, message: 'Failed to delete user', severity: 'error' });
    }
  };

  return (
    <Container>
      <Typography variant="h4" gutterBottom>
        User Management
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Username</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell>{user.id}</TableCell>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.username}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>
                  <IconButton onClick={() => setFormDialog({ open: true, user })}>
                    <Edit />
                  </IconButton>
                  <IconButton onClick={() => setDeleteDialog({ open: true, user })}>
                    <Delete />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={deleteDialog.open} onClose={() => setDeleteDialog({ open: false, user: null })}>
        <DialogTitle>Delete User</DialogTitle>
        <DialogContent>Are you sure you want to delete user "{deleteDialog.user?.name}"?</DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialog({ open: false, user: null })}>Cancel</Button>
          <Button onClick={handleDelete} variant="contained" color="error">
            Delete
          </Button>
        </DialogActions>
      </Dialog>

      <Fab color="primary" aria-label="add" sx={{ position: 'fixed', bottom: 16, right: 16 }} onClick={() => setFormDialog({ open: true, user: null })}>
        <Add />
      </Fab>

      <UserForm open={formDialog.open} onClose={() => setFormDialog({ open: false, user: null })} user={formDialog.user} onSave={fetchUsers} />

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

export default UserList;
