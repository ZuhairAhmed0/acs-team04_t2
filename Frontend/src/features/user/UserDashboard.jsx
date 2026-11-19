import { Button, Card, Stack, Typography } from '@mui/material';
import { Link } from 'react-router-dom';
import PageHeader from '../../components/common/PageHeader';
import Reveal from '../../components/common/Reveal';
import { ErrorState, PageSkeleton } from '../../components/common/States';
import TransactionList from '../../components/common/TransactionList';
import { ROLES } from '../../config/constants';
import { useAuth } from '../../context/AuthContext';
import useRequest from '../../hooks/useRequest';
import { fetchUserDashboard } from '../../services/user.service';
import LoyaltyCard from './LoyaltyCard';
import QuickActions from './QuickActions';
import WalletCard from './WalletCard';

export default function UserDashboard() {
  const { user } = useAuth();
  const { data, error, reload } = useRequest(fetchUserDashboard);

  if (error) return <ErrorState message={error} onRetry={reload} />;
  if (!data) return <PageSkeleton />;

  return (
    <>
      <PageHeader title={`أهلاً، ${user.name}`} subtitle="هذه نظرة سريعة على محفظتك." />
      <Stack spacing={3}>
        <Reveal><WalletCard wallet={data.wallet} /></Reveal>
        {user.role === ROLES.USER && (
          <Reveal delay={25}>
            <LoyaltyCard points={data.wallet.loyalty_points} gallons={data.wallet.free_gas_gallons} />
          </Reveal>
        )}
        <Reveal delay={50}><QuickActions /></Reveal>
        <Reveal delay={100}>
          <Card>
            <Stack direction="row" justifyContent="space-between" alignItems="center" flexWrap="wrap" gap={1} sx={{ px: { xs: 1.75, sm: 2.5 }, pt: 2, pb: 0.5 }}>
              <Typography variant="h6">آخر المعاملات</Typography>
              <Button component={Link} to="/user/transactions" size="small">عرض كل المعاملات</Button>
            </Stack>
            <TransactionList transactions={data.recent_transactions} />
          </Card>
        </Reveal>
      </Stack>
    </>
  );
}
