import { useMemo, useState } from 'react';
import { Box, ButtonBase, Card, Chip, Dialog, DialogContent, DialogTitle, Divider, Stack, Typography } from '@mui/material';
import LocalGasStationRounded from '@mui/icons-material/LocalGasStationRounded';
import CongestionChip from '../../components/common/CongestionChip';
import FilterChips from '../../components/common/FilterChips';
import FuelStatusChip from '../../components/common/FuelStatusChip';
import PageHeader from '../../components/common/PageHeader';
import ReceiptDialog from '../../components/common/ReceiptDialog';
import Reveal from '../../components/common/Reveal';
import { EmptyState, ErrorState, PageSkeleton } from '../../components/common/States';
import { CITIES, ROLES, cityLabel } from '../../config/constants';
import { useAuth } from '../../context/AuthContext';
import { useToast } from '../../context/ToastContext';
import useRequest from '../../hooks/useRequest';
import { fetchTransfers } from '../../services/transfer.service';
import { fetchUserStations, fetchUserWallet } from '../../services/user.service';
import { tokens } from '../../theme/tokens';
import { FuelPurchaseForm } from './TransferForms';

/** The API returns one row per branch; group them under their station. */
function groupByStation(rows) {
  const stations = new Map();
  rows.forEach((row) => {
    if (!stations.has(row.station_id)) stations.set(row.station_id, { id: row.station_id, name: row.station_name, branches: [] });
    stations.get(row.station_id).branches.push(row);
  });
  return [...stations.values()];
}

function StationCard({ station, onSelectBranch }) {
  return (
    <Card sx={{ p: 2.5, height: '100%' }}>
      <Stack direction="row" alignItems="center" spacing={1.5} sx={{ mb: 1.5 }}>
        <LocalGasStationRounded color="primary" />
        <Typography variant="h6">{station.name}</Typography>
      </Stack>
      {station.branches.map((branch, index) => {
        const isEmpty = branch.fuel_status === 'empty';
        return (
          <Box key={branch.branch_id}>
            {index > 0 && <Divider />}
            <ButtonBase
              disabled={isEmpty}
              onClick={() => onSelectBranch({ station, branch })}
              sx={{
                width: '100%',
                display: 'block',
                textAlign: 'start',
                borderRadius: 2,
                px: 1,
                py: 1.5,
                transition: 'background-color .2s',
                opacity: isEmpty ? 0.55 : 1,
                '&:hover': isEmpty ? undefined : { bgcolor: tokens.redSoft },
              }}
            >
              <Stack spacing={1}>
                <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={1}>
                  <Typography fontWeight={700} sx={{ minWidth: 0 }}>{cityLabel(branch.branch_city)}</Typography>
                  {!isEmpty && <Chip label="شراء وقود" size="small" color="primary" variant="outlined" sx={{ flexShrink: 0 }} />}
                </Stack>
                <Typography variant="body2" color="text.secondary" sx={{ overflowWrap: 'break-word' }}>
                  {branch.branch_address}
                </Typography>
                <Stack direction="row" alignItems="center" spacing={1} flexWrap="wrap" useFlexGap>
                  <FuelStatusChip status={branch.fuel_status} />
                  <CongestionChip percentage={branch.congestion_percentage} />
                </Stack>
              </Stack>
            </ButtonBase>
          </Box>
        );
      })}
    </Card>
  );
}

export default function UserStations() {
  const { user } = useAuth();
  const toast = useToast();
  const [city, setCity] = useState('');
  const [selected, setSelected] = useState(null);
  const [receipt, setReceipt] = useState(null);
  const [sessionPurchases, setSessionPurchases] = useState([]);
  const { data, error, reload } = useRequest(() => fetchUserStations(city || undefined), [city]);
  const { data: wallet, reload: reloadWallet } = useRequest(fetchUserWallet);
  const { data: transfers, reload: reloadTransfers } = useRequest(fetchTransfers);

  const recentPurchases = useMemo(
    () => [
      ...sessionPurchases,
      ...(transfers ?? []).filter((tx) => tx.type === 'fuel_purchase' && tx.status === 'completed'),
    ],
    [sessionPurchases, transfers],
  );

  const handleSuccess = (nextReceipt, submittedPayload) => {
    const branchId = submittedPayload?.branch_id ?? selected?.branch?.branch_id;
    if (branchId && nextReceipt?.reference && !nextReceipt.idempotent_replay) {
      setSessionPurchases((prev) => [
        {
          type: 'fuel_purchase',
          status: nextReceipt.status ?? 'completed',
          branch_id: branchId,
          amount: nextReceipt.amount,
          loyalty_points_used: nextReceipt.loyalty_points_used ?? 0,
          reference: nextReceipt.reference,
          created_at: new Date().toISOString(),
        },
        ...prev,
      ]);
    }
    if (nextReceipt?.idempotent_replay) {
      toast.warning('تم استرجاع نتيجة العملية السابقة دون خصم مبلغ إضافي');
    } else if (nextReceipt?.duplicate_purchase) {
      toast.warning('تنبيه: هذه عملية شراء مطابقة لعملية سابقة حديثة');
    }
    setSelected(null);
    setReceipt(nextReceipt);
    reloadWallet();
    reloadTransfers();
  };

  return (
    <>
      <PageHeader title="المحطات" subtitle="شوف المحطات وحالة توفر الوقود في كل فرع، واضغط على أي فرع متاح لشراء الوقود." />

      <FilterChips value={city} options={CITIES} allLabel="كل المدن" ariaLabel="تصفية حسب المدينة" onChange={setCity} />

      {error && <ErrorState message={error} onRetry={reload} />}
      {!error && !data && <PageSkeleton blocks={2} />}
      {data && !data.length && (
        <EmptyState icon={LocalGasStationRounded} title="لا توجد محطات في هذه المدينة" description="جرّب مدينة أخرى أو اعرض كل المدن." />
      )}
      {data && data.length > 0 && (
        <Box sx={{ display: 'grid', gap: 2.5, gridTemplateColumns: 'repeat(auto-fill, minmax(min(100%, 340px), 1fr))' }}>
          {groupByStation(data).map((station, index) => (
            <Reveal key={station.id} delay={Math.min(index, 5) * 50}>
              <StationCard station={station} onSelectBranch={setSelected} />
            </Reveal>
          ))}
        </Box>
      )}

      <Dialog open={Boolean(selected)} onClose={() => setSelected(null)} fullWidth maxWidth="xs">
        {selected && (
          <>
            <DialogTitle sx={{ fontWeight: 800 }}>شراء وقود</DialogTitle>
            <DialogContent>
              <Stack spacing={2} sx={{ pt: 0.5 }}>
                <Box>
                  <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={1} flexWrap="wrap" useFlexGap>
                    <Typography fontWeight={700}>{selected.station.name}</Typography>
                    <Stack direction="row" spacing={1} alignItems="center">
                      <FuelStatusChip status={selected.branch.fuel_status} />
                      <CongestionChip percentage={selected.branch.congestion_percentage} />
                    </Stack>
                  </Stack>
                  <Typography variant="body2" color="text.secondary">
                    {cityLabel(selected.branch.branch_city)} — {selected.branch.branch_address}
                  </Typography>
                  {user?.role === ROLES.USER && (
                    <Typography variant="caption" color="primary.main" sx={{ display: 'block', mt: 0.5 }}>
                      +5 نقاط ولاء عند الشراء من المحفظة، أو استبدل نقاطك بمضاعفات 10 (كل 10 نقاط = 1,000 ج.س)
                    </Typography>
                  )}
                </Box>
                <FuelPurchaseForm
                  key={selected.branch.branch_id}
                  branch={selected.branch}
                  recentPurchases={recentPurchases}
                  availableBalance={wallet?.available_balance}
                  allocatedBalance={wallet?.allocated_balance}
                  loyaltyPoints={wallet?.loyalty_points}
                  loyaltyPointValue={wallet?.loyalty_point_value}
                  isRegularUser={user?.role === ROLES.USER}
                  onCancel={() => setSelected(null)}
                  onSuccess={handleSuccess}
                />
              </Stack>
            </DialogContent>
          </>
        )}
      </Dialog>

      <ReceiptDialog receipt={receipt} onClose={() => setReceipt(null)} />
    </>
  );
}
