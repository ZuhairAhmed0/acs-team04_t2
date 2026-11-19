import { Box, Stack, Typography } from '@mui/material';
import CheckRounded from '@mui/icons-material/CheckRounded';
import { tokens } from '../../theme/tokens';
import Brand from '../brand/Brand';
import FuelGauge from '../brand/FuelGauge';

const HIGHLIGHTS = ['اعرف المحطات وحالة توفر الوقود', 'خصص رصيدك وحوّله بين المحطات بسهولة', 'أدر فروع محطتك وطلبات السحب'];

function BrandPanel() {
  return (
    <Box
      sx={{
        display: { xs: 'none', md: 'flex' },
        position: 'relative',
        overflow: 'hidden',
        flexDirection: 'column',
        justifyContent: 'center',
        gap: 5,
        p: { md: 6, lg: 9 },
        color: '#fff',
        background: tokens.gradientHero,
      }}
    >
      <Box
        component="svg"
        viewBox="0 0 400 400"
        aria-hidden
        sx={{ position: 'absolute', insetInlineStart: -140, bottom: -160, width: 440, height: 440 }}
      >
        {[190, 150, 110].map((radius) => (
          <circle key={radius} cx="200" cy="200" r={radius} fill="none" stroke="rgba(255,255,255,0.1)" strokeWidth="1.5" />
        ))}
      </Box>

      <Box sx={{ position: 'relative' }}>
        <FuelGauge />
      </Box>

      <Box sx={{ position: 'relative', maxWidth: 480 }}>
        <Typography variant="h3" sx={{ fontSize: { md: '2.2rem', lg: '2.8rem' }, lineHeight: 1.45 }}>
          وقودك أقرب،
          <br />
          ومحفظتك أذكى.
        </Typography>
        <Typography sx={{ mt: 2, opacity: 0.85, lineHeight: 2 }}>
          منصة لتنظيم شراء الوقود وإدارة الأرصدة بين العملاء ومحطات الوقود، من مكان واحد.
        </Typography>
        <Stack spacing={1.5} sx={{ mt: 3 }}>
          {HIGHLIGHTS.map((text) => (
            <Stack key={text} direction="row" spacing={1.5} alignItems="center">
              <Box sx={{ width: 24, height: 24, borderRadius: '50%', display: 'grid', placeItems: 'center', bgcolor: 'rgba(255,255,255,0.18)' }}>
                <CheckRounded sx={{ fontSize: 15 }} />
              </Box>
              <Typography variant="body2">{text}</Typography>
            </Stack>
          ))}
        </Stack>
      </Box>
    </Box>
  );
}

export default function AuthLayout({ children }) {
  return (
    <Box sx={{ minHeight: '100dvh', display: 'grid', gridTemplateColumns: { xs: '1fr', md: 'repeat(2, minmax(0, 1fr))' } }}>
      <Box sx={{ display: 'flex', alignItems: { xs: 'flex-start', md: 'center' }, justifyContent: 'center', p: { xs: 2, sm: 5 }, py: { xs: 3, sm: 6 } }}>
        <Box sx={{ width: '100%', maxWidth: 440 }}>
          <Box sx={{ mb: { xs: 3, sm: 4 } }}>
            <Brand />
          </Box>
          {children}
        </Box>
      </Box>
      <BrandPanel />
    </Box>
  );
}
