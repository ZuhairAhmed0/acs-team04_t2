import { Chip } from '@mui/material';
import TrafficRounded from '@mui/icons-material/TrafficRounded';
import useTone from '../../hooks/useTone';

const HIGH = 70;
const MEDIUM = 40;
const levelColor = (percentage) => (percentage >= HIGH ? 'error' : percentage >= MEDIUM ? 'warning' : 'success');

export default function CongestionChip({ percentage = 0, label = 'نسبة الازدحام', size = 'small' }) {
  const safeValue = Number.isFinite(percentage) ? percentage : 0;
  const tone = useTone(levelColor(safeValue));

  return (
    <Chip
      size={size}
      icon={<TrafficRounded sx={{ fontSize: 14 }} />}
      label={`${label}: ${safeValue}%`}
      sx={{
        bgcolor: tone.soft,
        color: tone.dark,
        fontWeight: 600,
        '& .MuiChip-icon': { color: tone.main },
      }}
    />
  );
}
