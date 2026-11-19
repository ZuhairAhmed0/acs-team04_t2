import { Button, Dialog, DialogActions, DialogContent, DialogTitle, MenuItem, Stack, TextField } from '@mui/material';
import { CITIES } from '../../config/constants';
import { useToast } from '../../context/ToastContext';
import useFormState from '../../hooks/useFormState';
import useSubmit from '../../hooks/useSubmit';
import { addBranch } from '../../services/station.service';

/** A branch needs only a city and an address; the station name and bank details are inherited. */
export default function AddBranchDialog({ open, onClose, onAdded }) {
  const toast = useToast();
  const { values, errors, handleChange, validateRequired, reset } = useFormState({ city: '', address: '' });
  const { submit, submitting } = useSubmit(addBranch);

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!validateRequired(['city', 'address'])) return;
    const branch = await submit({ city: values.city, address: values.address.trim() });
    if (branch) {
      toast.success('تمت إضافة الفرع');
      handleClose();
      onAdded();
    }
  };

  return (
    <Dialog open={open} onClose={submitting ? undefined : handleClose} fullWidth maxWidth="xs" component="form" onSubmit={handleSubmit} noValidate>
      <DialogTitle sx={{ fontWeight: 800 }}>إضافة فرع جديد</DialogTitle>
      <DialogContent>
        <Stack spacing={2} sx={{ pt: 1 }}>
          <TextField select name="city" label="المدينة" value={values.city} onChange={handleChange} error={Boolean(errors.city)} helperText={errors.city}>
            {CITIES.map(({ value, label }) => <MenuItem key={value} value={value}>{label}</MenuItem>)}
          </TextField>
          <TextField name="address" label="العنوان" value={values.address} onChange={handleChange} error={Boolean(errors.address)} helperText={errors.address} />
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2.5 }}>
        <Button color="inherit" onClick={handleClose} disabled={submitting}>إلغاء</Button>
        <Button type="submit" variant="contained" disabled={submitting}>حفظ الفرع</Button>
      </DialogActions>
    </Dialog>
  );
}
