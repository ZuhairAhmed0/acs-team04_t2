import { alpha, useTheme } from '@mui/material/styles';

/** Resolves a semantic palette colour ("success", "error"...) into main / dark / soft shades. */
export default function useTone(color) {
  const { palette } = useTheme();
  const tone = palette[color];
  const main = tone?.main ?? palette.grey[600];
  return { main, dark: tone?.dark ?? palette.grey[800], soft: alpha(main, 0.12) };
}
