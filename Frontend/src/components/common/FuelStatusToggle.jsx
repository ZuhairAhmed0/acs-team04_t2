import { Chip, Stack } from '@mui/material';
import { FUEL_STATUS } from '../../config/constants';
import useTone from '../../hooks/useTone';

function FuelOption({ label, color, selected, disabled, onSelect }) {
  const tone = useTone(color);
  return (
    <Chip
      clickable
      role="radio"
      aria-checked={selected}
      label={label}
      disabled={disabled}
      onClick={onSelect}
      variant={selected ? 'filled' : 'outlined'}
      sx={
        selected
          ? { bgcolor: tone.main, color: '#fff', '&:hover': { bgcolor: tone.dark } }
          : { borderColor: tone.main, color: tone.dark }
      }
    />
  );
}

export default function FuelStatusToggle({ value, onChange, disabled = false }) {
  return (
    <Stack direction="row" flexWrap="wrap" useFlexGap spacing={1} role="radiogroup" aria-label="حالة الوقود">
      {Object.entries(FUEL_STATUS).map(([key, meta]) => (
        <FuelOption
          key={key}
          label={meta.label}
          color={meta.color}
          selected={key === value}
          disabled={disabled}
          onSelect={() => key !== value && onChange(key)}
        />
      ))}
    </Stack>
  );
}
