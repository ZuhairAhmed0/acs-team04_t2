import { Route, Routes } from 'react-router-dom';
import { ROLES } from '../config/constants';
import AppShell from '../components/layout/AppShell';
import { HomeRedirect, ProtectedRoute, PublicOnlyRoute } from '../components/layout/RouteGuards';
import AdminDashboard from '../features/admin/AdminDashboard';
import AdminStations from '../features/admin/AdminStations';
import AdminUsers from '../features/admin/AdminUsers';
import AdminWithdrawals from '../features/admin/AdminWithdrawals';
import LoginPage from '../features/auth/LoginPage';
import RegisterPage from '../features/auth/RegisterPage';
import StationBranches from '../features/station/StationBranches';
import StationDashboard from '../features/station/StationDashboard';
import StationTransactions from '../features/station/StationTransactions';
import UserDashboard from '../features/user/UserDashboard';
import UserStations from '../features/user/UserStations';
import UserTransactions from '../features/user/UserTransactions';
import UserTransfers from '../features/user/UserTransfers';

export default function App() {
  return (
    <Routes>
      <Route element={<PublicOnlyRoute />}>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
      </Route>

      <Route element={<ProtectedRoute roles={[ROLES.USER]} />}>
        <Route element={<AppShell />}>
          <Route path="/user" element={<UserDashboard />} />
          <Route path="/user/stations" element={<UserStations />} />
          <Route path="/user/transfers" element={<UserTransfers />} />
          <Route path="/user/transactions" element={<UserTransactions />} />
        </Route>
      </Route>

      <Route element={<ProtectedRoute roles={[ROLES.STATION]} />}>
        <Route element={<AppShell />}>
          <Route path="/station" element={<StationDashboard />} />
          <Route path="/station/branches" element={<StationBranches />} />
          <Route path="/station/transactions" element={<StationTransactions />} />
        </Route>
      </Route>

      <Route element={<ProtectedRoute roles={[ROLES.ADMIN]} />}>
        <Route element={<AppShell />}>
          <Route path="/admin" element={<AdminDashboard />} />
          <Route path="/admin/users" element={<AdminUsers />} />
          <Route path="/admin/stations" element={<AdminStations />} />
          <Route path="/admin/withdrawals" element={<AdminWithdrawals />} />
        </Route>
      </Route>

      <Route path="*" element={<HomeRedirect />} />
    </Routes>
  );
}
