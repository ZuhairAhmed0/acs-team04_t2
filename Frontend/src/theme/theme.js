import { createTheme } from '@mui/material/styles';
import { tokens } from './tokens';

const theme = createTheme({
  direction: 'rtl',
  palette: {
    mode: 'light',
    primary: { main: tokens.red, dark: tokens.redDark, light: '#E64A61', contrastText: '#fff' },
    secondary: { main: tokens.redDeep },
    success: { main: tokens.success },
    warning: { main: tokens.warning },
    error: { main: '#B42318' },
    info: { main: '#56677A' },
    background: { default: tokens.canvas, paper: tokens.surface },
    text: { primary: tokens.ink, secondary: tokens.muted },
    divider: tokens.line,
  },
  shape: { borderRadius: 14 },
  typography: {
    fontFamily: '"Cairo", system-ui, sans-serif',
    h3: { fontWeight: 800 },
    h4: {
      fontWeight: 800,
      fontSize: '1.375rem',
      '@media (min-width:600px)': { fontSize: '1.65rem' },
    },
    h5: {
      fontWeight: 800,
      fontSize: '1.2rem',
      '@media (min-width:600px)': { fontSize: '1.35rem' },
    },
    h6: { fontWeight: 700, fontSize: '1.05rem', '@media (min-width:600px)': { fontSize: '1.1rem' } },
    subtitle1: { fontWeight: 700 },
    subtitle2: { fontWeight: 600 },
    button: { textTransform: 'none', fontWeight: 700 },
  },
  components: {
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: { borderRadius: 12, paddingBlock: 9, minHeight: 40 },
        sizeLarge: { paddingBlock: 12, minHeight: 48, fontSize: '1rem' },
      },
    },
    MuiIconButton: {
      styleOverrides: {
        root: { minWidth: 42, minHeight: 42 },
      },
    },
    MuiPaper: { styleOverrides: { root: { backgroundImage: 'none' } } },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 20,
          border: `1px solid ${tokens.line}`,
          boxShadow: tokens.shadowSoft,
        },
      },
    },
    MuiTextField: { defaultProps: { fullWidth: true } },
    MuiOutlinedInput: {
      styleOverrides: {
        root: { borderRadius: 12, backgroundColor: '#fff', '& fieldset': { borderColor: '#E7D8D6' } },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: { fontWeight: 700 },
        clickable: { minHeight: 36 },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          borderRadius: 20,
          margin: 16,
          width: 'calc(100% - 32px)',
          maxHeight: 'calc(100dvh - 32px)',
          '@media (min-width:600px)': { borderRadius: 24, margin: 32 },
        },
      },
    },
    MuiTabs: { styleOverrides: { indicator: { height: 3, borderRadius: 3 } } },
    MuiTab: { styleOverrides: { root: { fontWeight: 700, minHeight: 48 } } },
    MuiToggleButton: { styleOverrides: { root: { fontWeight: 700, borderRadius: 12 } } },
  },
});

export default theme;
