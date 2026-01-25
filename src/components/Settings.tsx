import React, { useState } from 'react';
import { Container, Box, Typography, Grid, Card, CardContent, Switch, FormControlLabel, Divider, Button, TextField, Select, MenuItem, FormControl, InputLabel, Chip, Alert, Stack } from '@mui/material';
import { Settings as SettingsIcon, Notifications as NotificationsIcon, Security as SecurityIcon, Palette as PaletteIcon, Language as LanguageIcon, Save as SaveIcon } from '@mui/icons-material';
import AppLayout from './AppLayout.tsx';

const Settings: React.FC = () => {
  const [settings, setSettings] = useState({
    emailNotifications: true,
    pushNotifications: false,
    darkMode: false,
    language: 'en',
    timezone: 'UTC',
    autoSave: true,
  });

  const [isSaved, setIsSaved] = useState(false);

  const handleToggle = (setting: string) => {
    setSettings((prev) => ({
      ...prev,
      [setting]: !prev[setting as keyof typeof prev],
    }));
    setIsSaved(false);
  };

  const handleChange = (setting: string, value: string) => {
    setSettings((prev) => ({
      ...prev,
      [setting]: value,
    }));
    setIsSaved(false);
  };

  const handleSave = () => {
    // TODO: Implement save to backend
    console.log('Saving settings:', settings);
    setIsSaved(true);
    setTimeout(() => setIsSaved(false), 3000);
  };

  return (
    <AppLayout>
      <Container maxWidth="lg">
        <Box sx={{ mb: 3 }}>
          <Typography variant="h4" sx={{ fontWeight: 'bold', mb: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
            <SettingsIcon fontSize="large" />
            Settings
          </Typography>
          <Typography variant="body2" color="textSecondary">
            Manage your application preferences and settings
          </Typography>
        </Box>

        {isSaved && (
          <Alert severity="success" sx={{ mb: 3 }}>
            Settings saved successfully!
          </Alert>
        )}

        <Grid container spacing={3}>
          {/* Notifications Settings */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2, gap: 1 }}>
                  <NotificationsIcon color="primary" />
                  <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                    Notifications
                  </Typography>
                </Box>
                <Divider sx={{ mb: 2 }} />
                <Stack spacing={2}>
                  <FormControlLabel control={<Switch checked={settings.emailNotifications} onChange={() => handleToggle('emailNotifications')} color="primary" />} label="Email Notifications" />
                  <Typography variant="caption" color="textSecondary" sx={{ ml: 4, mt: -1 }}>
                    Receive email updates about your account activity
                  </Typography>
                  <FormControlLabel control={<Switch checked={settings.pushNotifications} onChange={() => handleToggle('pushNotifications')} color="primary" />} label="Push Notifications" />
                  <Typography variant="caption" color="textSecondary" sx={{ ml: 4, mt: -1 }}>
                    Get browser notifications for important updates
                  </Typography>
                </Stack>
              </CardContent>
            </Card>
          </Grid>

          {/* Appearance Settings */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2, gap: 1 }}>
                  <PaletteIcon color="primary" />
                  <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                    Appearance
                  </Typography>
                </Box>
                <Divider sx={{ mb: 2 }} />
                <Stack spacing={2}>
                  <FormControlLabel control={<Switch checked={settings.darkMode} onChange={() => handleToggle('darkMode')} color="primary" />} label="Dark Mode" />
                  <Typography variant="caption" color="textSecondary" sx={{ ml: 4, mt: -1 }}>
                    Enable dark theme across the application
                    <Chip label="Coming Soon" size="small" color="info" sx={{ ml: 1 }} />
                  </Typography>
                </Stack>
              </CardContent>
            </Card>
          </Grid>

          {/* Language & Region Settings */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2, gap: 1 }}>
                  <LanguageIcon color="primary" />
                  <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                    Language & Region
                  </Typography>
                </Box>
                <Divider sx={{ mb: 2 }} />
                <Stack spacing={2}>
                  <FormControl fullWidth>
                    <InputLabel id="language-label">Language</InputLabel>
                    <Select labelId="language-label" value={settings.language} label="Language" onChange={(e) => handleChange('language', e.target.value)}>
                      <MenuItem value="en">English</MenuItem>
                      <MenuItem value="id">Bahasa Indonesia</MenuItem>
                      <MenuItem value="es">Español</MenuItem>
                      <MenuItem value="fr">Français</MenuItem>
                    </Select>
                  </FormControl>
                  <FormControl fullWidth>
                    <InputLabel id="timezone-label">Timezone</InputLabel>
                    <Select labelId="timezone-label" value={settings.timezone} label="Timezone" onChange={(e) => handleChange('timezone', e.target.value)}>
                      <MenuItem value="UTC">UTC (Coordinated Universal Time)</MenuItem>
                      <MenuItem value="Asia/Jakarta">Asia/Jakarta (WIB)</MenuItem>
                      <MenuItem value="America/New_York">America/New York (EST)</MenuItem>
                      <MenuItem value="Europe/London">Europe/London (GMT)</MenuItem>
                      <MenuItem value="Asia/Tokyo">Asia/Tokyo (JST)</MenuItem>
                    </Select>
                  </FormControl>
                </Stack>
              </CardContent>
            </Card>
          </Grid>

          {/* Security Settings */}
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2, gap: 1 }}>
                  <SecurityIcon color="primary" />
                  <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                    Security
                  </Typography>
                </Box>
                <Divider sx={{ mb: 2 }} />
                <Stack spacing={2}>
                  <FormControlLabel control={<Switch checked={settings.autoSave} onChange={() => handleToggle('autoSave')} color="primary" />} label="Auto-save Changes" />
                  <Typography variant="caption" color="textSecondary" sx={{ ml: 4, mt: -1 }}>
                    Automatically save your work periodically
                  </Typography>
                  <Button variant="outlined" color="secondary" fullWidth sx={{ mt: 2 }}>
                    Change Password
                  </Button>
                  <Button variant="outlined" color="error" fullWidth>
                    Enable Two-Factor Authentication
                  </Button>
                </Stack>
              </CardContent>
            </Card>
          </Grid>

          {/* Advanced Settings */}
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2 }}>
                  Advanced Settings
                </Typography>
                <Divider sx={{ mb: 2 }} />
                <Grid container spacing={2}>
                  <Grid item xs={12} sm={6}>
                    <TextField fullWidth label="API Endpoint" defaultValue="https://api.example.com" disabled helperText="Contact administrator to change" />
                  </Grid>
                  <Grid item xs={12} sm={6}>
                    <TextField fullWidth label="Session Timeout (minutes)" type="number" defaultValue="30" disabled helperText="Managed by server configuration" />
                  </Grid>
                </Grid>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Save Button */}
        <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end', gap: 2 }}>
          <Button variant="outlined" color="secondary">
            Reset to Defaults
          </Button>
          <Button variant="contained" startIcon={<SaveIcon />} onClick={handleSave}>
            Save Changes
          </Button>
        </Box>
      </Container>
    </AppLayout>
  );
};

export default Settings;
