import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { useUser } from '../providers/UserContext';
import { useRouter } from 'next/router';

export default function Navbar() {
  const { user, logout } = useUser();
  const router = useRouter();

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          OpenTelemetry Demo
        </Typography>
        <Box>
          {user ? (
            <>
              <Typography component="span" sx={{ mr: 2 }}>
                欢迎, {user.username}
              </Typography>
              <Button color="inherit" onClick={() => logout()}>
                登出
              </Button>
            </>
          ) : (
            <Button color="inherit" onClick={() => router.push('/login')}>
              登录
            </Button>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  );
} 