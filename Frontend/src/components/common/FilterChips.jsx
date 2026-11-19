import { Chip, Stack } from '@mui/material';

/**
 * Responsive filter chip bar: horizontally swipeable on narrow mobile screens
 * and wrapped across rows on wider screens.
 */
export default function FilterChips({ value, options, allLabel = 'الكل', ariaLabel, onChange }) {
  return (
    <Stack
      direction="row"
      spacing={1}
      role="group"
      aria-label={ariaLabel}
      sx={{
        mb: { xs: 2.5, sm: 3 },
        flexWrap: { xs: 'nowrap', sm: 'wrap' },
        overflowX: { xs: 'auto', sm: 'visible' },
        gap: 1,
        pb: { xs: 0.5, sm: 0 },
        scrollbarWidth: 'none',
        '&::-webkit-scrollbar': { display: 'none' },
        '& > *': { marginInlineStart: '0 !important', flexShrink: 0 },
      }}
    >
      <Chip
        label={allLabel}
        clickable
        color={value ? 'default' : 'primary'}
        variant={value ? 'outlined' : 'filled'}
        onClick={() => onChange('')}
      />
      {options.map((option) => (
        <Chip
          key={option.value}
          label={option.label}
          clickable
          color={value === option.value ? 'primary' : 'default'}
          variant={value === option.value ? 'filled' : 'outlined'}
          onClick={() => onChange(option.value)}
        />
      ))}
    </Stack>
  );
}
