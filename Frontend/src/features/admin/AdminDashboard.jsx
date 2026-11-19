import { Box } from '@mui/material';
import PaymentsRounded from '@mui/icons-material/PaymentsRounded';
import PeopleAltRounded from '@mui/icons-material/PeopleAltRounded';
import PendingActionsRounded from '@mui/icons-material/PendingActionsRounded';
import StorefrontRounded from '@mui/icons-material/StorefrontRounded';
import VerifiedUserRounded from '@mui/icons-material/VerifiedUserRounded';
import PageHeader from '../../components/common/PageHeader';
import Reveal from '../../components/common/Reveal';
import StatCard from '../../components/common/StatCard';
import { ErrorState, PageSkeleton } from '../../components/common/States';
import useRequest from '../../hooks/useRequest';
import { fetchAdminDashboard } from '../../services/admin.service';
import { formatNumber } from '../../utils/format';

const STATS = [
  { key: 'users_count', label: 'عدد المستخدمين', icon: PeopleAltRounded, to: '/admin/users' },
  { key: 'stations_count', label: 'المحطات المسجلة', icon: StorefrontRounded, to: '/admin/stations' },
  { key: 'approved_stations_count', label: 'المحطات المعتمدة', icon: VerifiedUserRounded },
  { key: 'pending_stations_count', label: 'طلبات تسجيل بانتظار المراجعة', icon: PendingActionsRounded, to: '/admin/stations' },
  { key: 'pending_withdrawals_count', label: 'طلبات سحب بانتظار المراجعة', icon: PaymentsRounded, to: '/admin/withdrawals' },
];

export default function AdminDashboard() {
  const { data, error, reload } = useRequest(fetchAdminDashboard);

  if (error) return <ErrorState message={error} onRetry={reload} />;
  if (!data) return <PageSkeleton />;

  return (
    <>
      <PageHeader title="لوحة التحكم" subtitle="ملخص عام عن النظام." />
      <Box sx={{ display: 'grid', gap: 2.5, gridTemplateColumns: 'repeat(auto-fill, minmax(min(100%, 250px), 1fr))' }}>
        {STATS.map(({ key, label, icon, to }, index) => (
          <Reveal key={key} delay={index * 50}>
            <StatCard label={label} value={formatNumber(data[key])} icon={icon} to={to} />
          </Reveal>
        ))}
      </Box>
    </>
  );
}
