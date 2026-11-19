import { useEffect, useState } from 'react';
import { Alert, Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField } from '@mui/material';
import DetailRow from '../../components/common/DetailRow';
import { bankLabel } from '../../config/constants';
import { formatMoney } from '../../utils/format';

/**
 * Shown before marking a withdrawal request as completed.
 * Allows recording the real external transfer reference from the banking app.
 */
export default function ExternalTransferDialog({ withdrawal, loading, onConfirm, onClose }) {
  const [externalRef, setExternalRef] = useState('');

  useEffect(() => {
    setExternalRef('');
  }, [withdrawal?.id]);

  return (
    <Dialog open={Boolean(withdrawal)} onClose={loading ? undefined : onClose} fullWidth maxWidth="xs">
      {withdrawal && (
        <>
          <DialogTitle sx={{ fontWeight: 800 }}>تنفيذ التحويل</DialogTitle>
          <DialogContent>
            <Alert severity="info" sx={{ mb: 2, borderRadius: '14px' }}>
              نفّذ التحويل في التطبيق المالي الخارجي، ثم أكّد إتمام الطلب هنا.
            </Alert>
            <DetailRow label="المحطة" value={withdrawal.station_name} />
            <DetailRow label="المبلغ" value={formatMoney(withdrawal.amount)} />
            <DetailRow label="البنك" value={bankLabel(withdrawal.bank_name)} />
            <DetailRow label="رقم الحساب" value={withdrawal.bank_account_number} />
            <DetailRow label="اسم صاحب الحساب" value={withdrawal.bank_account_name} />
            <TextField
              fullWidth
              margin="normal"
              label="رقم مرجع التحويل الخارجي (اختياري)"
              value={externalRef}
              onChange={(event) => setExternalRef(event.target.value)}
              disabled={loading}
            />
          </DialogContent>
          <DialogActions sx={{ px: 3, pb: 2.5 }}>
            <Button color="inherit" onClick={onClose} disabled={loading}>إلغاء</Button>
            <Button variant="contained" disabled={loading} onClick={() => onConfirm(withdrawal, externalRef.trim())}>
              تم التحويل، إكمال الطلب
            </Button>
          </DialogActions>
        </>
      )}
    </Dialog>
  );
}
