import { useState } from 'react';
import { Box, IconButton, Stack, Typography } from '@mui/material';
import VisibilityOffRounded from '@mui/icons-material/VisibilityOffRounded';
import VisibilityRounded from '@mui/icons-material/VisibilityRounded';
import GradientCard from '../../components/common/GradientCard';
import { formatMoney } from '../../utils/format';

const MASK = '••••••';

function Legend({ label, value, solid }) {
  return (
    <Box>
      <Stack direction="row" alignItems="center" spacing={1}>
        <Box sx={{ width: 10, height: 10, borderRadius: '50%', bgcolor: solid ? '#fff' : 'rgba(255,255,255,0.35)' }} />
        <Typography variant="body2" sx={{ opacity: 0.85 }}>{label}</Typography>
      </Stack>
      <Typography variant="h6" sx={{ mt: 0.25 }}>{value}</Typography>
    </Box>
  );
}

/** Total balance with the split between allocated and available funds. */
export default function WalletCard({ wallet }) {
  const [hidden, setHidden] = useState(false);
  const { total_balance: total, allocated_balance: allocated, available_balance: available } = wallet;
  const allocatedShare = total > 0 ? (allocated / total) * 100 : 0;
  const show = (value) => (hidden ? MASK : formatMoney(value));

  return (
    <GradientCard>
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Typography sx={{ opacity: 0.85 }}>إجمالي الرصيد</Typography>
        <IconButton onClick={() => setHidden((prev) => !prev)} aria-label={hidden ? 'إظهار الرصيد' : 'إخفاء الرصيد'} sx={{ color: '#fff' }}>
          {hidden ? <VisibilityOffRounded /> : <VisibilityRounded />}
        </IconButton>
      </Stack>
      <Typography variant="h3" sx={{ my: 1, fontSize: { xs: '2.1rem', sm: '2.9rem' } }}>{show(total)}</Typography>

      <Box
        role="img"
        aria-label={`المخصص للمحطات ${Math.round(allocatedShare)}% من الإجمالي`}
        sx={{ mt: 3, height: 8, borderRadius: 4, overflow: 'hidden', bgcolor: 'rgba(255,255,255,0.25)' }}
      >
        <Box sx={{ width: `${allocatedShare}%`, height: '100%', bgcolor: '#fff', borderRadius: 4 }} />
      </Box>
      <Stack direction="row" justifyContent="space-between" flexWrap="wrap" gap={1.5} sx={{ mt: 2 }}>
        <Legend solid label="المخصص للمحطات" value={show(allocated)} />
        <Legend label="المتاح" value={show(available)} />
      </Stack>
    </GradientCard>
  );
}
