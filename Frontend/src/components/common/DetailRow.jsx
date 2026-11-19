import { Stack, Typography } from '@mui/material';

export default function DetailRow({ label, value }) {
  return (
    <Stack
      direction="row"
      justifyContent="space-between"
      alignItems="center"
      spacing={2}
      sx={{ py: 1.25, borderBottom: '1px solid', borderColor: 'divider', '&:last-of-type': { borderBottom: 0 } }}
    >
      <Typography color="text.secondary" variant="body2" sx={{ flexShrink: 0 }}>
        {label}
      </Typography>
      <Typography component="div" fontWeight={700} sx={{ textAlign: 'end', minWidth: 0, overflowWrap: 'break-word' }}>
        {value}
      </Typography>
    </Stack>
  );
}
