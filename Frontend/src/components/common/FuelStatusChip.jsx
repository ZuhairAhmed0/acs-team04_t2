import { Chip } from '@mui/material';
import FiberManualRecordRounded from '@mui/icons-material/FiberManualRecordRounded';
import { FUEL_STATUS } from '../../config/constants';
import useTone from '../../hooks/useTone';

export default function FuelStatusChip({ status, size = 'small' }) {
  const { label, color } = FUEL_STATUS[status];
  const tone = useTone(color);
  return (
    <Chip
      size={size}
      label={label}
      icon={<FiberManualRecordRounded sx={{ fontSize: 10 }} />}
      sx={{ bgcolor: tone.soft, color: tone.dark, '& .MuiChip-icon': { color: tone.main } }}
    />
  );
}
