import { Chip } from '@mui/material';
import { statusMeta } from '../../config/constants';
import useTone from '../../hooks/useTone';

export default function StatusChip({ status, size = 'small' }) {
  const { label, color } = statusMeta(status);
  const tone = useTone(color);
  return <Chip size={size} label={label} sx={{ bgcolor: tone.soft, color: tone.dark }} />;
}
