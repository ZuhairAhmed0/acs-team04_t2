import { useMemo, useState } from 'react';
import { Alert, Box, Card, Chip, Stack, Tab, Tabs, Typography } from '@mui/material';
import StarsRounded from '@mui/icons-material/StarsRounded';
import { useSearchParams } from 'react-router-dom';
import PageHeader from '../../components/common/PageHeader';
import ReceiptDialog from '../../components/common/ReceiptDialog';
import Reveal from '../../components/common/Reveal';
import { ROLES } from '../../config/constants';
import { useAuth } from '../../context/AuthContext';
import useRequest from '../../hooks/useRequest';
import { fetchUserStations, fetchUserWallet } from '../../services/user.service';
import { formatNumber } from '../../utils/format';
import { AllocateForm, StationTransferForm, TopUpForm } from './TransferForms';

const TABS = [
  { value: 'topup', label: 'شحن المحفظة', hint: 'أضف رصيداً إلى محفظتك.' },
  { value: 'allocate', label: 'تخصيص لمحطة', hint: 'يُضاف المبلغ إلى رصيد المحطة المشترك ويصلح لكل فروعها.' },
  { value: 'transfer', label: 'تحويل بين المحطات', hint: 'انقل رصيدك المخصص من محطة إلى أخرى. لا يلزم اختيار فرع.' },
];

/** The stations endpoint lists branches; the selects need each station once. */
const uniqueStations = (rows) => [...new Map((rows ?? []).map(({ station_id: id, station_name: name }) => [id, { id, name }])).values()];

export default function UserTransfers() {
  const { user } = useAuth();
  const [searchParams, setSearchParams] = useSearchParams();
  const [receipt, setReceipt] = useState(null);
  const requestedTab = searchParams.get('tab');
  const activeTab = TABS.some((tab) => tab.value === requestedTab) ? requestedTab : TABS[0].value;
  const { data: rows, error } = useRequest(() => fetchUserStations());
  const { data: wallet, reload: reloadWallet } = useRequest(() => fetchUserWallet());
  const stations = useMemo(() => uniqueStations(rows), [rows]);
  const { hint } = TABS.find((tab) => tab.value === activeTab);

  const handleSuccess = (nextReceipt) => {
    setReceipt(nextReceipt);
    reloadWallet();
  };

  return (
    <>
      <PageHeader title="التحويلات" subtitle="اشحن محفظتك، خصّص رصيداً لمحطة، أو حوّل بين المحطات." />
      {user?.role === ROLES.USER && wallet && (
        <Stack direction="row" spacing={1} alignItems="center" flexWrap="wrap" useFlexGap sx={{ mb: 2 }}>
          <Chip
            icon={<StarsRounded />}
            color="primary"
            variant="outlined"
            label={`نقاط الولاء: ${formatNumber(wallet.loyalty_points)} نقطة (${formatNumber(wallet.free_gas_gallons)} جالون وقود)`}
          />
          <Typography variant="caption" color="text.secondary">
            +5 نقاط بعد كل شحن للمحفظة أو شراء وقود (10 نقاط = جالون وقود)
          </Typography>
        </Stack>
      )}
      <Reveal>
        <Card sx={{ maxWidth: 640, width: '100%' }}>
          <Tabs
            value={activeTab}
            onChange={(_, value) => setSearchParams({ tab: value })}
            variant="scrollable"
            scrollButtons="auto"
            allowScrollButtonsMobile
            sx={{ px: { xs: 0.5, sm: 1 }, borderBottom: '1px solid', borderColor: 'divider' }}
          >
            {TABS.map(({ value, label }) => <Tab key={value} value={value} label={label} />)}
          </Tabs>
          <Typography color="text.secondary" variant="body2" sx={{ px: { xs: 2, sm: 3 }, pt: { xs: 2, sm: 2.5 } }}>{hint}</Typography>
          {error && <Alert severity="error" sx={{ m: { xs: 2, sm: 3 }, mb: 0, borderRadius: '14px' }}>تعذّر تحميل قائمة المحطات</Alert>}
          <Box sx={{ p: { xs: 2, sm: 3 } }}>
            {activeTab === 'topup' && <TopUpForm onSuccess={handleSuccess} />}
            {activeTab === 'allocate' && <AllocateForm stations={stations} availableBalance={wallet?.available_balance} onSuccess={handleSuccess} />}
            {activeTab === 'transfer' && <StationTransferForm stations={stations} allocatedBalance={wallet?.allocated_balance} onSuccess={handleSuccess} />}
          </Box>
        </Card>
      </Reveal>
      <ReceiptDialog receipt={receipt} onClose={() => setReceipt(null)} />
    </>
  );
}
