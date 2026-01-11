import { Box, Button, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';

export default function Landing() {
  const navigate = useNavigate();

  return (
    <Box sx={{ p: 4, textAlign: 'center' }}>
      <Typography variant="h3" gutterBottom>
        Welcome to the CMS
      </Typography>
      <Typography variant="body1" sx={{ mb: 3 }}>
        Manage your content efficiently.
      </Typography>
      <Button variant="contained" onClick={() => navigate('/login')}>
        Get Started
      </Button>
    </Box>
  );
}
