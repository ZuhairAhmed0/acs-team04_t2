import { Box, Card, Stack, Typography } from '@mui/material';
import { cityLabel } from '../../config/constants';
import useTone from '../../hooks/useTone';

const HIGH = 70;
const MEDIUM = 40;
const levelColor = (percentage) => (percentage >= HIGH ? 'error' : percentage >= MEDIUM ? 'warning' : 'success');

function CongestionBar({ branch }) {
  const percentage = branch.congestion_percentage;
  const tone = useTone(levelColor(percentage));

  return (
    <Box>
      <Stack direction="row" justifyContent="space-between" sx={{ mb: 0.75 }}>
        <Typography variant="body2" fontWeight={700}>{cityLabel(branch.city)}</Typography>
        <Typography variant="body2" fontWeight={800} sx={{ color: tone.dark }}>{percentage}%</Typography>
      </Stack>
      <Box role="progressbar" aria-valuenow={percentage} aria-valuemin={0} aria-valuemax={100} aria-label={`ازدحام فرع ${cityLabel(branch.city)}`} sx={{ height: 10, borderRadius: 5, bgcolor: tone.soft, overflow: 'hidden' }}>
        <Box sx={{ width: `${percentage}%`, height: '100%', borderRadius: 5, bgcolor: tone.main, transition: 'width .8s cubic-bezier(.2,.8,.2,1)' }} />
      </Box>
    </Box>
  );
}

export default function CongestionChart({ branches }) {
  return (
    <Card sx={{ p: 2.5, height: '100%' }}>
      <Typography variant="h6">نسبة الازدحام في الفروع</Typography>
      <Typography variant="body2" color="text.secondary" sx={{ mb: 2.5 }}>الأعلى نسبةً هو الأكثر ازدحاماً حالياً.</Typography>
      <Stack spacing={2.5}>
        {branches.map((branch) => <CongestionBar key={branch.id} branch={branch} />)}
      </Stack>
    </Card>
  );
}
