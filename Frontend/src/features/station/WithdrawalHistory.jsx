import { Fragment } from 'react';
import { Box, Card, Divider, Stack, Typography } from '@mui/material';
import PaymentsRounded from '@mui/icons-material/PaymentsRounded';
import { EmptyState } from '../../components/common/States';
import StatusChip from '../../components/common/StatusChip';
import { formatDate, formatMoney } from '../../utils/format';

export default function WithdrawalHistory({ withdrawals }) {
  const ordered = [...withdrawals].sort(
    (a, b) => new Date(b.created_at) - new Date(a.created_at) || (b.id ?? 0) - (a.id ?? 0),
  );

  return (
    <Card>
      <Typography variant="h6" sx={{ px: 2.5, pt: 2, pb: 0.5 }}>طلبات السحب</Typography>
      {!ordered.length && <EmptyState icon={PaymentsRounded} title="لا توجد طلبات سحب" description="الطلبات التي ترسلها ستظهر هنا مع حالتها." />}
      {ordered.map((withdrawal, index) => (
        <Fragment key={withdrawal.id}>
          {index > 0 && <Divider />}
          <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={2} sx={{ px: 2.5, py: 1.75 }}>
            <Box>
              <Typography fontWeight={800}>{formatMoney(withdrawal.amount)}</Typography>
              <Typography variant="caption" color="text.secondary">{formatDate(withdrawal.created_at)}</Typography>
              {withdrawal.external_reference && (
                <Typography variant="body2" color="text.secondary">مرجع التحويل: {withdrawal.external_reference}</Typography>
              )}
              {withdrawal.rejection_reason && (
                <Typography variant="body2" color="error.main">السبب: {withdrawal.rejection_reason}</Typography>
              )}
            </Box>
            <StatusChip status={withdrawal.status} />
          </Stack>
        </Fragment>
      ))}
    </Card>
  );
}
