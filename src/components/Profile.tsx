import React, { useEffect, useState } from 'react';
import { Container, Paper, Box, Typography, TextField, Button, Avatar, Grid, Card, CardContent, Snackbar, Alert, CircularProgress, Divider } from '@mui/material';
import { Edit as EditIcon, Save as SaveIcon, Cancel as CancelIcon, Person as PersonIcon } from '@mui/icons-material';
import { useAuth } from '../context/AuthContext';
import { getProfile, updateProfile, UpdateProfileRequest } from '../api/auth';
import AppLayout from './AppLayout';

interface ProfileData {
  id: string;
  username: string;
  email: string;
  name: string;
  phone: string;
  created_at: string;
  updated_at: string;
  roles: string[];
}

const Profile: React.FC = () => {
  const { user } = useAuth();
  const [profileData, setProfileData] = useState<ProfileData | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [formData, setFormData] = useState<UpdateProfileRequest>({});
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' as 'success' | 'error' });

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      setIsLoading(true);
      const response = await getProfile();
      setProfileData(response.data);
      setFormData({
        name: response.data.name,
        email: response.data.email,
        phone: response.data.phone,
      });
    } catch (error) {
      setSnackbar({ open: true, message: 'Failed to load profile', severity: 'error' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleCancel = () => {
    setIsEditing(false);
    if (profileData) {
      setFormData({
        name: profileData.name,
        email: profileData.email,
        phone: profileData.phone,
      });
    }
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);
      await updateProfile(formData);
      await fetchProfile();
      setIsEditing(false);
      setSnackbar({ open: true, message: 'Profile updated successfully', severity: 'success' });
    } catch (error: any) {
      setSnackbar({
        open: true,
        message: error.response?.data?.message || 'Failed to update profile',
        severity: 'error',
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleInputChange = (field: keyof UpdateProfileRequest) => (event: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({
      ...prev,
      [field]: event.target.value,
    }));
  };

  if (isLoading) {
    return (
      <AppLayout>
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '400px' }}>
          <CircularProgress />
        </Box>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <Container maxWidth="md">
        <Grid container spacing={3}>
          {/* Profile Header Card */}
          <Grid item xs={12}>
            <Card sx={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', color: 'white' }}>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
                  <Avatar
                    sx={{
                      width: 100,
                      height: 100,
                      bgcolor: 'rgba(255, 255, 255, 0.3)',
                      border: '3px solid white',
                      fontSize: '2.5rem',
                    }}
                  >
                    <PersonIcon sx={{ fontSize: '3rem' }} />
                  </Avatar>
                  <Box sx={{ flex: 1 }}>
                    <Typography variant="h4" sx={{ fontWeight: 'bold', mb: 1 }}>
                      {profileData?.name}
                    </Typography>
                    <Typography variant="body1" sx={{ opacity: 0.9, mb: 1 }}>
                      @{profileData?.username}
                    </Typography>
                    <Box sx={{ display: 'flex', gap: 2 }}>
                      {profileData?.roles?.map((role) => (
                        <Typography
                          key={role}
                          variant="body2"
                          sx={{
                            bgcolor: 'rgba(255, 255, 255, 0.2)',
                            px: 1.5,
                            py: 0.5,
                            borderRadius: '12px',
                            fontWeight: 'bold',
                          }}
                        >
                          {role.toUpperCase()}
                        </Typography>
                      ))}
                    </Box>
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>

          {/* Profile Details Card */}
          <Grid item xs={12}>
            <Paper sx={{ p: 3 }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                  Profile Information
                </Typography>
                <Button variant={isEditing ? 'outlined' : 'contained'} startIcon={isEditing ? <CancelIcon /> : <EditIcon />} onClick={isEditing ? handleCancel : handleEdit} color={isEditing ? 'error' : 'primary'}>
                  {isEditing ? 'Cancel' : 'Edit Profile'}
                </Button>
              </Box>

              <Divider sx={{ mb: 3 }} />

              <Grid container spacing={3}>
                {/* Name */}
                <Grid item xs={12} sm={6}>
                  <TextField fullWidth label="Name" value={formData.name || ''} onChange={handleInputChange('name')} disabled={!isEditing} variant={isEditing ? 'outlined' : 'filled'} />
                </Grid>

                {/* Email */}
                <Grid item xs={12} sm={6}>
                  <TextField fullWidth label="Email" type="email" value={formData.email || ''} onChange={handleInputChange('email')} disabled={!isEditing} variant={isEditing ? 'outlined' : 'filled'} />
                </Grid>

                {/* Phone */}
                <Grid item xs={12} sm={6}>
                  <TextField fullWidth label="Phone" value={formData.phone || ''} onChange={handleInputChange('phone')} disabled={!isEditing} variant={isEditing ? 'outlined' : 'filled'} />
                </Grid>

                {/* Username (Read-only) */}
                <Grid item xs={12} sm={6}>
                  <TextField fullWidth label="Username" value={profileData?.username || ''} disabled variant="filled" />
                </Grid>

                {/* Account Details */}
                <Grid item xs={12}>
                  <Divider sx={{ my: 2 }} />
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" color="textSecondary">
                    Account Created
                  </Typography>
                  <Typography variant="body1">
                    {profileData?.created_at
                      ? new Date(profileData.created_at).toLocaleDateString('id-ID', {
                          year: 'numeric',
                          month: 'long',
                          day: 'numeric',
                        })
                      : '-'}
                  </Typography>
                </Grid>

                <Grid item xs={12} sm={6}>
                  <Typography variant="body2" color="textSecondary">
                    Last Updated
                  </Typography>
                  <Typography variant="body1">
                    {profileData?.updated_at
                      ? new Date(profileData.updated_at).toLocaleDateString('id-ID', {
                          year: 'numeric',
                          month: 'long',
                          day: 'numeric',
                        })
                      : '-'}
                  </Typography>
                </Grid>
              </Grid>

              {/* Save Button */}
              {isEditing && (
                <Box sx={{ display: 'flex', gap: 2, mt: 3, pt: 2, borderTop: '1px solid #e0e0e0' }}>
                  <Button variant="contained" startIcon={<SaveIcon />} onClick={handleSave} disabled={isSaving}>
                    {isSaving ? 'Saving...' : 'Save Changes'}
                  </Button>
                  <Button variant="outlined" onClick={handleCancel} disabled={isSaving}>
                    Cancel
                  </Button>
                </Box>
              )}
            </Paper>
          </Grid>

          {/* Account Statistics */}
          <Grid item xs={12}>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent sx={{ textAlign: 'center' }}>
                    <Typography color="textSecondary" gutterBottom>
                      Account Status
                    </Typography>
                    <Typography variant="h6" sx={{ fontWeight: 'bold', color: 'success.main' }}>
                      Active
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent sx={{ textAlign: 'center' }}>
                    <Typography color="textSecondary" gutterBottom>
                      Roles
                    </Typography>
                    <Typography variant="h6" sx={{ fontWeight: 'bold', color: 'primary.main' }}>
                      {profileData?.roles?.length || 0}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Container>

      {/* Snackbar */}
      <Snackbar open={snackbar.open} autoHideDuration={4000} onClose={() => setSnackbar({ ...snackbar, open: false })} anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}>
        <Alert severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </AppLayout>
  );
};

export default Profile;
