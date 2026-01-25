import React, { useState } from 'react';
import { AppBar, Toolbar, Typography, IconButton, Drawer, List, ListItem, ListItemIcon, ListItemText, Box, Avatar, Menu, MenuItem, Divider, Badge } from '@mui/material';
import { Menu as MenuIcon, Dashboard as DashboardIcon, Link as LinkIcon, People as PeopleIcon, Person as PersonIcon, Settings as SettingsIcon, Logout as LogoutIcon, Home as HomeIcon } from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface AppLayoutProps {
  children: React.ReactNode;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const navigate = useNavigate();
  const location = useLocation();
  const { logout, user } = useAuth();

  const handleDrawerToggle = () => {
    setDrawerOpen(!drawerOpen);
  };

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleProfileMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const isUserAdmin = () => {
    return user?.roles?.some((r) => (typeof r === 'string' ? r === 'admin' : r.name === 'admin'));
  };

  const menuItems = [{ label: 'Dashboard', icon: DashboardIcon, path: '/dashboard' }, { label: 'My Links', icon: LinkIcon, path: '/links' }, ...(isUserAdmin() ? [{ label: 'User Management', icon: PeopleIcon, path: '/users' }] : [])];

  const drawerContent = (
    <Box sx={{ width: 250, pt: 2 }}>
      <Box sx={{ px: 2, pb: 2 }}>
        <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
          DigiLearn
        </Typography>
      </Box>
      <Divider />
      <List>
        {menuItems.map((item) => (
          <ListItem
            key={item.path}
            component="button"
            onClick={() => {
              navigate(item.path);
              setDrawerOpen(false);
            }}
            sx={{
              backgroundColor: location.pathname === item.path ? 'rgba(33, 150, 243, 0.1)' : 'transparent',
              '&:hover': { backgroundColor: 'rgba(33, 150, 243, 0.05)' },
            }}
          >
            <ListItemIcon sx={{ color: location.pathname === item.path ? 'primary.main' : 'inherit' }}>
              <item.icon />
            </ListItemIcon>
            <ListItemText primary={item.label} sx={{ color: location.pathname === item.path ? 'primary.main' : 'inherit' }} />
          </ListItem>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      {/* AppBar */}
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          <IconButton color="inherit" aria-label="open drawer" onClick={handleDrawerToggle} sx={{ mr: 2 }}>
            <MenuIcon />
          </IconButton>

          <HomeIcon sx={{ mr: 1 }} />
          <Typography variant="h6" sx={{ flexGrow: 1, fontWeight: 'bold' }}>
            DigiLearn
          </Typography>

          {/* Profile Menu */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography variant="body2" sx={{ mr: 1 }}>
              {user?.username}
            </Typography>
            <IconButton onClick={handleProfileMenuOpen} sx={{ p: 0 }}>
              <Badge
                overlap="circular"
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'right',
                }}
                variant="dot"
                sx={{
                  '& .MuiBadge-badge': {
                    backgroundColor: '#44b700',
                    color: '#44b700',
                    boxShadow: '0 0 0 2px white',
                    '&::after': {
                      position: 'absolute',
                      top: 0,
                      left: 0,
                      width: '100%',
                      height: '100%',
                      borderRadius: '50%',
                      animation: 'ripple 1.2s infinite ease-in-out',
                      border: '1px solid currentColor',
                    },
                  },
                  '@keyframes ripple': {
                    '0%': {
                      transform: 'scale(.8)',
                      opacity: 1,
                    },
                    '100%': {
                      transform: 'scale(2.4)',
                      opacity: 0,
                    },
                  },
                }}
              >
                <Avatar sx={{ width: 36, height: 36, bgcolor: 'secondary.main' }}>{user?.username?.[0]?.toUpperCase()}</Avatar>
              </Badge>
            </IconButton>
          </Box>

          {/* Profile Dropdown Menu */}
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleProfileMenuClose}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right',
            }}
            transformOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
          >
            <MenuItem disabled sx={{ color: 'text.secondary' }}>
              <Typography variant="body2">
                {user?.username} ({user?.roles?.map((r) => (typeof r === 'string' ? r : r.name)).join(', ')})
              </Typography>
            </MenuItem>
            <Divider />
            <MenuItem
              onClick={() => {
                navigate('/profile');
                handleProfileMenuClose();
              }}
            >
              <PersonIcon sx={{ mr: 1 }} />
              Profile
            </MenuItem>
            <MenuItem
              onClick={() => {
                navigate('/settings');
                handleProfileMenuClose();
              }}
            >
              <SettingsIcon sx={{ mr: 1 }} />
              Settings
            </MenuItem>
            <Divider />
            <MenuItem
              onClick={() => {
                handleLogout();
                handleProfileMenuClose();
              }}
            >
              <LogoutIcon sx={{ mr: 1 }} />
              Logout
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* Drawer */}
      <Drawer
        anchor="left"
        open={drawerOpen}
        onClose={handleDrawerToggle}
        sx={{
          width: 250,
          flexShrink: 0,
          '& .MuiDrawer-paper': {
            width: 250,
            boxSizing: 'border-box',
            marginTop: '64px',
          },
        }}
      >
        {drawerContent}
      </Drawer>

      {/* Main Content */}
      <Box
        sx={{
          flexGrow: 1,
          p: 3,
          mt: 8,
          width: '100%',
          backgroundColor: '#f5f5f5',
          minHeight: '100vh',
        }}
      >
        {children}
      </Box>
    </Box>
  );
};

export default AppLayout;
