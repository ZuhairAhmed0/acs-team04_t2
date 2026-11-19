import HomeRounded from '@mui/icons-material/HomeRounded';
import LocalGasStationRounded from '@mui/icons-material/LocalGasStationRounded';
import PaymentsRounded from '@mui/icons-material/PaymentsRounded';
import PeopleAltRounded from '@mui/icons-material/PeopleAltRounded';
import ReceiptLongRounded from '@mui/icons-material/ReceiptLongRounded';
import SpaceDashboardRounded from '@mui/icons-material/SpaceDashboardRounded';
import StorefrontRounded from '@mui/icons-material/StorefrontRounded';
import SwapHorizRounded from '@mui/icons-material/SwapHorizRounded';
import { ROLES } from './constants';

export const NAVIGATION = {
  [ROLES.USER]: [
    { to: '/user', label: 'الرئيسية', icon: HomeRounded, end: true },
    { to: '/user/stations', label: 'المحطات', icon: LocalGasStationRounded },
    { to: '/user/transfers', label: 'التحويلات', icon: SwapHorizRounded },
    { to: '/user/transactions', label: 'المعاملات', icon: ReceiptLongRounded },
  ],
  [ROLES.STATION]: [
    { to: '/station', label: 'لوحة المحطة', icon: SpaceDashboardRounded, end: true },
    { to: '/station/branches', label: 'الفروع', icon: StorefrontRounded },
    { to: '/station/transactions', label: 'المعاملات', icon: ReceiptLongRounded },
  ],
  [ROLES.ADMIN]: [
    { to: '/admin', label: 'لوحة التحكم', icon: SpaceDashboardRounded, end: true },
    { to: '/admin/users', label: 'المستخدمون', icon: PeopleAltRounded },
    { to: '/admin/stations', label: 'المحطات', icon: LocalGasStationRounded },
    { to: '/admin/withdrawals', label: 'طلبات السحب', icon: PaymentsRounded },
  ],
};

