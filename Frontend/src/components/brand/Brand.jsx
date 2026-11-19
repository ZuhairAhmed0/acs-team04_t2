import { Box, Stack, Typography } from '@mui/material';
import LocalGasStationRounded from '@mui/icons-material/LocalGasStationRounded';
import { tokens } from '../../theme/tokens';

export default function Brand({ light = false, small = false }) {
  const boxSize = small ? 34 : 42;
  return (
    <Stack direction="row" alignItems="center" spacing={1.25}>
      <Box
        sx={{
          width: boxSize,
          height: boxSize,
          borderRadius: '12px',
          display: 'grid',
          placeItems: 'center',
          color: '#fff',
          background: light ? 'rgba(255,255,255,0.16)' : tokens.gradientHero,
        }}
      >
        <LocalGasStationRounded fontSize={small ? 'small' : 'medium'} />
      </Box>
      <Typography component="span" sx={{ fontWeight: 800, fontSize: small ? '1.15rem' : '1.4rem', color: light ? '#fff' : 'text.primary' }}>
        وين الوقود؟
      </Typography>
    </Stack>
  );
}
