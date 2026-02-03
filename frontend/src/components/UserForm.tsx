import React, { useState } from 'react';
import { Dialog, DialogTitle, DialogContent, DialogActions, Button, TextField, FormControlLabel, Switch } from '@mui/material';
import { create, update, CreateUserRequest } from '../api/user.ts';

interface User {
  id?: number;
  name: string;
  username: string;
  email: string;
  phone?: string;
  is_active?: boolean;
  password?: string;
}

interface UserFormProps {
  open: boolean;
  onClose: () => void;
  user?: User | null;
  onSave: () => void;
}

const UserForm: React.FC<UserFormProps> = ({ open, onClose, user, onSave }) => {
  const [formData, setFormData] = useState<User>({
    name: user?.name || '',
    username: user?.username || '',
    email: user?.email || '',
    phone: user?.phone || '',
    is_active: user?.is_active ?? true,
    password: '',
  });
  const [errors, setErrors] = useState<{ [key: string]: string }>({});

  const validate = () => {
    const newErrors: { [key: string]: string } = {};
    if (!formData.name.trim()) newErrors.name = 'Name is required';
    if (!formData.username.trim()) newErrors.username = 'Username is required';
    if (!formData.email.trim()) newErrors.email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(formData.email)) newErrors.email = 'Email is invalid';
    if (!user && !formData.password?.trim()) newErrors.password = 'Password is required';
    return newErrors;
  };

  const handleSubmit = async () => {
    const validationErrors = validate();
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }
    setErrors({});
    try {
      if (user) {
        await update(user.id!.toString(), formData);
      } else {
        await create(formData as CreateUserRequest);
      }
      onSave();
      onClose();
    } catch (err: any) {
      // Handle error
    }
  };

  const handleChange = (field: keyof User, value: any) => {
    setFormData({ ...formData, [field]: value });
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{user ? 'Edit User' : 'Create User'}</DialogTitle>
      <DialogContent>
        <TextField fullWidth label="Name" value={formData.name} onChange={(e) => handleChange('name', e.target.value)} margin="normal" error={!!errors.name} helperText={errors.name} />
        <TextField fullWidth label="Username" value={formData.username} onChange={(e) => handleChange('username', e.target.value)} margin="normal" error={!!errors.username} helperText={errors.username} />
        <TextField fullWidth label="Email" value={formData.email} onChange={(e) => handleChange('email', e.target.value)} margin="normal" error={!!errors.email} helperText={errors.email} />
        <TextField fullWidth label="Phone" value={formData.phone} onChange={(e) => handleChange('phone', e.target.value)} margin="normal" />
        <FormControlLabel control={<Switch checked={formData.is_active} onChange={(e) => handleChange('is_active', e.target.checked)} />} label="Active" />
        {!user && <TextField fullWidth label="Password" type="password" value={formData.password} onChange={(e) => handleChange('password', e.target.value)} margin="normal" error={!!errors.password} helperText={errors.password} />}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSubmit} variant="contained">
          {user ? 'Update' : 'Create'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default UserForm;
