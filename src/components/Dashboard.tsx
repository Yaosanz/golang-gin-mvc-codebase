import React, { useEffect, useState, useCallback } from 'react';
import { Container, Box, Typography, Grid, Card, CardContent, CardMedia, Button, Paper, Divider, Alert, CircularProgress, Stack, Chip } from '@mui/material';
import { Link as LinkIcon, Dashboard as DashboardIcon, People as PeopleIcon, MoreHorizontal as MoreHorizontalIcon, TrendingUp as TrendingUpIcon, Clock as ClockIcon } from '@mui/icons-material';
import { useAuth } from '../context/AuthContext.tsx';
import { useNavigate } from 'react-router-dom';
import AppLayout from './AppLayout.tsx';
import { getAll as getAllLinks } from '../api/shortenlink.ts';

interface DashboardStats {
  totalLinks: number;
  totalUsers: number;
  recentLinks: any[];
  isLoading: boolean;
}

export default function Dashboard() {
  const { user, token } = useAuth();
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats>({
    totalLinks: 0,
    totalUsers: 0,
    recentLinks: [],
    isLoading: true,
  });

  const fetchStats = useCallback(async () => {
    try {
      setStats((prev) => ({ ...prev, isLoading: true }));
      const linksData = await getAllLinks();
      setStats((prev) => ({
        ...prev,
        totalLinks: Array.isArray(linksData.data) ? linksData.data.length : linksData.data.contents?.length || 0,
        recentLinks: Array.isArray(linksData.data) ? linksData.data.slice(0, 5) : linksData.data.contents?.slice(0, 5) || [],
        isLoading: false,
      }));
    } catch (err) {
      console.error('Failed to fetch stats:', err);
      setStats((prev) => ({ ...prev, isLoading: false }));
    }
  }, []);

  useEffect(() => {
    if (token) {
      fetchStats();
    }
  }, [token, fetchStats]);

  if (!token) {
    navigate('/login');
    return null;
  }

  return (
    <AppLayout>
      <Container maxWidth="lg">
        {/* Welcome Section */}
        <Box sx={{ mb: 4 }}>
          <Box
            sx={{
              background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
              borderRadius: 3,
              p: 4,
              color: 'white',
            }}
          >
            <Typography variant="h4" sx={{ fontWeight: 'bold', mb: 1 }}>
              Welcome back, {user?.name}! 👋
            </Typography>
            <Typography variant="body1" sx={{ opacity: 0.9 }}>
              Here's what's happening with your links today.
            </Typography>
          </Box>
        </Box>

        {/* Stats Cards */}
        <Grid container spacing={3} sx={{ mb: 4 }}>
          {/* Total Links Card */}
          <Grid item xs={12} sm={6} md={4}>
            <Card
              sx={{
                height: '100%',
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                color: 'white',
                cursor: 'pointer',
                transition: 'transform 0.2s, box-shadow 0.2s',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: '0 12px 24px rgba(102, 126, 234, 0.4)',
                },
              }}
              onClick={() => navigate('/links')}
            >
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <Box>
                    <Typography color="inherit" sx={{ fontSize: '0.9rem', opacity: 0.9, mb: 1 }}>
                      Total Links
                    </Typography>
                    <Typography variant="h4" sx={{ fontWeight: 'bold' }}>
                      {stats.isLoading ? '-' : stats.totalLinks}
                    </Typography>
                  </Box>
                  <LinkIcon sx={{ fontSize: '2.5rem', opacity: 0.3 }} />
                </Box>
                <Typography variant="caption" sx={{ opacity: 0.8, mt: 1, display: 'block' }}>
                  Click to view all links →
                </Typography>
              </CardContent>
            </Card>
          </Grid>

          {/* Admin Only - Total Users Card */}
          {user?.roles?.some((r) => r.name === 'admin' || r === 'admin') && (
            <Grid item xs={12} sm={6} md={4}>
              <Card
                sx={{
                  height: '100%',
                  background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
                  color: 'white',
                  cursor: 'pointer',
                  transition: 'transform 0.2s, box-shadow 0.2s',
                  '&:hover': {
                    transform: 'translateY(-4px)',
                    boxShadow: '0 12px 24px rgba(245, 87, 108, 0.4)',
                  },
                }}
                onClick={() => navigate('/users')}
              >
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <Box>
                      <Typography color="inherit" sx={{ fontSize: '0.9rem', opacity: 0.9, mb: 1 }}>
                        Total Users
                      </Typography>
                      <Typography variant="h4" sx={{ fontWeight: 'bold' }}>
                        -
                      </Typography>
                    </Box>
                    <PeopleIcon sx={{ fontSize: '2.5rem', opacity: 0.3 }} />
                  </Box>
                  <Typography variant="caption" sx={{ opacity: 0.8, mt: 1, display: 'block' }}>
                    Click to manage users →
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          )}

          {/* Quick Actions Card */}
          <Grid item xs={12} sm={6} md={user?.roles?.some((r) => r.name === 'admin') ? 4 : 6}>
            <Card
              sx={{
                height: '100%',
                background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
                color: 'white',
                cursor: 'pointer',
                transition: 'transform 0.2s, box-shadow 0.2s',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: '0 12px 24px rgba(79, 172, 254, 0.4)',
                },
              }}
              onClick={() => navigate('/links/create')}
            >
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <Box>
                    <Typography color="inherit" sx={{ fontSize: '0.9rem', opacity: 0.9, mb: 1 }}>
                      Create New Link
                    </Typography>
                    <Typography variant="h4" sx={{ fontWeight: 'bold' }}>
                      +
                    </Typography>
                  </Box>
                  <TrendingUpIcon sx={{ fontSize: '2.5rem', opacity: 0.3 }} />
                </Box>
                <Typography variant="caption" sx={{ opacity: 0.8, mt: 1, display: 'block' }}>
                  Start shortening →
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Recent Links Section */}
        <Paper sx={{ p: 3, borderRadius: 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6" sx={{ fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: 1 }}>
              <ClockIcon /> Recent Links
            </Typography>
            <Button variant="text" onClick={() => navigate('/links')}>
              View All →
            </Button>
          </Box>

          <Divider sx={{ mb: 2 }} />

          {stats.isLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 3 }}>
              <CircularProgress />
            </Box>
          ) : stats.recentLinks.length === 0 ? (
            <Alert severity="info">No links yet. Create your first shortened link!</Alert>
          ) : (
            <Stack spacing={2}>
              {stats.recentLinks.map((link, index) => (
                <Box
                  key={link.id || index}
                  sx={{
                    p: 2,
                    border: '1px solid',
                    borderColor: 'divider',
                    borderRadius: 1.5,
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    '&:hover': {
                      backgroundColor: 'action.hover',
                    },
                  }}
                >
                  <Box sx={{ flex: 1, minWidth: 0 }}>
                    <Chip label={link.code} variant="outlined" sx={{ fontFamily: 'monospace', fontWeight: 'bold', mb: 0.5 }} />
                    <Typography
                      variant="body2"
                      sx={{
                        color: 'textSecondary',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                      }}
                    >
                      {link.url}
                    </Typography>
                    <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mt: 0.5 }}>
                      {new Date(link.created_at).toLocaleDateString('id-ID', {
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', gap: 1 }}>
                    <Button size="small" variant="outlined">
                      Copy
                    </Button>
                    <Button size="small" variant="outlined" color="info">
                      View
                    </Button>
                  </Box>
                </Box>
              ))}
            </Stack>
          )}
        </Paper>

        {/* User Info Section */}
        <Box sx={{ mt: 4, p: 3, backgroundColor: 'background.paper', borderRadius: 2, border: '1px solid', borderColor: 'divider' }}>
          <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2 }}>
            Your Account Information
          </Typography>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="textSecondary">
                Username
              </Typography>
              <Typography variant="body1" sx={{ fontWeight: 500 }}>
                {user?.username}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="textSecondary">
                Email
              </Typography>
              <Typography variant="body1" sx={{ fontWeight: 500 }}>
                {user?.email}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="textSecondary">
                Full Name
              </Typography>
              <Typography variant="body1" sx={{ fontWeight: 500 }}>
                {user?.name}
              </Typography>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Typography variant="body2" color="textSecondary">
                Roles
              </Typography>
              <Box sx={{ display: 'flex', gap: 1, mt: 0.5, flexWrap: 'wrap' }}>
                {user?.roles?.map((role, idx) => (
                  <Chip key={idx} label={role.name || role} size="small" color={role.name === 'admin' || role === 'admin' ? 'error' : 'default'} variant="outlined" />
                ))}
              </Box>
            </Grid>
          </Grid>
          <Button variant="outlined" sx={{ mt: 2 }} onClick={() => navigate('/profile')}>
            Edit Profile
          </Button>
        </Box>
      </Container>
    </AppLayout>
  );
}
