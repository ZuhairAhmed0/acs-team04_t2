import { Fragment, useState } from 'react';
import { Box, Button, Card, Divider, Stack, Typography } from '@mui/material';
import PaymentsRounded from '@mui/icons-material/PaymentsRounded';
import ActionDialog from '../../components/common/ActionDialog';
import FilterChips from '../../components/common/FilterChips';
import PageHeader from '../../components/common/PageHeader';
import Reveal from '../../components/common/Reveal';
import { EmptyState, ErrorState, PageSkeleton } from '../../components/common/States';
import StatusChip from '../../components/common/StatusChip';
import { WITHDRAWAL_STATUSES, bankLabel, statusMeta } from '../../config/constants';
import { useToast } from '../../context/ToastContext';
import useConfirmedAction from '../../hooks/useConfirmedAction';
import useRequest from '../../hooks/useRequest';
import useSubmit from '../../hooks/useSubmit';
import { approveWithdrawal, completeWithdrawal, fetchWithdrawals, rejectWithdrawal } from '../../services/admin.service';
import { formatDate, formatMoney } from '../../utils/format';
import ExternalTransferDialog from './ExternalTransferDialog';

const ACTIONS = {
  approve: {
    title: 'الموافقة على طلب السحب', confirmLabel: 'موافقة', success: 'تمت الموافقة على الطلب',
    message: (w) => `الموافقة على سحب ${formatMoney(w.amount)} لصالح ${w.station_name}؟`,
    run: (w) => approveWithdrawal(w.id),
  },
  reject: {
    title: 'رفض طلب السحب', confirmLabel: 'رفض الطلب', tone: 'error', requireReason: true, success: 'تم رفض الطلب',
    message: (w) => `اكتب سبب رفض طلب ${w.station_name} بمبلغ ${formatMoney(w.amount)}.`,
    run: (w, reason) => rejectWithdrawal(w.id, reason),
  },
};

function WithdrawalRow({ withdrawal, onAction, onTransfer }) {
  const { status } = withdrawal;
  return (
    <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} alignItems={{ xs: 'flex-start', md: 'center' }} sx={{ px: 2.5, py: 2 }}>
      <Box sx={{ flex: 1 }}>
        <Typography fontWeight={700}>{withdrawal.station_name}</Typography>
        <Typography variant="body2" color="text.secondary">
          {bankLabel(withdrawal.bank_name)} — {withdrawal.bank_account_number}
        </Typography>
        <Typography variant="caption" color="text.secondary">{formatDate(withdrawal.created_at)}</Typography>
        {withdrawal.external_reference && (
          <Typography variant="body2" color="text.secondary">مرجع التحويل: {withdrawal.external_reference}</Typography>
        )}
        {withdrawal.rejection_reason && <Typography variant="body2" color="error.main">السبب: {withdrawal.rejection_reason}</Typography>}
      </Box>
      <Typography variant="h6" sx={{ whiteSpace: 'nowrap' }}>{formatMoney(withdrawal.amount)}</Typography>
      <StatusChip status={status} />
      <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ width: { xs: '100%', md: 'auto' }, minWidth: { md: 170 } }}>
        {status === 'pending' && (
          <>
            <Button size="small" variant="contained" onClick={() => onAction(ACTIONS.approve, withdrawal)}>موافقة</Button>
            <Button size="small" variant="outlined" color="error" onClick={() => onAction(ACTIONS.reject, withdrawal)}>رفض</Button>
          </>
        )}
        {status === 'approved' && <Button size="small" variant="contained" onClick={() => onTransfer(withdrawal)}>تنفيذ التحويل</Button>}
      </Stack>
    </Stack>
  );
}

export default function AdminWithdrawals() {
  const toast = useToast();
  const [statusFilter, setStatusFilter] = useState('');
  const [transferring, setTransferring] = useState(null);
  const { data, error, reload } = useRequest(fetchWithdrawals);
  const controller = useConfirmedAction(reload);
  const { submit: complete, submitting: completing } = useSubmit(async (withdrawal, externalRef) => {
    await completeWithdrawal(withdrawal.id, externalRef);
    return true;
  });

  const handleComplete = async (withdrawal, externalRef) => {
    if (await complete(withdrawal, externalRef)) {
      toast.success('تم إكمال طلب السحب');
      setTransferring(null);
      reload();
    }
  };

  const withdrawals = [...(data ?? [])]
    .filter((entry) => !statusFilter || entry.status === statusFilter)
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at) || (b.id ?? 0) - (a.id ?? 0));

  return (
    <>
      <PageHeader title="طلبات السحب" subtitle="راجع طلبات السحب، ثم أكمل الطلب بعد تنفيذ التحويل." />
      <FilterChips
        value={statusFilter}
        options={WITHDRAWAL_STATUSES.map((value) => ({ value, label: statusMeta(value).label }))}
        ariaLabel="تصفية حسب الحالة"
        onChange={setStatusFilter}
      />

      {error && <ErrorState message={error} onRetry={reload} />}
      {!error && !data && <PageSkeleton blocks={2} />}
      {data && (
        <Reveal>
          <Card>
            {!withdrawals.length && <EmptyState icon={PaymentsRounded} title="لا توجد طلبات" description="غيّر الفلتر لعرض طلبات بحالة أخرى." />}
            {withdrawals.map((withdrawal, index) => (
              <Fragment key={withdrawal.id}>
                {index > 0 && <Divider />}
                <WithdrawalRow withdrawal={withdrawal} onAction={controller.request} onTransfer={setTransferring} />
              </Fragment>
            ))}
          </Card>
        </Reveal>
      )}

      <ActionDialog controller={controller} />
      <ExternalTransferDialog withdrawal={transferring} loading={completing} onConfirm={handleComplete} onClose={() => setTransferring(null)} />
    </>
  );
}
