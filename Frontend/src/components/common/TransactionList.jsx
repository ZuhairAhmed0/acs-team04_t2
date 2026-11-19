import { Fragment, useEffect, useState } from 'react';
import { Box, Button, ButtonBase, Dialog, DialogActions, DialogContent, Divider, LinearProgress, Stack, Typography } from '@mui/material';
import AddCardRounded from '@mui/icons-material/AddCardRounded';
import ChevronLeftRounded from '@mui/icons-material/ChevronLeftRounded';
import LocalGasStationRounded from '@mui/icons-material/LocalGasStationRounded';
import PaymentsRounded from '@mui/icons-material/PaymentsRounded';
import ReceiptLongRounded from '@mui/icons-material/ReceiptLongRounded';
import SwapHorizRounded from '@mui/icons-material/SwapHorizRounded';
import { ROLES, TRANSACTION_TYPES } from '../../config/constants';
import { useAuth } from '../../context/AuthContext';
import { fetchTransferById } from '../../services/transfer.service';
import { tokens } from '../../theme/tokens';
import { formatDateTime, formatMoney, formatNumber } from '../../utils/format';
import DetailRow from './DetailRow';
import { EmptyState } from './States';
import StatusChip from './StatusChip';

const TYPE_ICONS = {
  topup: AddCardRounded,
  allocation: LocalGasStationRounded,
  station_transfer: SwapHorizRounded,
  fuel_purchase: LocalGasStationRounded,
  withdrawal: PaymentsRounded,
  cash_withdrawal: PaymentsRounded,
};

function TransactionRow({ transaction, onClick }) {
  const {
    type,
    status,
    from_name: from,
    to_name: to,
    branch_name: branch,
    reference,
    created_at: createdAt,
    amount,
    loyalty_points_used: pointsUsed,
    loyalty_discount: loyaltyDiscount,
  } = transaction;
  const Icon = TYPE_ICONS[type] ?? ReceiptLongRounded;
  const parties = from && to ? `من ${from} إلى ${to}` : from ? `من ${from}` : to ? `إلى ${to}` : '';

  return (
    <ButtonBase
      onClick={() => onClick(transaction)}
      sx={{
        width: '100%',
        display: 'block',
        textAlign: 'start',
        transition: 'background-color .2s',
        '&:hover': { bgcolor: tokens.redSoft },
      }}
    >
      <Stack direction="row" spacing={{ xs: 1.25, sm: 2 }} alignItems="center" sx={{ py: { xs: 1.5, sm: 1.75 }, px: { xs: 1.75, sm: 2.5 } }}>
        <Box
          sx={{
            width: { xs: 38, sm: 44 },
            height: { xs: 38, sm: 44 },
            flexShrink: 0,
            borderRadius: { xs: '12px', sm: '14px' },
            display: 'grid',
            placeItems: 'center',
            bgcolor: tokens.redSoft,
            color: 'primary.main',
          }}
        >
          <Icon fontSize="small" />
        </Box>
        <Box sx={{ flex: 1, minWidth: 0 }}>
          <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap>
            <Typography fontWeight={700} variant="body2" sx={{ fontSize: { sm: '1rem' } }}>
              {TRANSACTION_TYPES[type] ?? type}
            </Typography>
            {status !== 'completed' && <StatusChip status={status} />}
          </Stack>
          {(parties || branch) && (
            <Typography variant="body2" color="text.secondary" sx={{ fontSize: { xs: '0.8rem', sm: '0.875rem' }, overflowWrap: 'break-word' }}>
              {parties}
              {branch ? `${parties ? '، ' : ''}${branch}` : ''}
            </Typography>
          )}
          {pointsUsed > 0 && (
            <Typography variant="caption" color="primary.main" sx={{ display: 'block', fontWeight: 600 }}>
              تم استخدام {formatNumber(pointsUsed)} نقطة ولاء (خصم {formatMoney(loyaltyDiscount)})
            </Typography>
          )}
          <Typography variant="caption" color="text.secondary" sx={{ display: 'block', overflowWrap: 'break-word' }}>
            {reference} — {formatDateTime(createdAt)}
          </Typography>
        </Box>
        <Stack direction="row" spacing={0.5} alignItems="center" sx={{ flexShrink: 0 }}>
          <Typography fontWeight={800} variant="body2" sx={{ whiteSpace: 'nowrap', fontSize: { sm: '1rem' } }}>
            {amount === 0 && pointsUsed > 0 ? `${formatNumber(pointsUsed)} نقطة` : formatMoney(amount)}
          </Typography>
          <ChevronLeftRounded fontSize="small" color="action" />
        </Stack>
      </Stack>
    </ButtonBase>
  );
}

function TransactionDetailDialog({ transaction, loading, onClose }) {
  const Icon = transaction ? (TYPE_ICONS[transaction.type] ?? ReceiptLongRounded) : ReceiptLongRounded;
  const pointsUsed = transaction?.loyalty_points_used ?? 0;
  const loyaltyDiscount = transaction?.loyalty_discount ?? 0;
  const totalValue = (transaction?.amount ?? 0) + loyaltyDiscount;

  return (
    <Dialog open={Boolean(transaction)} onClose={onClose} fullWidth maxWidth="xs">
      {transaction && (
        <>
          <DialogContent sx={{ pt: 4 }}>
            {loading && <LinearProgress sx={{ mb: 2, borderRadius: 2 }} />}
            <Box sx={{ textAlign: 'center', mb: 2.5 }}>
              <Box
                sx={{
                  width: 56,
                  height: 56,
                  mx: 'auto',
                  mb: 1.5,
                  borderRadius: '50%',
                  display: 'grid',
                  placeItems: 'center',
                  bgcolor: tokens.redSoft,
                  color: 'primary.main',
                }}
              >
                <Icon />
              </Box>
              <Typography variant="subtitle1" fontWeight={700}>
                {TRANSACTION_TYPES[transaction.type] ?? transaction.type}
              </Typography>
              <Typography variant="h4" sx={{ mt: 0.5 }}>
                {pointsUsed > 0 ? formatMoney(totalValue) : formatMoney(transaction.amount)}
              </Typography>
            </Box>

            <DetailRow label="رقم المرجع" value={transaction.reference} />
            <DetailRow label="نوع العملية" value={TRANSACTION_TYPES[transaction.type] ?? transaction.type} />
            <DetailRow label="الحالة" value={<StatusChip status={transaction.status} />} />
            <DetailRow label="التاريخ والوقت" value={formatDateTime(transaction.created_at)} />
            {transaction.from_name && <DetailRow label="من" value={transaction.from_name} />}
            {transaction.to_name && <DetailRow label="إلى" value={transaction.to_name} />}
            {transaction.branch_name && <DetailRow label="الفرع" value={transaction.branch_name} />}
            {pointsUsed > 0 && (
              <>
                <DetailRow label="نقاط الولاء المستخدمة" value={`${formatNumber(pointsUsed)} نقطة`} />
                <DetailRow label="خصم نقاط الولاء" value={`-${formatMoney(loyaltyDiscount)}`} />
                <DetailRow label="المدفوع من المحفظة" value={formatMoney(transaction.amount)} />
              </>
            )}
          </DialogContent>
          <DialogActions sx={{ px: 3, pb: 2.5 }}>
            <Button fullWidth variant="contained" size="large" onClick={onClose}>
              إغلاق
            </Button>
          </DialogActions>
        </>
      )}
    </Dialog>
  );
}

const timestampOf = (value) => {
  const time = new Date(value).getTime();
  return Number.isFinite(time) ? time : 0;
};

export default function TransactionList({ transactions, emptyTitle = 'لا توجد معاملات بعد', emptyDescription }) {
  const { user } = useAuth();
  const [selected, setSelected] = useState(null);
  const [detailLoading, setDetailLoading] = useState(false);

  useEffect(() => {
    if (!selected?.id || user?.role !== ROLES.USER || selected.type === 'withdrawal') {
      setDetailLoading(false);
      return;
    }
    let cancelled = false;
    setDetailLoading(true);
    fetchTransferById(selected.id)
      .then((fresh) => {
        if (!cancelled && fresh) setSelected(fresh);
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setDetailLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [selected?.id, selected?.type, user?.role]);

  if (!transactions.length) return <EmptyState title={emptyTitle} description={emptyDescription} icon={ReceiptLongRounded} />;

  const ordered = [...transactions].sort(
    (a, b) => timestampOf(b.created_at) - timestampOf(a.created_at) || (b.id ?? 0) - (a.id ?? 0),
  );

  return (
    <>
      <Box>
        {ordered.map((transaction, index) => (
          <Fragment key={transaction.id}>
            {index > 0 && <Divider />}
            <TransactionRow transaction={transaction} onClick={setSelected} />
          </Fragment>
        ))}
      </Box>
      <TransactionDetailDialog transaction={selected} loading={detailLoading} onClose={() => setSelected(null)} />
    </>
  );
}
