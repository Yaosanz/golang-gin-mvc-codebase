import * as React from 'react';
import { styled, createTheme, ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import MuiDrawer from '@mui/material/Drawer';
import Box from '@mui/material/Box';
import MuiAppBar, { AppBarProps as MuiAppBarProps } from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Badge from '@mui/material/Badge';
import Container from '@mui/material/Container';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import Link from '@mui/material/Link';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import NotificationsIcon from '@mui/icons-material/Notifications';
import { mainListItems, secondaryListItems } from './listItems.tsx';
import { Button, ListItem, ListItemText, Dialog, DialogTitle, DialogContent, DialogActions, TextField, Alert } from '@mui/material';
import { Delete, Edit, Add } from '@mui/icons-material';
import axios from 'axios';
import { useAuth } from '../context/AuthContext.tsx';
import { Navigate } from 'react-router-dom';

function Copyright(props: any) {
  return (
    <Typography variant="body2" color="text.secondary" align="center" {...props}>
      {'Copyright © '}
      <Link color="inherit" href="https://mui.com/">
        Your Website
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const drawerWidth: number = 240;

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== 'open',
})<AppBarProps>(({ theme, open }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  ...(open && {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  }),
}));

const Drawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(({ theme, open }) => ({
  '& .MuiDrawer-paper': {
    position: 'relative',
    whiteSpace: 'nowrap',
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
    boxSizing: 'border-box',
    ...(!open && {
      overflowX: 'hidden',
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
      }),
      width: theme.spacing(7),
      [theme.breakpoints.up('sm')]: {
        width: theme.spacing(9),
      },
    }),
  },
}));

const mdTheme = createTheme();

interface ShortenLink {
  id: number;
  original_url: string;
  short_code: string;
  created_at: string;
}

function DashboardContent() {
  const { token, logout } = useAuth();

  const [open, setOpen] = React.useState(true);
  const [links, setLinks] = React.useState<ShortenLink[]>([]);
  const [dialogOpen, setDialogOpen] = React.useState(false);
  const [editing, setEditing] = React.useState<ShortenLink | null>(null);
  const [originalUrl, setOriginalUrl] = React.useState('');
  const [error, setError] = React.useState('');

  const toggleDrawer = () => {
    setOpen(!open);
  };

  const fetchLinks = React.useCallback(async () => {
    try {
      const response = await axios.get('/api/v1/shorten-links/', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setLinks(response.data.data);
      setError('');
    } catch (err: any) {
      if (err.response?.status === 401) {
        setError('Unauthorized. Please login again.');
        logout();
      } else {
        setError('Failed to fetch links.');
      }
      console.error(err);
    }
  }, [token, logout]);

  React.useEffect(() => {
    if (token) {
      fetchLinks();
    }
  }, [fetchLinks, token]);

  const handleCreate = () => {
    setEditing(null);
    setOriginalUrl('');
    setDialogOpen(true);
  };

  const handleEdit = (link: ShortenLink) => {
    setEditing(link);
    setOriginalUrl(link.original_url);
    setDialogOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await axios.delete(`/api/v1/shorten-links/${id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      fetchLinks();
    } catch (err: any) {
      if (err.response?.status === 401) {
        setError('Unauthorized. Please login again.');
        logout();
      } else {
        setError('Failed to delete link.');
      }
      console.error(err);
    }
  };

  const handleSave = async () => {
    try {
      if (editing) {
        await axios.patch(
          `/api/v1/shorten-links/${editing.id}`,
          { original_url: originalUrl },
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
      } else {
        await axios.post(
          '/api/v1/shorten-links',
          { original_url: originalUrl },
          {
            headers: { Authorization: `Bearer ${token}` },
          }
        );
      }
      setDialogOpen(false);
      fetchLinks();
    } catch (err: any) {
      if (err.response?.status === 401) {
        setError('Unauthorized. Please login again.');
        logout();
      } else {
        setError('Failed to save link.');
      }
      console.error(err);
    }
  };

  if (!token) {
    return <Navigate to="/login" />;
  }

  return (
    <ThemeProvider theme={mdTheme}>
      <Box sx={{ display: 'flex' }}>
        <CssBaseline />
        <AppBar position="absolute" open={open}>
          <Toolbar
            sx={{
              pr: '24px', // keep right padding when drawer closed
            }}
          >
            <IconButton
              edge="start"
              color="inherit"
              aria-label="open drawer"
              onClick={toggleDrawer}
              sx={{
                marginRight: '36px',
                ...(open && { display: 'none' }),
              }}
            >
              <MenuIcon />
            </IconButton>
            <Typography component="h1" variant="h6" color="inherit" noWrap sx={{ flexGrow: 1 }}>
              CMS Dashboard
            </Typography>
            <IconButton color="inherit">
              <Badge badgeContent={4} color="secondary">
                <NotificationsIcon />
              </Badge>
            </IconButton>
          </Toolbar>
        </AppBar>
        <Drawer variant="permanent" open={open}>
          <Toolbar
            sx={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'flex-end',
              px: [1],
            }}
          >
            <IconButton onClick={toggleDrawer}>
              <ChevronLeftIcon />
            </IconButton>
          </Toolbar>
          <Divider />
          <List component="nav">
            {mainListItems}
            <Divider sx={{ my: 1 }} />
            {secondaryListItems}
          </List>
        </Drawer>
        <Box
          component="main"
          sx={{
            backgroundColor: (theme) => (theme.palette.mode === 'light' ? theme.palette.grey[100] : theme.palette.grey[900]),
            flexGrow: 1,
            height: '100vh',
            overflow: 'auto',
          }}
        >
          <Toolbar />
          <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Grid container spacing={3}>
              {/* Chart */}
              <Grid item xs={12} md={8} lg={9}>
                <Paper
                  sx={{
                    p: 2,
                    display: 'flex',
                    flexDirection: 'column',
                    height: 240,
                  }}
                >
                  <Typography variant="h6">Shorten Links Overview</Typography>
                  <Typography>Total Links: {links.length}</Typography>
                </Paper>
              </Grid>
              {/* Recent Deposits */}
              <Grid item xs={12} md={4} lg={3}>
                <Paper
                  sx={{
                    p: 2,
                    display: 'flex',
                    flexDirection: 'column',
                    height: 240,
                  }}
                >
                  <Typography variant="h6">Actions</Typography>
                  <Button variant="contained" startIcon={<Add />} onClick={handleCreate} sx={{ mt: 2 }}>
                    Create Link
                  </Button>
                  <Button variant="outlined" onClick={logout} sx={{ mt: 2 }}>
                    Logout
                  </Button>
                </Paper>
              </Grid>
              {/* Recent Orders */}
              <Grid item xs={12}>
                <Paper sx={{ p: 2, display: 'flex', flexDirection: 'column' }}>
                  <Typography variant="h6">Recent Links</Typography>
                  {error && (
                    <Alert severity="error" sx={{ mt: 2 }}>
                      {error}
                    </Alert>
                  )}
                  <List>
                    {links.map((link) => (
                      <ListItem
                        key={link.id}
                        secondaryAction={
                          <>
                            <IconButton onClick={() => handleEdit(link)}>
                              <Edit />
                            </IconButton>
                            <IconButton onClick={() => handleDelete(link.id)}>
                              <Delete />
                            </IconButton>
                          </>
                        }
                      >
                        <ListItemText primary={link.original_url} secondary={`Code: ${link.short_code}`} />
                      </ListItem>
                    ))}
                  </List>
                </Paper>
              </Grid>
            </Grid>
            <Copyright sx={{ pt: 4 }} />
          </Container>
        </Box>
      </Box>
      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
        <DialogTitle>{editing ? 'Edit' : 'Create'} Shorten Link</DialogTitle>
        <DialogContent>
          <TextField autoFocus margin="dense" label="Original URL" fullWidth value={originalUrl} onChange={(e) => setOriginalUrl(e.target.value)} />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSave}>Save</Button>
        </DialogActions>
      </Dialog>
    </ThemeProvider>
  );
}

export default function Dashboard() {
  return <DashboardContent />;
}
