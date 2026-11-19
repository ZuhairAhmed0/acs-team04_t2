import { Alert, Box, Button, Dialog, DialogActions, DialogContent, Typography } from '@mui/material';
import CheckRounded from '@mui/icons-material/CheckRounded';
import { formatMoney, formatNumber } from '../../utils/format';
import DetailRow from './DetailRow';
import StatusChip from './StatusChip';

/** Simple notification shown after a transfer. The branch row appears only when the API returns one. */
export default function ReceiptDialog({ receipt, onClose }) {
  const duplicateNotice = receipt?.idempotent_replay
    ? 'تنبيه: تم إرسال هذا الطلب مسبقاً، وعرضت النتيجة السابقة دون خصم مبلغ إضافي.'
    : receipt?.duplicate_purchase
      ? `تنبيه: هذه عملية شراء مطابقة لعملية سابقة حديثة${receipt.previous_reference ? ` (المرجع السابق: ${receipt.previous_reference})` : ''}.`
      : '';

  return (
    <Dialog open={Boolean(receipt)} onClose={onClose} fullWidth maxWidth="xs">
      {receipt && (
        <>
          <DialogContent sx={{ pt: 4 }}>
            <Box sx={{ textAlign: 'center', mb: 2 }}>
              <Box sx={{ width: 56, height: 56, mx: 'auto', mb: 1.5, borderRadius: '50%', display: 'grid', placeItems: 'center', bgcolor: 'success.main', color: '#fff' }}>
                <CheckRounded />
              </Box>
              <Typography variant="subtitle1">تمت العملية بنجاح</Typography>
              <Typography variant="h4" sx={{ mt: 0.5 }}>
                {formatMoney(receipt.total_amount || receipt.amount)}
              </Typography>
            </Box>
            {duplicateNotice && (
              <Alert severity="warning" sx={{ mb: 2, borderRadius: '14px' }}>
                {duplicateNotice}
              </Alert>
            )}
            <DetailRow label="رقم المرجع" value={receipt.reference} />
            {receipt.from_name && <DetailRow label="من" value={receipt.from_name} />}
            {receipt.to_name && <DetailRow label="إلى" value={receipt.to_name} />}
            {receipt.branch_name && <DetailRow label="الفرع" value={receipt.branch_name} />}
            {receipt.loyalty_points_used > 0 && (
              <>
                <DetailRow
                  label="خصم نقاط الولاء"
                  value={`-${formatMoney(receipt.loyalty_discount)} (${formatNumber(receipt.loyalty_points_used)} نقطة)`}
                />
                <DetailRow label="المدفوع من المحفظة" value={formatMoney(receipt.amount)} />
              </>
            )}
            <DetailRow label="الحالة" value={<StatusChip status={receipt.status} />} />
            {receipt.points_earned > 0 && (
              <DetailRow label="نقاط الولاء المكتسبة" value={`+${formatNumber(receipt.points_earned)} نقاط`} />
            )}
            {typeof receipt.loyalty_points === 'number' && (
              <DetailRow
                label="إجمالي نقاط الولاء"
                value={`${formatNumber(receipt.loyalty_points)} نقطة (${formatNumber(Math.floor(receipt.loyalty_points / 10))} جالون)`}
              />
            )}
          </DialogContent>
          <DialogActions sx={{ px: 3, pb: 2.5 }}>
            <Button fullWidth variant="contained" size="large" onClick={onClose}>
              تم
            </Button>
          </DialogActions>
        </>
      )}
    </Dialog>
  );
}
