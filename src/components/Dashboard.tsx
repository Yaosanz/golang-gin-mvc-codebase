import * as React from 'react';
import { Box, Typography, Button, Tabs, Tab } from '@mui/material';
import { useAuth } from '../context/AuthContext.tsx';
import { useNavigate } from 'react-router-dom';
import ShortenLinkCreate from './ShortenLinkCreate.tsx';
import ShortenLinkList from './ShortenLinkList.tsx';
import UserList from './UserList.tsx';

export default function Dashboard() {
  const { logout, user } = useAuth();
  const navigate = useNavigate();
  const [tabValue, setTabValue] = React.useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  return (
    <Box sx={{ p: 4 }}>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      <Typography variant="body1" sx={{ mb: 3 }}>
        Welcome! Manage your shorten links here.
      </Typography>

      <Tabs value={tabValue} onChange={handleTabChange} sx={{ mb: 3 }}>
        <Tab label="Create Link" />
        <Tab label="My Links" />
        {user?.role === 'admin' && <Tab label="User Management" />}
      </Tabs>

      {tabValue === 0 && <ShortenLinkCreate />}
      {tabValue === 1 && <ShortenLinkList />}
      {tabValue === 2 && user?.role === 'admin' && <UserList />}

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
