import { Box, ButtonBase, Typography } from '@mui/material';
import AccountBalanceWalletRounded from '@mui/icons-material/AccountBalanceWalletRounded';
import AddCardRounded from '@mui/icons-material/AddCardRounded';
import SwapHorizRounded from '@mui/icons-material/SwapHorizRounded';
import { Link } from 'react-router-dom';
import { tokens } from '../../theme/tokens';

const ACTIONS = [
  { tab: 'topup', label: 'شحن المحفظة', icon: AddCardRounded },
  { tab: 'allocate', label: 'تخصيص لمحطة', icon: AccountBalanceWalletRounded },
  { tab: 'transfer', label: 'تحويل بين المحطات', icon: SwapHorizRounded },
];

export default function QuickActions() {
  return (
    <Box sx={{ display: 'grid', gap: { xs: 1, sm: 1.5 }, gridTemplateColumns: 'repeat(3, minmax(0, 1fr))' }}>
      {ACTIONS.map(({ tab, label, icon: Icon }) => (
        <ButtonBase
          key={tab}
          component={Link}
          to={`/user/transfers?tab=${tab}`}
          sx={{
            flexDirection: 'column',
            gap: 1,
            px: { xs: 1, sm: 2 },
            py: { xs: 1.75, sm: 2 },
            borderRadius: '18px',
            bgcolor: 'background.paper',
            border: '1px solid',
            borderColor: 'divider',
            transition: 'background-color .2s, border-color .2s',
            '&:hover': { bgcolor: tokens.redSoft, borderColor: 'primary.light' },
          }}
        >
          <Box sx={{ width: 44, height: 44, borderRadius: '50%', display: 'grid', placeItems: 'center', bgcolor: tokens.redSoft, color: 'primary.main' }}>
            <Icon />
          </Box>
          <Typography variant="subtitle2" sx={{ textAlign: 'center', fontSize: { xs: '0.78rem', sm: '0.875rem' }, lineHeight: 1.35 }}>
            {label}
          </Typography>
        </ButtonBase>
      ))}
    </Box>
  );
}
