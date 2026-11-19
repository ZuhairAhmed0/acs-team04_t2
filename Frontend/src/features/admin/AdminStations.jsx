import { Fragment, useState } from 'react';
import { Box, Button, Card, Dialog, DialogActions, DialogContent, DialogTitle, Divider, Stack, Tab, Tabs, Typography } from '@mui/material';
import PendingActionsRounded from '@mui/icons-material/PendingActionsRounded';
import StorefrontRounded from '@mui/icons-material/StorefrontRounded';
import ActionDialog from '../../components/common/ActionDialog';
import CongestionChip from '../../components/common/CongestionChip';
import PageHeader from '../../components/common/PageHeader';
import Reveal from '../../components/common/Reveal';
import { EmptyState, ErrorState, PageSkeleton } from '../../components/common/States';
import StatusChip from '../../components/common/StatusChip';
import useConfirmedAction from '../../hooks/useConfirmedAction';
import useRequest from '../../hooks/useRequest';
import { approveStation, deleteStation, fetchPendingStations, fetchStations, rejectStation, revokeStation } from '../../services/admin.service';
import StationDetails from './StationDetails';

const ACTIONS = {
  approve: {
    title: 'الموافقة على المحطة', confirmLabel: 'موافقة', success: 'تمت الموافقة على المحطة',
    message: (station) => `ستتمكن ${station.name} من الدخول وإدارة فروعها بعد الموافقة.`,
    run: (station) => approveStation(station.id),
  },
  reject: {
    title: 'رفض طلب التسجيل', confirmLabel: 'رفض الطلب', tone: 'error', requireReason: true, success: 'تم رفض الطلب',
    message: (station) => `اكتب سبب رفض طلب ${station.name}.`,
    run: (station, reason) => rejectStation(station.id, reason),
  },
  revoke: {
    title: 'إلغاء الموافقة', confirmLabel: 'إلغاء الموافقة', tone: 'error', success: 'تم إلغاء الموافقة على المحطة',
    message: (station) => `لن تتمكن ${station.name} من الدخول حتى تتم الموافقة عليها من جديد.`,
    run: (station) => revokeStation(station.id),
  },
  remove: {
    title: 'حذف المحطة', confirmLabel: 'حذف نهائي', tone: 'error', success: 'تم حذف المحطة',
    message: (station) => `سيتم حذف ${station.name} وكل فروعها نهائياً. لا يمكن التراجع عن هذا الإجراء.`,
    run: (station) => deleteStation(station.id),
  },
};

function PendingStationCard({ station, onAction }) {
  return (
    <Card sx={{ p: 2.5, height: '100%' }}>
      <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={1} sx={{ mb: 1 }}>
        <Typography variant="h6">{station.name}</Typography>
        <StatusChip status={station.status} />
      </Stack>
      <StationDetails station={station} />
      <Stack direction="row" spacing={1.5} sx={{ mt: 2.5 }}>
        <Button fullWidth variant="contained" onClick={() => onAction(ACTIONS.approve, station)}>موافقة</Button>
        <Button fullWidth variant="outlined" color="error" onClick={() => onAction(ACTIONS.reject, station)}>رفض</Button>
      </Stack>
    </Card>
  );
}

function PendingTab({ request, onAction }) {
  const { data, error, reload } = request;
  if (error) return <ErrorState message={error} onRetry={reload} />;
  if (!data) return <PageSkeleton blocks={2} />;
  if (!data.length) return <EmptyState icon={PendingActionsRounded} title="لا توجد طلبات جديدة" description="طلبات تسجيل المحطات الجديدة ستظهر هنا للمراجعة." />;

  return (
    <Box sx={{ display: 'grid', gap: 2.5, gridTemplateColumns: 'repeat(auto-fill, minmax(min(100%, 380px), 1fr))' }}>
      {data.map((station, index) => (
        <Reveal key={station.id} delay={Math.min(index, 5) * 50}>
          <PendingStationCard station={station} onAction={onAction} />
        </Reveal>
      ))}
    </Box>
  );
}

function ExistingTab({ request, onAction, onDetails }) {
  const { data, error, reload } = request;
  if (error) return <ErrorState message={error} onRetry={reload} />;
  if (!data) return <PageSkeleton blocks={2} />;
  const stations = data.filter((station) => station.status !== 'pending');
  if (!stations.length) return <EmptyState icon={StorefrontRounded} title="لا توجد محطات بعد" />;

  return (
    <Reveal>
      <Card>
        {stations.map((station, index) => {
          const branches = station.branches ?? [];
          const avgCongestion = branches.length
            ? Math.round(branches.reduce((sum, b) => sum + (b.congestion_percentage ?? 0), 0) / branches.length)
            : 0;
          return (
            <Fragment key={station.id}>
              {index > 0 && <Divider />}
              <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ xs: 'flex-start', md: 'center' }} sx={{ px: 2.5, py: 2 }}>
                <Box sx={{ flex: 1, minWidth: 0 }}>
                  <Typography fontWeight={700}>{station.name}</Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ overflowWrap: 'break-word' }}>{station.location} — {branches.length} فروع</Typography>
                </Box>
                <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap>
                  <StatusChip status={station.status} />
                  {branches.length > 0 && <CongestionChip percentage={avgCongestion} label="متوسط الازدحام" />}
                </Stack>
                <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ width: { xs: '100%', md: 'auto' } }}>
                  <Button size="small" variant="outlined" onClick={() => onDetails(station)}>التفاصيل</Button>
                  {station.status === 'active' && (
                    <Button size="small" variant="outlined" color="error" onClick={() => onAction(ACTIONS.revoke, station)}>إلغاء الموافقة</Button>
                  )}
                  <Button size="small" color="error" onClick={() => onAction(ACTIONS.remove, station)}>حذف</Button>
                </Stack>
              </Stack>
            </Fragment>
          );
        })}
      </Card>
    </Reveal>
  );
}

export default function AdminStations() {
  const [tab, setTab] = useState('pending');
  const [detailsOf, setDetailsOf] = useState(null);
  const pending = useRequest(fetchPendingStations);
  const existing = useRequest(fetchStations);
  const controller = useConfirmedAction(() => {
    pending.reload();
    existing.reload();
  });

  return (
    <>
      <PageHeader title="المحطات" subtitle="راجع طلبات التسجيل الجديدة وأدر المحطات الحالية." />
      <Tabs
        value={tab}
        onChange={(_, value) => setTab(value)}
        variant="scrollable"
        scrollButtons="auto"
        allowScrollButtonsMobile
        sx={{ mb: 3, borderBottom: '1px solid', borderColor: 'divider' }}
      >
        <Tab value="pending" label={`طلبات التسجيل${pending.data ? ` (${pending.data.length})` : ''}`} />
        <Tab value="existing" label="المحطات الحالية" />
      </Tabs>

      {tab === 'pending' ? (
        <PendingTab request={pending} onAction={controller.request} />
      ) : (
        <ExistingTab request={existing} onAction={controller.request} onDetails={setDetailsOf} />
      )}

      <Dialog open={Boolean(detailsOf)} onClose={() => setDetailsOf(null)} fullWidth maxWidth="sm">
        {detailsOf && (
          <>
            <DialogTitle sx={{ fontWeight: 800 }}>
              <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={1}>
                <span>{detailsOf.name}</span>
                <StatusChip status={detailsOf.status} />
              </Stack>
            </DialogTitle>
            <DialogContent><StationDetails station={detailsOf} /></DialogContent>
            <DialogActions sx={{ px: 3, pb: 2.5 }}><Button onClick={() => setDetailsOf(null)}>إغلاق</Button></DialogActions>
          </>
        )}
      </Dialog>
      <ActionDialog controller={controller} />
    </>
  );
}
