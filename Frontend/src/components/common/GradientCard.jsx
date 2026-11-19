import { Box, Card } from '@mui/material';
import { tokens } from '../../theme/tokens';

/** Concentric arcs echo the fuel-gauge motif used on the sign-in screen. */
function Arcs() {
  return (
    <Box
      component="svg"
      viewBox="0 0 320 320"
      aria-hidden
      sx={{ position: 'absolute', insetInlineStart: -90, bottom: -130, width: 320, height: 320, pointerEvents: 'none' }}
    >
      {[150, 110, 70].map((radius) => (
        <circle key={radius} cx="160" cy="160" r={radius} fill="none" stroke="rgba(255,255,255,0.12)" strokeWidth="1.5" />
      ))}
    </Box>
  );
}

export default function GradientCard({ children, sx, ...props }) {
  return (
    <Card
      {...props}
      sx={{ position: 'relative', overflow: 'hidden', border: 0, color: '#fff', background: tokens.gradientHero, p: { xs: 2.5, sm: 3.5 }, ...sx }}
    >
      <Arcs />
      <Box sx={{ position: 'relative' }}>{children}</Box>
    </Card>
  );
}
