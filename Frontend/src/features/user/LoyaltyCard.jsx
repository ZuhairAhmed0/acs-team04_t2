import { Box, Card, Chip, Stack, Typography } from '@mui/material';
import LocalGasStationRounded from '@mui/icons-material/LocalGasStationRounded';
import StarsRounded from '@mui/icons-material/StarsRounded';
import { tokens } from '../../theme/tokens';
import { formatNumber } from '../../utils/format';

const POINTS_PER_GALLON = 10;

/** Displays user loyalty points, free gas gallons (10 points = 1 gallon), and progress toward the next gallon. */
export default function LoyaltyCard({ points = 0, gallons }) {
  const safePoints = Math.max(0, Number(points) || 0);
  const freeGallons = typeof gallons === 'number' ? gallons : Math.floor(safePoints / POINTS_PER_GALLON);
  const towardNext = safePoints % POINTS_PER_GALLON;
  const remaining = POINTS_PER_GALLON - towardNext;
  const progress = (towardNext / POINTS_PER_GALLON) * 100;

  return (
    <Card sx={{ p: { xs: 2, sm: 2.5 } }}>
      <Stack direction="row" justifyContent="space-between" alignItems="flex-start" flexWrap="wrap" gap={1.5}>
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Box
            sx={{
              width: 44,
              height: 44,
              flexShrink: 0,
              borderRadius: '14px',
              display: 'grid',
              placeItems: 'center',
              bgcolor: tokens.redSoft,
              color: 'primary.main',
            }}
          >
            <StarsRounded />
          </Box>
          <Box>
            <Typography variant="body2" color="text.secondary">نقاط الولاء</Typography>
            <Typography variant="h5" fontWeight={800}>
              {formatNumber(safePoints)} نقطة
            </Typography>
          </Box>
        </Stack>

        <Chip
          icon={<LocalGasStationRounded />}
          color="primary"
          variant="outlined"
          label={`${formatNumber(freeGallons)} جالون وقود`}
          sx={{ fontWeight: 700 }}
        />
      </Stack>

      <Box sx={{ mt: 2 }}>
        <Stack direction="row" justifyContent="space-between" alignItems="center" flexWrap="wrap" gap={1} sx={{ mb: 0.75 }}>
          <Typography variant="caption" color="text.secondary">
            كل 10 نقاط = جالون وقود مجاني
          </Typography>
          <Typography variant="caption" fontWeight={700} color="primary.main">
            تبقّى {formatNumber(remaining)} نقاط للجالون التالي ({formatNumber(towardNext)}/10)
          </Typography>
        </Stack>
        <Box
          role="progressbar"
          aria-valuenow={towardNext}
          aria-valuemin={0}
          aria-valuemax={POINTS_PER_GALLON}
          aria-label="التقدم نحو الجالون المجاني التالي"
          sx={{ height: 8, borderRadius: 4, bgcolor: tokens.redSoft, overflow: 'hidden' }}
        >
          <Box sx={{ width: `${progress}%`, height: '100%', borderRadius: 4, bgcolor: 'primary.main', transition: 'width .4s ease' }} />
        </Box>
      </Box>

      <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1.5 }}>
        +10 نقاط عند التسجيل • +5 نقاط عند الشحن أو الشراء من المحفظة • الاستبدال بمضاعفات 10 نقاط (10، 20، 30... = 1,000 ج.س للجالون)
      </Typography>
    </Card>
  );
}
