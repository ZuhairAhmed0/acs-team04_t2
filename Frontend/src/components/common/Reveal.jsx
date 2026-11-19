import { Box } from '@mui/material';

/** Scroll-reveal wrapper powered by AOS. `delay` must be a multiple of 50 (ms). */
export default function Reveal({ children, delay = 0, ...boxProps }) {
  return (
    <Box data-aos="fade-up" data-aos-delay={delay} {...boxProps}>
      {children}
    </Box>
  );
}
