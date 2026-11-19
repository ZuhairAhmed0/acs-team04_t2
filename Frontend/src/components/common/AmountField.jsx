import { InputAdornment, TextField } from '@mui/material';

export default function AmountField({ label = 'المبلغ', min = 1, ...props }) {
  return (
    <TextField
      type="number"
      label={label}
      placeholder="مثلاً 500000"
      slotProps={{
        htmlInput: { min, step: 1, inputMode: 'numeric' },
        input: { endAdornment: <InputAdornment position="end">ج.س</InputAdornment> },
      }}
      {...props}
    />
  );
}
