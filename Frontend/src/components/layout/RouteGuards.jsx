import { Box, CircularProgress } from '@mui/material';
import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { ROLE_HOME } from '../../config/constants';
import { useAuth } from '../../context/AuthContext';

function AuthLoading() {
  return (
    <Box sx={{ minHeight: '100vh', display: 'grid', placeItems: 'center' }}>
      <CircularProgress />
    </Box>
  );
}

/** Only lets signed-in accounts with one of `roles` through. */
export function ProtectedRoute({ roles }) {
  const { user, isAuthenticated, loading } = useAuth();
  const location = useLocation();

  if (loading) return <AuthLoading />;
  if (!isAuthenticated) return <Navigate to="/login" replace state={{ from: location }} />;
  if (roles && !roles.includes(user.role)) return <Navigate to={ROLE_HOME[user.role]} replace />;
  return <Outlet />;
}

/** Login and registration are hidden from accounts that are already signed in. */
export function PublicOnlyRoute() {
  const { user, isAuthenticated, loading } = useAuth();
  if (loading) return <AuthLoading />;
  return isAuthenticated ? <Navigate to={ROLE_HOME[user.role]} replace /> : <Outlet />;
}

export function HomeRedirect() {
  const { user, isAuthenticated, loading } = useAuth();
  if (loading) return <AuthLoading />;
  return <Navigate to={isAuthenticated ? ROLE_HOME[user.role] : '/login'} replace />;
}
