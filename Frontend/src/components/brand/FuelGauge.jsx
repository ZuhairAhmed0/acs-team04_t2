import { Box, useMediaQuery } from '@mui/material';

const CENTER = 120;
const FILL_PERCENT = 70;
const NEEDLE_SWEEP = (FILL_PERCENT / 100) * 180;
const TICK_ANGLES = Array.from({ length: 9 }, (_, index) => 180 - index * 22.5);

const polar = (radius, degrees) => {
  const radians = (degrees * Math.PI) / 180;
  return { x: CENTER + radius * Math.cos(radians), y: CENTER - radius * Math.sin(radians) };
};

/** Decorative fuel gauge: the needle and the filled arc sweep to 70% once on load. */
export default function FuelGauge({ width = 320 }) {
  const reduceMotion = useMediaQuery('(prefers-reduced-motion: reduce)');
  const arc = 'M 20 120 A 100 100 0 0 1 220 120';

  return (
    <Box component="svg" viewBox="0 0 240 150" role="img" aria-label="مؤشر مستوى الوقود" sx={{ width, maxWidth: '100%' }}>
      <path d={arc} fill="none" stroke="rgba(255,255,255,0.16)" strokeWidth="14" strokeLinecap="round" />
      <path d={arc} pathLength="100" fill="none" stroke="#fff" strokeWidth="14" strokeLinecap="round" strokeDasharray={`${reduceMotion ? FILL_PERCENT : 0} 100`}>
        {!reduceMotion && (
          <animate attributeName="stroke-dasharray" from="0 100" to={`${FILL_PERCENT} 100`} begin="0.4s" dur="1.8s" fill="freeze" calcMode="spline" keyTimes="0;1" keySplines="0.2 0.8 0.2 1" />
        )}
      </path>
      {TICK_ANGLES.map((angle) => {
        const start = polar(80, angle);
        const end = polar(90, angle);
        return <line key={angle} x1={start.x} y1={start.y} x2={end.x} y2={end.y} stroke="rgba(255,255,255,0.45)" strokeWidth="2" strokeLinecap="round" />;
      })}
      <g transform={`rotate(${reduceMotion ? NEEDLE_SWEEP : 0} ${CENTER} ${CENTER})`}>
        <line x1={CENTER} y1={CENTER} x2="46" y2={CENTER} stroke="#fff" strokeWidth="4" strokeLinecap="round" />
        {!reduceMotion && (
          <animateTransform attributeName="transform" type="rotate" from={`0 ${CENTER} ${CENTER}`} to={`${NEEDLE_SWEEP} ${CENTER} ${CENTER}`} begin="0.4s" dur="1.8s" fill="freeze" calcMode="spline" keyTimes="0;1" keySplines="0.2 0.8 0.2 1" />
        )}
      </g>
      <circle cx={CENTER} cy={CENTER} r="9" fill="#fff" />
      <text x="20" y="144" textAnchor="middle" fontSize="11" fill="rgba(255,255,255,0.75)" fontFamily="Cairo, sans-serif">فارغ</text>
      <text x="220" y="144" textAnchor="middle" fontSize="11" fill="rgba(255,255,255,0.75)" fontFamily="Cairo, sans-serif">ممتلئ</text>
    </Box>
  );
}
