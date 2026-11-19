import { Box, Button, Card, Stack, Typography } from '@mui/material';
import StorefrontRounded from '@mui/icons-material/StorefrontRounded';
import { Link } from 'react-router-dom';
import GradientCard from '../../components/common/GradientCard';
import PageHeader from '../../components/common/PageHeader';
import Reveal from '../../components/common/Reveal';
import StatCard from '../../components/common/StatCard';
import { ErrorState, PageSkeleton } from '../../components/common/States';
import TransactionList from '../../components/common/TransactionList';
import { useAuth } from '../../context/AuthContext';
import useRequest from '../../hooks/useRequest';
import { fetchStationDashboard, fetchStationWithdrawals } from '../../services/station.service';
import { formatMoney } from '../../utils/format';
import CongestionChart from './CongestionChart';
import WithdrawalForm from './WithdrawalForm';
import WithdrawalHistory from './WithdrawalHistory';

const twoColumns = { display: 'grid', gap: 2.5, gridTemplateColumns: { xs: '1fr', md: 'repeat(2, minmax(0, 1fr))' } };

export default function StationDashboard() {
  const { user } = useAuth();
  const dashboard = useRequest(fetchStationDashboard);
  const withdrawals = useRequest(fetchStationWithdrawals);

  if (dashboard.error) return <ErrorState message={dashboard.error} onRetry={dashboard.reload} />;
  if (!dashboard.data) return <PageSkeleton />;
  const { balance, branches_count: branchesCount, branches, recent_transactions: transactions } = dashboard.data;

  return (
    <>
      <PageHeader title={user.name} subtitle="لوحة المحطة" />
      <Stack spacing={2.5}>
        <Reveal sx={{ display: 'grid', gap: 2.5, gridTemplateColumns: { xs: '1fr', md: '2fr 1fr' } }}>
          <GradientCard>
            <Typography sx={{ opacity: 0.85 }}>رصيد المحطة المشترك</Typography>
            <Typography variant="h3" sx={{ my: 1, fontSize: { xs: '2.1rem', sm: '2.9rem' } }}>{formatMoney(balance)}</Typography>
            <Typography variant="body2" sx={{ opacity: 0.85 }}>يستخدمه كل فروع المحطة، ولا يُقسَّم بينها.</Typography>
          </GradientCard>
          <StatCard label="عدد الفروع" value={branchesCount} icon={StorefrontRounded} hint="أدر الفروع وأضف فروعاً جديدة" to="/station/branches" />
        </Reveal>

        <Reveal delay={50} sx={twoColumns}>
          <CongestionChart branches={branches} />
          <WithdrawalForm balance={balance} onSubmitted={withdrawals.reload} />
        </Reveal>

        <Reveal delay={100}>
          <Card>
            <Stack direction="row" justifyContent="space-between" alignItems="center" flexWrap="wrap" gap={1} sx={{ px: { xs: 1.75, sm: 2.5 }, pt: 2, pb: 0.5 }}>
              <Typography variant="h6">آخر المعاملات</Typography>
              <Button component={Link} to="/station/transactions" size="small">عرض كل المعاملات</Button>
            </Stack>
            <TransactionList transactions={transactions} />
          </Card>
        </Reveal>

        <Reveal delay={150}>
          <Box>{withdrawals.data && <WithdrawalHistory withdrawals={withdrawals.data} />}</Box>
        </Reveal>
      </Stack>
    </>
  );
}
