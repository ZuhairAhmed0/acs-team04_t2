import { useEffect, useState } from 'react';
import { Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, TextField } from '@mui/material';

/** Generic confirmation. With `requireReason`, `onConfirm` receives the typed reason. */
export default function ConfirmDialog({
  open,
  title,
  message,
  confirmLabel = 'تأكيد',
  tone = 'primary',
  requireReason = false,
  reasonLabel = 'سبب القرار',
  loading = false,
  onConfirm,
  onClose,
}) {
  const [reason, setReason] = useState('');

  useEffect(() => {
    if (open) setReason('');
  }, [open]);

  const canConfirm = !requireReason || reason.trim().length > 0;

  return (
    <Dialog open={open} onClose={loading ? undefined : onClose} fullWidth maxWidth="xs">
      <DialogTitle sx={{ fontWeight: 800 }}>{title}</DialogTitle>
      <DialogContent>
        <DialogContentText>{message}</DialogContentText>
        {requireReason && (
          <TextField autoFocus multiline minRows={2} label={reasonLabel} value={reason} onChange={(event) => setReason(event.target.value)} sx={{ mt: 2 }} />
        )}
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2.5 }}>
        <Button color="inherit" onClick={onClose} disabled={loading}>
          إلغاء
        </Button>
        <Button variant="contained" color={tone} disabled={!canConfirm || loading} onClick={() => onConfirm(reason.trim())}>
          {confirmLabel}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
