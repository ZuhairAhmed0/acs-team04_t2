import { Avatar, BottomNavigation, BottomNavigationAction, Box, Container, Divider, IconButton, List, ListItemButton, ListItemIcon, ListItemText, Paper, Stack, Typography, useMediaQuery } from '@mui/material';
import { useTheme } from '@mui/material/styles';
import LogoutRounded from '@mui/icons-material/LogoutRounded';
import { Link, NavLink, Outlet, useLocation } from 'react-router-dom';
import { ROLE_LABELS } from '../../config/constants';
import { NAVIGATION } from '../../config/navigation';
import { useAuth } from '../../context/AuthContext';
import { tokens } from '../../theme/tokens';
import Brand from '../brand/Brand';

const SIDEBAR_WIDTH = 272;

function UserBadge({ user, onLogout }) {
  return (
    <Stack direction="row" alignItems="center" spacing={1.5}>
      <Avatar sx={{ bgcolor: tokens.redSoft, color: 'primary.main', fontWeight: 800 }}>{user.name[0]}</Avatar>
      <Box sx={{ flex: 1, minWidth: 0 }}>
        <Typography fontWeight={700} noWrap>{user.name}</Typography>
        <Typography variant="caption" color="text.secondary">{ROLE_LABELS[user.role]}</Typography>
      </Box>
      <IconButton aria-label="تسجيل الخروج" onClick={onLogout}>
        <LogoutRounded />
      </IconButton>
    </Stack>
  );
}

function Sidebar({ items, user, onLogout }) {
  return (
    <Box
      component="aside"
      sx={{
        position: 'fixed',
        insetBlock: 0,
        insetInlineStart: 0,
        width: SIDEBAR_WIDTH,
        display: 'flex',
        flexDirection: 'column',
        p: 2.5,
        bgcolor: 'background.paper',
        borderInlineEnd: '1px solid',
        borderColor: 'divider',
      }}
    >
      <Box sx={{ px: 1, py: 1, mb: 3 }}>
        <Brand />
      </Box>
      <List component="nav" sx={{ flex: 1 }}>
        {items.map(({ to, label, icon: Icon, end }) => (
          <ListItemButton
            key={to}
            component={NavLink}
            to={to}
            end={end}
            sx={{
              borderRadius: '14px',
              mb: 0.5,
              py: 1.2,
              color: 'text.secondary',
              '&.active': { bgcolor: tokens.redSoft, color: 'primary.main', '& .MuiListItemIcon-root': { color: 'primary.main' } },
            }}
          >
            <ListItemIcon sx={{ minWidth: 40, color: 'inherit' }}>
              <Icon />
            </ListItemIcon>
            <ListItemText primary={label} primaryTypographyProps={{ fontWeight: 700 }} />
          </ListItemButton>
        ))}
      </List>
      <Divider sx={{ mb: 2 }} />
      <UserBadge user={user} onLogout={onLogout} />
    </Box>
  );
}

function MobileNavigation({ items }) {
  const { pathname } = useLocation();
  const activeIndex = items.findIndex(({ to, end }) => (end ? pathname === to : pathname.startsWith(to)));

  return (
    <Paper
      elevation={0}
      sx={{
        position: 'fixed',
        insetInline: 0,
        bottom: 0,
        zIndex: 10,
        borderTop: '1px solid',
        borderColor: 'divider',
        borderRadius: 0,
        pb: 'env(safe-area-inset-bottom)',
      }}
    >
      <BottomNavigation showLabels value={activeIndex} sx={{ height: 64 }}>
        {items.map(({ to, label, icon: Icon }) => (
          <BottomNavigationAction
            key={to}
            component={Link}
            to={to}
            label={label}
            icon={<Icon />}
            sx={{
              minWidth: 0,
              px: 0.5,
              '& .MuiBottomNavigationAction-label': {
                fontSize: '0.72rem',
                lineHeight: 1.2,
                whiteSpace: 'nowrap',
                mt: 0.25,
                '&.Mui-selected': { fontSize: '0.75rem', fontWeight: 700 },
              },
            }}
          />
        ))}
      </BottomNavigation>
    </Paper>
  );
}

export default function AppShell() {
  const theme = useTheme();
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'));
  const { user, logout } = useAuth();
  const items = NAVIGATION[user.role];

  return (
    <Box sx={{ minHeight: '100dvh', overflowX: 'clip' }}>
      {isDesktop ? (
        <Sidebar items={items} user={user} onLogout={logout} />
      ) : (
        <>
          <Stack
            component="header"
            direction="row"
            justifyContent="space-between"
            alignItems="center"
            spacing={1}
            sx={{
              position: 'sticky',
              top: 0,
              zIndex: 9,
              px: 2,
              pt: 'calc(10px + env(safe-area-inset-top))',
              pb: 1.25,
              bgcolor: 'background.paper',
              borderBottom: '1px solid',
              borderColor: 'divider',
            }}
          >
            <Brand small />
            <Stack direction="row" alignItems="center" spacing={0.5} sx={{ minWidth: 0 }}>
              <Typography variant="caption" color="text.secondary" noWrap sx={{ maxWidth: 130, fontWeight: 700 }}>
                {user.name}
              </Typography>
              <IconButton aria-label="تسجيل الخروج" onClick={logout}>
                <LogoutRounded />
              </IconButton>
            </Stack>
          </Stack>
          <MobileNavigation items={items} />
        </>
      )}

      <Box
        component="main"
        sx={{
          marginInlineStart: isDesktop ? `${SIDEBAR_WIDTH}px` : 0,
          pb: isDesktop ? 4 : 'calc(84px + env(safe-area-inset-bottom))',
          minWidth: 0,
        }}
      >
        <Container maxWidth="lg" sx={{ py: { xs: 2, sm: 3, md: 4 }, px: { xs: 1.75, sm: 3 } }}>
          <Outlet />
        </Container>
      </Box>
    </Box>
  );
}
